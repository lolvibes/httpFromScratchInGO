package main

import (
	"fmt"
	"net"

	"github.httpfromtcp/internal/request"
)

func main() {
	//`file, _ := os.Open("messages.txt")
	//	channel := getLineChannel(file)
	lister, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, _ := lister.Accept()
		if conn != nil {
			fmt.Println("connection has been aceepted")
		}
		// for i := range channel {
		//	fmt.Println("read:", i)
		//channel := getLineChannel(conn)
		//
		//for i := range channel {
		//	fmt.Printf("%s", i)
		//}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Errorf("don't know what happened", err)
		}
		fmt.Printf("Request line :\n -Method: %s\n -Target: %s\n -Version:%s", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

	}
}
