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

type javascriptContext struct {
	URL           string
	Method        string
	HeadersSingle []headerField
	HeadersMulti  []headerField
	Body          string
}

const jsTemplateStr = `const url = {{.URL}};
{{if .HeadersMulti}}const headers = new Headers();
{{range .HeadersMulti}}headers.append({{.Key}}, {{.Value}});
{{end}}
{{end}}const options = {
  method: {{.Method}}{{if .HeadersSingle}},
  headers: {
{{range $i, $h := .HeadersSingle}}{{if $i}},
{{end}}    {{$h.Key}}: {{$h.Value}}{{end}}
  }{{else if .HeadersMulti}},
  headers: headers{{end}}{{if .Body}},
  body: {{.Body}}{{end}}
};

fetch(url, options)
  .then(res => {
    console.log(res.status);
    return res.text();
  })
  .then(text => console.log(text))
  .catch(err => console.error(err));
`

var javascriptTemplate = template.Must(template.New("js").Parse(jsTemplateStr))

// generateJavascript builds a JavaScript fetch snippet using text/template.
func generateJavascript(data *request.Data) (string, error) {
	method := defaultMethod(data.Method)
	hasMulti := hasMultiValueHeaders(data.Headers)
	keys := sortedHeaderKeys(data.Headers)

	var headersSingle []headerField
	var headersMulti []headerField

	for _, k := range keys {
		vals := data.Headers[k]
		if hasMulti {
			for _, v := range vals {
				headersMulti = append(headersMulti, headerField{Key: strconv.Quote(k), Value: strconv.Quote(v)})
			}
		} else {
			if len(vals) > 0 {
				headersSingle = append(headersSingle, headerField{Key: strconv.Quote(k), Value: strconv.Quote(vals[0])})
			}
		}
	}

	bodyVal := ""
	if data.Body != "" {
		bodyVal = strconv.Quote(data.Body)
	}

	ctx := javascriptContext{
		URL:           strconv.Quote(data.URL),
		Method:        strconv.Quote(method),
		HeadersSingle: headersSingle,
		HeadersMulti:  headersMulti,
		Body:          bodyVal,
	}

	var sb strings.Builder
	if err := javascriptTemplate.Execute(&sb, ctx); err != nil {
		return "", err
	}
	return sb.String(), nil
}
