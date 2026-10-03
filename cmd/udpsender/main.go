package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Print(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Print(err)
	}
	defer conn.Close()
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		lin, err := input.ReadString('\n')
		if err == io.EOF {
			os.Exit(1)
		}
		conn.Write([]byte(lin))
	}
}
