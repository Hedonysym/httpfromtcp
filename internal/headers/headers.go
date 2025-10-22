package headers

import (
	"strings"
	"unicode"
)

var validUnicode = map[rune]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'|':  true,
	'~':  true,
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	s := string(data)

	// End of headers: starts with CRLF
	if strings.HasPrefix(s, "\r\n") {
		return 2, true, nil
	}

	// Need at least one full line
	i := strings.Index(s, "\r\n")
	if i == -1 {
		return 0, false, nil
	}
	line := s[:i]
	if line == "" {
		// Safety: if blank line shows up not at start, still consume
		return 2, true, nil
	}
	if strings.Contains(line, " :") {
		return 0, false, InvaidHeaderSpacingError
	}
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 || parts[0] == "" || invalidName(parts[0]) {
		return 0, false, InvalidHeaderNameError
	}
	key := strings.TrimSpace(strings.ToLower(parts[0]))
	value := strings.TrimSpace(strings.ToLower(parts[1]))
	if h[key] != "" {
		h[key] += ", " + value
	} else {
		h[key] = value
	}
	return i + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (string, bool) {
	if h[strings.ToLower(key)] == "" {
		return "", false
	}
	return h[strings.ToLower(key)], true
}

func invalidName(data string) bool {
	stripped := strings.TrimSpace(data)
	for _, c := range stripped {
		if !unicode.IsDigit(rune(c)) && !unicode.IsLetter(rune(c)) && !validUnicode[rune(c)] {
			return true
		}
	}
	return false
}
