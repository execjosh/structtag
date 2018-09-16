package structtag

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// StructTag represents a struct tag.
type StructTag interface {
	Set(key, val string)
	Get(key string) (string, bool)
	Keys() []string
	String() string
}

type structTagImpl map[string]string

// FromMap creates a new StructTag.
func FromMap(m map[string]string) StructTag {
	return structTagImpl(m)
}

// Union creates a new StructTag with all keys from a and b.
// Values for keys in a will be overwritten if set in b.
func Union(a, b StructTag) StructTag {
	tag := structTagImpl(map[string]string{})

	tag.append(a)
	tag.append(b)

	return tag
}

func (t structTagImpl) Set(key, val string) {
	t[key] = val
}

func (t structTagImpl) Get(key string) (string, bool) {
	val, ok := t[key]
	return val, ok
}

func (t structTagImpl) Keys() []string {
	keys := []string{}
	for key := range t {
		keys = append(keys, key)
	}

	// Ensure reproducible ordering
	sort.Strings(keys)

	return keys
}

func (t structTagImpl) String() string {
	s := []string{}
	for _, key := range t.Keys() {
		val := strconv.Quote(t[key])
		s = append(s, fmt.Sprintf("%s:%s", key, val))
	}

	return fmt.Sprintf("`%s`", strings.Join(s, " "))
}

func (t structTagImpl) append(tag StructTag) {
	for _, key := range tag.Keys() {
		val, _ := tag.Get(key)
		t.Set(key, val)
	}
}
