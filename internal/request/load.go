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
 * along with this program.  If not, see <https: //www.gnu.org/licenses/>.
 */

package request

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"sigs.k8s.io/yaml"
)

var reqCache = cache.New(time.Minute*5, time.Second*30)

func LoadAll(dataDir string) ([]*Group, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("reading data directory: %w", err)
	}

	requestDir, err := findRequestDir(entries)
	if err != nil {
		return nil, errors.New("requests directory not found")
	}

	requestsPath := filepath.Join(dataDir, requestDir)
	groups := make([]*Group, 0)

	err = filepath.WalkDir(requestsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) != ".yaml" {
			return nil
		}
		group, err := loadRequests(path)
		if err != nil {
			return fmt.Errorf("loading request: %w", err)
		}
		groups = append(groups, group)
		return nil
	})
	if err != nil {
		return nil, err
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

	for _, r := range g.Requests {
		reqCache.Set(r.ID, r, cache.DefaultExpiration)
	}

	return g, nil
}

func LoadGroupFromDir(gid string, dataDir string) (*Group, error) {
	targetFile, err := groupFilePath(gid, dataDir)
	if err != nil {
		return nil, fmt.Errorf("finding requests: %w", err)
	}

	group, err := loadRequests(targetFile)
	if err != nil {
		return nil, fmt.Errorf("loading requests: %w", err)
	}

	return group, nil
}

func LoadGroupIDsFromDir(dataDir string) ([]string, error) {
	var uids []string
	d := filepath.Join(dataDir, "requests")
	err := filepath.WalkDir(
		d,
		func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			if filepath.Ext(entry.Name()) != ".yaml" {
				return nil
			}

			id := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			uids = append(uids, id)

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return uids, nil
}

func groupFilePath(id string, dataDir string) (string, error) {
	dirs, err := os.ReadDir(dataDir)
	if err != nil {
		return "", fmt.Errorf("data directory: %w", err)
	}

	requestsDir, err := findRequestDir(dirs)
	if err != nil {
		return "", errors.New("requests directory not found")
	}

	requestsPath := filepath.Join(dataDir, requestsDir)
	var targetFile string

	err = filepath.WalkDir(requestsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == fmt.Sprintf("%s.yaml", id) {
			targetFile = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("requests directory: %w", err)
	} else if targetFile == "" {
		return "", errors.New("request group not found")
	}

	return targetFile, nil
}

func GroupExists(id string, dataDir string) (bool, error) {
	fp, err := groupFilePath(id, dataDir)
	return fp != "", err
}
