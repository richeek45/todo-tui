package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up              key.Binding
	Down            key.Binding
	PreviousSection key.Binding
	NextSection     key.Binding
	TogglePreview   key.Binding
	Search          key.Binding
	Filter          key.Binding
	AddTask         key.Binding
	EditTask        key.Binding
	DeleteTask      key.Binding
	Quit            key.Binding
	Enter           key.Binding
	PageUp          key.Binding
	PageDown        key.Binding
	NextPage        key.Binding
	PrevPage        key.Binding
	Help            key.Binding
	Esc             key.Binding
}

var Keys = &KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↑/j", "move down"),
	),
	PreviousSection: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous section"),
	),
	NextSection: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next section"),
	),
	AddTask: key.NewBinding(
		key.WithKeys("a", "+"),
		key.WithHelp("a/+", "add task"),
	),
	EditTask: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	DeleteTask: key.NewBinding(
		key.WithKeys("d", "backspace"),
		key.WithHelp("d/del", "delete"),
	),
	TogglePreview: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Toggle View"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Search"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f", "F"),
		key.WithHelp("f/F", "Switch Filter"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc", "q"),
		key.WithHelp("q", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Selected"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("Page up", "Sidebar Page Up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("Page down", "Sidebar Page Down"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "Next Page"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "PrevPage"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.TogglePreview, k.PageUp, k.PageDown, k.Enter, k.Help}}
}
