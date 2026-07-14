/*
 * go-rest - A TUI for a REST client
 * Copyright (C) 2026  Corbin Staaben
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package snippet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/request"
	"github.com/cstaaben/go-rest/internal/snippet"
)

func TestGenerate(t *testing.T) {
	reqStandard := &request.Request{
		Data: &request.Data{
			URL:    "http://example.com/api?foo=bar",
			Method: "POST",
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"X-Custom":     {"val1", "val2"},
			},
			Body: `{"key": 'value'}`,
		},
	}

	reqSimple := &request.Request{
		Data: &request.Data{
			URL:    "http://example.com",
			Method: "GET",
		},
	}

	reqPlaceholders := &request.Request{
		Data: &request.Data{
			URL:    "{{host}}/api?foo={{param}}",
			Method: "POST",
			Headers: map[string][]string{
				"Authorization": {"Bearer {{token}}"},
			},
			Body: `{"key": "{{body_val}}"}`,
		},
	}

	tests := []struct {
		name     string
		request  *request.Request
		target   snippet.Target
		expected string
		contains []string
		wantErr  bool
	}{
		{
			name:    "nil request returns error",
			request: nil,
			target:  snippet.TargetCurl,
			wantErr: true,
		},
		{
			name:    "nil data returns error",
			request: &request.Request{},
			target:  snippet.TargetCurl,
			wantErr: true,
		},
		{
			name:    "unsupported target returns error",
			request: reqSimple,
			target:  snippet.Target("invalid"),
			wantErr: true,
		},
		{
			name:     "curl target - standard request",
			request:  reqStandard,
			target:   snippet.TargetCurl,
			expected: `curl -X 'POST' -H 'Content-Type: application/json' -H 'X-Custom: val1' -H 'X-Custom: val2' -d '{"key": '\''value'\''}' 'http://example.com/api?foo=bar'`,
		},
		{
			name:    "go target - standard request",
			request: reqStandard,
			target:  snippet.TargetGo,
			contains: []string{
				"package main",
				`url := "http://example.com/api?foo=bar"`,
				`body := "{\"key\": 'value'}"`,
				`req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, strings.NewReader(body))`,
				`req.Header.Set("Content-Type", "application/json")`,
				`req.Header.Set("X-Custom", "val1")`,
				`req.Header.Add("X-Custom", "val2")`,
				`log.Println(err)`,
				`return`,
			},
		},
		{
			name:    "python target - standard request",
			request: reqStandard,
			target:  snippet.TargetPython,
			contains: []string{
				"import requests",
				`url = "http://example.com/api?foo=bar"`,
				`    "Content-Type": "application/json",`,
				`    "X-Custom": "val1, val2",`,
				`payload = "{\"key\": 'value'}"`,
				`response = requests.request("POST", url, headers=headers, data=payload)`,
			},
		},
		{
			name:    "javascript target - standard request",
			request: reqStandard,
			target:  snippet.TargetJavascript,
			contains: []string{
				`const url = "http://example.com/api?foo=bar";`,
				`const headers = new Headers();`,
				`headers.append("Content-Type", "application/json");`,
				`headers.append("X-Custom", "val1");`,
				`headers.append("X-Custom", "val2");`,
				`  body: "{\"key\": 'value'}"`,
			},
		},
		{
			name:     "curl target - placeholders request",
			request:  reqPlaceholders,
			target:   snippet.TargetCurl,
			expected: `curl -X 'POST' -H 'Authorization: Bearer {{token}}' -d '{"key": "{{body_val}}"}' '{{host}}/api?foo={{param}}'`,
		},
		{
			name:    "go target - placeholders request",
			request: reqPlaceholders,
			target:  snippet.TargetGo,
			contains: []string{
				`url := "{{host}}/api?foo={{param}}"`,
				`body := "{\"key\": \"{{body_val}}\"}"`,
				`req.Header.Set("Authorization", "Bearer {{token}}")`,
			},
		},
		{
			name:    "python target - placeholders request",
			request: reqPlaceholders,
			target:  snippet.TargetPython,
			contains: []string{
				`url = "{{host}}/api?foo={{param}}"`,
				`    "Authorization": "Bearer {{token}}",`,
				`payload = "{\"key\": \"{{body_val}}\"}"`,
			},
		},
		{
			name:    "javascript target - placeholders request",
			request: reqPlaceholders,
			target:  snippet.TargetJavascript,
			contains: []string{
				`const url = "{{host}}/api?foo={{param}}";`,
				`    "Authorization": "Bearer {{token}}"`,
				`  body: "{\"key\": \"{{body_val}}\"}"`,
			},
		},
		{
			name: "curl target - shell special characters",
			request: &request.Request{
				Data: &request.Data{
					URL:    "http://example.com/api?foo=$bar&baz=`echo 1`",
					Method: "GET",
					Headers: map[string][]string{
						"X-Shell": {`"double"`, `val's`},
					},
				},
			},
			target:   snippet.TargetCurl,
			expected: `curl -X 'GET' -H 'X-Shell: "double"' -H 'X-Shell: val'\''s' 'http://example.com/api?foo=$bar&baz=` + "`echo 1`" + `'`,
		},
		{
			name:     "curl target - simple empty request",
			request:  reqSimple,
			target:   snippet.TargetCurl,
			expected: `curl -X 'GET' 'http://example.com'`,
		},
		{
			name:    "go target - simple empty request",
			request: reqSimple,
			target:  snippet.TargetGo,
			contains: []string{
				`url := "http://example.com"`,
				`req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)`,
			},
		},
		{
			name:    "python target - simple empty request",
			request: reqSimple,
			target:  snippet.TargetPython,
			contains: []string{
				`url = "http://example.com"`,
				`response = requests.request("GET", url)`,
			},
		},
		{
			name:    "javascript target - simple empty request",
			request: reqSimple,
			target:  snippet.TargetJavascript,
			contains: []string{
				`const url = "http://example.com";`,
				`  method: "GET"`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := snippet.Generate(tc.request, tc.target)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.expected != "" {
					assert.Equal(t, tc.expected, res)
				}
				for _, part := range tc.contains {
					assert.Contains(t, res, part)
				}
			}
		})
	}
}
