package response

import (
	"fmt"
	"io"
	"strconv"

	"github.httpfromtcp/internal/headers"
)

type StatusCode int

const (
	ServerOk            StatusCode = 200
	BadReq              StatusCode = 400
	InternalServerError            = 500
)

func WriteStatusLine(w io.Writer, StatusCode StatusCode) error {
	switch StatusCode {
	case ServerOk:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case BadReq:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case InternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d \r\n", StatusCode)))
		return err
	}
}

func GetDefaultHeaders(contentlen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentlen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n")) // blank line ends the headers
	return err
}
