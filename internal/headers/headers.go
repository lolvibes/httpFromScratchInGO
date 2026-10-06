package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == 0 { // this means we just got zero headers
		return 0, true, nil
	}

	parts := bytes.Split(data, []byte("\r\n")) // this gives the n number of headers
	// check or all the headers is those are valid or not
	for _, data2 := range parts {
		subParts := bytes.Split(data2, []byte(":"))
		for i, subdata := range subParts {
			if bytes.Contains(subdata, []byte(" ")) { // not thats the invalid headers
				return 0, false, fmt.Errorf("invalid header , header : %s", string(subdata))
			}
			subdata = append(subdata, ':')
			h[string(subParts[i])] = string(subParts[i+1])
			break
		}

	}
}
