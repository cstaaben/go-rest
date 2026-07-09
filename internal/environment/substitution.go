package environment

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cstaaben/go-rest/internal/request"
)

var placeholderRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Substitute replaces {{variable_name}} placeholders in URL, HTTP method, header keys, header values, and body
// using the environment's variables. If any placeholder cannot be resolved, it returns an error listing the missing variables.
func (e *Environment) Substitute(data *request.Data) (*request.Data, error) {
	if data == nil {
		return nil, nil
	}

	variables := make(map[string]any)
	if e != nil && e.Variables != nil {
		variables = e.Variables
	}

	var missingVars []string
	seenMissing := make(map[string]bool)

	substituteStr := func(input string) string {
		return placeholderRegex.ReplaceAllStringFunc(input, func(match string) string {
			key := match[2 : len(match)-2]
			key = strings.TrimSpace(key)
			if val, ok := variables[key]; ok {
				return fmt.Sprintf("%v", val)
			}
			if !seenMissing[key] {
				seenMissing[key] = true
				missingVars = append(missingVars, key)
			}
			return match
		})
	}

	subURL := substituteStr(data.URL)
	subMethod := substituteStr(data.Method)
	subBody := substituteStr(data.Body)

	var subHeaders map[string][]string
	if data.Headers != nil {
		subHeaders = make(map[string][]string)
		for k, vals := range data.Headers {
			subK := substituteStr(k)
			var subVals []string
			for _, val := range vals {
				subVals = append(subVals, substituteStr(val))
			}
			subHeaders[subK] = subVals
		}
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing variables: %s", strings.Join(missingVars, ", "))
	}

	return &request.Data{
		URL:     subURL,
		Method:  subMethod,
		Headers: subHeaders,
		Proto:   data.Proto,
		Body:    subBody,
	}, nil
}
