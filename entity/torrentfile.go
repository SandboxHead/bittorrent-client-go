package entity

import (
	"fmt"
	"net/url"
	"strconv"
)

type TorrentFile struct {
	Announce	string
	InfoHash	[20]byte
	PieceHashes	[][20]byte
	PieceLength	int
	Length		int 
	Name		string
}

func (t *TorrentFile) BuildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		fmt.Errorf("Error while parsing the announce: %d", err)
		return "", err 
	}
	params := url.Values{
		"info_hash" : []string{string(t.InfoHash[:])},
		"peer_id" : []string{string(peerID[:])},
		"port" : []string{strconv.Itoa(int(port))},
		"uploaded" : []string{"0"},
		"downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) CalculatePieceSize(index int) (int) {
	start := index * t.PieceLength
	end := start + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return end - start
}