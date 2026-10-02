package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, _ := os.Open("messages.txt")
	channel := getLineChannel(file)
	for i := range channel {
		fmt.Println("read:", i)
	}
}

func getLineChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)

		buf := make([]byte, 123460)
		b, err := f.Read(buf)
		//fmt.Println("bytes read : ", b, "err: ", err)
		currentLine := ""
		for i := 0; i < b; i += 8 {
			j := min(i+8, b)
			str := string(buf[i:j])
			//fmt.Printf("read: %s\n", string(buf[i:j]))
			parts := strings.Split(str, "\n")
			for k := 0; k < len(parts)-1; k++ {
				ch <- currentLine + parts[k]
				currentLine = ""

			}
			// this piece of block never runs bcz err is calculated outside of loop and never checked again
			// cuz of reading entirefile at once we need to check for currentLine is empty or not  outside of loop
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
					currentLine = ""
				}
				if err == io.EOF {
					break
				}
			}
			currentLine += parts[len(parts)-1]
		}
		if currentLine != "" {
			ch <- currentLine
		}
	}()
	return ch
}
