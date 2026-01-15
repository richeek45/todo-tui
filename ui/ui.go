package ui

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
	"github.com/richeek45/todo-tui/components/footer"
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/components/sidebar"
	"github.com/richeek45/todo-tui/components/tabs"
	"github.com/richeek45/todo-tui/components/tasksection"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/keys"
	"github.com/richeek45/todo-tui/theme"
)

const (
// bullet   = "•"
)

type Item struct {
	title, desc string
}

type Model struct {
	sidebar       sidebar.Model
	tabs          tabs.Model
	footer        footer.Model
	taskSpinner   spinner.Model
	currSectionId int
	keys          *keys.KeyMap
	priorityTasks []section.Section
	statusTasks   []section.Section
	ctx           *context.ProgramContext
	Tasks         map[string]context.Task
}

type initMsg struct {
	Config config.Config
}

func NewModel(items []Item) Model {
	m := Model{
		keys:        keys.Keys,
		sidebar:     sidebar.NewModel(),
		taskSpinner: spinner.Model{Spinner: spinner.Ellipsis},
	}

	m.ctx = &context.ProgramContext{}

	m.tabs = tabs.NewModel(m.ctx)
	m.footer = footer.NewModel(m.ctx)

	return m
}

func (m Model) InitScreen() tea.Msg {
	cfg := config.ParseConfig()
	return initMsg{Config: cfg}
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.InitScreen, tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd        tea.Cmd
		listCmd    tea.Cmd
		sidebarCmd tea.Cmd
		// currSection = m.GetCurrSection()
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.PreviousSection):
			log.Print("PreviousSection")

			prevSection := m.getSectionAt(m.getPreviousSectionId())
			if prevSection != nil {
				m.setCurrSectionId(prevSection.GetId())
				// cmd = m.onViewedRowChanged()
			}
		case key.Matches(msg, m.keys.NextSection):
			log.Print("NextSection")
			nextSection := m.getSectionAt(m.getNextSectionId())
			if nextSection != nil {
				m.setCurrSectionId(nextSection.GetId())
				// cmd = m.onViewedRowChanged()
			}

		// case key.Matches(msg, m.keys.Enter):

		case key.Matches(msg, m.keys.TogglePreview):
			// _, ok := m.list.SelectedItem().(item)
			// if ok {
			// 	m.sidebar.IsOpen = !m.sidebar.IsOpen
			// 	m.SyncMainContentWidth()
			// }
		case key.Matches(msg, m.keys.Help):
			_, v := docStyle.GetFrameSize()
			if !m.footer.ShowAll {
				m.ctx.MainContentHeight = m.ctx.ScreenHeight - v - context.ExpandedHelpHeight
			} else {
				m.ctx.MainContentHeight = m.ctx.ScreenHeight - v - context.FooterHeight
			}

			m.footer.ShowAll = !m.footer.ShowAll
		}
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)
	case initMsg:
		m.ctx.Config = &msg.Config
		m.ctx.Theme = theme.ParseTheme()
		m.ctx.Styles = context.InitStyles(m.ctx.Theme)
		m.ctx.View = m.ctx.Config.Defaults.View

		m.sidebar.IsOpen = m.ctx.Config.Defaults.Preview.Open
		m.SyncMainContentWidth()

		newSections, fetchSectionsCmds := m.fetchAllViewSections()
		m.setCurrentViewSections(newSections)
		m.setCurrSectionId(1)

		cmds = append(cmds, fetchSectionsCmds)
	}
	m.SyncProgramContext(m.ctx)
	m.SyncSideBar()

	m.sidebar, sidebarCmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd, listCmd, sidebarCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.ctx.Config == nil {
		return lipgloss.Place(m.ctx.ScreenWidth, m.ctx.ScreenHeight, lipgloss.Center, lipgloss.Center, "Reading config...")
	}

	s := strings.Builder{}

	s.WriteString(m.tabs.View())

	// currSection := m.GetCurrSection()

	// if currSection != nil {
	// 	s.WriteString(lipgloss.JoinHorizontal(
	// 		lipgloss.Top,
	// 		m.GetCurrSection().View(),
	// 		m.sidebar.View(),
	// 	))
	// }
	// s.WriteString("\n")
	s.WriteString(m.footer.View())

	return s.String()
}

