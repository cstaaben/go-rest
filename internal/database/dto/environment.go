/*
 * go-rest - A TUI for a REST client
 * Copyright (C) 2024  Corbin Staaben
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

package dto

import (
	"bytes"
	"reflect"
)

// Environment represents an execution environment containing environment variables.
type Environment struct {
	UUID      string                 `json:"uuid,omitempty" db:"uuid"`
	Name      string                 `json:"name,omitempty" db:"name"`
	Variables []*EnvironmentVariable `json:"variables,omitempty" db:"variables"`
}

// NewEnvironment creates a new Environment instance with initialized variables slice.
func NewEnvironment(name string) *Environment {
	return &Environment{
		Name:      name,
		Variables: make([]*EnvironmentVariable, 0),
	}
}

// FilterValue returns the environment name for list filtering.
func (e *Environment) FilterValue() string {
	return e.Name
}

// Title returns the environment name for list display.
func (e *Environment) Title() string {
	return e.Name
}

// Description returns the environment name description.
func (e *Environment) Description() string {
	return e.Name
}

// Equal compares structural equality between two Environment instances.
func (e *Environment) Equal(other *Environment) bool {
	return reflect.DeepEqual(e, other)
}

// EnvironmentVariable represents a key-value variable associated with an environment.
type EnvironmentVariable struct {
	EnvironmentUUID string `json:"environment_uuid,omitempty" db:"environment_uuid"`
	Name            string `json:"name,omitempty" db:"name"`
	Value           []byte `json:"value,omitempty" db:"value"`
}

// Equal compares equality between two EnvironmentVariable instances.
func (ev *EnvironmentVariable) Equal(other *EnvironmentVariable) bool {
	if ev == nil || other == nil {
		return ev == other
	}
	return ev.EnvironmentUUID == other.EnvironmentUUID &&
		ev.Name == other.Name &&
		bytes.Equal(ev.Value, other.Value)
}
