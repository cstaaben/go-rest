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

package snippet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedHeaderKeys(t *testing.T) {
	headers := map[string][]string{
		"Z-Header": {"val1"},
		"A-Header": {"val2"},
	}
	res := sortedHeaderKeys(headers)
	assert.Equal(t, []string{"A-Header", "Z-Header"}, res)
}

func TestHasMultiValueHeaders(t *testing.T) {
	t.Run("no multi-value", func(t *testing.T) {
		headers := map[string][]string{
			"A-Header": {"val1"},
		}
		assert.False(t, hasMultiValueHeaders(headers))
	})

	t.Run("has multi-value", func(t *testing.T) {
		headers := map[string][]string{
			"A-Header": {"val1", "val2"},
		}
		assert.True(t, hasMultiValueHeaders(headers))
	})
}

func TestDefaultMethod(t *testing.T) {
	t.Run("empty method defaults to GET", func(t *testing.T) {
		assert.Equal(t, "GET", defaultMethod(""))
	})

	t.Run("preserves non-empty method", func(t *testing.T) {
		assert.Equal(t, "POST", defaultMethod("POST"))
	})
}
