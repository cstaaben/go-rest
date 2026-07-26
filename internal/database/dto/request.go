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

// Package dto defines the data types used by go-rest.
package dto

import (
	"database/sql"
	"reflect"
)

type Request struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name,omitempty"`
	Desc      string `json:"desc,omitempty"`
	GroupUUID string `json:"group_uuid,omitempty"`
	Data      *Data  `json:"data,omitempty"`
}

type Data struct {
	URL      string              `json:"url,omitempty"`
	Headers  map[string][]string `json:"headers,omitempty"`
	Method   string              `json:"method,omitempty"`
	Proto    string              `json:"proto,omitempty"`
	Body     string              `json:"body,omitempty"`
	Response *Data               `json:"response,omitempty"`
}

// FilterValue is the value used when filtering against this item when
// filtering the list.
func (r *Request) FilterValue() string {
	return r.Name
}

func (r *Request) Title() string {
	return r.Name
}

func (r *Request) Description() string {
	return r.Desc
}

func (r *Request) Equal(other *Request) bool {
	return reflect.DeepEqual(r, other)
}

type RequestListData struct {
	RequestID     string         `db:"request_id"`
	RequestName   string         `db:"request_name"`
	RequestDesc   sql.NullString `db:"request_desc"`
	RequestMethod string         `db:"request_method"`
	GroupID       string         `db:"group_id"`
	GroupName     string         `db:"group_name"`
	GroupDesc     sql.NullString `db:"group_desc"`
}

func (r *RequestListData) RequestDescription() string {
	return r.RequestDesc.String
}

func (r *RequestListData) GroupDescription() string {
	return r.GroupDesc.String
}
