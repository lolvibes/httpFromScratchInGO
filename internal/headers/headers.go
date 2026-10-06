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
	if idx == -1 {
		return 0, false, nil // need data no full line yet
	}
	if idx == 0 {
		return 2, true, nil // no header to parse
	}
	line := data[:idx]
	hname, hvalue, found := bytes.Cut(line, []byte(":"))
	if !found {
		return 0, false, fmt.Errorf("malformed header: no colon in %q", line)
	}
	hname = bytes.TrimLeft(hname, " ")
	if len(hname) == 0 || bytes.Contains(hname, []byte(" ")) {
		return 0, false, fmt.Errorf("invalid header name %q", hname)
	}
	h[string(hname)] = string(bytes.TrimSpace(hvalue))
	return idx + 2, false, nil
}
