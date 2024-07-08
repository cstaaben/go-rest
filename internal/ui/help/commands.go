package help

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cstaaben/go-rest/internal/model/keymap"
)

type UpdateKeymapMsg keymap.KeyMap

func UpdateKeymap(keys keymap.KeyMap) tea.Cmd {
	return func() tea.Msg {
		return UpdateKeymapMsg(keys)
	}
}
