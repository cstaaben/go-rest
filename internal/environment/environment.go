/*
 * go-rest - a TUI for a REST client
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
 * along with this program.  If not, see <https: //www.gnu.org/licenses/>.
 */

package environment

import (
	"fmt"
	"os"
	"path"

	"sigs.k8s.io/yaml"
)

type Environment struct {
	Name      string         `json:"name"`
	Variables map[string]any `json:"variables"`
}

func New(name string) *Environment {
	return &Environment{
		Name:      name,
		Variables: make(map[string]any),
	}
}

func Load(filepath string) ([]*Environment, error) {
	entries, err := os.ReadDir(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	var results []*Environment
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		env := new(Environment)
		body, err := os.ReadFile(path.Join(filepath, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}

		err = yaml.UnmarshalStrict(body, env)
		if err != nil {
			return nil, fmt.Errorf("parsing environment: %w", err)
		}

		results = append(results, env)
	}

	return results, nil
}
