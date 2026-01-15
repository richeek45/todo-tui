package constants

import tea "github.com/charmbracelet/bubbletea"

type Dimensions struct {
	Width  int
	Height int
}

const (
	Ellipsis = "..."
	Logo     = `TODO`
)

type TaskFinishedMsg struct {
	TaskId      string
	SectionId   int
	SectionType string
	Err         error
	Msg         tea.Msg
}
