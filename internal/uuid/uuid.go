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

// Package uuid provides a simple utility for creating v7 UUIDs.
package uuid

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Must returns a new v7 UUID as a string. If there is any error creating the UUID, it panics.
func Must() string {
	id, err := uuid.NewV7()
	if err != nil {
		slog.Error("error creating new v7 UUID", slog.String("error", err.Error()))
		panic(fmt.Errorf("creating v7 uuid: %w", err))
	}

	return id.String()
}
