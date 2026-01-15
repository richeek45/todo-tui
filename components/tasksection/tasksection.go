package tasksection

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/components/table"
	"github.com/richeek45/todo-tui/components/taskrow"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
)

type Model struct {
	section.BaseModel
	Tasks []context.Task
}

type SectionTaskDataFetchedMsg struct {
	Tasks      []context.Task
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
				Title:   cfg.Title,
				Filters: cfg.Filters,
				Limit:   cfg.Limit,
				Type:    cfg.Type,
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
			Title: "Created At",
			Width: layout.CreatedAt.Width,
		},
		{
			Title: "Updated At",
			Width: layout.UpdatedAt.Width,
		},
		{
			Title: "Title",
			Width: layout.Title.Width,
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

		// fetch the task data from database
		// 		res, err := data.FetchTaskData(m.GetFilters(), *limit, m.PageInfo)
		// if err != nil {
		// 	return constants.TaskFinishedMsg{
		// 		SectionId:   m.Id,
		// 		SectionType: m.Type,
		// 		TaskId:      taskId,
		// 		Err:         err,
		// 	}
		// }

		tasks := make([]context.Task, 0)

		return constants.TaskFinishedMsg{
			SectionId:   m.Id,
			SectionType: m.Type,
			TaskId:      taskId,
			Msg: SectionTaskDataFetchedMsg{
				Tasks:      tasks,
				TotalCount: 0,
				PageInfo:   *m.PageInfo,
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
) (sections []section.Section, fetchAllCmd tea.Cmd) {

	fetchPRsCmds := make([]tea.Cmd, 0, len(ctx.Config.TaskSections))
	sections = make([]section.Section, 0, len(ctx.Config.TaskSections))
	for i, sectionConfig := range ctx.Config.TaskSections {
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
		// fetchPRsCmds = append(
		// 	fetchPRsCmds,
		// 	sectionModel.FetchNextPageSectionRows()...)
	}
	return sections, tea.Batch(fetchPRsCmds...)
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
	m.Table.SetIsLoading(val)
}
