package structtag_test

import (
	"testing"

	"github.com/execjosh/structtag"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type kvs map[string]string

	testcases := map[string]struct {
		errExpected bool
		orig        string
		want        kvs
	}{
		"empty string":                    {false, "", kvs{}},
		"empty backquoted string":         {false, "``", kvs{}},
		"space":                           {false, " ", kvs{}},
		"space backquoted string":         {false, "` `", kvs{}},
		"spaces":                          {false, "        ", kvs{}},
		"spaces backquoted string":        {false, "`      `", kvs{}},
		"one pair, no quotes":             {false, "json:\"-\"", kvs{"json": "-"}},
		"one pair, backquoted":            {false, "`json:\"-\"`", kvs{"json": "-"}},
		"two pair, no quotes":             {false, "json:\"-\" valid:\"-\"", kvs{"json": "-", "valid": "-"}},
		"two pair, backquoted":            {false, "`json:\"-\" valid:\"-\"`", kvs{"json": "-", "valid": "-"}},
		"escape in value":                 {false, "json:\"\\n\"", kvs{"json": "\n"}},
		"duplicate keys, last value wins": {false, "json:\"-\" json:\"omg\" json:\"required\"", kvs{"json": "required"}},

		"just colon":       {true, ":", kvs{}},
		"no value":         {true, "json:", kvs{}},
		"no key":           {true, ":\"-\"", kvs{}},
		"no start quote":   {true, "json:-\"", kvs{}},
		"no end quote":     {true, "json:\"-", kvs{}},
		"newline in value": {true, "json:\"\n\"", kvs{}},
	}
	for _, tc := range testcases {
		tag, err := structtag.Parse(tc.orig)
		if tc.errExpected {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		for key, wantVal := range tc.want {
			val, ok := tag.Get(key)
			assert.True(t, ok)
			assert.Equal(t, wantVal, val)
		}
	}
}
