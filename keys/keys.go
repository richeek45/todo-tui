package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	TogglePreview key.Binding
	Quit          key.Binding
	Enter         key.Binding
	PageUp        key.Binding
	PageDown      key.Binding
}

var Keys = &KeyMap{
	TogglePreview: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Toggle View"),
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
}
