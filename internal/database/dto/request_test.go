package dto_test

import (
	"database/sql"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/database/dto"
)

var (
	_ list.Item        = (*dto.Request)(nil)
	_ list.DefaultItem = (*dto.Request)(nil)
	_ list.Item        = (*dto.Group)(nil)
	_ list.DefaultItem = (*dto.Group)(nil)
)

func TestRequest_MetadataMethods(t *testing.T) {
	tests := []struct {
		name            string
		req             *dto.Request
		wantFilterValue string
		wantTitle       string
		wantDesc        string
	}{
		{
			name: "populated request",
			req: &dto.Request{
				UUID: "r-1",
				Name: "Get Users",
				Desc: "Fetch all user accounts",
			},
			wantFilterValue: "Get Users",
			wantTitle:       "Get Users",
			wantDesc:        "Fetch all user accounts",
		},
		{
			name:            "empty request",
			req:             &dto.Request{},
			wantFilterValue: "",
			wantTitle:       "",
			wantDesc:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFilterValue, tt.req.FilterValue())
			assert.Equal(t, tt.wantTitle, tt.req.Title())
			assert.Equal(t, tt.wantDesc, tt.req.Description())
		})
	}
}

func TestRequest_Equal(t *testing.T) {
	baseReq := &dto.Request{
		UUID:      "r-100",
		Name:      "Create Order",
		Desc:      "POST request to create order",
		GroupUUID: "g-50",
		Data: &dto.Data{
			URL:     "https://api.example.com/orders",
			Headers: map[string][]string{"Content-Type": {"application/json"}},
			Method:  "POST",
			Proto:   "HTTP/1.1",
			Body:    `{"item_id": 42}`,
		},
	}

	tests := []struct {
		name      string
		req       *dto.Request
		other     *dto.Request
		wantEqual bool
	}{
		{
			name:      "identical requests",
			req:       baseReq,
			other:     baseReq,
			wantEqual: true,
		},
		{
			name: "equal copy of request",
			req:  baseReq,
			other: &dto.Request{
				UUID:      "r-100",
				Name:      "Create Order",
				Desc:      "POST request to create order",
				GroupUUID: "g-50",
				Data: &dto.Data{
					URL:     "https://api.example.com/orders",
					Headers: map[string][]string{"Content-Type": {"application/json"}},
					Method:  "POST",
					Proto:   "HTTP/1.1",
					Body:    `{"item_id": 42}`,
				},
			},
			wantEqual: true,
		},
		{
			name: "different UUID",
			req:  baseReq,
			other: &dto.Request{
				UUID:      "r-101",
				Name:      baseReq.Name,
				Desc:      baseReq.Desc,
				GroupUUID: baseReq.GroupUUID,
				Data:      baseReq.Data,
			},
			wantEqual: false,
		},
		{
			name: "different data body",
			req:  baseReq,
			other: &dto.Request{
				UUID:      baseReq.UUID,
				Name:      baseReq.Name,
				Desc:      baseReq.Desc,
				GroupUUID: baseReq.GroupUUID,
				Data: &dto.Data{
					URL:     baseReq.Data.URL,
					Headers: baseReq.Data.Headers,
					Method:  baseReq.Data.Method,
					Proto:   baseReq.Data.Proto,
					Body:    `{"item_id": 99}`,
				},
			},
			wantEqual: false,
		},
		{
			name:      "nil comparison",
			req:       baseReq,
			other:     nil,
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEqual, tt.req.Equal(tt.other))
		})
	}
}

func TestRequestListData_Descriptions(t *testing.T) {
	tests := []struct {
		name                 string
		data                 *dto.RequestListData
		wantReqDescription   string
		wantGroupDescription string
	}{
		{
			name: "valid descriptions",
			data: &dto.RequestListData{
				RequestDesc: sql.NullString{String: "Request details", Valid: true},
				GroupDesc:   sql.NullString{String: "Group details", Valid: true},
			},
			wantReqDescription:   "Request details",
			wantGroupDescription: "Group details",
		},
		{
			name: "null descriptions",
			data: &dto.RequestListData{
				RequestDesc: sql.NullString{String: "", Valid: false},
				GroupDesc:   sql.NullString{String: "", Valid: false},
			},
			wantReqDescription:   "",
			wantGroupDescription: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantReqDescription, tt.data.RequestDescription())
			assert.Equal(t, tt.wantGroupDescription, tt.data.GroupDescription())
		})
	}
}
