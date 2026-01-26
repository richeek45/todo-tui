package tasksection

import (
	ctx "context"
	"fmt"
	"log"
	"time"

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
	Tasks      []models.Task
	TotalCount int
	PageInfo   config.PageInfo
	TaskId     string
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
				return m, tea.Batch(m.FetchNextPageSectionRows()...)
			}
		}
	case SectionTaskDataFetchedMsg:
		if m.LastFetchTaskId == msg.TaskId {
			if m.PageInfo != nil {
				m.Tasks = append(m.Tasks, msg.Tasks...)
			} else {
				m.Tasks = msg.Tasks
			}
			m.TotalCount = msg.TotalCount
			m.PageInfo = &msg.PageInfo
			m.SetIsLoading(false)
			m.Table.SetRows(m.BuildRows())
			// m.Table.UpdateLastUpdated(time.Now())
			m.UpdateTotalItemsCount(m.TotalCount)
		}
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
			Title: "Description",
			Width: layout.Description.Width,
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

func (m *Model) FetchNextPageSectionRows() []tea.Cmd {
	if m == nil {
		return nil
	}

	if m.PageInfo != nil && !m.PageInfo.HasNextPage {
		return nil
	}

	var cmds []tea.Cmd

	startCursor := time.Now().String()
	isFirstFetch := m.LastFetchTaskId == ""
	taskId := fmt.Sprintf("fetching_tasks_%d_%s", m.Id, startCursor)
	m.LastFetchTaskId = taskId

	fetchCmd := func() tea.Msg {
		limit := m.Config.Limit

		if limit == nil {
			limit = &m.Ctx.Config.Defaults.TaskLimit
		}

		var filter models.TodoFilter
		var pagination models.CursorPagination

		if m.BaseModel.Type != "" {
			filter.Status = models.Status(m.BaseModel.Config.FilterValue)
		}

		pagination.OrderBy = "ASC"
		pagination.OrderDir = "next"

		paginatedTodos, err := m.Ctx.Repo.GetTodosWithCursor(ctx.TODO(), filter, pagination)

		if err != nil {
			log.Fatal(err)
		}

		log.Print("paginatedTodos= ", paginatedTodos.HasNext)

		tasks := make([]models.Task, 0)
		for _, task := range paginatedTodos.Todos {
			log.Print(task.Description)
			tasks = append(tasks, task)
		}

		return constants.TaskFinishedMsg{
			SectionId:   m.Id,
			SectionType: m.Type,
			TaskId:      taskId,
			Msg: SectionTaskDataFetchedMsg{
				Tasks:      tasks,
				TotalCount: 3,                                                                     // res.TotalCount
				PageInfo:   config.PageInfo{HasNextPage: false, StartCursor: "1", EndCursor: "3"}, // res.PageInfo
				TaskId:     taskId,
			},
		}
	}

	cmds = append(cmds, fetchCmd)

	m.IsLoading = true
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

	// Create different Sections for this
	switch filterType {
	case config.Category:
		sectionConfig = ctx.Config.TaskSections
	case config.Priority:
		sectionConfig = ctx.Config.TaskSections
	case config.Status:
		sectionConfig = ctx.Config.TaskSections
	}

	fetchPRsCmds := make([]tea.Cmd, 0, len(sectionConfig))
	sections = make([]section.Section, 0, len(sectionConfig))
	for i, sectionConfig := range sectionConfig {
		sectionModel := NewModel(
			i+1, // 0 is the search section
			ctx,
			sectionConfig,
		)
		if len(taskSections) > 0 && len(taskSections) >= i+1 && taskSections[i+1] != nil {
			oldSection := taskSections[i+1].(*Model)
			sectionModel.Tasks = oldSection.Tasks
			sectionModel.LastFetchTaskId = oldSection.LastFetchTaskId
		}
		sections = append(sections, &sectionModel)
		fetchPRsCmds = append(
			fetchPRsCmds,
			sectionModel.FetchNextPageSectionRows()...)
	}
	return sections, tea.Batch(fetchPRsCmds...)
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
	m.Table.SetIsLoading(val)
}
