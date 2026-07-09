package variables

import "sync"

// Session represents a temporary, thread-safe in-memory key-value store for session variables.
type Session struct {
	mu        sync.RWMutex
	variables map[string]any
}

// NewSession initializes and returns a new Session.
func NewSession() *Session {
	return &Session{
		variables: make(map[string]any),
	}
}

// Get retrieves a session variable value by its key.
func (s *Session) Get(key string) (any, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.variables == nil {
		return nil, false
	}
	val, ok := s.variables[key]
	return val, ok
}

// Set stores or updates a session variable value.
func (s *Session) Set(key string, val any) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.variables == nil {
		s.variables = make(map[string]any)
	}
	s.variables[key] = val
}
