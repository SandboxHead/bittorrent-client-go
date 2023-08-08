package main

import (
	"bittorrent-client-go/entity"
	"bittorrent-client-go/helper"
	"fmt"
	"os"
)

func DownloadTorrent(torrentFile string) (error) {
	f, _ := os.Open(torrentFile)
	bto, _ := helper.Open(f)
	peerId := helper.PeerIdGenerator()
	url, _ := bto.BuildTrackerURL(peerId, 20102)
	fmt.Println(url)
	response := query(url)
	parsed, err := helper.ParseTrackerResponse(response)
	fmt.Println(err)
	if err != nil {
		return err
	}
	fmt.Println(parsed)
	
	peers, err := entity.UnmarshalPeer([]byte(parsed.Peers))
	if err != nil {
		return err
	}

	workQueue := make(chan *entity.PieceWork, len(bto.PieceHashes))
	resultQueue := make(chan *entity.PieceResult)

	for index, hash := range bto.PieceHashes {
		length := bto.CalculatePieceSize(index)
		piecework := entity.PieceWork{Index: uint(index), Hash: hash, Length: uint(length)}
		workQueue <- &piecework
	}
	fmt.Println("HEre")

	for _, peer := range peers {
		go startDownload(peer, workQueue, resultQueue, bto, peerId)
	}
	// buf := make([]byte, bto.Length)
	donePieces := 0
	for donePieces < len(bto.PieceHashes) {
		res := <-resultQueue
		begin := bto.PieceLength * res.Index
		end := begin + bto.PieceLength
		if end > bto.Length {
			end = bto.Length
		}
		// copy(buf[begin:end], res.Buf)
		donePieces++
	}
	close(workQueue)
	close(resultQueue)
	// fmt.Println(x)
	return nil
}

func startDownload(peer entity.Peer, workQueue chan *entity.PieceWork, resultQueue chan *entity.PieceResult, bto entity.TorrentFile, peerId [20]byte) (error) {
	client, err := entity.New(peer, bto.InfoHash, peerId)
	fmt.Println(client)

	if err != nil {
		return err
	}
	defer client.Conn.Close()
	err = client.SendUnchoke()
	if err != nil {
		return err
	}
	err = client.SendInterested()
	if err != nil {
		return err
	}

	for pw := range workQueue {
		if !client.BitField.HasPiece(int(pw.Index)) {
			workQueue <- pw
			continue
		}
		buf, err := client.DownloadPiece(pw)
		if err != nil {
			fmt.Printf("Piece #%d failed integrity check\n", pw.Index)
			workQueue <- pw
			continue
		}
		client.SendHave(int(pw.Index))
		resultQueue <- &entity.PieceResult{Index: int(pw.Index), Buf: buf}
		fmt.Println("Downloaded Piece %d", pw.Index)
	}
	return nil

}