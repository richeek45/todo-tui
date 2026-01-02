package context

import (
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
}
