package request

import (
	"fmt"
	"io"
	"strings"
)

// Request represent the full parsed Http request
type Request struct {
	RequestLine RequestLine
}
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// escaping the pointer value aren't dead
func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rl, err := PareseRequestLine(string(data[:]))
	if err != nil {
		return nil, err
	}
	// composit literal
	return &Request{RequestLine: rl}, nil
}

func PareseRequestLine(data string) (RequestLine, error) {
	parts := strings.Split(data, "\r\n")

	reqLine := parts[0]
	rqSlice := strings.Split(reqLine, " ")
	if len(rqSlice) != 3 {
		fmt.Errorf("bad req nedded 3 parts gave %d", len(rqSlice))
	}
	v := make(map[int]string)
	for i, l := range rqSlice {
		v[i] = l
	}
	var rl RequestLine
	method := v[0]
	switch method {
	case "GET":
		rl.Method = method
	case "POST":
		rl.Method = method

	case "PUT":
		rl.Method = method

	case "DELETE":
		rl.Method = method

	case "HEAD":
		rl.Method = method

	case "OPTIONS":
		rl.Method = method
	case "TRACE":
		rl.Method = method

	case "CONNECT":
		rl.Method = method

	default:
		return RequestLine{}, fmt.Errorf("invalid method %s", method)
	}
	reqTarget := v[1]
	if reqTarget[0] == '/' {
		rl.RequestTarget = reqTarget
	}
	httpversion := v[2]
	if httpversion != "HTTP/1.1" {
		return RequestLine{}, nil
	}
	rl.HttpVersion = strings.TrimPrefix(httpversion, "HTTP/")
	return rl, nil
}
