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
	"log/slog"
	"slices"

	"github.com/cstaaben/go-rest/internal/uuid"
)

const UnsortedGroup = "unsorted"

type Group struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Desc     string     `json:"desc"`
	Requests []*Request `json:"requests"`
}

func NewGroup() *Group {
	return &Group{
		ID:       uuid.Must(),
		Requests: make([]*Request, 0),
	}
}

// FilterValue is the value we use when filtering against this item when
// we're filtering the list.
func (group *Group) FilterValue() string {
	return group.Name
}

// Title returns the name of the request group.
func (group *Group) Title() string {
	return group.Name
}

// Description returns a brief description of the group.
func (group *Group) Description() string {
	return group.Desc
}

func (group *Group) AddRequest(r *Request) {
	group.Requests = append(group.Requests, r)
}

func (group *Group) RemoveRequest(r *Request) {
	idx := -1
	for i, req := range group.Requests {
		if !req.Equals(r) {
			continue
		}

		idx = i
		break
	}
	if idx == -1 {
		slog.Error("Failed to find request for deletion", slog.String("name", group.Name))
		return
	}

	group.Requests = slices.Delete(group.Requests, idx, idx+1)
}
