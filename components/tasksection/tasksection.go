package tasksection

import (
	ctx "context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/components/table"
	"github.com/richeek45/todo-tui/components/taskrow"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/models"
)

type Model struct {
	section.BaseModel
	Tasks []models.Task
}

type SectionTaskDataFetchedMsg struct {
	Tasks          []models.Task
	TotalCount     int
	direction      string
	PaginatedTodos *models.PaginatedTodos
}

// need to pass it as prop with the FilterType in NewModel func
// Check updateSection of ui.go
const SectionType = "task"

func NewModel(
	id int,
	ctx *context.ProgramContext,
	cfg config.SectionConfig,
) Model {
	m := Model{}
	m.BaseModel = section.NewModel(
		ctx,
		section.NewSectionOptions{
			Id: id,
			Config: config.SectionConfig{
				Title:       cfg.Title,
				FilterType:  cfg.FilterType,
				FilterValue: cfg.FilterValue,
				Limit:       cfg.Limit,
				Type:        cfg.Type,
			},
			Type:    SectionType,
			Columns: GetSectionColumns(ctx),
		},
	)

	return m
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.IsSearchFocused() {
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.SearchBar.SetValue(m.SearchValue)
				blinkCmd := m.SetIsSearching(false)
				return m, blinkCmd
			case tea.KeyEnter:
				m.SearchValue = m.SearchBar.Value()
				m.SetIsSearching(false)
				m.ResetRows()
				m.PaginatedTodos = nil
				m.CurrentPage = 0
				return m, tea.Batch(m.FetchNextPageSectionRows("next")...)
			}
		}
	case SectionTaskDataFetchedMsg:
		m.PaginatedTodos = msg.PaginatedTodos
		if msg.PaginatedTodos != nil {
			m.Tasks = msg.PaginatedTodos.Todos
			if msg.direction == "next" {
				m.CurrentPage += 1
			} else {
				m.CurrentPage = max(m.CurrentPage-1, 0)
			}
		} else {
			m.Tasks = []models.Task{}
		}

		m.TotalCount = msg.TotalCount
		m.SetIsLoading(false)
		m.Ctx.Loading = false
		m.Table.SetRows(m.BuildRows())
		// m.Table.UpdateLastUpdated(time.Now())
		m.UpdateTotalItemsCount(m.TotalCount)
	}

	search, searchCmd := m.SearchBar.Update(msg)
	m.Table.SetRows(m.BuildRows())
	m.SearchBar = search

	table, tableCmd := m.Table.Update(msg)
	m.Table = table

	return m, tea.Batch(cmd, searchCmd, tableCmd)
}

func GetSectionColumns(ctx *context.ProgramContext) []table.Column {
	layout := ctx.Config.Defaults.Layout

	return []table.Column{
		{
			Title: "Task ID",
			Width: layout.TaskId.Width,
		},
		{
			Title: "Title",
			Width: layout.Title.Width,
		},
		{
			Title: "Status",
			Width: layout.ReviewStatus.Width,
		},
		{
			Title: "Priority",
			Width: layout.Priority.Width,
		},
	}
}

func (m *Model) BuildRows() []table.Row {
	var rows []table.Row
	// currItem := m.Table.GetCurrItem()

	for _, task := range m.Tasks {
		taskModel := taskrow.TaskRow{
			Ctx:     m.Ctx,
			Data:    task,
			Columns: m.Table.Columns,
		}

		rows = append(rows, taskModel.ToTableRow())
	}

	if rows == nil {
		rows = []table.Row{}
	}
	return rows
}

func (m *Model) NumRows() int {
	return len(m.Tasks)
}

func (m *Model) GetCurrRow() *models.Task {
	if len(m.Tasks) == 0 {
		return nil
	}

	i := m.Table.GetCurrItem()

	return &m.Tasks[i]
}

func (m *Model) FetchNextPageSectionRows(direction string) []tea.Cmd {
	if m == nil {
		return nil
	}

	if m.PaginatedTodos != nil {
		if direction == "next" && m.PaginatedTodos.HasNext {
			m.Pagination.Cursor = m.PaginatedTodos.NextCursor
		}

		if direction == "prev" && m.PaginatedTodos.HasPrev {
			m.Pagination.Cursor = m.PaginatedTodos.PrevCursor
		}
	}

	var cmds []tea.Cmd

	fetchCmd := func() tea.Msg {
		limit := m.Config.Limit

		if limit == nil {
			limit = &m.Ctx.Config.Defaults.TaskLimit
		}

		var filter models.TodoFilter

		if m.Config.FilterType == string(config.Priority) {
			filter.Priority = models.Priority(m.Config.FilterValue)
		}

		if m.Config.FilterType == string(config.Status) {
			filter.Status = models.Status(m.Config.FilterValue)
		}

		if m.SearchValue != "" {
			filter.Search = m.SearchValue
		}

		paginatedTodos, err := m.Ctx.Repo.GetTodosWithCursor(ctx.TODO(), filter, m.Pagination, direction)

		if err != nil {
			log.Fatal(err)
		}

		tasks := make([]models.Task, 0)
		for _, task := range paginatedTodos.Todos {
			tasks = append(tasks, task)
		}

		return constants.TaskFinishedMsg{
			SectionId: m.Id,
			Msg: SectionTaskDataFetchedMsg{
				Tasks:          tasks,
				TotalCount:     paginatedTodos.TotalCount,
				PaginatedTodos: paginatedTodos,
				direction:      direction,
			},
		}
	}

	cmds = append(cmds, fetchCmd)

	m.Ctx.Loading = true
	m.IsLoading = true
	isFirstFetch := m.PaginatedTodos == nil
	if isFirstFetch {
		m.SetIsLoading(true)
		cmds = append(cmds, m.Table.StartLoadingSpinner())
	}

	return cmds
}

func FetchAllSections(
	ctx *context.ProgramContext,
	taskSections []section.Section,
	filterType config.FilterType,
) (sections []section.Section, fetchAllCmd tea.Cmd) {

	var sectionConfig []config.SectionConfig

	switch filterType {
	case config.Category:
		sectionConfig = ctx.Config.TaskSections
	case config.Priority:
		sectionConfig = ctx.Config.PrioritySections
	case config.Status:
		sectionConfig = ctx.Config.TaskSections
	}

	fetchTaskCmds := make([]tea.Cmd, 0, len(sectionConfig))
	sections = make([]section.Section, 0, len(sectionConfig))
	for i, sectionConfig := range sectionConfig {
		sectionModel := NewModel(
			i,
			ctx,
			sectionConfig,
		)
		if len(taskSections) > 0 && len(taskSections) >= i && taskSections[i] != nil {
			oldSection := taskSections[i].(*Model)
			sectionModel.Tasks = oldSection.Tasks
			sectionModel.LastFetchTaskId = oldSection.LastFetchTaskId
		}
		sections = append(sections, &sectionModel)

		fetchTaskCmds = append(
			fetchTaskCmds,
			sectionModel.FetchNextPageSectionRows("next")...)
	}
	return sections, tea.Batch(fetchTaskCmds...)
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
	m.Table.SetIsLoading(val)
}

func (m *Model) FetchNextPage() tea.Cmd {
	return tea.Batch(m.FetchNextPageSectionRows("next")...)
}

func (m *Model) FetchPrevPage() tea.Cmd {
	return tea.Batch(m.FetchNextPageSectionRows("prev")...)
}
