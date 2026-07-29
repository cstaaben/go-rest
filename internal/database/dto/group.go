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

package dto

import (
	"charm.land/bubbles/v2/list"
)

// UnsortedName is the identifier reserved for requests that do not belong to a custom group.
const UnsortedName = "unsorted"

// Group represents a logical grouping of related HTTP requests.
type Group struct {
	UUID     string     `json:"uuid,omitempty" db:"uuid"`
	Name     string     `json:"name,omitempty" db:"name"`
	Desc     string     `json:"desc,omitempty" db:"desc"`
	Requests []*Request `json:"requests,omitempty" db:"requests"`
}

// NewGroup creates and initializes a new Group instance with the specified name.
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

	for _, req := range group.Requests {
		items = append(items, req)
	}

	return items
}

// AddRequest appends a request to the group's collection.
func (group *Group) AddRequest(req *Request) {
	group.Requests = append(group.Requests, req)
}

// RemoveRequest removes a specific request from the group's collection.
func (group *Group) RemoveRequest(targetReq *Request) {
	for i, req := range group.Requests {
		if req == targetReq {
			group.Requests = append(group.Requests[:i], group.Requests[i+1:]...)
			break
		}
	}
}
