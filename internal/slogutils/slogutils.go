// Package slogutils provides helper utilities for structured logging with log/slog.
package slogutils

import (
	"log/slog"
)

// Err returns a slog.Attr for an error with key "error".
// If targetError is nil, it returns a slog.Attr with value "<nil>".
func Err(targetError error) slog.Attr {
	if targetError == nil {
		return slog.String("error", "<nil>")
	}
	return slog.String("error", targetError.Error())
}
