package section

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/common"
	"github.com/richeek45/todo-tui/components/search"
	"github.com/richeek45/todo-tui/components/table"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
)

type BaseModel struct {
	Id              int
	Ctx             *context.ProgramContext
	Config          config.SectionConfig
	IsSearching     bool
	IsLoading       bool
	PageInfo        *config.PageInfo
	Spinner         spinner.Model
	SearchBar       search.Model
	Table           table.Model
	SearchValue     string
	Type            string
	Columns         []table.Column
	TotalCount      int
	LastFetchTaskId string
}

type NewSectionOptions struct {
	Id      int
	ctx     *context.ProgramContext
	Config  config.SectionConfig
	Type    string
	Columns []table.Column
}

func NewModel(
	ctx *context.ProgramContext,
	options NewSectionOptions,
) BaseModel {
	m := BaseModel{
		Ctx:      ctx,
		Config:   options.Config,
		PageInfo: nil,
		Spinner:  spinner.Model{Spinner: spinner.Ellipsis},
		SearchBar: search.NewModel(ctx, search.SearchOptions{
			Prefix:       "",
			InitialValue: "Search Todo List",
			Placeholder:  "Search Todo List",
		}),
		Id:          options.Id,
		Type:        options.Type,
		Columns:     options.Columns,
		IsSearching: false,
		TotalCount:  0,
	}

	emptyMsg := m.Ctx.Styles.Section.EmptyStateStyle.Render(
		"No were found that match the given filters",
	)

	m.Table = table.NewModel(
		ctx,
		m.GetDimensions(),
		"Label",
		m.Columns,
		nil,
		false,
		"Loading...",
		&emptyMsg,
	)

	return m
}

func (m *BaseModel) GetDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  max(0, m.Ctx.MainContentWidth-m.Ctx.Styles.Section.ContainerStyle.GetHorizontalPadding()),
		Height: max(0, m.Ctx.MainContentHeight-common.SearchHeight),
	}
}

type Section interface {
	Identifier
	Component
	Table
	Search
	GetConfig() config.SectionConfig
	UpdateProgramContext(ctx *context.ProgramContext)
}

type Identifier interface {
	GetId() int
	GetType() string
}

type Component interface {
	Update(msg tea.Msg) (Section, tea.Cmd)
	View() string
}

type Table interface {
	NumRows() int
	ResetRows()
	CurrRow() int
	NextRow() int
	PrevRow() int
	FirstItem() int
	LastItem() int
	GetIsLoading() bool
	BuildRows() []table.Row
	FetchNextPageSectionRows() []tea.Cmd
}

type Search interface {
	IsSearchFocused() bool
	SetIsSearching(val bool) tea.Cmd
}

func (m *BaseModel) GetMainContent() string {
	if m.Table.Rows == nil {
		d := m.GetDimensions()
		return lipgloss.Place(
			d.Width,
			d.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf(
				"%s you can search query by pressing %s and submitting it with %s",
				lipgloss.NewStyle().Bold(true).Render(" Tip:"),
				m.Ctx.Styles.Section.KeyStyle.Render("/"),
				m.Ctx.Styles.Section.KeyStyle.Render("Enter"),
			),
		)
	}

	return m.Table.View()
}

func (m BaseModel) View() string {
	search := m.SearchBar.View()

	return m.Ctx.Styles.Section.ContainerStyle.
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				search,
				m.GetMainContent(),
			))
}

type SectionRowsFetchedMsg struct {
	SectionId int
	// TodoList    []data.RowData
}

func (m *BaseModel) ResetPageInfo() {
	m.PageInfo = nil
}

func (m *BaseModel) GetConfig() config.SectionConfig {
	return m.Config
}

func (msg SectionRowsFetchedMsg) GetSectionId() int {
	return msg.SectionId
}

func (m *BaseModel) GetId() int {
	return m.Id
}

func (m *BaseModel) GetType() string {
	return m.Type
}

func (m *BaseModel) ResetRows() {
	m.Table.Rows = nil
	m.ResetPageInfo()
	m.Table.ResetCurrItem()
}

func (m *BaseModel) CurrRow() int {
	return m.Table.GetCurrItem()
}

func (m *BaseModel) NextRow() int {
	return m.Table.NextItem()
}

func (m *BaseModel) PrevRow() int {
	return m.Table.PrevItem()
}

func (m *BaseModel) FirstItem() int {
	return m.Table.FirstItem()
}

func (m *BaseModel) LastItem() int {
	return m.Table.LastItem()
}

func (m *BaseModel) IsSearchFocused() bool {
	return m.IsSearching
}

func (m *BaseModel) GetIsLoading() bool {
	return m.IsLoading
}

func (m *BaseModel) SetIsSearching(val bool) tea.Cmd {
	m.IsSearching = val
	if val {
		m.SearchBar.Focus()
		return m.SearchBar.Init()
	} else {
		m.SearchBar.Blur()
		return nil
	}
}

func (m *BaseModel) UpdateTotalItemsCount(count int) {
	m.Table.UpdateTotalItemsCount(count)
}

func (m *BaseModel) UpdateProgramContext(ctx *context.ProgramContext) {
	m.Ctx = ctx
	newDimensions := m.GetDimensions()
	tableDimensions := constants.Dimensions{
		Height: max(0, newDimensions.Height-2),
		Width:  max(0, newDimensions.Width),
	}
	m.Table.SetDimensions(tableDimensions)
	m.Table.UpdateProgramContext(ctx)
	m.Table.SyncViewportContent()
	m.SearchBar.UpdateProgramContext(ctx)
}
