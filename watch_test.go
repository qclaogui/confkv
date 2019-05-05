package confkv

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAppendPrefix(t *testing.T) {
	prefix := "/app"
	keys := []string{"/redis/addr", "kafka/addr"}
	want := []string{"/app/redis/addr", "/app/kafka/addr"}
	got := appendPrefix(prefix, keys)
	if df := cmp.Diff(got, want); df != "" {
		t.Errorf("👉 \x1b[92m%s\x1b[39m", df)
	}
}
