package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
	"github.com/richeek45/todo-tui/components/footer"
	"github.com/richeek45/todo-tui/components/form"
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/components/sidebar"
	"github.com/richeek45/todo-tui/components/tabs"
	"github.com/richeek45/todo-tui/components/tasksection"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/database"
	"github.com/richeek45/todo-tui/keys"
	"github.com/richeek45/todo-tui/models"
	"github.com/richeek45/todo-tui/theme"
)

const (
// bullet   = "•"
)

type Item struct {
	title, desc string
}

type state int

const (
	stateBrowsing state = iota
	stateFiltering
	stateAdding
	stateEditing
	stateDeleting
)

type Model struct {
	sidebar       sidebar.Model
	addTaskForm   form.Model
	tabs          tabs.Model
	footer        footer.Model
	taskSpinner   spinner.Model
	currSectionId int
	keys          *keys.KeyMap
	priorityTasks []section.Section
	statusTasks   []section.Section
	ctx           *context.ProgramContext
	Tasks         map[string]models.Task
	CurrentState  state
	pagination    models.CursorPagination
}

type initMsg struct {
	Config config.Config
}

func NewModel(items []Item) Model {

	db, err := database.NewDB("test.db")
	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		database.SeedSampleData(db)
	}

	taskRepo := database.NewTodoRepository(db)

	m := Model{
		keys:         keys.Keys,
		sidebar:      sidebar.NewModel(),
		taskSpinner:  spinner.Model{Spinner: spinner.Ellipsis},
		CurrentState: stateBrowsing,
		addTaskForm:  form.NewModel(),
		pagination:   models.CursorPagination{Limit: 10, OrderBy: "created_at", OrderDir: "DESC"},
	}

	m.ctx = &context.ProgramContext{
		Repo: taskRepo,
	}

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
	return tea.Batch(m.InitScreen, tea.EnterAltScreen, textarea.Blink, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd         tea.Cmd
		listCmd     tea.Cmd
		sidebarCmd  tea.Cmd
		currSection = m.GetCurrSection()
		cmds        []tea.Cmd
	)
	switch msg := msg.(type) {

	case tea.KeyMsg:

		if currSection != nil && (currSection.IsSearchFocused()) {
			cmd = m.updateSection(currSection.GetId(), currSection.GetType(), msg)
			return m, cmd
		}

		switch m.CurrentState {
		case stateBrowsing:
			m, cmd = m.handleBrowsingKeys(msg, currSection, cmd, cmds)
		case stateAdding:
			m, cmd = m.handleAddingTaskKeys(msg, cmd)
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

	case constants.TaskFinishedMsg:
		sectionCmd := m.updateSection(msg.SectionId, msg.SectionType, msg.Msg)
		cmds = append(cmds, sectionCmd)

		// syncCmd := m.SyncSideBar()
		// cmds = append(cmds, syncCmd)
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

	content := strings.Builder{}

	switch m.CurrentState {
	case stateBrowsing:
		content.WriteString(m.renderBrowsingView())
	// case stateFiltering:
	// 	content.WriteString(m.renderFilterView())
	case stateAdding:
		content.WriteString(m.addTaskForm.View())
		// case stateEditing:
		// 	content.WriteString(m.renderEditView())
	}

	content.WriteString("\n\n")
	content.WriteString(m.footer.View())

	return lipgloss.NewStyle().Padding(1, 2).Render(content.String())
}

func (m Model) renderBrowsingView() string {
	content := strings.Builder{}

	content.WriteString(m.tabs.View())

	currSection := m.GetCurrSection()

	if currSection != nil {
		content.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.GetCurrSection().View(),
			m.sidebar.View(),
		))
	}
	return content.String()
}

