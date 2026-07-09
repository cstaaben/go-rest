package environment_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/environment"
	"github.com/cstaaben/go-rest/internal/request"
)

func TestSubstitute(t *testing.T) {
	tests := []struct {
		name         string
		env          *environment.Environment
		inputData    *request.Data
		expectedData *request.Data
		expectedErr  string
	}{
		{
			name: "URL success",
			env: &environment.Environment{
				Name: "test",
				Variables: map[string]any{
					"host": "https://api.example.com",
				},
			},
			inputData: &request.Data{
				URL: "{{host}}/users",
			},
			expectedData: &request.Data{
				URL: "https://api.example.com/users",
			},
		},
		{
			name: "all fields success",
			env: &environment.Environment{
				Name: "test",
				Variables: map[string]any{
					"host":       "https://api.example.com",
					"method":     "POST",
					"auth_key":   "Authorization",
					"token":      12345,
					"is_admin":   true,
					"body_param": "hello-world",
				},
			},
			inputData: &request.Data{
				URL:    "{{host}}/users",
				Method: "{{method}}",
				Headers: map[string][]string{
					"{{auth_key}}": {"Bearer {{token}}"},
					"X-Admin":      {"{{is_admin}}"},
				},
				Body: `{"data": "{{body_param}}"}`,
			},
			expectedData: &request.Data{
				URL:    "https://api.example.com/users",
				Method: "POST",
				Headers: map[string][]string{
					"Authorization": {"Bearer 12345"},
					"X-Admin":       {"true"},
				},
				Body: `{"data": "hello-world"}`,
			},
		},
		{
			name: "spaces in placeholders",
			env: &environment.Environment{
				Name: "test",
				Variables: map[string]any{
					"host": "https://api.example.com",
				},
			},
			inputData: &request.Data{
				URL: "{{ host }}/users",
			},
			expectedData: &request.Data{
				URL: "https://api.example.com/users",
			},
		},
		{
			name: "missing variables error",
			env: &environment.Environment{
				Name: "test",
				Variables: map[string]any{
					"host": "https://api.example.com",
				},
			},
			inputData: &request.Data{
				URL: "{{host}}/users/{{id}}/details/{{token}}",
			},
			expectedErr: "missing variables: id, token",
		},
		{
			name: "nil environment succeeds when no placeholders",
			env:  nil,
			inputData: &request.Data{
				URL: "https://api.example.com/users",
			},
			expectedData: &request.Data{
				URL: "https://api.example.com/users",
			},
		},
		{
			name: "nil environment fails when placeholders present",
			env:  nil,
			inputData: &request.Data{
				URL: "{{host}}/users",
			},
			expectedErr: "missing variables: host",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.env.Substitute(tc.inputData)
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, res)
			}
		})
	}
}
