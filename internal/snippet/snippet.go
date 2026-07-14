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

// Package snippet handles code snippet generation for requests.
package snippet

import (
	"fmt"
	"sort"

	"github.com/cstaaben/go-rest/internal/request"
)

// Target represents a supported target code language format.
type Target string

const (
	// TargetCurl represents the curl command-line utility target.
	TargetCurl Target = "curl"
	// TargetGo represents the Go net/http library target.
	TargetGo Target = "go"
	// TargetPython represents the Python requests library target.
	TargetPython Target = "python"
	// TargetJavascript represents the JavaScript fetch API target.
	TargetJavascript Target = "javascript"
)

type headerField struct {
	Key   string
	Value string
}

// Generate compiles a Request definition into equivalent code for the specified target.
func Generate(req *request.Request, target Target) (string, error) {
	if req == nil {
		return "", fmt.Errorf("request cannot be nil")
	}
	if req.Data == nil {
		return "", fmt.Errorf("request data cannot be nil")
	}

	switch target {
	case TargetCurl:
		return generateCurl(req.Data)
	case TargetGo:
		return generateGo(req.Data)
	case TargetPython:
		return generatePython(req.Data)
	case TargetJavascript:
		return generateJavascript(req.Data)
	default:
		return "", fmt.Errorf("unsupported snippet target: %s", target)
	}
}

// sortedHeaderKeys returns the keys of the headers map sorted alphabetically.
func sortedHeaderKeys(headers map[string][]string) []string {
	var keys []string
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// hasMultiValueHeaders returns true if any header key has more than one value.
func hasMultiValueHeaders(headers map[string][]string) bool {
	for _, vals := range headers {
		if len(vals) > 1 {
			return true
		}
	}
	return false
}

// defaultMethod returns the HTTP method, defaulting to "GET" if empty.
func defaultMethod(method string) string {
	if method == "" {
		return "GET"
	}
	return method
}
