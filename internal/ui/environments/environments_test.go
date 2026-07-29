package environments

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModel_View(t *testing.T) {
	m := New(nil)
	require.NotNil(t, m)

	assert.NotPanics(t, func() {
		viewStr := m.View()
		assert.NotEmpty(t, viewStr)
	})
}

func TestModel_Update_WindowSizeMsg(t *testing.T) {
	m := New(nil)
	require.NotNil(t, m)

	winMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	updatedModel, cmd := m.Update(winMsg)

	assert.NotNil(t, updatedModel)
	assert.Nil(t, cmd)
	assert.Greater(t, m.List.Width(), 0)
	assert.Greater(t, m.List.Height(), 0)
}
