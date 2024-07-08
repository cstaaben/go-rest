package target_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/model/target"
)

func TestNextTarget(t *testing.T) {
	testCases := []struct {
		name          string
		currentTarget target.Target
		view          target.View
		expected      target.Target
	}{
		{
			name:          "Cycle client focusedTarget",
			view:          target.ClientView,
			currentTarget: target.ResponseTarget,
			expected:      target.RequestsTarget,
		},
		{
			name:          "Next client focusedTarget",
			view:          target.ClientView,
			currentTarget: target.RequestsTarget,
			expected:      target.EditorTarget,
		},
		{
			name:          "Cycle environments focusedTarget",
			view:          target.EnvironmentView,
			currentTarget: target.EnvEditorTarget,
			expected:      target.EnvironmentsTarget,
		},
		{
			name:          "Next environments focusedTarget",
			view:          target.EnvironmentView,
			currentTarget: target.EnvironmentsTarget,
			expected:      target.EnvEditorTarget,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(
			tc.name, func(t *testing.T) {
				actual := target.NextTarget(tc.view, tc.currentTarget)
				assert.Equal(t, tc.expected, actual)
			},
		)
	}
}

func TestPreviousTarget(t *testing.T) {
	testCases := []struct {
		name          string
		currentTarget target.Target
		view          target.View
		expected      target.Target
	}{
		{
			name:          "Cycle client focusedTarget",
			view:          target.ClientView,
			currentTarget: target.RequestsTarget,
			expected:      target.ResponseTarget,
		},
		{
			name:          "Previous client focusedTarget",
			view:          target.ClientView,
			currentTarget: target.EditorTarget,
			expected:      target.RequestsTarget,
		},
		{
			name:          "Cycle environments focusedTarget",
			view:          target.EnvironmentView,
			currentTarget: target.EnvironmentsTarget,
			expected:      target.EnvEditorTarget,
		},
		{
			name:          "Previous environments focusedTarget",
			view:          target.EnvironmentView,
			currentTarget: target.EnvEditorTarget,
			expected:      target.EnvironmentsTarget,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(
			tc.name, func(t *testing.T) {
				actual := target.PrevTarget(tc.view, tc.currentTarget)
				assert.Equal(t, tc.expected, actual)
			},
		)
	}
}

func TestTarget_Equal(t *testing.T) {
	assert.False(t, target.RequestsTarget == target.EnvironmentsTarget)
}

func TestNextView(t *testing.T) {
	testCases := []struct {
		name     string
		current  target.View
		expected target.View
	}{
		{
			name:     "Next focusedView",
			current:  target.ClientView,
			expected: target.EnvironmentView,
		},
		{
			name:     "Cycle focusedView",
			current:  target.EnvironmentView,
			expected: target.ClientView,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(
			tc.name, func(t *testing.T) {
				actual := target.NextView(tc.current)
				assert.Equal(t, tc.expected, actual)
			},
		)
	}
}

func TestPrevView(t *testing.T) {
	testCases := []struct {
		name     string
		current  target.View
		expected target.View
	}{
		{
			name:     "Previous focusedView",
			current:  target.EnvironmentView,
			expected: target.ClientView,
		},
		{
			name:     "Cycle focusedView",
			current:  target.ClientView,
			expected: target.EnvironmentView,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(
			tc.name, func(t *testing.T) {
				actual := target.PrevView(tc.current)
				assert.Equal(t, tc.expected, actual)
			},
		)
	}
}

func TestChangeFocus(t *testing.T) {
	testCases := []struct {
		name            string
		focusedView     target.View
		focusedTarget   target.Target
		unfocusedView   target.View
		unfocusedTarget target.Target
		expectedMsg     target.FocusMsg
		expectingCmd    bool
	}{
		{
			name:            "Change focus to requests",
			focusedView:     target.ClientView,
			focusedTarget:   target.RequestsTarget,
			unfocusedView:   target.ClientView,
			unfocusedTarget: target.ResponseTarget,
			expectedMsg: target.FocusMsg{
				FocusedView:     target.ClientView,
				FocusedTarget:   target.RequestsTarget,
				UnfocusedView:   target.ClientView,
				UnfocusedTarget: target.ResponseTarget,
			},
			expectingCmd: true,
		},
		{
			name:            "Change focus to environments",
			focusedView:     target.EnvironmentView,
			focusedTarget:   target.EnvironmentsTarget,
			unfocusedView:   target.ClientView,
			unfocusedTarget: target.RequestsTarget,
			expectedMsg: target.FocusMsg{
				FocusedView:     target.EnvironmentView,
				FocusedTarget:   target.EnvironmentsTarget,
				UnfocusedView:   target.ClientView,
				UnfocusedTarget: target.RequestsTarget,
			},
			expectingCmd: true,
		},
		{
			name:            "Invalid focusedTarget for focusedView",
			focusedView:     target.ClientView,
			focusedTarget:   target.EnvironmentsTarget,
			unfocusedView:   target.ClientView,
			unfocusedTarget: target.RequestsTarget,
			expectingCmd:    false,
		},
		{
			name:            "Invalid focusedView",
			focusedView:     target.View(-1),
			focusedTarget:   target.Target(-1),
			unfocusedView:   target.ClientView,
			unfocusedTarget: target.RequestsTarget,
			expectingCmd:    false,
		},
		{
			name:            "Invalid unfocusedView",
			focusedView:     target.ClientView,
			focusedTarget:   target.RequestsTarget,
			unfocusedView:   target.View(-1),
			unfocusedTarget: target.Target(-1),
			expectingCmd:    false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(
			tc.name, func(t *testing.T) {
				cmd := target.ChangeFocus(tc.focusedView, tc.focusedTarget, tc.unfocusedView, tc.unfocusedTarget)
				if !tc.expectingCmd {
					assert.Nil(t, cmd, "unexpected command present")
					return
				}

				actual := cmd()
				assert.Equal(t, tc.expectedMsg, actual)
			},
		)
	}
}
