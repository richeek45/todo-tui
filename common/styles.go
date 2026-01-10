package common

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/theme"
)

type CommonStyles struct {
	FaintTextStyle lipgloss.Style
}

var (
	TabsContentHeight = 2
	HeaderHeight      = 1
	SingleRuneWidth   = 4
	SearchHeight      = 3
)

func BuildStyles(theme theme.Theme) CommonStyles {
	var s CommonStyles

	s.FaintTextStyle = lipgloss.NewStyle().Foreground(theme.FaintText)

	return s
}
