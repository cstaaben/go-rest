package uuid

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func NewString() string {
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("error creating new v7 UUID", slog.String("error", err.Error()))
		panic(fmt.Errorf("creating v7 uuid: %w", err))
	}

	return id.String()
}
