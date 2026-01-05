package sidebar

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/keys"
)

type Model struct {
	viewport viewport.Model
	ctx      *context.ProgramContext

	IsOpen     bool
	data       string
	emptyState string
}

func NewModel() Model {
	return Model{
		IsOpen: false,
		data:   "",
		viewport: viewport.Model{
			Height: 0,
			Width:  0,
		},
		ctx:        nil,
		emptyState: "Nothing selected...",
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.PageUp):
			m.viewport.HalfPageUp()
		case key.Matches(msg, keys.Keys.PageDown):
			m.viewport.HalfPageDown()
		}
	}
	return m, nil
}

func (m Model) View() string {
	if !m.IsOpen {
		return ""
	}

	height := m.ctx.MainContentHeight
	style := m.ctx.Styles.Sidebar.Root.
		Height(height).
		Width(m.ctx.Config.Defaults.Preview.Width).
		MaxWidth(m.ctx.Config.Defaults.Preview.Width)

	if m.data == "" {
		return style.Align(lipgloss.Center).
			Render(lipgloss.PlaceVertical(height, lipgloss.Center, m.emptyState))
	}

	return style.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		m.ctx.Styles.Sidebar.PagerStyle.
			Render(fmt.Sprintf("%d%%", int(m.viewport.ScrollPercent()*100)))))

}

func (m Model) GetData() string {
	return m.data
}

func (m *Model) SetContent(data string) {
	m.data = data
	m.viewport.SetContent(data)
}

func (m *Model) GetSidebarContentWidth() int {
	if m.ctx.Config == nil {
		return 0
	}

	return m.ctx.Config.Defaults.Preview.Width
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	if ctx == nil {
		return
	}

	m.ctx = ctx
	m.viewport.Height = m.ctx.MainContentHeight - m.ctx.Styles.Sidebar.PagerHeight
	m.viewport.Width = m.GetSidebarContentWidth()
}
