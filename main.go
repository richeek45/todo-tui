package main

import (
	// "database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/richeek45/todo-tui/components/footer"
	"github.com/richeek45/todo-tui/components/sidebar"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/keys"
	"github.com/richeek45/todo-tui/theme"
)

const (
	// bullet   = "•"
	ellipsis = "…"
)

type item struct {
	title, desc string
}

type Model struct {
	choice   string
	list     list.Model
	sidebar  sidebar.Model
	footer   footer.Model
	quitting bool

	keys *keys.KeyMap
	ctx  *context.ProgramContext
}

type initMsg struct {
	Config config.Config
}

type itemDelegate struct {
	ShowDescription bool
}

func NewModel(items []list.Item) Model {
	m := Model{
		list:    list.New(items, itemDelegate{ShowDescription: true}, 0, 0),
		keys:    keys.Keys,
		sidebar: sidebar.NewModel(),
	}

	m.ctx = &context.ProgramContext{}

	m.footer = footer.NewModel(m.ctx)
	m.list.SetShowHelp(false)

	return m
}

func (m Model) InitScreen() tea.Msg {
	cfg := config.ParseConfig()
	return initMsg{Config: cfg}
}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		title, desc, indexValue string
		matchedRunes            []int
		s                       = list.NewDefaultItemStyles()
	)

	if i, ok := listItem.(item); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	textwidth := m.Width() / 3
	title = ansi.Truncate(title, textwidth, ellipsis)
	if d.ShowDescription {
		desc = ansi.Truncate(strings.Split(desc, "\n")[0], textwidth, ellipsis)
	}

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.PaddingLeft(6).Render(desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		indexValue = s.SelectedTitle.Render(strconv.Itoa(index + 1))
		title = s.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).Render(title)
		desc = s.SelectedDesc.PaddingLeft(5).Render(desc)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		indexValue = s.NormalTitle.Render(strconv.Itoa(index + 1))
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.PaddingLeft(6).Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s.%s\n%s", indexValue, title, desc) //nolint: errcheck
		return
	}
	fmt.Fprintf(w, "%s", title) //nolint: errcheck
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.InitScreen, tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		listCmd    tea.Cmd
		sidebarCmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.Title())
			}
		case key.Matches(msg, m.keys.TogglePreview):
			_, ok := m.list.SelectedItem().(item)
			if ok {
				m.sidebar.IsOpen = !m.sidebar.IsOpen
				m.SyncMainContentWidth()
			}
		case key.Matches(msg, m.keys.Help):
			_, v := docStyle.GetFrameSize()
			if !m.footer.ShowAll {
				m.ctx.MainContentHeight = m.ctx.ScreenHeight - v - context.ExpandedHelpHeight
			} else {
				m.ctx.MainContentHeight = m.ctx.ScreenHeight - v - context.FooterHeight
			}
			m.list.SetHeight(m.ctx.MainContentHeight)

			m.footer.ShowAll = !m.footer.ShowAll
		}
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)
	case initMsg:
		m.ctx.Config = &msg.Config
		m.ctx.Theme = theme.ParseTheme()
		m.ctx.Styles = context.InitStyles(m.ctx.Theme)

		m.sidebar.IsOpen = m.ctx.Config.Defaults.Preview.Open
		m.SyncMainContentWidth()
	}
	m.SyncProgramContext(m.ctx)
	m.SyncSideBar()

	var cmds []tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	m.sidebar, sidebarCmd = m.sidebar.Update(msg)
	cmds = append(cmds, listCmd, sidebarCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.ctx.Config == nil {
		return lipgloss.Place(m.ctx.ScreenWidth, m.ctx.ScreenHeight, lipgloss.Center, lipgloss.Center, "Reading config...")
	}

	s := strings.Builder{}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.sidebar.View()))
	s.WriteString("\n")
	s.WriteString(m.footer.View())

	return docStyle.Render(s.String())
}

func main() {
	// notify()

	// db, err := sql.Open("sqlite3", "./test.db")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// sqlStatement := `
	// 	CREATE TABLE IF NOT EXISTS users (
	// 	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	// 	name TEXT
	// 	)
	// `

	// _, err = db.Exec(sqlStatement)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Table 'users' created successfully")

	// _, err = db.Exec("INSERT INTO users(name) VALUES(?)", "John Doe")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err := db.Query("SELECT id, name FROM users")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// if err = rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		item{title: "Linux", desc: "Pretty much the best OS"},
		item{title: "Business school", desc: "Just kidding"},
		item{title: "Pottery", desc: "Wet clay is a great feeling Wet clay is a great feelingWet clay is a great feelingWet \nclay is a great feelingWet clay is a great feelingWet clay is a great feeling"},
		item{title: "Shampoo", desc: "Nothing like clean hair"},
		item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		item{title: "Stickers", desc: "The thicker the vinyl the better"},
		item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := NewModel(items)
	// m.list.SetSize()
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

func (m *Model) SyncProgramContext(ctx *context.ProgramContext) {
	m.sidebar.UpdateProgramContext(m.ctx)
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	log.Println("window size changed", "width", msg.Width, "height", msg.Height)
	h, v := docStyle.GetFrameSize()
	log.Println("frame width and height", h, v)
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	if m.footer.ShowAll {
		m.ctx.MainContentHeight = msg.Height - v - context.ExpandedHelpHeight
	} else {
		m.ctx.MainContentHeight = msg.Height - v - context.FooterHeight
	}
	m.list.SetSize(msg.Width-h, m.ctx.MainContentHeight)
	m.ctx.MainContentWidth = msg.Width
	m.footer.SetWidth(msg.Width)
}

func (m *Model) SyncSideBar() {
	// width := m.GetSidebarContentWidth()
	m.sidebar.SetContent(fmt.Sprintf("Whatever content I want to put here.... %d", 1000))
}

func (m *Model) SyncMainContentWidth() {
	sideBarOffset := 0
	if m.sidebar.IsOpen {
		sideBarOffset = m.ctx.Config.Defaults.Preview.Width
	}
	m.ctx.ScreenWidth = m.ctx.ScreenWidth - sideBarOffset
}
