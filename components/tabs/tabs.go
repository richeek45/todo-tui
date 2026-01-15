package tabs

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/common"
	"github.com/richeek45/todo-tui/components/carousel"
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
)

type SectionTab struct {
	section section.Section
	spinner spinner.Model
}

type Model struct {
	carousel    carousel.Model
	ctx         *context.ProgramContext
	sectionTabs []SectionTab
}

func NewModel(ctx *context.ProgramContext) Model {
	c := carousel.New(carousel.WithHeight(1), carousel.WithOverflowIndicators("←", "→"), carousel.WithSeparators())
	m := Model{
		carousel: c,
	}

	m.UpdateProgramContext(ctx)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case spinner.TickMsg:
		for i, tab := range m.sectionTabs {
			// if section is loading
			var spinnerCmd tea.Cmd
			m.sectionTabs[i].spinner, spinnerCmd = tab.spinner.Update(msg)
			cmds = append(cmds, spinnerCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	c := m.carousel.View()
	logo := m.ViewLogo()
	return m.ctx.Styles.Tabs.TabsRow.
		Width(m.ctx.ScreenWidth).
		Height(common.HeaderHeight).
		Render(lipgloss.JoinHorizontal(lipgloss.Bottom, lipgloss.NewStyle().
			Width(m.ctx.ScreenWidth-lipgloss.Width(logo)).
			Render(c), logo))

}

func (m *Model) CurrentSectionId() int {
	return m.carousel.Cursor()
}

func (m *Model) SetCurrentSectionId(id int) {
	m.carousel.SetCursor(id)
}

func (m *Model) SetSections(sections []section.Section) {
	sectionTabs := make([]SectionTab, 0)

	for _, s := range sections {
		tab := SectionTab{section: s, spinner: spinner.New(
			spinner.WithSpinner(spinner.Ellipsis),
			spinner.WithStyle(
				lipgloss.NewStyle().Foreground(m.ctx.Theme.FaintText).Padding(2),
			),
		)}

		sectionTabs = append(sectionTabs, tab)
	}

	m.sectionTabs = sectionTabs
	m.UpdateTabTitles()
}

func (m *Model) UpdateTabTitles() {
	titles := make([]string, 0)

	for _, tab := range m.sectionTabs {
		cfg := tab.section.GetConfig()
		title := cfg.Title

		titles = append(titles, title)
	}

	log.Print("titles = ", titles)

	oldCursor := m.carousel.Cursor()
	m.carousel.SetItems(titles)
	m.carousel.SetCursor(oldCursor)
}

func (m *Model) SetAllLoading() []tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for i := range m.sectionTabs {
		cmds = append(cmds, m.sectionTabs[i].spinner.Tick)
	}

	return cmds
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.carousel.SetStyles(carousel.Styles{
		Item:              ctx.Styles.Tabs.Tab,
		Selected:          ctx.Styles.Tabs.ActiveTab,
		OverflowIndicator: ctx.Styles.Tabs.OverflowIndicator,
		Separator:         ctx.Styles.Tabs.TabsSeparator,
	})
	m.carousel.SetWidth(ctx.ScreenWidth - lipgloss.Width(m.ViewLogo()))
}

func (m *Model) ViewLogo() string {

	return lipgloss.NewStyle().
		Padding(0, 1, 0, 2).
		Height(2).
		Render(lipgloss.NewStyle().Foreground(context.LogoColor).Render(constants.Logo))
}
