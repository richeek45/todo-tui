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
	Id          int
	ctx         *context.ProgramContext
	Config      config.SectionConfig
	IsSearching bool
	IsLoading   bool
	PageInfo    *config.PageInfo
	Spinner     spinner.Model
	SearchBar   search.Model
	Table       table.Model
	Type        string
	Columns     []table.Column
	TotalCount  int
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
		ctx:      ctx,
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

	emptyMsg := m.ctx.Styles.Section.EmptyStateStyle.Render(
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
		Width:  max(0, m.ctx.Styles.Section.ContainerStyle.GetHorizontalPadding()),
		Height: max(0, m.ctx.MainContentHeight-common.SearchHeight),
	}
}

type Section interface {
	Identifier
	Table
	Search
	GetConfig() config.SectionConfig
	UpdateProgramContext(ctx *context.ProgramContext)
	GetTotalCount() int
}

type Identifier interface {
	GetId()
	GetType()
}

type Table interface {
	CurrRow()
	NextRow()
	PrevRow()
	FirstItem()
	LastItem()
}

type Search interface {
	IsSearchFocused()
	SetIsSearching()
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
				m.ctx.Styles.Section.KeyStyle.Render("/"),
				m.ctx.Styles.Section.KeyStyle.Render("Enter"),
			),
		)
	}

	return m.Table.View()
}

func (m BaseModel) View() string {
	search := m.SearchBar.View()

	return m.ctx.Styles.Section.ContainerStyle.
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
