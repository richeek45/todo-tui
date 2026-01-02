package context

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/theme"
)

type Styles struct {
	Sidebar struct {
		Root        lipgloss.Style
		PagerStyle  lipgloss.Style
		PagerHeight int
	}
}

func InitStyles(theme theme.Theme) Styles {
	var s Styles

	s.Sidebar.Root.
		BorderLeft(true).
		BorderStyle(lipgloss.Border{
			Top:         "",
			Bottom:      "",
			Left:        "|",
			Right:       "",
			TopLeft:     "",
			TopRight:    "",
			BottomLeft:  "",
			BottomRight: "",
		}).BorderForeground(theme.PrimaryBorder)

	return s
}
