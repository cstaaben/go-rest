package environment

import (
	"fmt"
	"regexp"

	"github.com/cstaaben/go-rest/internal/variables"
)

var placeholderRegex = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.-]+)\s*\}\}`)

func stringifyPrimitive(val any) (string, bool) {
	switch v := val.(type) {
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprintf("%v", v), true
	default:
		return "", false
	}
}

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
		sub := placeholderRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		key := sub[1]

		// Check session first
		if session != nil {
			if val, ok := session.Get(key); ok {
				if strVal, ok := stringifyPrimitive(val); ok {
					return strVal
				}
			}
		}

		// Fallback to environment variables
		val, ok := variablesMap[key]
		if ok {
			if strVal, ok := stringifyPrimitive(val); ok {
				return strVal
			}
		}

		if !seenMissing[key] {
			seenMissing[key] = true
			missingVars = append(missingVars, key)
		}
		return match
	})

	return result, missingVars
}
