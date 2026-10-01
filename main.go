package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 123460)
	b, err := file.Read(buf)

	for i := 0; i < b; i += 8 {
		j := min(i+8, b)
		fmt.Printf("read: %s\n", string(buf[i:j]))
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
}
