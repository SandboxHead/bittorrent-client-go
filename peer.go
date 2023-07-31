package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Peer struct {
	IP net.IP
	Port uint16
}

func Unmarshal(peerBin []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peerBin) / peerSize
	if len(peerBin)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
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

