package request_test

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/cstaaben/go-rest/internal/request"
)

var (
	_ list.Item        = (*request.Request)(nil)
	_ list.DefaultItem = (*request.Request)(nil)
	_ list.Item        = (*request.Group)(nil)
	_ list.DefaultItem = (*request.Group)(nil)
)
