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
	"strings"
	"text/template"

	"github.com/cstaaben/go-rest/internal/request"
)

type curlContext struct {
	Method  string
	Headers []headerField
	Body    string
	URL     string
}

// Wrap all values in single quotes to prevent command injection and shell expansion.
const curlTemplateStr = `curl -X '{{.Method}}' {{range .Headers}}-H '{{.Key}}: {{.Value}}' {{end}}{{if .Body}}-d '{{.Body}}' {{end}}'{{.URL}}'`

var curlTemplate = template.Must(template.New("curl").Parse(curlTemplateStr))

// generateCurl builds a curl command snippet using text/template.
func generateCurl(data *request.Data) (string, error) {
	escapeShellString := func(s string) string {
		return strings.ReplaceAll(s, "'", `'\''`)
	}

	method := escapeShellString(defaultMethod(data.Method))

	var headers []headerField
	keys := sortedHeaderKeys(data.Headers)

	for _, k := range keys {
		for _, v := range data.Headers[k] {
			headers = append(headers, headerField{
				Key:   escapeShellString(k),
				Value: escapeShellString(v),
			})
		}
	}

	body := escapeShellString(data.Body)
	url := escapeShellString(data.URL)

	ctx := curlContext{
		Method:  method,
		Headers: headers,
		Body:    body,
		URL:     url,
	}

	var sb strings.Builder
	if err := curlTemplate.Execute(&sb, ctx); err != nil {
		return "", err
	}
	return sb.String(), nil
}
