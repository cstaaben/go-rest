package slogutils_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cstaaben/go-rest/internal/slogutils"
)

func TestErr(t *testing.T) {
	t.Run("non-nil error", func(t *testing.T) {
		err := errors.New("database connection failed")
		attr := slogutils.Err(err)

		assert.Equal(t, "error", attr.Key)
		assert.Equal(t, slog.KindString, attr.Value.Kind())
		assert.Equal(t, "database connection failed", attr.Value.String())
	})

	t.Run("nil error", func(t *testing.T) {
		attr := slogutils.Err(nil)

		assert.Equal(t, "error", attr.Key)
		assert.Equal(t, slog.KindString, attr.Value.Kind())
		assert.Equal(t, "<nil>", attr.Value.String())
	})
}
