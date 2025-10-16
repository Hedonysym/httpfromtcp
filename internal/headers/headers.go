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
	split := strings.Split(string(data), "\r\n")
	if len(split) < 2 {
		return 0, false, nil
	}
	if split[0] == "" {
		return 0, true, nil
	}
	if strings.Contains(split[0], " :") {
		return 0, false, InvaidHeaderSpacingError
	}

	split2 := strings.SplitN(split[0], ":", 2)
	if split2[0] == "" || invalidName(split2[0]) {
		return 0, false, InvalidHeaderNameError
	}
	key := strings.TrimSpace(strings.ToLower(split2[0]))
	value := strings.TrimSpace(strings.ToLower(split2[1]))
	if h[key] != "" {
		h[key] += ", " + value
	} else {
		h[key] = value
	}
	return len(split[0]) + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
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
