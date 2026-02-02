package context

import (
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/database"
	"github.com/richeek45/todo-tui/theme"
)

type State int

const (
	StateBrowsing State = iota
	StateFiltering
	StateAdding
	StateEditing
	StateDeleting
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
	Repo              *database.TodoRepository
	Loading           bool
	CurrentState      State
}
