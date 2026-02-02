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
}

type initMsg struct {
	Config config.Config
}

func NewModel() Model {

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
		keys:        keys.Keys,
		sidebar:     sidebar.NewModel(),
		taskSpinner: spinner.Model{Spinner: spinner.Ellipsis},
		addTaskForm: form.NewModel(),
	}

	m.ctx = &context.ProgramContext{
		Repo:         taskRepo,
		Loading:      true,
		CurrentState: context.StateBrowsing,
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
			cmd = m.updateSection(currSection.GetId(), msg)
			return m, cmd
		}

		switch m.ctx.CurrentState {
		case context.StateBrowsing:
			m, cmd = m.handleBrowsingKeys(msg, currSection, cmd, cmds)
		case context.StateAdding:
			m, cmd = m.handleAddingTaskKeys(msg, cmd)
		case context.StateEditing:
			m, cmd = m.handleEditingKeys(msg, cmd)

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
		sectionCmd := m.updateSection(msg.SectionId, msg.Msg)
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

	if m.ctx.Loading {
		content.WriteString(m.taskSpinner.View())
		return content.String()
	}

	switch m.ctx.CurrentState {
	case context.StateBrowsing:
		content.WriteString(m.renderBrowsingView())
	case context.StateAdding:
		content.WriteString(m.addTaskForm.View())
	case context.StateEditing:
		content.WriteString(m.addTaskForm.View())
	}

	content.WriteString("\n\n")
	content.WriteString(m.footer.View())

	return lipgloss.NewStyle().Render(content.String())
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

	m := NewModel()

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
			cmds = append(cmds, currSection.FetchNextPageSectionRows("next")...)
		}
		// cmd = m.onViewedRowChanged()
	case key.Matches(msg, m.keys.Up):
		currSection.PrevRow()
	case key.Matches(msg, m.keys.Search):
		if currSection != nil {
			cmd = currSection.SetIsSearching(true)
			return m, cmd
		}
	case key.Matches(msg, m.keys.Filter):
		if m.ctx.View == config.Status {
			m.ctx.View = config.Priority
		} else {
			m.ctx.View = config.Status
		}

		newSections, cmd := m.fetchAllViewSections()
		m.setCurrentViewSections(newSections)
		m.setCurrSectionId(1)
		return m, cmd

	case key.Matches(msg, m.keys.NextPage):
		section := m.GetCurrSection()
		fetchCmd := section.FetchNextPage()
		cmds = append(cmds, fetchCmd)

	case key.Matches(msg, m.keys.PrevPage):
		section := m.GetCurrSection()
		fetchCmd := section.FetchPrevPage()
		cmds = append(cmds, fetchCmd)

	case key.Matches(msg, m.keys.AddTask):
		m.ctx.CurrentState = context.StateAdding
	case key.Matches(msg, m.keys.EditTask):
		m.ctx.CurrentState = context.StateEditing
		currRowData := m.getCurrRowData()
		m.addTaskForm.SetTaskValues(currRowData)
	case key.Matches(msg, m.keys.DeleteTask):
		currRowData := m.getCurrRowData()
		m.addTaskForm.DeleteTask(currRowData.Id)
		m.ctx.CurrentState = context.StateBrowsing
	case key.Matches(msg, m.keys.TogglePreview):
		m.sidebar.IsOpen = !m.sidebar.IsOpen
		m.SyncMainContentWidth()
	case key.Matches(msg, m.keys.Help):
		m.footer.ShowAll = !m.footer.ShowAll
		if m.footer.ShowAll {
			m.ctx.MainContentHeight = m.ctx.MainContentHeight +
				context.FooterHeight - context.ExpandedHelpHeight
		} else {
			m.ctx.MainContentHeight = m.ctx.MainContentHeight +
				context.ExpandedHelpHeight - context.FooterHeight
		}

	}

	return m, cmd
}

func (m Model) handleAddingTaskKeys(
	msg tea.KeyMsg,
	cmd tea.Cmd,
) (Model, tea.Cmd) {

	switch msg.String() {
	case "esc":
		m.ctx.CurrentState = context.StateBrowsing
		m.addTaskForm.ResetAddForm()
		return m, nil

	}
	m.addTaskForm, _ = m.addTaskForm.Update(msg, cmd)

	return m, cmd
}

func (m Model) handleEditingKeys(
	msg tea.KeyMsg,
	cmd tea.Cmd,
) (Model, tea.Cmd) {

	m.addTaskForm, _ = m.addTaskForm.Update(msg, cmd)

	switch msg.String() {
	case "esc":
		m.ctx.CurrentState = context.StateBrowsing
		m.addTaskForm.ResetAddForm()
		return m, cmd
	case "ctrl+s":
		newSections, fetchSectionsCmds := m.fetchAllViewSections()
		m.setCurrentViewSections(newSections)
		m.setCurrSectionId(1)
		return m, tea.Batch(cmd, fetchSectionsCmds)
	}

	return m, cmd
}

func (m *Model) SyncProgramContext(ctx *context.ProgramContext) {

	for _, section := range m.getCurrentViewSections() {
		section.UpdateProgramContext(m.ctx)
	}
	m.sidebar.UpdateProgramContext(m.ctx)
	m.tabs.UpdateProgramContext(m.ctx)
	m.footer.UpdateProgramContext(m.ctx)
	m.addTaskForm.UpdateProgramContext(m.ctx)
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
	m.SyncMainContentWidth()
	m.addTaskForm.DescInput.SetWidth(msg.Width / 2)

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

	switch m.ctx.View {
	case config.Priority:
		sections, taskCmds1 := tasksection.FetchAllSections(m.ctx, m.priorityTasks, m.ctx.View)
		cmds = append(cmds, taskCmds1)
		return sections, tea.Batch(cmds...)
	case config.Status:
		sections, taskCmds2 := tasksection.FetchAllSections(m.ctx, m.statusTasks, m.ctx.View)
		cmds = append(cmds, taskCmds2)
		return sections, tea.Batch(cmds...)
	default:
		sections, taskCmds2 := tasksection.FetchAllSections(m.ctx, m.statusTasks, m.ctx.View)
		cmds = append(cmds, taskCmds2)
		return sections, tea.Batch(cmds...)
	}
}

func (m *Model) setCurrentViewSections(newSections []section.Section) {
	if newSections == nil {
		return
	}

	s := make([]section.Section, 0)

	if m.ctx.View == config.Status {
		m.statusTasks = append(s, newSections...)
		newSections = m.statusTasks
	}

	if m.ctx.View == config.Priority {
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

func (m *Model) updateSection(id int, msg tea.Msg) (cmd tea.Cmd) {
	var updatedSection section.Section
	switch m.ctx.View {
	case config.Status:
		updatedSection, cmd = m.statusTasks[id].Update(msg)
		m.statusTasks[id] = updatedSection
	case config.Priority:
		updatedSection, cmd = m.priorityTasks[id].Update(msg)
		m.priorityTasks[id] = updatedSection
	default:
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

	return m.updateSection(section.GetId(), msg)
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
