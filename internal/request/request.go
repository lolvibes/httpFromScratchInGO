package request

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.httpfromtcp/internal/headers"
)

type parseState int

const (
	stateint parseState = iota
	stateHeaders
	statebody
	stateDone
)

// Request represent the full parsed Http request
type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parseState
}
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// escaping the pointer value aren't dead
func RequestFromReader(reader io.Reader) (*Request, error) {
	// trying to read from buffer
	buf := make([]byte, 8)
	readToIdx := 0
	req := &Request{
		state:   stateint,
		Headers: headers.NewHeaders(),
	}
	//	data, err := io.ReadAll(reader)
	for req.state != stateDone {
		if readToIdx == len(buf) {
			bigbuf := make([]byte, len(buf)*2)
			copy(bigbuf, buf[:readToIdx])
			buf = bigbuf
		}
		n, err := reader.Read(buf[readToIdx:])
		readToIdx += n
		consumed, parseError := req.parse(buf[:readToIdx])
		if parseError != nil {
			return nil, parseError
		}
		copy(buf, buf[consumed:readToIdx])
		readToIdx -= consumed
		if err != nil {
			return nil, err
		}
		if err == io.EOF {
			return nil, fmt.Errorf("connection close before the request line complete")
			break
		}

	}
	return req, nil
}

func PareseRequestLine(data string) (RequestLine, error, int) {
	idx := strings.Index(data, "\r\n")
	if idx == -1 {
		return RequestLine{}, nil, 0
	}

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
		return RequestLine{}, fmt.Errorf("invalid method %s", method), 0
	}
	reqTarget := v[1]
	if reqTarget[0] == '/' {
		rl.RequestTarget = reqTarget
	}
	httpversion := v[2]
	if httpversion != "HTTP/1.1" {
		return RequestLine{}, fmt.Errorf("wrong htt version need %s", httpversion), 0
	}
	rl.HttpVersion = strings.TrimPrefix(httpversion, "HTTP/")
	return rl, nil, idx + 2
}

func (r *Request) parse(data []byte) (int, error) {
	total := 0
	for r.state != stateDone {
		n, err := r.parseSingle(data[total:])
		if err != nil {
			return 0, nil
		}
		if n == 0 {
			break
		}
		total += n
	}
	return total, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case stateint:
		rl, err, n := PareseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 { // simply means not enough byte
			return 0, nil
		}
		r.RequestLine = rl
		r.state = stateHeaders
		return n, nil
	case stateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			cl, ok := r.Headers["content-length"]
			if !ok {
				r.state = stateDone
			} else {
				length, err := strconv.Atoi(cl)
				if err != nil || length < 0 {
					return 0, fmt.Errorf("invalid content lenght  %q", cl)
				}
				if length == 0 {
					r.state = stateDone
				} else {
					r.state = statebody
				}
			}
		}
		return n, nil
	case statebody:
		length, err := strconv.Atoi(r.Headers["content-length"])
		if err != nil {
			return 0, fmt.Errorf("invalid content-length %q ")
		}
		remaining := length - len(r.Body)
		take := min(remaining, len(data))
		r.Body = append(r.Body, data[:take]...)
		if len(r.Body) == length {
			r.state = stateDone
		}

		return take, nil
	default:
		return 0, fmt.Errorf("parser called unfimiliar state %d", r.state)

	}
}
