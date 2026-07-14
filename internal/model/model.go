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

// Package model defines utilities for the current state of the UI.
package model

import (
	"sync"

	"github.com/cstaaben/go-rest/internal/request"
)

var (
	mu            sync.Mutex
	activeRequest *request.Request
	groups        []*request.Group
)

func ActiveRequest() *request.Request {
	mu.Lock()
	defer mu.Unlock()

	return activeRequest
}

func Groups() []*request.Group {
	mu.Lock()
	defer mu.Unlock()

	return groups
}

func AddGroup(g *request.Group) {
	mu.Lock()
	defer mu.Unlock()

	groups = append(groups, g)
}

func AddRequest(gid string, r *request.Request) {
	mu.Lock()
	defer mu.Unlock()

	for _, g := range groups {
		if gid != g.ID {
			continue
		}

		g.Requests = append(g.Requests, r)
	}
}

func SetActiveRequest(r *request.Request) {
	mu.Lock()
	defer mu.Unlock()

	activeRequest = r
}

func GroupChildren(gid string) []string {
	mu.Lock()
	defer mu.Unlock()

	var group *request.Group
	for _, g := range groups {
		if g.ID == gid {
			group = g
		}
	}

	var childUIDs []string
	for _, r := range group.Requests {
		childUIDs = append(childUIDs, r.ID)
	}

	return childUIDs
}

func IsGroup(uid string) bool {
	mu.Lock()
	defer mu.Unlock()

	for _, g := range groups {
		if uid == g.ID {
			return true
		}
	}

	return false
}
