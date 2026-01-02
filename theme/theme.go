package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	PrimaryBorder      lipgloss.AdaptiveColor
	SelectedBackground lipgloss.AdaptiveColor
}

var DefaultTheme = &Theme{
	PrimaryBorder:      lipgloss.AdaptiveColor{Light: "013", Dark: "008"},
	SelectedBackground: lipgloss.AdaptiveColor{Light: "006", Dark: "008"},
}

func ParseTheme() Theme {

	return *DefaultTheme
}
