package main

import (
	"bytes"
	"fmt"

	"github.com/jackpal/bencode-go"
	"crypto/sha1"
	"io"
)

type bencodeInfo struct {
	Pieces		string	`bencode:"pieces"`
	PieceLength	int		`bencode:"piece length"`
	Length		int 	`bencode:"length"`
	Name		string	`bencode:"name"`
}

type bencodeTorrent struct {
	Announce 	string		`bencode:"announce"`
	Info 		bencodeInfo	`bencode:"Info"`
}

func createInfoHash(i bencodeInfo) ([]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, i)

	if err != nil {
		return nil, err
	}
	hash := sha1.Sum(buf.Bytes())
	return hash[:], nil
}

func toByteChunks(s string) ([][20]byte, error) {
	hashLen := 20
	buff := []byte(s)
	
	if len(buff) % hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buff))
		return nil, err
	}

	output := make([][20]byte, len(buff)/hashLen)
	for i := 0; i<len(buff)/hashLen; i++ {
		copy(output[i][:], buff[i*hashLen:(i+1)*hashLen])
	}
	return output, nil;
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := createInfoHash(bto.Info)
	if err != nil {
		return TorrentFile{}, err
	}

	byteChunk, err := toByteChunks(bto.Info.Pieces)
	if err != nil {
		return TorrentFile{}, err
	}
	torrentFile := TorrentFile{
		Announce: bto.Announce,
		InfoHash: [20]byte(infoHash),
		PieceHashes: byteChunk,
		PieceLength: bto.Info.PieceLength,
		Length: bto.Info.Length,
		Name: bto.Info.Name,
	}
	return torrentFile, nil
}

func Open(r io.Reader) (TorrentFile, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.toTorrentFile()
}

func ParseTrackerResponse(input []byte) (TrackerResponse, error) {
	tr := TrackerResponse{}
	err := bencode.Unmarshal(bytes.NewReader(input), &tr)

	if err != nil {
		fmt.Errorf("Error Parsing the input")
		return tr, err
	}

	return tr, nil
}

