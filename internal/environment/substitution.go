package environment

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cstaaben/go-rest/internal/variables"
)

var placeholderRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Substitute replaces {{variable_name}} placeholders in the input string using the environment's variables
// and the provided session variables (which take precedence). It returns the substituted string and a slice
// of unique unresolved placeholder names.
func (e *Environment) Substitute(input string, session *variables.Session) (string, []string) {
	variablesMap := make(map[string]any)
	if e != nil && e.Variables != nil {
		for k, v := range e.Variables {
			variablesMap[k] = v
		}
	}

	var missingVars []string
	seenMissing := make(map[string]bool)

	result := placeholderRegex.ReplaceAllStringFunc(input, func(match string) string {
		key := match[2 : len(match)-2]
		key = strings.TrimSpace(key)

		// Check session first
		if session != nil {
			if val, ok := session.Get(key); ok {
				return fmt.Sprintf("%v", val)
			}
		}

		// Fallback to environment variables
		val, ok := variablesMap[key]
		if !ok {
			if !seenMissing[key] {
				seenMissing[key] = true
				missingVars = append(missingVars, key)
			}
			return match
		}
		return fmt.Sprintf("%v", val)
	})

	return result, missingVars
}
