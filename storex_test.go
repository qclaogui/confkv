package confkv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStoreX_Get(t *testing.T) {
	var tests = map[string]struct {
		key   string
		value string
		err   error
		want  KVPair
	}{
		"case1": {"/db/user", "admin", nil, KVPair{"/db/user", "admin"}},
		"case2": {"/db/pass", "foo", nil, KVPair{"/db/pass", "foo"}},
		"case3": {"/missing", "", ErrNotExist, KVPair{}},
	}
	db := NewStoreX()
	for name, test := range tests {
		// set first
		if test.err == nil {
			db.Set(test.key, test.value)
		}

		t.Run(name, func(t *testing.T) {
			got, err := db.Get(test.key)
			if df := cmp.Diff(err, test.err); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
			if df := cmp.Diff(got, test.want); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
		})
	}
}

func TestStoreX_GetValue(t *testing.T) {
	var tests = map[string]struct {
		key   string
		value string
		err   error
		want  string
	}{
		"case1": {"/db/user", "admin", nil, "admin"},
		"case2": {"/db/pass", "foo", nil, "foo"},
		"case3": {"/missing", "", ErrNotExist, ""},
	}

	db := NewStoreX()
	for name, test := range tests {
		// set first
		if test.err == nil {
			db.Set(test.key, test.value)
		}
		t.Run(name, func(t *testing.T) {
			got, err := db.GetValue(test.key)
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
	db := NewStoreX()

	got, err := db.GetValue("/db/user", "defaultValue")
	if df := cmp.Diff(err, nil); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
	if df := cmp.Diff(got, want); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
}

func TestGetValueWithEmptyDefault(t *testing.T) {
	want := ""
	db := NewStoreX()

	got, err := db.GetValue("/db/user", "")
	if df := cmp.Diff(err, nil); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
	if df := cmp.Diff(got, want); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
}

var getalltestinput = map[string]string{
	"/app/db/pass":               "foo",
	"/app/db/user":               "admin",
	"/app/port":                  "443",
	"/app/url":                   "app.example.com",
	"/app/vhosts/host1":          "app.example.com",
	"/app/upstream/host1":        "203.0.113.0.1:8080",
	"/app/upstream/host1/domain": "app.example.com",
	"/app/upstream/host2":        "203.0.113.0.2:8080",
	"/app/upstream/host2/domain": "app.example.com",
}

var getalltests = map[string]struct {
	pattern string
	err     error
	want    KVPairs
}{
	"case1": {"/app/db/*", nil,
		KVPairs{
			KVPair{"/app/db/pass", "foo"},
			KVPair{"/app/db/user", "admin"}}},
	"case2": {"/app/*/host1", nil,
		KVPairs{
			KVPair{"/app/upstream/host1", "203.0.113.0.1:8080"},
			KVPair{"/app/vhosts/host1", "app.example.com"}}},

	"case3": {"/app/upstream/*", nil,
		KVPairs{
			KVPair{"/app/upstream/host1", "203.0.113.0.1:8080"},
			KVPair{"/app/upstream/host2", "203.0.113.0.2:8080"}}},
	"case4": {"[]a]", ErrNoMatch, nil},
	"case5": {"/app/missing/*", ErrNoMatch, nil},
}

func TestStoreX_GetAll(t *testing.T) {
	db := NewStoreX()
	for key, value := range getalltestinput {
		db.Set(key, value)
	}

	for name, test := range getalltests {
		t.Run(name, func(t *testing.T) {
			got, err := db.GetAll(test.pattern)
			if df := cmp.Diff(err, test.err); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
			if df := cmp.Diff(got, test.want); df != "" {
				t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
			}
		})
	}
}

func TestDel(t *testing.T) {
	db := NewStoreX()
	db.Set("/app/port", "8080")
	want := KVPair{"/app/port", "8080"}
	got, err := db.Get("/app/port")
	if df := cmp.Diff(err, nil); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
	if df := cmp.Diff(got, want); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}

	db.Del("/app/port")
	want = KVPair{}
	got, err = db.Get("/app/port")
	if df := cmp.Diff(err, ErrNotExist); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
}
