package request

import (
	"io"
	"unicode"
)

// Request represent the full parsed Http request
type Request struct {
	RequestLine RequestLine
}
type RequestLine struct {
	HTTPversion    string
	RequestTargert string
	Method         string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	var httpMessage Request

	PareseRequestLine(req, &httpMessage)

	return &httpMessage
}

func PareseRequestLine(req []byte, httpMessage *Request) {
	var j int = 0
	for i := range string(req) {
		if string(req[i]) == " " {
			method := string(req[j:i])
			switch method {
			case "Get":
				httpMessage.RequestLine.Method = method
			case "POST":
				httpMessage.RequestLine.Method = method

			case "PUT":
				httpMessage.RequestLine.Method = method

			case "DELETE":
				httpMessage.RequestLine.Method = method

			case "HEAD":
				httpMessage.RequestLine.Method = method

			case "OPTIONS":
				httpMessage.RequestLine.Method = method

			case "TRACE":
				httpMessage.RequestLine.Method = method

			case "CONNECT":
				httpMessage.RequestLine.Method = method
			default:
				httpMessage.RequestLine.Method = ""
			}
			if string(req[i]) == "/" {
				for k := i; k < len(req); k++ {
					if string(req[k]) == " " {
						httpMessage.RequestLine.RequestTargert = string(req[i:k])
					}
				}
			}

		}

		if string(req[i]) == "" {
			httpMessage.RequestLine.Method = ""
		}
		r := rune(req[i])
		if unicode.IsLower(r) {
			httpMessage.RequestLine.Method = ""
		}
	}
}
