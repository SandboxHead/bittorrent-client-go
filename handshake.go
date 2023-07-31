package main

import "io"

type HandShake struct {
	Pstr 		string
	InfoHash	[20]byte
	PeerID		[20]byte
}


func (h *HandShake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.InfoHash[:])

	return buf
}

func Deserialize(r io.Reader) (*HandShake, error) {
	var lengthIdentifier [1]byte
	r.Read(lengthIdentifier[:])
	pstr := make([]byte, int(lengthIdentifier[0]))

	var infoHash [20]byte
	var peerId [20]byte
	var reserveByte [8]byte
	r.Read(pstr)
	r.Read(reserveByte[:])
	r.Read(infoHash[:])
	r.Read(peerId[:])

	handShake := HandShake{
		Pstr: string(pstr),
		InfoHash: infoHash,
		PeerID: peerId,
	}
	return &handShake, nil
}
