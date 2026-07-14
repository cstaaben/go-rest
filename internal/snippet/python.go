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

package snippet

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/cstaaben/go-rest/internal/request"
)

type pythonContext struct {
	URL     string
	Method  string
	Headers []headerField
	Body    string
}

const pythonTemplateStr = `import requests

url = {{.URL}}
{{if .Headers}}
headers = {
{{range .Headers}}    {{.Key}}: {{.Value}},
{{end}}}
{{end}}{{if .Body}}
payload = {{.Body}}
{{end}}
response = requests.request({{.Method}}, url{{if .Headers}}, headers=headers{{end}}{{if .Body}}, data=payload{{end}})

print(response.status_code)
print(response.text)
`

var pythonTemplate = template.Must(template.New("py").Parse(pythonTemplateStr))

// generatePython builds a Python requests snippet using text/template.
func generatePython(data *request.Data) (string, error) {
	method := defaultMethod(data.Method)
	keys := sortedHeaderKeys(data.Headers)

	var headers []headerField

	for _, k := range keys {
		vals := data.Headers[k]
		if len(vals) > 0 {
			// Comma-join multi-value headers to conform with python requests dict expectations
			joinedVal := strings.Join(vals, ", ")
			headers = append(headers, headerField{Key: strconv.Quote(k), Value: strconv.Quote(joinedVal)})
		}
	}

	bodyVal := ""
	if data.Body != "" {
		bodyVal = strconv.Quote(data.Body)
	}

	ctx := pythonContext{
		URL:     strconv.Quote(data.URL),
		Method:  strconv.Quote(method),
		Headers: headers,
		Body:    bodyVal,
	}

	var sb strings.Builder
	if err := pythonTemplate.Execute(&sb, ctx); err != nil {
		return "", err
	}
	return sb.String(), nil
}
