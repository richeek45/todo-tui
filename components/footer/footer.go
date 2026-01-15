package footer

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/keys"
)

type Model struct {
	ctx     *context.ProgramContext
	help    help.Model
	ShowAll bool
}

func NewModel(ctx *context.ProgramContext) Model {
	help := help.New()
	help.ShowAll = true
	help.Styles = ctx.Styles.Help.BubbleStyles

	return Model{
		ctx:  ctx,
		help: help,
	}
}

func (m *Model) View() string {
	var footer string

	helpIndicator := lipgloss.NewStyle().Foreground(m.ctx.Theme.SelectedBackground).Padding(0, 1).Render("? help")

	footer = m.ctx.Styles.FooterStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, helpIndicator))

	if m.ShowAll {
		fullHelp := m.help.View(keys.Keys)
		return lipgloss.JoinVertical(lipgloss.Top, footer, fullHelp)
	}

	return footer
}

func (m *Model) SetWidth(width int) {
	m.help.Width = width
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.help.Styles = ctx.Styles.Help.BubbleStyles
}
