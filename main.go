package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("messages.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//
	eightByteSlice := make([]byte, 8)
	for {
		n, err := f.Read(eightByteSlice)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("read: %s\n", string(eightByteSlice[:n]))
		// eightByteSlice = eightByteSlice[:0]
		// n, err = f.ReadAt(eightByteSlice, int64(n))
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
}
