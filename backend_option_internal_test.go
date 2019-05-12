package kv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseNodes(t *testing.T) {
	var tests = map[string]struct {
		nodes []string
		want  []string
	}{
		"case1": {[]string{}, nil},
		"case2": {nil, nil},
		"case3": {[]string{"127.0.0.1"}, []string{"127.0.0.1:2181"}},
		"case4": {[]string{"127.0.0.1:2181"}, []string{"127.0.0.1:2181"}},
		"case5": {[]string{"127.0.0.1:2182"}, []string{"127.0.0.1:2182"}},
		"case6": {
			[]string{"127.0.0.1,127.0.0.1:2182,127.0.0.1:2183"},
			[]string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"},
		},
		"case7": {
			[]string{"127.0.0.1,127.0.0.1:2182", "127.0.0.1:2183,127.0.0.1:2184", "127.0.0.1:2185"},
			[]string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183", "127.0.0.1:2184", "127.0.0.1:2185"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseNodes(test.nodes...)
			if df := cmp.Diff(got, test.want); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
		})
	}
}
