package request

import (
	"fmt"
	"strings"

	"github.com/cstaaben/go-rest/internal/environment"
	"github.com/cstaaben/go-rest/internal/variables"
)

// Substitute resolves all placeholders in the Request using the active environment and session variables.
// It returns a new copied Request. If any placeholders are unresolved, it returns a validation error.
func (r *Request) Substitute(env *environment.Environment, session *variables.Session) (*Request, error) {
	if r == nil {
		return nil, nil
	}

	subData, err := r.Data.Substitute(env, session)
	if err != nil {
		return nil, err
	}

	return &Request{
		ID:       r.ID,
		Name:     r.Name,
		Desc:     r.Desc,
		Data:     subData,
		Response: r.Response,
	}, nil
}

// Substitute resolves all placeholders in the request data using the active environment and session variables.
// It returns a new copied Data struct. If any placeholders are unresolved, it returns a validation error.
func (d *Data) Substitute(env *environment.Environment, session *variables.Session) (*Data, error) {
	if d == nil {
		return nil, nil
	}

	var missingVars []string
	seenMissing := make(map[string]bool)

	addMissing := func(vars []string) {
		for _, v := range vars {
			if !seenMissing[v] {
				seenMissing[v] = true
				missingVars = append(missingVars, v)
			}
		}
	}

	subURL, missing := env.Substitute(d.URL, session)
	addMissing(missing)

	subMethod, missing := env.Substitute(d.Method, session)
	addMissing(missing)

	subBody, missing := env.Substitute(d.Body, session)
	addMissing(missing)

	var subHeaders map[string][]string
	if d.Headers != nil {
		subHeaders = make(map[string][]string)
		for k, vals := range d.Headers {
			subK, missing := env.Substitute(k, session)
			addMissing(missing)

			var subVals []string
			for _, val := range vals {
				subV, missing := env.Substitute(val, session)
				addMissing(missing)
				subVals = append(subVals, subV)
			}
			subHeaders[subK] = subVals
		}
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing variables: %s", strings.Join(missingVars, ", "))
	}

	return &Data{
		URL:     subURL,
		Method:  subMethod,
		Headers: subHeaders,
		Proto:   d.Proto,
		Body:    subBody,
	}, nil
}
