package main 

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open(os.Args[1])
	bto, _ := Open(f)
	peerId := peerIdGenerator()
	url, _ := bto.buildTrackerURL(peerId, 20102)
	fmt.Println(url)
	response := query(url)
	parsed, _ := ParseTrackerResponse(response)
	peers, _ := Unmarshal([]byte(parsed.Peers))
	fmt.Println(peers)
}