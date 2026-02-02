package context

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/common"
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
	Tabs struct {
		Tab               lipgloss.Style
		ActiveTab         lipgloss.Style
		OverflowIndicator lipgloss.Style
		TabsSeparator     lipgloss.Style
		TabsRow           lipgloss.Style
	}
	Table struct {
		CellStyle                lipgloss.Style
		SelectedCellStyle        lipgloss.Style
		TitleCellStyle           lipgloss.Style
		SingleRuneTitleCellStyle lipgloss.Style
		HeaderStyle              lipgloss.Style
		RowStyle                 lipgloss.Style
	}
	Section struct {
		ContainerPadding int
		ContainerStyle   lipgloss.Style
		SpinnerStyle     lipgloss.Style
		EmptyStateStyle  lipgloss.Style
		KeyStyle         lipgloss.Style
	}
	FooterStyle lipgloss.Style
	Common      common.CommonStyles
}

var (
	FooterHeight       = 1
	ExpandedHelpHeight = 10
	MainContentPadding = 1
	TabsHeight         = 3
	TableHeaderHeight  = 2
)

const (
	LogoColor = lipgloss.Color("#09e750ff")
)

func InitStyles(theme theme.Theme) Styles {
	var s Styles

	s.Common = common.BuildStyles(theme)

	s.Tabs.Tab = lipgloss.NewStyle().Faint(true).Padding(0, 2)
	s.Tabs.ActiveTab = s.Tabs.Tab.
		Faint(false).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("50")).
		Foreground(theme.PrimaryText)

	s.Tabs.OverflowIndicator = s.Common.FaintTextStyle.Bold(true).Padding(0, 1)
	s.Tabs.TabsSeparator = lipgloss.NewStyle().Foreground(theme.SecondaryBorder)
	s.Tabs.TabsRow = lipgloss.NewStyle().
		Height(common.TabsContentHeight).
		BorderBottom(true).
		BorderBottomForeground(theme.PrimaryBorder)

	s.Table.CellStyle = lipgloss.NewStyle().PaddingLeft(1).
		PaddingRight(1).
		MaxHeight(1)
	s.Table.SelectedCellStyle = s.Table.CellStyle.
		Background(theme.SelectedBackground)
	s.Table.TitleCellStyle = s.Table.CellStyle.
		Bold(true).
		Foreground(theme.PrimaryText)
	s.Table.SingleRuneTitleCellStyle = s.Table.TitleCellStyle.
		Width(common.SingleRuneWidth)
	s.Table.HeaderStyle = lipgloss.NewStyle()
	s.Table.RowStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.FaintBorder)

	s.Section.ContainerPadding = 1
	s.Section.ContainerStyle = lipgloss.NewStyle().
		Padding(0, s.Section.ContainerPadding)
	s.Section.SpinnerStyle = lipgloss.NewStyle().Padding(0, 1)
	s.Section.EmptyStateStyle = lipgloss.NewStyle().
		Faint(true).
		PaddingLeft(1).
		MarginBottom(1)
	s.Section.KeyStyle = lipgloss.NewStyle().
		Foreground(theme.PrimaryText).
		Background(theme.SelectedBackground).
		Padding(0, 1)

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
			Right:       "|",
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