func Run() {
	// notify()

	// db, err := sql.Open("sqlite3", "./test.db")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// sqlStatement := `
	// 	CREATE TABLE IF NOT EXISTS users (
	// 	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	name TEXT
	// 	)
	// `

	// _, err = db.Exec(sqlStatement)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Table 'users' created successfully")

	// _, err = db.Exec("INSERT INTO users(name) VALUES(?)", "John Doe")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err := db.Query("SELECT id, name FROM users")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// if err = rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	items := []Item{
		{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		{title: "Nutella", desc: "It's good on toast"},
		{title: "Bitter melon", desc: "It cools you down"},
		{title: "Nice socks", desc: "And by that I mean socks without holes"},
		{title: "Eight hours of sleep", desc: "I had this once"},
		{title: "Cats", desc: "Usually"},
		{title: "Plantasia, the album", desc: "My plants love it too"},
		{title: "Pour over coffee", desc: "It takes forever to make though"},
		{title: "VR", desc: "Virtual reality...what is there to say?"},
		{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		{title: "Linux", desc: "Pretty much the best OS"},
		{title: "Business school", desc: "Just kidding"},
		{title: "Pottery", desc: "Wet clay is a great feeling Wet clay is a great feelingWet clay is a great feelingWet \nclay is a great feelingWet clay is a great feelingWet clay is a great feeling"},
		{title: "Shampoo", desc: "Nothing like clean hair"},
		{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		{title: "Stickers", desc: "The thicker the vinyl the better"},
		{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		{title: "Warm light", desc: "Like around 2700 Kelvin"},
		{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := NewModel(items)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func (m *Model) SyncProgramContext(ctx *context.ProgramContext) {

	for _, section := range m.getCurrentViewSections() {
		section.UpdateProgramContext(m.ctx)
	}
	m.sidebar.UpdateProgramContext(m.ctx)
	m.tabs.UpdateProgramContext(m.ctx)
	m.footer.UpdateProgramContext(m.ctx)
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	log.Println("window size changed", "width", msg.Width, "height", msg.Height)
	_, v := docStyle.GetFrameSize()
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	log.Println("ScreenWidth=", m.ctx.ScreenWidth, " ScreenHeight=", m.ctx.ScreenHeight)
	if m.footer.ShowAll {
		m.ctx.MainContentHeight = msg.Height - v - context.ExpandedHelpHeight
	} else {
		m.ctx.MainContentHeight = msg.Height - v - context.FooterHeight
	}
	m.ctx.MainContentWidth = msg.Width
	m.footer.SetWidth(msg.Width)
}

func (m *Model) SyncSideBar() {
	// width := m.GetSidebarContentWidth()
	m.sidebar.SetContent(fmt.Sprintf("Whatever content I want to put here.... %d", 1000))
}

func (m *Model) SyncMainContentWidth() {
	sideBarOffset := 0
	if m.sidebar.IsOpen {
		sideBarOffset = m.ctx.Config.Defaults.Preview.Width
	}
	m.ctx.ScreenWidth = m.ctx.ScreenWidth - sideBarOffset
	log.Println("m.ctx.ScreenWidth=", m.ctx.ScreenWidth, " ScreenHeight=", m.ctx.ScreenHeight)
}

func (m *Model) fetchAllViewSections() ([]section.Section, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.tabs.SetAllLoading()...)

	log.Print("CTX VIEW =", m.ctx.View)

	switch m.ctx.View {
	case config.Priority:
		section, taskCmds1 := tasksection.FetchAllSections(m.ctx, m.priorityTasks)
		cmds = append(cmds, taskCmds1)
		return section, tea.Batch(cmds...)
	case config.Status:
		section, taskCmds2 := tasksection.FetchAllSections(m.ctx, m.statusTasks)
		cmds = append(cmds, taskCmds2)
		return section, tea.Batch(cmds...)
	default:
		section, taskCmds2 := tasksection.FetchAllSections(m.ctx, m.statusTasks)
		cmds = append(cmds, taskCmds2)
		return section, tea.Batch(cmds...)
	}
}

func (m *Model) setCurrentViewSections(newSections []section.Section) {
	if newSections == nil {
		return
	}

	missingSearchSection := true
	s := make([]section.Section, 0)

	if m.ctx.View == config.Status {
		if missingSearchSection {
			search := tasksection.NewModel(
				0,
				m.ctx,
				config.SectionConfig{
					Title:   "",
					Filters: "archived:false",
				},
			)
			s = append(s, &search)
		}

		m.statusTasks = append(s, newSections...)
		newSections = m.statusTasks
	}

	if m.ctx.View == config.Priority {
		if missingSearchSection {
			search := tasksection.NewModel(
				0,
				m.ctx,
				config.SectionConfig{
					Title:   "",
					Filters: "archived:false",
				},
			)
			s = append(s, &search)
		}

		m.priorityTasks = append(s, newSections...)
		newSections = m.priorityTasks
	}

	m.tabs.SetSections(newSections)
}

func (m *Model) getCurrentViewSections() []section.Section {
	switch m.ctx.View {
	case config.Status:
		return m.statusTasks
	case config.Priority:
		return m.priorityTasks
	default:
		return m.priorityTasks
	}
}

func (m *Model) updateSection(id int, sType string, msg tea.Msg) (cmd tea.Cmd) {
	var updatedSection section.Section
	switch sType {
	case tasksection.SectionType:
		updatedSection, cmd = m.statusTasks[id].Update(msg)
		m.statusTasks[id] = updatedSection
	}

	// currSection := m.GetCurrSection()
	// if currSection != nil && id == currSection.GetId() {
	// 	if _, ok := msg.(prssection.SectionPullRequestsFetchedMsg); ok {
	// 		cmd = m.onViewedRowChanged()
	// 	}
	// }

	return cmd
}

func (m *Model) UpdateCurrentSection(msg tea.Msg) tea.Cmd {
	section := m.GetCurrSection()

	if section == nil {
		return nil
	}

	return m.updateSection(section.GetId(), section.GetType(), msg)
}

func (m *Model) setCurrSectionId(newSectionId int) {
	m.currSectionId = newSectionId
	m.tabs.SetCurrentSectionId(newSectionId)
}

func (m *Model) onViewedRowChanged() tea.Cmd {
	var cmd tea.Cmd

	m.SyncSideBar()
	m.sidebar.ScrollToTop()

	return cmd
}
