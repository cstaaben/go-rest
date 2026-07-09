package environment_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/environment"
	"github.com/cstaaben/go-rest/internal/variables"
)

func TestSubstitute(t *testing.T) {
	tests := []struct {
		name         string
		env          *environment.Environment
		session      *variables.Session
		input        string
		expected     string
		expectedMiss []string
	}{
		{
			name: "plain URL substitution",
			env: &environment.Environment{
				Variables: map[string]any{
					"host": "https://api.example.com",
				},
			},
			input:    "{{host}}/users",
			expected: "https://api.example.com/users",
		},
		{
			name: "spacing in brackets",
			env: &environment.Environment{
				Variables: map[string]any{
					"host": "https://api.example.com",
				},
			},
			input:    "{{   host   }}/users",
			expected: "https://api.example.com/users",
		},
		{
			name: "primitive values casting",
			env: &environment.Environment{
				Variables: map[string]any{
					"token":    12345,
					"is_admin": true,
				},
			},
			input:    "Bearer {{token}} admin={{is_admin}}",
			expected: "Bearer 12345 admin=true",
		},
		{
			name: "session variables take precedence",
			env: &environment.Environment{
				Variables: map[string]any{
					"token": "env-token",
					"host":  "https://api.env.com",
				},
			},
			session: func() *variables.Session {
				s := variables.NewSession()
				s.Set("token", "session-token")
				return s
			}(),
			input:    "{{host}} with {{token}}",
			expected: "https://api.env.com with session-token",
		},
		{
			name: "missing variables collected",
			env: &environment.Environment{
				Variables: map[string]any{
					"host": "https://api.env.com",
				},
			},
			input:        "{{host}}/users/{{id}}/details/{{token}}",
			expected:     "https://api.env.com/users/{{id}}/details/{{token}}",
			expectedMiss: []string{"id", "token"},
		},
		{
			name:         "nil environment with placeholders fails",
			env:          nil,
			input:        "{{host}}",
			expected:     "{{host}}",
			expectedMiss: []string{"host"},
		},
		{
			name:     "nil environment without placeholders succeeds",
			env:      nil,
			input:    "https://api.example.com",
			expected: "https://api.example.com",
		},
		{
			name:     "malformed placeholder ignored",
			env:      &environment.Environment{},
			input:    "{{a}b}}",
			expected: "{{a}b}}",
		},
		{
			name:     "out-of-scope placeholder ignored",
			env:      &environment.Environment{},
			input:    "{{$randomInt}}",
			expected: "{{$randomInt}}",
		},
		{
			name: "complex variable type fails",
			env: &environment.Environment{
				Variables: map[string]any{
					"complex_map": map[string]string{"foo": "bar"},
				},
			},
			input:        "{{complex_map}}",
			expected:     "{{complex_map}}",
			expectedMiss: []string{"complex_map"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, missing := tc.env.Substitute(tc.input, tc.session)
			assert.Equal(t, tc.expected, res)
			if len(tc.expectedMiss) > 0 {
				assert.Equal(t, tc.expectedMiss, missing)
			} else {
				assert.Empty(t, missing)
			}
		})
	}
}
