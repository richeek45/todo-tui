package context

import (
	"time"

	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/theme"
)

type ProgramContext struct {
	ScreenWidth       int
	ScreenHeight      int
	MainContentWidth  int
	MainContentHeight int
	Styles            Styles
	Config            *config.Config
	Theme             theme.Theme
	View              config.FilterType
}

type Status string

const (
	Completed  Status = "completed"
	NotStarted Status = "not-started"
	InProgress Status = "in-progress"
)

type Task struct {
	Id           string
	StartText    string
	FinishedText string
	Status       Status
	Error        error
	StartTime    time.Time
	FinishedTime *time.Time
}
