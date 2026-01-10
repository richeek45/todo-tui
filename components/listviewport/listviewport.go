package listviewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/constants"
	"github.com/richeek45/todo-tui/context"
)

type Model struct {
	ctx             *context.ProgramContext
	viewport        viewport.Model
	topBoundId      int
	bottomBoundId   int
	currId          int
	NumItems        int
	NumCurrentItems int
	NumTotalItems   int
	ListItemHeight  int
	ItemTypeLabel   string
}

func NewModel(
	ctx *context.ProgramContext,
	dimensions constants.Dimensions,
	itemTypeLable string,
	numItems, listItemHeight int,
) Model {
	model := Model{
		ctx:    ctx,
		currId: 0,
		viewport: viewport.Model{
			Width:  dimensions.Width,
			Height: dimensions.Height,
		},
		topBoundId:     0,
		ItemTypeLabel:  itemTypeLable,
		NumItems:       numItems,
		ListItemHeight: listItemHeight,
	}
	model.bottomBoundId = min(model.NumCurrentItems-1, model.getNumTodosPerPage()-1)

	return model
}

func (m *Model) SetTotalItems(total int) {
	m.NumTotalItems = total
}

func (m *Model) SetNumItems(numItems int) {
	m.NumCurrentItems = numItems
	m.bottomBoundId = min(m.NumCurrentItems-1, m.getNumTodosPerPage()-1)
}

func (m *Model) getNumTodosPerPage() int {
	if m.ListItemHeight == 0 {
		return 0
	}
	return m.viewport.Height / m.ListItemHeight
}

func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) ResetCurrItem() {
	m.currId = 0
	m.viewport.GotoTop()
}

func (m *Model) GetCurrItem() int {
	return m.currId
}

func (m *Model) FirstItem() int {
	m.currId = 0
	m.viewport.GotoTop()
	return m.currId
}

func (m *Model) LastItem() int {
	m.currId = m.NumCurrentItems - 1
	m.viewport.GotoBottom()
	return m.currId
}

func (m *Model) NextItem() int {
	atBottomOfViewport := m.currId >= m.bottomBoundId
	if atBottomOfViewport {
		m.topBoundId += 1
		m.bottomBoundId += 1
		m.viewport.ScrollDown(m.ListItemHeight)
	}

	newCurrId := min(m.currId+1, m.NumCurrentItems-1)
	newCurrId = max(newCurrId, 0)
	m.currId = newCurrId
	return m.currId
}

func (m *Model) PrevItem() int {
	atTopOfViewport := m.currId <= m.topBoundId
	if atTopOfViewport {
		m.topBoundId -= 1
		m.bottomBoundId -= 1
		m.viewport.ScrollUp(m.ListItemHeight)
	}
	newCurrId := min(m.currId-1, 0)
	m.currId = newCurrId
	return m.currId
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.viewport.Height = max(0, dimensions.Height)
	m.viewport.Width = max(0, dimensions.Width)
}

func (m Model) View() string {
	viewport := m.viewport.View()
	return lipgloss.NewStyle().
		Width(m.viewport.Width).
		MaxWidth(m.viewport.Width).
		Render(viewport)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
}
