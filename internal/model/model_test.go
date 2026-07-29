package model

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cstaaben/go-rest/internal/model/target"
	notify "github.com/cstaaben/go-rest/internal/ui/notify"
)

func TestModel_Update_NotifyMsg(t *testing.T) {
	m := New(nil)
	require.NotNil(t, m)

	assert.NotPanics(t, func() {
		_, _ = m.Update(notify.Msg{Message: "test message", Level: notify.Info})
	})

	assert.NotPanics(t, func() {
		_, _ = m.Update(&notify.Msg{Message: "test message ptr", Level: notify.Info})
	})
}

func TestModel_Update_WindowSizeMsg(t *testing.T) {
	m := New(nil)
	require.NotNil(t, m)

	winMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	assert.NotPanics(t, func() {
		updatedModel, cmd := m.Update(winMsg)
		assert.NotNil(t, updatedModel)
		_ = cmd
	})
}

func TestModel_View_EnvironmentView(t *testing.T) {
	m := New(nil)
	require.NotNil(t, m)
	m.CurrentView = target.EnvironmentView

	assert.NotPanics(t, func() {
		view := m.View()
		assert.NotNil(t, view)
	})
}
