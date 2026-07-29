package help

import (
	tea "charm.land/bubbletea/v2"

	"github.com/cstaaben/go-rest/internal/model/keymap"
)

type UpdateKeymapMsg keymap.KeyMap

func UpdateKeymap(keys keymap.KeyMap) tea.Cmd {
	return func() tea.Msg {
		return UpdateKeymapMsg(keys)
	}
}
