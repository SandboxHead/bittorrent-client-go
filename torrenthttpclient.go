package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func query(url string) []byte {
	resp, err  := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status: ", resp.Status)
	
	scanner := bufio.NewScanner(resp.Body)
	output := []byte{}
	for i:=0; scanner.Scan(); i++ {
		output = append(output, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return output
}