package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/environment"
	"github.com/cstaaben/go-rest/internal/request"
	"github.com/cstaaben/go-rest/internal/variables"
)

func TestData_Substitute(t *testing.T) {
	env := &environment.Environment{
		Variables: map[string]any{
			"host":     "https://api.example.com",
			"method":   "POST",
			"auth_key": "Authorization",
			"token":    "env-token",
		},
	}

	session := variables.NewSession()
	session.Set("token", "session-token")
	session.Set("body_param", "hello-world")

	tests := []struct {
		name            string
		input           *request.Data
		expected        *request.Data
		expectedMissing []string
	}{
		{
			name: "full data substitution success",
			input: &request.Data{
				URL:    "{{host}}/users",
				Method: "{{method}}",
				Headers: map[string][]string{
					"{{auth_key}}": {"Bearer {{token}}"},
				},
				Body: `{"data": "{{body_param}}"}`,
			},
			expected: &request.Data{
				URL:    "https://api.example.com/users",
				Method: "POST",
				Headers: map[string][]string{
					"Authorization": {"Bearer session-token"},
				},
				Body: `{"data": "hello-world"}`,
			},
		},
		{
			name: "missing variables validation error",
			input: &request.Data{
				URL:    "{{host}}/users/{{missing_id}}",
				Method: "GET",
				Headers: map[string][]string{
					"X-Header-{{missing_header}}": {"some-value"},
				},
				Body: `{"data": "{{missing_body}}"}`,
			},
			expectedMissing: []string{"missing_id", "missing_header", "missing_body"},
		},
		{
			name:     "nil data returns nil",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.input.Substitute(env, session)
			if len(tc.expectedMissing) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "missing variables")
				for _, m := range tc.expectedMissing {
					assert.Contains(t, err.Error(), m)
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}

func TestRequest_Substitute(t *testing.T) {
	env := &environment.Environment{
		Variables: map[string]any{
			"host": "https://api.example.com",
		},
	}

	tests := []struct {
		name            string
		input           *request.Request
		expected        *request.Request
		expectedMissing []string
	}{
		{
			name: "request substitution success",
			input: &request.Request{
				ID:   "req-1",
				Name: "Get Users",
				Desc: "Test request",
				Data: &request.Data{
					URL: "{{host}}/users",
				},
			},
			expected: &request.Request{
				ID:   "req-1",
				Name: "Get Users",
				Desc: "Test request",
				Data: &request.Data{
					URL: "https://api.example.com/users",
				},
			},
		},
		{
			name:     "nil request returns nil",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.input.Substitute(env, nil)
			if len(tc.expectedMissing) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "missing variables")
				for _, m := range tc.expectedMissing {
					assert.Contains(t, err.Error(), m)
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}
