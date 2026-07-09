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

// Package request defines the a request wrapper around an HTTP request.
package request

import (
	"reflect"

	"github.com/cstaaben/go-rest/internal/uuid"
)

// Request contains all the information to create an HTTP request.
type Request struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Data     *Data  `json:"data,omitempty"`
	Response *Data  `json:"response,omitempty"`
}

type Data struct {
	URL     string              `json:"url,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
	Method  string              `json:"method,omitempty"`
	Proto   string              `json:"proto,omitempty"`
	Body    string              `json:"body,omitempty"`
}

// New returns a new request with a new ID.
func New() *Request {
	return &Request{
		ID: uuid.Must(),
	}
}

// FilterValue is the value we use when filtering against this item when
// we're filtering the list.
func (r *Request) FilterValue() string {
	return r.Name
}

// Title returns the request's name.
func (r *Request) Title() string {
	return r.Name
}

// Description returns the requests description.
func (r *Request) Description() string {
	return r.Desc
}

// Equals returns true if the Request `other` is deeply equal to the current request and false otherwise.
func (r *Request) Equals(other *Request) bool {
	return reflect.DeepEqual(r, other)
}
