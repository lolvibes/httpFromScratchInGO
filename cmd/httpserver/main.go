package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.httpfromtcp/internal/request"
	"github.httpfromtcp/internal/response"
	"github.httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(w io.Writer, req *request.Request) *server.HandleError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandleError{
				StatusCode: response.BadReq,
				Message:    "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandleError{
				StatusCode: response.InternalServerError,
				Message:    "Woopsie, my bad\n",
			}
		default:
			w.Write([]byte("All good, frfr\n"))
			return nil
		}
	}
	server, err := server.Serv(handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
