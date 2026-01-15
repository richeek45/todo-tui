package carousel

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/richeek45/todo-tui/constants"
)

type Model struct {
	KeyMap                  KeyMap
	items                   []string
	cursor                  int
	width                   int
	height                  int
	focus                   bool
	showOverflowIndicators  bool
	leftOverflowIndicators  string
	rightOverflowIndicators string
	separator               string
	showSeparator           bool
	Styles                  Styles
	content                 string
	start                   int
	end                     int
}

type KeyMap struct {
	SelectLeft  key.Binding
	SelectRight key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		SelectLeft: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←/h", "h"),
		),
		SelectRight: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("←/l", "l"),
		),
	}
}

type Styles struct {
	Item              lipgloss.Style
	Selected          lipgloss.Style
	OverflowIndicator lipgloss.Style
	Separator         lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Item: lipgloss.NewStyle().Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Padding(0, 1).Foreground(lipgloss.Color("212")),
	}
}

func (m *Model) SetStyles(s Styles) {
	m.Styles = s
	m.UpdateSize()
}

func WithItems(items []string) Option {
	return func(m *Model) {
		m.SetItems(items)
	}
}

func WithHeight(h int) Option {
	return func(m *Model) {
		m.height = h
	}
}

func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

func WithOverflowIndicators(indicators ...string) Option {
	return func(m *Model) {
		m.showOverflowIndicators = true
		if len(indicators) > 0 {
			m.leftOverflowIndicators = indicators[0]
		}
		if len(indicators) > 1 {
			m.rightOverflowIndicators = indicators[1]
		}
	}
}

func WithSeparators(separators ...string) Option {
	return func(m *Model) {
		m.showSeparator = true
		if len(separators) > 0 {
			m.separator = separators[0]
		}
	}
}

func WithFocused(f bool) Option {
	return func(m *Model) {
		m.focus = f
	}
}

func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.Styles = s
	}
}

func WithKeyMap(km KeyMap) Option {
	return func(m *Model) {
		m.KeyMap = km
	}
}

type Option func(m *Model)

func New(opts ...Option) Model {
	m := Model{
		cursor: 0,

		KeyMap: DefaultKeyMap(),
		Styles: DefaultStyles(),

		leftOverflowIndicators:  "<",
		rightOverflowIndicators: ">",
		separator:               "|",
	}

	for _, opt := range opts {
		opt(&m)
	}

	m.UpdateSize()

	return m
}

func (m *Model) Focused() bool {
	return m.focus
}

func (m *Model) Focus() {
	m.focus = true
	m.UpdateSize()
}

func (m *Model) Blur() {
	m.focus = false
	m.UpdateSize()
}

func (m *Model) UpdateSize() {
	// we are starting the render from cursor and render left and right from that based on the available width
	leftOver := m.width

	itemsContent := ""
	currDirection := -1

	left := m.cursor
	lastLeft := left
	right := min(m.cursor+1, len(m.items))
	lastRight := right

	for len(m.items) > 0 && leftOver > 0 && (left >= 0 || right < len(m.items)) {
		if currDirection < 0 && left >= 0 {
			lItem := m.renderItem(left, leftOver)
			leftOver -= lipgloss.Width(lItem)
			itemsContent = lipgloss.JoinHorizontal(lipgloss.Top, lItem, itemsContent)
			lastLeft = left
			left--
		} else if currDirection > 0 && right < len(m.items) {
			rItem := m.renderItem(right, leftOver)
			leftOver -= lipgloss.Width(rItem)
			itemsContent = lipgloss.JoinHorizontal(lipgloss.Top, itemsContent, rItem)
			lastRight = right
			right++
		}

		if left < 0 {
			currDirection = 1
		} else if right > len(m.items)-1 {
			currDirection = -1
		} else {
			currDirection = currDirection * -1
		}
	}

	lastRight = min(lastRight, len(m.items)-1)

	m.start = lastLeft
	m.end = lastRight

	l := m.width
	loIndicator, roIndicator := "", ""

	if m.showOverflowIndicators && lastLeft != 0 {
		loIndicator = m.Styles.OverflowIndicator.Render(m.leftOverflowIndicators)
		l -= lipgloss.Width(loIndicator)
	}
	if m.showOverflowIndicators && lastRight != len(m.items)-1 {
		roIndicator = m.Styles.OverflowIndicator.Render(m.rightOverflowIndicators)
		l -= lipgloss.Width(roIndicator)
	}

	if loIndicator != "" {
		truncate := lipgloss.Width(itemsContent) - l + 1
		itemsContent = ansi.TruncateLeft(itemsContent, truncate, "")
		if truncate > 0 {
			itemsContent = lipgloss.JoinHorizontal(lipgloss.Center, m.Styles.Item.Inline(true).Render(constants.Ellipsis), itemsContent)
		}
	} else {
		w := lipgloss.Width(itemsContent)
		if w > l {
			itemsContent = ansi.Truncate(itemsContent, l, "")
			itemsContent = lipgloss.JoinHorizontal(lipgloss.Center, itemsContent, m.Styles.Item.Inline(true).Render(constants.Ellipsis))
		}
	}

	m.content = lipgloss.NewStyle().Height(m.height).Render(lipgloss.JoinHorizontal(lipgloss.Center, loIndicator, itemsContent, roIndicator))

}

func (m *Model) renderItem(ItemId int, maxWidth int) string {
	var item string

	if ItemId == m.cursor {
		item = m.Styles.Selected.Render(m.items[ItemId])
	} else if ItemId < m.cursor {
		r := m.Styles.Item.Render(m.items[ItemId])
		truncate := lipgloss.Width(r) - maxWidth - 1
		item = ansi.TruncateLeft(r, truncate, "")
		if truncate > 0 {
			item = lipgloss.JoinHorizontal(lipgloss.Center,
				m.Styles.Item.Inline(true).Render(constants.Ellipsis), item)
		}
	} else {
		r := m.Styles.Item.Render(m.items[ItemId])
		item = ansi.Truncate(r, maxWidth, m.Styles.Item.Inline(true).Render(constants.Ellipsis))
	}

	if m.showSeparator && ItemId != len(m.items)-1 {
		return lipgloss.JoinHorizontal(lipgloss.Center, item, m.Styles.Separator.Render(m.separator))
	}

	return item
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.SelectLeft):
			m.MoveLeft()
		case key.Matches(msg, m.KeyMap.SelectRight):
			m.MoveRight()
		}
	}

	return m, nil
}

func (m *Model) View() string {
	return m.content
}

func (m *Model) SelectedItem() string {
	return m.items[m.cursor]
}

func (m *Model) Items() []string {
	return m.items
}

func (m *Model) SetItems(items []string) {
	m.items = items
}

func (m *Model) SetWidth(w int) {
	m.width = w
	m.UpdateSize()
}

func (m *Model) SetHeight(h int) {
	m.height = h
	m.UpdateSize()
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return m.height
}

func (m *Model) Cursor() int {
	return m.cursor
}

func (m *Model) HasRightItems() bool {
	return m.end < len(m.items)
}

func (m *Model) HasLeftItems() bool {
	return m.start > 0
}

func (m *Model) SetCursor(n int) {
	m.cursor = m.clamp(n, 0, len(m.items)-1)
	m.UpdateSize()
}

func (m *Model) MoveLeft() {
	m.cursor = m.clamp(m.cursor-1, 0, len(m.items)-1)
}

func (m *Model) MoveRight() {
	m.cursor = m.clamp(m.cursor+1, 0, len(m.items)-1)
}

func (m *Model) clamp(v, low, high int) int {
	return min(max(v, low), high)
}
