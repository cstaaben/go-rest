package request_test

import (
	"slices"
	"testing"

	"github.com/cstaaben/go-rest/internal/request"
)

func TestLoadGroupIDsFromDir(t *testing.T) {
	dir := "/home/corbinstaaben/code/src/github.com/cstaaben/go-rest/example/data"
	t.Logf("Loading group IDs from %s", dir)

	ids, err := request.LoadGroupIDsFromDir(dir)
	if err != nil {
		t.Fatalf("unexpected error loading group IDs: %v", err)
	}

	t.Logf("Loaded group IDs: %#v", ids)

	if len(ids) != 1 {
		t.Fatalf("expected 1 group, but received %d", len(ids))
	}

	if !slices.Contains(ids, "3901089e-b8a5-47c2-8504-064e4a6c3413") {
		t.Fatalf("missing id for example group in loaded IDs")
	}
}
