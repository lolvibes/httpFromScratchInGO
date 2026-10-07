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

func isUpper(s string) bool {
	for _, r := range s {
		if !(r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}

func isLower(s string) bool {
	for _, r := range s {
		if !(r >= 'a' && r <= 'z') {
			return false
		}
	}
	return true
}

func validSpeacialChar(s string) bool {
	specialchar := []string{"!", "#", "$", "%", "&", ","}
	for _, r := range s {
		if !((r >= 'A' && r <= 'z') || (r >= 'a' && r <= 'z')) { // then it must be a special character if it is a special character then is it from our slice
			for _, f := range specialchar {
				if string(r) == f {
					continue
				}
			}
			return true

		}
	}
	return false
}

func validityCheck(headername string) bool {
	if (isUpper(headername) || isLower(headername)) && validSpeacialChar(headername) {
		return true
	}
	return false
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
	if validityCheck(string(hname)) {

		hnameStr := string(hname)
		hnameStr = strings.ToLower(hnameStr)
		h[hnameStr] = string(bytes.TrimSpace(hvalue))
		return idx + 2, false, nil
	}
	return 0, false, fmt.Errorf("this is above my pay grade")
}
