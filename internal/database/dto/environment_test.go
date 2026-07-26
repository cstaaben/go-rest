package dto_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/database/dto"
)

var (
	_ list.Item        = (*dto.Environment)(nil)
	_ list.DefaultItem = (*dto.Environment)(nil)
)

func TestNewEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		envName string
	}{
		{
			name:    "create environment with name",
			envName: "Production",
		},
		{
			name:    "create environment with empty name",
			envName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := dto.NewEnvironment(tt.envName)
			assert.NotNil(t, env)
			assert.Equal(t, tt.envName, env.Name)
			assert.NotNil(t, env.Variables)
			assert.Empty(t, env.Variables)
		})
	}
}

func TestEnvironment_MetadataMethods(t *testing.T) {
	tests := []struct {
		name            string
		env             *dto.Environment
		wantFilterValue string
		wantTitle       string
		wantDesc        string
	}{
		{
			name: "populated environment",
			env: &dto.Environment{
				UUID: "env-1",
				Name: "Staging",
			},
			wantFilterValue: "Staging",
			wantTitle:       "Staging",
			wantDesc:        "Staging",
		},
		{
			name:            "empty environment",
			env:             &dto.Environment{},
			wantFilterValue: "",
			wantTitle:       "",
			wantDesc:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFilterValue, tt.env.FilterValue())
			assert.Equal(t, tt.wantTitle, tt.env.Title())
			assert.Equal(t, tt.wantDesc, tt.env.Description())
		})
	}
}

func TestEnvironment_Equal(t *testing.T) {
	baseEnv := &dto.Environment{
		UUID: "env-10",
		Name: "Production",
		Variables: []*dto.EnvironmentVariable{
			{EnvironmentUUID: "env-10", Name: "API_KEY", Value: []byte("secret")},
		},
	}

	tests := []struct {
		name      string
		env       *dto.Environment
		other     *dto.Environment
		wantEqual bool
	}{
		{
			name:      "identical environments",
			env:       baseEnv,
			other:     baseEnv,
			wantEqual: true,
		},
		{
			name: "equal copy",
			env:  baseEnv,
			other: &dto.Environment{
				UUID: "env-10",
				Name: "Production",
				Variables: []*dto.EnvironmentVariable{
					{EnvironmentUUID: "env-10", Name: "API_KEY", Value: []byte("secret")},
				},
			},
			wantEqual: true,
		},
		{
			name: "different UUID",
			env:  baseEnv,
			other: &dto.Environment{
				UUID:      "env-11",
				Name:      baseEnv.Name,
				Variables: baseEnv.Variables,
			},
			wantEqual: false,
		},
		{
			name:      "nil comparison",
			env:       baseEnv,
			other:     nil,
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEqual, tt.env.Equal(tt.other))
		})
	}
}

func TestEnvironmentVariable_Equal(t *testing.T) {
	baseVar := &dto.EnvironmentVariable{
		EnvironmentUUID: "env-10",
		Name:            "BASE_URL",
		Value:           []byte("https://api.example.com"),
	}

	tests := []struct {
		name      string
		envVar    *dto.EnvironmentVariable
		other     *dto.EnvironmentVariable
		wantEqual bool
	}{
		{
			name:      "identical variables",
			envVar:    baseVar,
			other:     baseVar,
			wantEqual: true,
		},
		{
			name: "equal copy",
			envVar: baseVar,
			other: &dto.EnvironmentVariable{
				EnvironmentUUID: "env-10",
				Name:            "BASE_URL",
				Value:           []byte("https://api.example.com"),
			},
			wantEqual: true,
		},
		{
			name: "different name",
			envVar: baseVar,
			other: &dto.EnvironmentVariable{
				EnvironmentUUID: "env-10",
				Name:            "OTHER_URL",
				Value:           []byte("https://api.example.com"),
			},
			wantEqual: false,
		},
		{
			name: "different value",
			envVar: baseVar,
			other: &dto.EnvironmentVariable{
				EnvironmentUUID: "env-10",
				Name:            "BASE_URL",
				Value:           []byte("https://other.example.com"),
			},
			wantEqual: false,
		},
		{
			name:      "both nil",
			envVar:    nil,
			other:     nil,
			wantEqual: true,
		},
		{
			name:      "one nil",
			envVar:    baseVar,
			other:     nil,
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEqual, tt.envVar.Equal(tt.other))
		})
	}
}
