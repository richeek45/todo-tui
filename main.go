package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Where should we begin this journey?"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		// render the row
		s += fmt.Sprintf("\n%s [%s] %s", cursor, checked, choice)
	}
	s += "\n Press q to quit.\n"

	return s
}

func main() {

	beeep.AppName = "My App Name"
	title := "Title"
	message := "Message body"
	// AppName := "My App Name"

	icon := "warning.png"

	// err := beeep.Notify(title, message, icon)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	beeep.AppName = "My App Name"

	err := beeep.Alert(title, message, icon)
	if err != nil {
		panic(err)
	}

	// cmd, err := exec.LookPath("terminal-notifier")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// args := []string{"-title", title, "-message", message, "-group", AppName, "-appIcon", img, "-sound", "default"}

	// c := exec.Command(cmd, args...)
	// c.Run()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
