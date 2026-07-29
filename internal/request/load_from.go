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
	"github.com/cstaaben/go-rest/internal/database"
	"github.com/cstaaben/go-rest/internal/database/dto"
)

const UnsortedName = "unsorted"

func LoadFrom(_ *database.DAO) ([]*dto.Group, error) {
	panic("request.LoadFrom not updated for sqlite") // TODO: update for sqlite
	// entries, err := os.ReadDir(dataDir)
	// if err != nil {
	// 	return nil, fmt.Errorf("reading data directory: %w", err)
	// }
	//
	// requestDir, err := findRequestDir(entries)
	// if err != nil {
	// 	return nil, errors.New("requests directory not found")
	// }
	//
	// files, err := os.ReadDir(path.Join(dataDir, requestDir))
	// if err != nil {
	// 	return nil, fmt.Errorf("reading requests directory: %w", err)
	// }
	//
	// groups := make([]*Group, 0)
	// for _, file := range files {
	// 	group, err := loadRequests(path.Join(dataDir, requestDir, file.Name()))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("loading request: %w", err)
	// 	}
	//
	// 	groups = append(groups, group)
	// }
	//
	// return groups, nil
}

