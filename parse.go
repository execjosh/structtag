package structtag

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse parses a string representing a struct tag.
func Parse(s string) (StructTag, error) {
	t := FromMap(map[string]string{})

	orig := s
	s = normalize(s)

	// This code is based on reflect.StructTag.Lookup

	i := 0
	for i < len(s) {
		// Skip leading space.
		j := i
		for j < len(s) && s[j] == ' ' {
			j++
		}
		if j >= len(s) {
			break
		}

		// Scan to colon.
		i = j
		for j < len(s) && s[j] > ' ' && s[j] != ':' && s[j] != '"' && s[j] != 0x7f {
			j++
		}
		if j == i || j+1 >= len(s) || s[j] != ':' || s[j+1] != '"' {
			return nil, fmt.Errorf("invalid key in pair in struct tag: %s:%d-%d", orig, i, j)
		}
		key := string(s[i:j])
		j++ // move past colon

		// Scan quoted string to find value.
		i = j
		j++ // move past beg quote
		for j < len(s) && s[j] != '"' {
			if s[j] == '\\' {
				j++
			}
			j++
		}
		if j >= len(s) {
			return nil, fmt.Errorf("invalid value in pair in struct tag: %s:%d-%d, key: %s", orig, i, j, key)
		}
		j++ // move to end quote
		qval := string(s[i:j])

		val, err := strconv.Unquote(qval)
		if err != nil {
			return nil, fmt.Errorf("invald value in pair in struct tag: %s:%d-%d, key: %s, val: %s, err: %v", orig, i, j, key, qval, err)
		}

		t.Set(key, val)

		j++ // move past end quote
		i = j
	}

	return t, nil
}

func normalize(s string) string {
	s = strings.Trim(s, "`")
	return s
}
