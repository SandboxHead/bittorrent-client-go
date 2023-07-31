package main 

import (
	"fmt"
	"crypto/rand"
)

func peerIdGenerator() [20]byte {
	output := make([]byte, 20)
	_, err := rand.Read(output)
	if err != nil {
		fmt.Errorf("Error while generating random bytes")
	}
	return [20]byte(output)
}
