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

	if len(ids) != 2 {
		t.Fatalf("expected 2 requests, but received %d", len(ids))
	}

	if !slices.Contains(ids, "33336ae0-0545-43a7-a6be-1dce501170c4") {
		t.Fatalf("missing id for example request 1 in loaded IDs")
	} else if !slices.Contains(ids, "f6b11dd8-db26-41b9-a93a-be19afc29e0f") {
		t.Fatalf("missing id for example request 2 in loaded IDs")
	}
}
