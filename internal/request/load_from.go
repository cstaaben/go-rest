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
 * along with this program.  If not, see <https: //www.gnu.org/licenses/>.
 */

package request

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"sigs.k8s.io/yaml"
)

func LoadFrom(dataDir string) ([]*Group, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("reading data directory: %w", err)
	}

	requestDir, err := findRequestDir(entries)
	if err != nil {
		return nil, errors.New("requests directory not found")
	}

	files, err := os.ReadDir(path.Join(dataDir, requestDir))
	if err != nil {
		return nil, fmt.Errorf("reading requests directory: %w", err)
	}

	groups := make([]*Group, 0)
	for _, file := range files {
		group, err := loadRequests(path.Join(dataDir, requestDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("loading request: %w", err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func findRequestDir(entries []fs.DirEntry) (string, error) {
	for _, entry := range entries {
		if entry.IsDir() && strings.EqualFold(entry.Name(), "requests") {
			return entry.Name(), nil
		}
	}

	return "", os.ErrNotExist
}

func loadRequests(filepath string) (*Group, error) {
	body, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	g := new(Group)
	if err = yaml.Unmarshal(body, g); err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	return g, nil
}
