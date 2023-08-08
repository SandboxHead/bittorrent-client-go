package entity

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

const MaxBlockSize = 16384
const MaxBackLog = 5

type Client struct {
	Peer Peer
	Choked bool
	peerId [20]byte
	infoHash [20]byte
	Conn net.Conn
	BitField BitField
}


func (client *Client) doHandShake() (error) {
	client.Conn.SetDeadline(time.Now().Add(5*time.Second))
	defer client.Conn.SetDeadline(time.Time{})

	handShake := HandShake{
		Pstr: "BitTorrent protocol",
		InfoHash: client.infoHash,
		PeerID: client.peerId,
	}

	message := handShake.Serialize()

	_, err := client.Conn.Write(message)
	if err != nil {
		fmt.Errorf("error sending message")
		return err
	}

	response := make([]byte, 1024)
	n, err := client.Conn.Read(response)
	if err != nil {
		fmt.Errorf("error receiving message")
		return err
	}
	handShakeResponse, err := Deserialize(bytes.NewReader(response[:n]))
	if err != nil {
		return err
	}
	if bytes.Equal(handShake.InfoHash[:], handShakeResponse.InfoHash[:]) {
		return nil
	} else {
		err = fmt.Errorf("InfoHash doesn't match in the handshake")
		return err
	}

}

func New(peer Peer, infoHash, peerId [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
	if err != nil {
		return nil, err 
	}
	client := Client{
		Peer: peer,
		peerId: peerId,
		infoHash: infoHash,
		Conn: conn,
		Choked: true,
	}
	err = client.doHandShake()
	if err != nil {
		client.Conn.Close()
		return nil, err
	}
	err = client.setBitField()
	if err != nil {
		client.Conn.Close()
		return nil, err
	}
	return &client, nil
}

func (client *Client) setBitField() (error) {
	client.Conn.SetDeadline(time.Now().Add(5*time.Second))
	defer client.Conn.SetDeadline(time.Time{})

	msg, err := Read(client.Conn)
	if err != nil {
		return err
	}
	if msg == nil {
		return fmt.Errorf("Expected BitField but got %s", msg)
	}
	if msg.ID != MsgBitfield {
		return fmt.Errorf("Expected BitField but got %s", msg)
	}
	client.BitField = msg.Payload
	return nil
}

func (client *Client) Read() (*Message, error) {
	msg, err := Read(client.Conn)
	if err != nil {
		return nil, err
	}
	return msg, err
	
}

func (client *Client) DownloadPiece(pw *PieceWork) ([]byte, error) {
	client.Conn.SetDeadline(time.Now().Add(30*time.Second))
	defer client.Conn.SetDeadline(time.Time{})

	downloaded := 0
	lengthDownloaded := 0
	currBacklog := 0
	requested := 0
	buf := make([]byte, pw.Length)
	for lengthDownloaded < int(pw.Length) {
		if !client.Choked {
			for currBacklog < MaxBackLog && requested < int(pw.Length) {
				blockSize := MaxBlockSize
				if int(pw.Length) - requested < blockSize {
					blockSize = int(pw.Length) - requested
				}

				err := client.sendRequest(int(pw.Index), requested, blockSize)
				if err != nil {
					return nil, err
				}
				currBacklog ++
				requested += blockSize
			}
		}
		n, err := client.readMessage(int(pw.Index), buf)
		
		if err != nil {
			fmt.Println("Error in readMessage")
			return nil, err 
		}
		if n == 0 {
			continue
		}
		downloaded += n
		currBacklog --
	}	
	return buf, nil
}

func (client *Client) sendRequest(index int, requested int, blockSize int) (error) {
	return nil
}

func (client *Client) readMessage(index int, buf []byte) (int, error) {
	msg, err := client.Read()
	n := 0
	if err != nil {
		return 0, err
	}
	switch msg.ID {
	case MsgUnchoke:
		client.Choked = false
	case MsgChoke:
		client.Choked = true
	case MsgHave:
		index, _ := msg.ParseHave()
		client.BitField.SetPiece(index)
	case MsgPiece:
		n, err = msg.ParsePiece(index, buf)
		if err != nil {
			return 0, err 
		}
	}
	return n, nil
}

func (client *Client) SendUnchoke() (error) {
	message := Message{ID: MsgUnchoke}
	_, err := client.Conn.Write(message.Serialize())
	return err
}

func (client *Client) SendInterested() (error) {
	message := Message{ID: MsgInterested}
	_, err := client.Conn.Write(message.Serialize())
	return err
}

func (client *Client) SendHave(index int) (error) {
	message := FormatHave(index)
	_, err := client.Conn.Write(message.Serialize())
	return err
}


