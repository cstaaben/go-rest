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

package request

import (
	"reflect"
)

type Request struct {
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
	Data *Data  `json:"data,omitempty"`
}

type Data struct {
	URL      string              `json:"url,omitempty"`
	Headers  map[string][]string `json:"headers,omitempty"`
	Method   string              `json:"method,omitempty"`
	Proto    string              `json:"proto,omitempty"`
	Body     string              `json:"body,omitempty"`
	Response *Data               `json:"response,omitempty"`
}

// FilterValue is the value we use when filtering against this item when
// we're filtering the list.
func (request *Request) FilterValue() string {
	return request.Name
}

func (request *Request) Title() string {
	return request.Name
}

func (request *Request) Description() string {
	return request.Desc
}

func (r *Request) Equal(other *Request) bool {
	return reflect.DeepEqual(r, other)
}
