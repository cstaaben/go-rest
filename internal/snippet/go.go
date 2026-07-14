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

type goContext struct {
	URL           string
	Method        string
	HeadersSingle []headerField
	HeadersMulti  []headerField
	Body          string
}

const goTemplateStr = `package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"{{if .Body}}
	"strings"{{end}}
)

func main() {
	url := {{.URL}}
{{if .Body}}	body := {{.Body}}
	req, err := http.NewRequestWithContext(context.Background(), {{.Method}}, url, strings.NewReader(body))
{{else}}	req, err := http.NewRequestWithContext(context.Background(), {{.Method}}, url, nil)
{{end}}	if err != nil {
		log.Println(err)
		return
	}

{{if or .HeadersSingle .HeadersMulti}}{{range .HeadersSingle}}	req.Header.Set({{.Key}}, {{.Value}})
{{end}}{{range .HeadersMulti}}	req.Header.Add({{.Key}}, {{.Value}})
{{end}}
{{end}}	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
}
`

var goTemplate = template.Must(template.New("go").Parse(goTemplateStr))

// generateGo builds a Go net/http snippet using text/template.
func generateGo(data *request.Data) (string, error) {
	var headersSingle []headerField
	var headersMulti []headerField

	keys := sortedHeaderKeys(data.Headers)

	for _, k := range keys {
		vals := data.Headers[k]
		for i, v := range vals {
			if i == 0 {
				headersSingle = append(headersSingle, headerField{Key: strconv.Quote(k), Value: strconv.Quote(v)})
			} else {
				headersMulti = append(headersMulti, headerField{Key: strconv.Quote(k), Value: strconv.Quote(v)})
			}
		}
	}

	bodyVal := ""
	if data.Body != "" {
		bodyVal = strconv.Quote(data.Body)
	}

	ctx := goContext{
		URL:           strconv.Quote(data.URL),
		Method:        getGoMethod(data.Method),
		HeadersSingle: headersSingle,
		HeadersMulti:  headersMulti,
		Body:          bodyVal,
	}

	var sb strings.Builder
	if err := goTemplate.Execute(&sb, ctx); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// getGoMethod returns the standard net/http constant string if applicable, or quotes the custom method.
func getGoMethod(method string) string {
	m := defaultMethod(method)
	switch strings.ToUpper(m) {
	case "GET":
		return "http.MethodGet"
	case "POST":
		return "http.MethodPost"
	case "PUT":
		return "http.MethodPut"
	case "DELETE":
		return "http.MethodDelete"
	case "PATCH":
		return "http.MethodPatch"
	case "HEAD":
		return "http.MethodHead"
	case "OPTIONS":
		return "http.MethodOptions"
	case "CONNECT":
		return "http.MethodConnect"
	case "TRACE":
		return "http.MethodTrace"
	default:
		return strconv.Quote(m)
	}
}
