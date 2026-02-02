package sidebar

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/common"
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

	height := m.ctx.MainContentHeight - context.FooterHeight - context.TabsHeight - common.SearchHeight
	width := m.ctx.ScreenWidth - m.ctx.MainContentWidth - 10

	style := m.ctx.Styles.Sidebar.Root.
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("37")).
		Margin(1, 2).
		Width(width)

	if m.data == "" {
		return style.Align(lipgloss.Center).
			Render(lipgloss.PlaceVertical(height, lipgloss.Center, m.emptyState))
	}

	return style.Render(lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		m.ctx.Styles.Sidebar.PagerStyle.
			Render(fmt.Sprintf("%d%%", int(m.viewport.ScrollPercent()*100))),
	))

}

func (m *Model) ScrollToTop() {
	m.viewport.GotoTop()
}

func (m *Model) ScrollToBottom() {
	m.viewport.GotoBottom()
}

func (m Model) GetData() string {
	return m.data
}

func (m *Model) SetContent(data string) {
	m.data = data
	m.viewport.SetContent("\n" + data)
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
	height := m.ctx.MainContentHeight - context.FooterHeight - context.TabsHeight - common.SearchHeight
	m.viewport.Height = height
	m.viewport.Width = m.GetSidebarContentWidth() - 10
}
