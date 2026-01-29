package form

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/models"
)

var (
	filterInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("69")).
				Padding(0, 1)

	filterLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("69")).
				Bold(true)
)

type Model struct {
	TitleInput    textinput.Model
	DescInput     textarea.Model
	PriorityInput textinput.Model
	StatusInput   textinput.Model
}

func NewModel() Model {
	addTitleInput := textinput.New()
	addTitleInput.Placeholder = "Todo title"
	addTitleInput.Focus()
	addTitleInput.CharLimit = 100
	addTitleInput.Width = 50

	addDescInput := textarea.New()
	addDescInput.Placeholder = "Description (optional)"
	addDescInput.CharLimit = 500

	addPriorityInput := textinput.New()
	addPriorityInput.Placeholder = "low/medium/high"
	addPriorityInput.CharLimit = 10
	addPriorityInput.SetValue(string(models.PriorityMedium))
	addPriorityInput.Width = 15

	addStatusInput := textinput.New()
	addStatusInput.Placeholder = "pending/in_progress/completed"
	addStatusInput.CharLimit = 12
	addStatusInput.SetValue(string(models.PriorityMedium))
	addTitleInput.Width = 50

	return Model{
		TitleInput:    addTitleInput,
		DescInput:     addDescInput,
		PriorityInput: addPriorityInput,
		StatusInput:   addStatusInput,
	}
}

func (m Model) Update(msg tea.KeyMsg, cmd tea.Cmd) (Model, tea.Cmd) {
	cmds := []tea.Cmd{cmd}
	switch msg.String() {
	case "tab", "shift+tab":
		if m.TitleInput.Focused() {
			m.TitleInput.Blur()
			m.DescInput.Focus()
		} else if m.DescInput.Focused() {
			m.DescInput.Blur()
			m.PriorityInput.Focus()
		} else if m.PriorityInput.Focused() {
			m.PriorityInput.Blur()
			m.StatusInput.Focus()
		} else {
			m.StatusInput.Blur()
			m.TitleInput.Focus()
		}
	// case "enter":
	// 	// save the task
	default:
		if m.TitleInput.Focused() {
			m.TitleInput, cmd = m.TitleInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.DescInput.Focused() {
			log.Print("Desc Focus")
			m.DescInput, cmd = m.DescInput.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.PriorityInput.Focused() {
			m.PriorityInput, cmd = m.PriorityInput.Update(msg)
			cmds = append(cmds, cmd)

		} else {
			m.StatusInput, cmd = m.StatusInput.Update(msg)
			cmds = append(cmds, cmd)

		}

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := strings.Builder{}

	content.WriteString(filterLabelStyle.Render("Add New Todo (Tab to navigate, Enter to save, Esc to cancel)"))
	content.WriteString("\n\n")

	log.Print("Placeholder= ", m.TitleInput.Placeholder)

	fmt.Fprintf(&content, "%s %s\n",
		filterLabelStyle.Render("Title:"),
		filterInputStyle.Render(m.TitleInput.View()))

	fmt.Fprintf(&content, "%s %s\n",
		filterLabelStyle.Render("Description:"),
		filterInputStyle.Render(m.DescInput.View()))

	fmt.Fprintf(&content, "%s %s\n",
		filterLabelStyle.Render("Priority:"),
		filterInputStyle.Render(m.PriorityInput.View()))

	fmt.Fprintf(&content, "%s %s\n",
		filterLabelStyle.Render("Status:"),
		filterInputStyle.Render(m.StatusInput.View()))

	return content.String()
}

func (m *Model) ResetAddForm() {
	m.TitleInput.SetValue("")
	m.DescInput.SetValue("")
	m.PriorityInput.SetValue(string(models.PriorityMedium))
	m.StatusInput.SetValue(string(models.NotStarted))
	m.TitleInput.Focus()
}
