package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/database/dto"
)

func TestNewGroup(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
	}{
		{
			name:      "create group with name",
			groupName: "Test Group",
		},
		{
			name:      "create group with empty name",
			groupName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := dto.NewGroup(tt.groupName)
			assert.NotNil(t, group)
			assert.Equal(t, tt.groupName, group.Name)
			assert.NotNil(t, group.Requests)
			assert.Empty(t, group.Requests)
		})
	}
}

func TestGroup_MetadataMethods(t *testing.T) {
	tests := []struct {
		name            string
		group           *dto.Group
		wantFilterValue string
		wantTitle       string
		wantDesc        string
	}{
		{
			name: "populated group",
			group: &dto.Group{
				UUID: "g-1",
				Name: "Auth Service",
				Desc: "Authentication endpoints",
			},
			wantFilterValue: "Auth Service",
			wantTitle:       "Auth Service",
			wantDesc:        "Authentication endpoints",
		},
		{
			name:            "empty group",
			group:           &dto.Group{},
			wantFilterValue: "",
			wantTitle:       "",
			wantDesc:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFilterValue, tt.group.FilterValue())
			assert.Equal(t, tt.wantTitle, tt.group.Title())
			assert.Equal(t, tt.wantDesc, tt.group.Description())
		})
	}
}

func TestGroup_AddRequest(t *testing.T) {
	tests := []struct {
		name          string
		initial       []*dto.Request
		toAdd         []*dto.Request
		wantLength    int
		wantLastTitle string
	}{
		{
			name:          "add to empty group",
			initial:       nil,
			toAdd:         []*dto.Request{{Name: "Req 1"}},
			wantLength:    1,
			wantLastTitle: "Req 1",
		},
		{
			name:          "add multiple requests",
			initial:       []*dto.Request{{Name: "Req 1"}},
			toAdd:         []*dto.Request{{Name: "Req 2"}, {Name: "Req 3"}},
			wantLength:    3,
			wantLastTitle: "Req 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := dto.NewGroup("Test")
			if tt.initial != nil {
				group.Requests = tt.initial
			}
			for _, r := range tt.toAdd {
				group.AddRequest(r)
			}
			assert.Len(t, group.Requests, tt.wantLength)
			if tt.wantLength > 0 {
				assert.Equal(t, tt.wantLastTitle, group.Requests[tt.wantLength-1].Name)
			}
		})
	}
}

func TestGroup_RemoveRequest(t *testing.T) {
	req1 := &dto.Request{UUID: "r1", Name: "Req 1"}
	req2 := &dto.Request{UUID: "r2", Name: "Req 2"}
	req3 := &dto.Request{UUID: "r3", Name: "Req 3"}

	tests := []struct {
		name       string
		initial    []*dto.Request
		toRemove   *dto.Request
		wantLength int
		wantNames  []string
	}{
		{
			name:       "remove existing request from middle",
			initial:    []*dto.Request{req1, req2, req3},
			toRemove:   req2,
			wantLength: 2,
			wantNames:  []string{"Req 1", "Req 3"},
		},
		{
			name:       "remove non-existent request",
			initial:    []*dto.Request{req1, req2},
			toRemove:   req3,
			wantLength: 2,
			wantNames:  []string{"Req 1", "Req 2"},
		},
		{
			name:       "remove from empty group",
			initial:    []*dto.Request{},
			toRemove:   req1,
			wantLength: 0,
			wantNames:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := dto.NewGroup("Test")
			group.Requests = append([]*dto.Request(nil), tt.initial...)
			group.RemoveRequest(tt.toRemove)

			assert.Len(t, group.Requests, tt.wantLength)
			gotNames := make([]string, len(group.Requests))
			for i, r := range group.Requests {
				gotNames[i] = r.Name
			}
			assert.Equal(t, tt.wantNames, gotNames)
		})
	}
}

func TestGroup_ListItems(t *testing.T) {
	req1 := &dto.Request{UUID: "r1", Name: "Req 1"}
	req2 := &dto.Request{UUID: "r2", Name: "Req 2"}

	tests := []struct {
		name       string
		requests   []*dto.Request
		wantLength int
	}{
		{
			name:       "empty requests list",
			requests:   nil,
			wantLength: 0,
		},
		{
			name:       "populated requests list",
			requests:   []*dto.Request{req1, req2},
			wantLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := dto.NewGroup("Test")
			group.Requests = tt.requests
			items := group.ListItems()

			assert.Len(t, items, tt.wantLength)
			for i, item := range items {
				req, ok := item.(*dto.Request)
				assert.True(t, ok)
				assert.Equal(t, tt.requests[i], req)
			}
		})
	}
}
