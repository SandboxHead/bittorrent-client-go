package main 

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open(os.Args[1])
	bto, _ := Open(f)
	fmt.Println(bto)
}