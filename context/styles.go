package context

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/theme"
)

type Styles struct {
	Sidebar struct {
		Root        lipgloss.Style
		PagerStyle  lipgloss.Style
		PagerHeight int
	}
	Help struct {
		Text         lipgloss.Style
		KeyText      lipgloss.Style
		BubbleStyles help.Styles
	}
	FooterStyle lipgloss.Style
}

var (
	FooterHeight       = 1
	ExpandedHelpHeight = 10
)

func InitStyles(theme theme.Theme) Styles {
	var s Styles

	s.Sidebar.PagerHeight = 1
	s.Sidebar.PagerStyle = lipgloss.NewStyle().
		Height(s.Sidebar.PagerHeight).
		Bold(true).
		Foreground(theme.FaintText)

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

	s.Help.Text = lipgloss.NewStyle().Foreground(theme.PrimaryText)
	s.Help.KeyText = lipgloss.NewStyle().Foreground(theme.SecondaryText)

	s.FooterStyle = lipgloss.NewStyle().
		Background(theme.SelectedBackground).
		Height(FooterHeight)

	return s
}
