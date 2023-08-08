package entity

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Peer struct {
	IP net.IP
	Port uint16
}


func (peer Peer) String() string {
	return fmt.Sprintf("%v:%d", peer.IP, peer.Port)
}

func UnmarshalPeer(peerBin []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peerBin) / peerSize
	if len(peerBin)%peerSize != 0 {
		err := fmt.Errorf("received malformed peers")
		return nil, err 
	}

	peers := make([]Peer, numPeers)
	for i:=0; i<numPeers; i++ {
		offset := i*peerSize
		peers[i].IP = net.IPv4(peerBin[offset], peerBin[offset + 1], peerBin[offset + 2], peerBin[offset + 3])
		peers[i].Port = binary.BigEndian.Uint16(peerBin[offset + 4 : offset + 6])
	}
	return peers, nil
}

func (p *Peer) DoHandShake(infoHash [20]byte, peerId [20]byte) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", p.String(), 5*time.Second)
	if err != nil {
		return nil, err
	}

	handShake := HandShake{
		Pstr: "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID: peerId,
	}

	message := handShake.Serialize()

	_, err = conn.Write(message)
	if err != nil {
		fmt.Errorf("error sending message")
		conn.Close()
		return nil, err
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Errorf("error receiving message")
		conn.Close()
		return nil, err
	}
	handShakeResponse, err := Deserialize(bytes.NewReader(response[:n]))
	if err != nil {
		conn.Close()
		return nil, err
	}
	if bytes.Equal(handShake.InfoHash[:], handShakeResponse.InfoHash[:]) {
		return conn, nil
	} else {
		err = fmt.Errorf("InfoHash doesn't match in the handshake")
		conn.Close()
		return nil, err
	}
} 


