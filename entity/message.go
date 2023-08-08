package entity

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	MsgChoke			messageID = 0
	MsgUnchoke      	messageID = 1
    MsgInterested    	messageID = 2
    MsgNotInterested	messageID = 3
    MsgHave          	messageID = 4
    MsgBitfield      	messageID = 5
    MsgRequest       	messageID = 6
    MsgPiece         	messageID = 7
    MsgCancel        	messageID = 8
)

type Message struct {
	ID messageID
	Payload []byte 
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4 + length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthbuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthbuf)
	if err != nil {
		return nil, err 
	}
	length := binary.BigEndian.Uint32(lengthbuf)
	if length == 0 {
		return nil, nil
	}

	messageBuff := make([]byte, length)
	_, err = io.ReadFull(r, messageBuff)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID: messageID(messageBuff[0]),
		Payload: messageBuff[1:],
	}
	return &m, nil
}

func (msg *Message) ParseHave() (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("Expected msg Id for Have as %d. Got %d", MsgHave, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload of size 4, Got %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func (msg *Message) ParsePiece(index int, buf []byte) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("Expected msg Id for Piece as %d. Got %d", MsgPiece, msg.ID)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Expected payload of size greater than 8, Got %d", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if len(data) + begin > len(buf) {
		return 0, fmt.Errorf("Length of buffer is not enough")
	}
	copy(buf[begin:], data)
	return len(data), nil
}

func FormatHave(index int) (*Message) {
	payload := make([]byte, 4)
	binary.BigEndian.AppendUint32(payload, uint32(index)) 
	msg := Message{ID: MsgHave, Payload: payload}
	return &msg
} 