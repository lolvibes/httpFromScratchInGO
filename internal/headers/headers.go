package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func checkvalid(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
		case strings.ContainsRune("!#$%&^*+_-`.|~^", r):
		default:
			return false

		}
	}
	return true
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
	if !checkvalid(string(hname)) {
		return 0, false, fmt.Errorf("invalid header name ; %q", hname)
	}

	hnameStr := string(hname)
	hnameStr = strings.ToLower(hnameStr)
	hvalueStr := string(bytes.TrimSpace(hvalue))
	if existing, ok := h[hnameStr]; ok {
		h[hnameStr] = existing + ", " + hvalueStr
	} else {
		h[hnameStr] = hvalueStr
	}
	return idx + 2, false, nil
}
