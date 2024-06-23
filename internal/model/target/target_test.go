package target_test

import (
	"testing"

	"github.com/cstaaben/go-rest/internal/model/target"
	"github.com/stretchr/testify/assert"
)

func TestNextTarget(t *testing.T) {
	testCases := []struct {
		name          string
		currentTarget target.Target
		view          target.View
		expected      target.Target
	}{
		{
			name:          "Cycle client target",
			view:          target.ClientView,
			currentTarget: target.ResponseTarget,
			expected:      target.RequestsTarget,
		},
		{
			name:          "Next client target",
			view:          target.ClientView,
			currentTarget: target.RequestsTarget,
			expected:      target.EditorTarget,
		},
		{
			name:          "Cycle environments target",
			view:          target.EnvironmentView,
			currentTarget: target.EnvEditorTarget,
			expected:      target.EnvironmentsTarget,
		},
		{
			name:          "Next environments target",
			view:          target.EnvironmentView,
			currentTarget: target.EnvironmentsTarget,
			expected:      target.EnvEditorTarget,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := target.NextTarget(tc.view, tc.currentTarget)
			assert.Equal(t, tc.expected, actual)
		})
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
			name:          "Cycle client target",
			view:          target.ClientView,
			currentTarget: target.RequestsTarget,
			expected:      target.ResponseTarget,
		},
		{
			name:          "Previous client target",
			view:          target.ClientView,
			currentTarget: target.EditorTarget,
			expected:      target.RequestsTarget,
		},
		{
			name:          "Cycle environments target",
			view:          target.EnvironmentView,
			currentTarget: target.EnvironmentsTarget,
			expected:      target.EnvEditorTarget,
		},
		{
			name:          "Previous environments target",
			view:          target.EnvironmentView,
			currentTarget: target.EnvEditorTarget,
			expected:      target.EnvironmentsTarget,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := target.PrevTarget(tc.view, tc.currentTarget)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestTarget_Equal(t *testing.T) {
	assert.False(t, target.RequestsTarget == target.EnvironmentsTarget)
}
