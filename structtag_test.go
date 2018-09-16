package structtag_test

import (
	"testing"

	"github.com/execjosh/structtag"
	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	type kvs map[string]string

	testcases := map[string]struct {
		orig kvs
		want []string
	}{
		"empty":    {kvs{}, []string{}},
		"one pair": {kvs{"json": "-"}, []string{"json"}},
		"two pair": {kvs{"valid": "-", "json": "-"}, []string{"json", "valid"}},
		"ordering": {kvs{"valid": "-", "json": "-", "xxx": "x", "abc": "a"}, []string{"abc", "json", "valid", "xxx"}},
	}
	for _, tc := range testcases {
		tag := structtag.FromMap(tc.orig)
		assert.Equal(t, tc.want, tag.Keys())
	}
}

func TestString(t *testing.T) {
	type kvs map[string]string

	testcases := map[string]struct {
		orig kvs
		want string
	}{
		"empty":    {kvs{}, "``"},
		"one pair": {kvs{"json": "-"}, "`json:\"-\"`"},
		"two pair": {kvs{"valid": "-", "json": "-"}, "`json:\"-\" valid:\"-\"`"},
		"ordering": {kvs{"valid": "-", "json": "-", "xxx": "x", "abc": "a"}, "`abc:\"a\" json:\"-\" valid:\"-\" xxx:\"x\"`"},
	}
	for _, tc := range testcases {
		tag := structtag.FromMap(tc.orig)
		assert.Equal(t, tc.want, tag.String())
	}
}

func TestUnion(t *testing.T) {
	testcases := map[string]struct {
		a    string
		b    string
		want string
	}{
		"second tag empty":            {`json:"-"`, "", "`json:\"-\"`"},
		"first tag empty":             {"", `json:"json"`, "`json:\"json\"`"},
		"duplicate keys, second wins": {`json:"-"`, `json:"json"`, "`json:\"json\"`"},
	}
	for _, tc := range testcases {
		a, err := structtag.Parse(tc.a)
		assert.NoError(t, err)
		b, err := structtag.Parse(tc.b)
		assert.NoError(t, err)
		assert.Equal(t, tc.want, structtag.Union(a, b).String())
	}
}
