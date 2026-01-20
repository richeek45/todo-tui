package table

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/richeek45/todo-tui/common"
	"github.com/richeek45/todo-tui/components/listviewport"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
)

type Model struct {
	ctx            *context.ProgramContext
	Columns        []Column
	Rows           []Row
	dimensions     constants.Dimensions
	EmptyState     *string
	loadingMessage string
	isLoading      bool
	loadingSpinner spinner.Model
	rowsViewport   listviewport.Model
}

type Column struct {
	Title         string
	Hidden        *bool
	Width         *int
	Grow          *bool
	ComputedWidth int
}

type Row []string

func NewModel(
	ctx *context.ProgramContext,
	dimensions constants.Dimensions,
	itemsTypeLable string,
	columns []Column,
	rows []Row,
	isLoading bool,
	loadingMessage string,
	emptyState *string,
) Model {
	itemHeight := 1
	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Ellipsis
	loadingSpinner.Style = lipgloss.NewStyle().Foreground(ctx.Theme.SecondaryText)

	m := Model{
		ctx:            ctx,
		Columns:        columns,
		Rows:           rows,
		dimensions:     dimensions,
		loadingSpinner: loadingSpinner,
		isLoading:      isLoading,
		EmptyState:     emptyState,
		loadingMessage: loadingMessage,
		rowsViewport: listviewport.NewModel(
			ctx,
			dimensions,
			itemsTypeLable,
			len(rows),
			itemHeight,
		),
	}

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.isLoading {
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	header := m.RenderHeader()
	body := m.RenderBody()

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

func (m Model) StartLoadingSpinner() tea.Cmd {
	return m.loadingSpinner.Tick
}

func (m *Model) getShownColumns() []Column {
	shownColumns := make([]Column, 0, len(m.Columns))

	for _, col := range m.Columns {
		if col.Hidden != nil && *col.Hidden {
			continue
		}

		shownColumns = append(shownColumns, col)
	}

	return shownColumns
}

func (m *Model) RenderHeaderColumns() []string {
	shownColumns := m.getShownColumns()
	renderedColumns := make([]string, len(shownColumns))
	takenWidth := 0
	numGrowingColumns := 0

	for i, column := range shownColumns {
		if column.Grow != nil && *column.Grow {
			numGrowingColumns += 1
			continue
		}

		if column.Width != nil {
			renderedColumns[i] = m.ctx.Styles.Table.TitleCellStyle.
				Width(*column.Width).
				MaxWidth(*column.Width).
				Render(column.Title)
			takenWidth += *column.Width
			continue
		}

		cell := m.ctx.Styles.Table.TitleCellStyle.Render(column.Title)
		renderedColumns[i] = cell
		takenWidth += lipgloss.Width(cell)
	}

	if numGrowingColumns == 0 {
		return renderedColumns
	}

	leftOverWidth := m.dimensions.Width - takenWidth
	growCellWidth := leftOverWidth / numGrowingColumns

	for i, column := range shownColumns {
		if column.Grow == nil && !*column.Grow {
			continue
		}

		renderedColumns[i] = m.ctx.Styles.Table.TitleCellStyle.
			Width(growCellWidth).
			MaxWidth(growCellWidth).
			Render(column.Title)
	}

	return renderedColumns
}

func (m *Model) RenderHeader() string {
	headerColumns := m.RenderHeaderColumns()
	header := ansi.Truncate(lipgloss.JoinHorizontal(lipgloss.Top, headerColumns...), m.dimensions.Width, constants.Ellipsis)
	return m.ctx.Styles.Table.HeaderStyle.
		Width(m.dimensions.Width).
		Height(common.HeaderHeight).Render(header)
}

func (m *Model) RenderBody() string {
	if m.isLoading {
		return lipgloss.Place(
			m.dimensions.Width,
			m.dimensions.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s%s", m.loadingSpinner.View(), m.loadingMessage),
		)
	}

	if len(m.Rows) == 0 && m.EmptyState != nil {
		return lipgloss.NewStyle().
			Width(m.dimensions.Width).
			Height(m.dimensions.Height).
			Render(*m.EmptyState)
	}

	return m.rowsViewport.View()
}

func (m *Model) renderRows(rowId int, headerColumns []string) string {
	var style lipgloss.Style

	if m.rowsViewport.GetCurrItem() == rowId {
		style = m.ctx.Styles.Table.SelectedCellStyle
	} else {
		style = m.ctx.Styles.Table.CellStyle
	}

	heaaderColumnId := 0
	renderedColumns := make([]string, 0, len(m.Columns))

	for i, column := range m.Columns {
		if column.Hidden != nil && *column.Hidden {
			continue
		}

		colWidth := lipgloss.Width(headerColumns[heaaderColumnId])
		col := m.Rows[rowId][i]
		colHeight := 1
		renderedCol := style.
			Width(colWidth).
			MaxWidth(colWidth).
			Height(colHeight).
			MaxHeight(colHeight).Render(col)

		renderedColumns = append(renderedColumns, renderedCol)
		heaaderColumnId++
	}

	return m.ctx.Styles.Table.RowStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...))

}

func (m *Model) cacheColumnWidth() {
	columns := m.RenderHeaderColumns()

	for i, column := range columns {
		if m.Columns[i].Hidden != nil && *m.Columns[i].Hidden {
			continue
		}

		m.Columns[i].ComputedWidth = lipgloss.Width(column)
	}
}

func (m *Model) SyncViewportContent() {
	headerColumns := m.RenderHeaderColumns()
	m.cacheColumnWidth()

	renderedRows := make([]string, 0, len(m.Rows))
	for i := range m.Rows {
		renderedRows = append(renderedRows, m.renderRows(i, headerColumns))
	}
	m.rowsViewport.SyncViewPort(lipgloss.JoinVertical(lipgloss.Left, renderedRows...))
}

func (m *Model) PrevItem() int {
	currItem := m.rowsViewport.PrevItem()
	m.SyncViewportContent()
	return currItem
}

func (m *Model) NextItem() int {
	nextItem := m.rowsViewport.NextItem()
	m.SyncViewportContent()
	return nextItem
}

func (m *Model) FirstItem() int {
	firstItem := m.rowsViewport.FirstItem()
	m.SyncViewportContent()
	return firstItem
}

func (m *Model) LastItem() int {
	lastItem := m.rowsViewport.LastItem()
	m.SyncViewportContent()
	return lastItem
}

func (m *Model) SetRows(rows []Row) {
	m.Rows = rows
	m.rowsViewport.SetNumItems(len(m.Rows))
	m.SyncViewportContent()
}

func (m *Model) OnLineDown() {
	m.rowsViewport.NextItem()
}

func (m *Model) OnLineUp() {
	m.rowsViewport.PrevItem()
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.rowsViewport.UpdateProgramContext(ctx)
}

func (m *Model) UpdateTotalItemsCount(count int) {
	m.rowsViewport.SetTotalItems(count)
}

func (m *Model) SetIsLoading(isLoading bool) {
	m.isLoading = isLoading
}

func (m *Model) IsLoading() bool {
	return m.isLoading
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.dimensions = dimensions
	m.rowsViewport.SetDimensions(constants.Dimensions{
		Width:  m.dimensions.Width,
		Height: m.dimensions.Height,
	})
}

func (m *Model) ResetCurrItem() {
	m.rowsViewport.ResetCurrItem()
}

func (m *Model) GetCurrItem() int {
	return m.rowsViewport.GetCurrItem()
}
