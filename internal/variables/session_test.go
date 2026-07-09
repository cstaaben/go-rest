package variables_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/variables"
)

func TestSession_Operations(t *testing.T) {
	s := variables.NewSession()

	tests := []struct {
		name          string
		action        string
		key           string
		value         any
		expectedVal   any
		expectedFound bool
	}{
		{
			name:          "get initially empty",
			action:        "get",
			key:           "token",
			expectedVal:   nil,
			expectedFound: false,
		},
		{
			name:          "set token",
			action:        "set",
			key:           "token",
			value:         "12345",
			expectedVal:   "12345",
			expectedFound: true,
		},
		{
			name:          "overwrite token",
			action:        "set",
			key:           "token",
			value:         "67890",
			expectedVal:   "67890",
			expectedFound: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.action == "set" {
				s.Set(tc.key, tc.value)
			}
			val, ok := s.Get(tc.key)
			assert.Equal(t, tc.expectedFound, ok)
			assert.Equal(t, tc.expectedVal, val)
		})
	}
}

func TestSession_Concurrency(t *testing.T) {
	s := variables.NewSession()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			s.Set("key", val)
		}(i)
		go func() {
			defer wg.Done()
			_, _ = s.Get("key")
		}()
	}

	wg.Wait()
}