func Run() {
	// notify()

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

func (m Model) handleBrowsingKeys(
	msg tea.KeyMsg,
	currSection section.Section,
	cmd tea.Cmd,
	cmds []tea.Cmd,
) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.PreviousSection):
		prevSection := m.getSectionAt(m.getPreviousSectionId())
		if prevSection != nil {
			m.setCurrSectionId(prevSection.GetId())
			// cmd = m.onViewedRowChanged()
		}
	case key.Matches(msg, m.keys.NextSection):
		nextSection := m.getSectionAt(m.getNextSectionId())
		if nextSection != nil {
			m.setCurrSectionId(nextSection.GetId())
			// cmd = m.onViewedRowChanged()
		}
	case key.Matches(msg, m.keys.Down):
		prevRow := currSection.CurrRow()
		nextRow := currSection.NextRow()

		if prevRow != nextRow && nextRow == currSection.NumRows()-1 {
			cmds = append(cmds, currSection.FetchNextPageSectionRows()...)
		}
		// cmd = m.onViewedRowChanged()
	case key.Matches(msg, m.keys.Up):
		currSection.PrevRow()

	case key.Matches(msg, m.keys.Search):
		if currSection != nil {
			cmd = currSection.SetIsSearching(true)
			return m, cmd
		}

	case key.Matches(msg, m.keys.AddTask):
		if m.CurrentState == stateBrowsing {
			m.CurrentState = stateAdding
		}
	// case key.Matches(msg, m.keys.Enter):

	case key.Matches(msg, m.keys.TogglePreview):
		m.sidebar.IsOpen = !m.sidebar.IsOpen
		m.SyncMainContentWidth()
	case key.Matches(msg, m.keys.Help):
		_, v := docStyle.GetFrameSize()
		if !m.footer.ShowAll {
			m.ctx.MainContentHeight = m.ctx.ScreenHeight - v - context.ExpandedHelpHeight
		} else {
			m.ctx.MainContentHeight = m.ctx.ScreenHeight - v - context.FooterHeight
		}

		m.footer.ShowAll = !m.footer.ShowAll
	}

	return m, cmd
}

func (m Model) handleAddingTaskKeys(
	msg tea.KeyMsg,
	cmd tea.Cmd,
) (Model, tea.Cmd) {

	switch msg.String() {
	case "esc":
		m.CurrentState = stateBrowsing
		m.addTaskForm.ResetAddForm()
		return m, nil

	}
	m.addTaskForm, cmd = m.addTaskForm.Update(msg, cmd)

	return m, cmd
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
	if m.footer.ShowAll {
		m.ctx.MainContentHeight = msg.Height - v - context.ExpandedHelpHeight
	} else {
		m.ctx.MainContentHeight = msg.Height - v - context.FooterHeight
	}
	m.addTaskForm.DescInput.SetWidth(msg.Width / 2)

	m.ctx.MainContentWidth = msg.Width
	m.footer.SetWidth(msg.Width)
}

func (m *Model) SyncSideBar() {
	// width := m.sidebar.GetSidebarContentWidth()
	currRowData := m.getCurrRowData()
	if currRowData == nil {
		m.sidebar.SetContent("")
		return
	}
	m.sidebar.SetContent(string(currRowData.Description))
}

func (m *Model) SyncMainContentWidth() {
	sideBarOffset := 0
	if m.sidebar.IsOpen {
		sideBarOffset = m.ctx.Config.Defaults.Preview.Width
	}
	m.ctx.MainContentWidth = m.ctx.ScreenWidth - sideBarOffset
}

func (m *Model) fetchAllViewSections() ([]section.Section, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, m.tabs.SetAllLoading()...)

	log.Print("CTX VIEW =", m.ctx.View)

	switch m.ctx.View {
	case config.Priority:
		section, taskCmds1 := tasksection.FetchAllSections(m.ctx, m.priorityTasks, m.ctx.View)
		cmds = append(cmds, taskCmds1)
		return section, tea.Batch(cmds...)
	case config.Status:
		section, taskCmds2 := tasksection.FetchAllSections(m.ctx, m.statusTasks, m.ctx.View)
		cmds = append(cmds, taskCmds2)
		return section, tea.Batch(cmds...)
	default:
		section, taskCmds2 := tasksection.FetchAllSections(m.ctx, m.statusTasks, m.ctx.View)
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
					Title:       "",
					FilterType:  "archived",
					FilterValue: "false",
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
					Title:       "",
					FilterType:  "archived",
					FilterValue: "false"},
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
