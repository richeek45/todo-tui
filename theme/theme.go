package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	PrimaryText        lipgloss.AdaptiveColor
	SecondaryText      lipgloss.AdaptiveColor
	FaintText          lipgloss.AdaptiveColor
	SelectedBackground lipgloss.AdaptiveColor
	PrimaryBorder      lipgloss.AdaptiveColor
	SecondaryBorder    lipgloss.AdaptiveColor
	FaintBorder        lipgloss.AdaptiveColor
}

var DefaultTheme = &Theme{
	PrimaryText:        lipgloss.AdaptiveColor{Light: "000", Dark: "015"},
	SecondaryText:      lipgloss.AdaptiveColor{Light: "244", Dark: "251"},
	FaintText:          lipgloss.AdaptiveColor{Light: "007", Dark: "245"},
	SelectedBackground: lipgloss.AdaptiveColor{Light: "006", Dark: "008"},
	PrimaryBorder:      lipgloss.AdaptiveColor{Light: "013", Dark: "008"},
	SecondaryBorder:    lipgloss.AdaptiveColor{Light: "008", Dark: "007"},
	FaintBorder:        lipgloss.AdaptiveColor{Light: "254", Dark: "000"},
}

func ParseTheme() Theme {

	return *DefaultTheme
}
