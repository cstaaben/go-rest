package styles

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/cstaaben/go-rest/internal/config"
)

var (
	Base          = lipgloss.NewStyle().Padding(1, 2).Foreground(Colors().Foreground)
	BorderPanel   = lipgloss.NewStyle().Inherit(Base).Border(lipgloss.RoundedBorder(), true)
	FocusedBorder = lipgloss.NewStyle().Inherit(BorderPanel).BorderForeground(Colors().FocusHighlight)
	Title         = lipgloss.NewStyle().Padding(0, 1).Bold(true).Align(lipgloss.Center).Inherit(Base)

	colors        *ColorScheme
	defaultColors = ColorScheme{
		FocusHighlight: lipgloss.AdaptiveColor{
			Light: "#5fffff",
			Dark:  "#5fffff",
		},
		Foreground: lipgloss.AdaptiveColor{
			Light: "#000000",
			Dark:  "#ffffff",
		},
	}
)

type ColorScheme struct {
	FocusHighlight lipgloss.AdaptiveColor
	Foreground     lipgloss.AdaptiveColor
}

func Colors() *ColorScheme {
	if colors != nil {
		return colors
	}

	switch config.ColorScheme() {
	case "default":
		colors = &defaultColors
	default:
		colors = &defaultColors
	}

	return colors
}
