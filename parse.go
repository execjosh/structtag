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
	tag := normalize(s)

	// This code is based on reflect.StructTag.Lookup

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			return nil, fmt.Errorf("invalid key in pair in struct tag: %s", orig)
		}
		key := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			return nil, fmt.Errorf("invalid value in pair in struct tag: %s, key: %s", orig, key)
		}
		qval := string(tag[:i+1])
		tag = tag[i+1:]

		val, err := strconv.Unquote(qval)
		if err != nil {
			return nil, fmt.Errorf("invald value in pair in struct tag: %s, key: %s, val: %s, err: %v", orig, key, qval, err)
		}

		t.Set(key, val)
	}

	return t, nil
}

func normalize(s string) string {
	s = strings.Trim(s, "`")
	return s
}
