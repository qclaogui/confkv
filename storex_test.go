package kv_test

import (
	"testing"

	"github.com/qclaogui/kv"

	"github.com/google/go-cmp/cmp"
)

func TestStoreX_GetValue(t *testing.T) {
	var tests = map[string]struct {
		key   string
		value string
		err   error
		want  string
	}{
		"case1": {"/db/user", "admin", nil, "admin"},
		"case2": {"/db/pass", "foo", nil, "foo"},
		"case3": {"/missing", "", kv.ErrNotExist, ""},
	}

	db := kv.NewDB()
	for name, test := range tests {
		// Set first
		if test.err == nil {
			db.Set(test.key, test.value)
		}
		t.Run(name, func(t *testing.T) {
			got, err := db.GetV(test.key)
			if df := cmp.Diff(err, test.err); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
			if df := cmp.Diff(got, test.want); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
		})
	}
}

func TestGetValueWithDefault(t *testing.T) {
	want := "defaultValue"
	db := kv.NewDB()

	got, err := db.GetV("/db/user", "defaultValue")
	if df := cmp.Diff(err, nil); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
	if df := cmp.Diff(got, want); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
}

func TestGetValueWithEmptyDefault(t *testing.T) {
	want := ""
	db := kv.NewDB()

	got, err := db.GetV("/db/user", "")
	if df := cmp.Diff(err, nil); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
	if df := cmp.Diff(got, want); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
}
