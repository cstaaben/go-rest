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
	"github.com/charmbracelet/bubbles/list"
)

const UnsortedName = "unsorted"

type Group struct {
	Name     string     `json:"name"`
	Desc     string     `json:"desc"`
	Requests []*Request `json:"requests"`
}

func NewGroup(name string) *Group {
	return &Group{
		Name:     name,
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

// ListItems returns the collection of requests as list.Item.
func (group *Group) ListItems() []list.Item {
	var items []list.Item

	for _, r := range group.Requests {
		items = append(items, r)
	}

	return items
}

func (group *Group) AddRequest(r *Request) {
	group.Requests = append(group.Requests, r)
}

func (group *Group) RemoveRequest(r *Request) {
	for i, req := range group.Requests {
		if req == r {
			group.Requests = append(group.Requests[:i], group.Requests[i+1:]...)
			break
		}
	}
}
