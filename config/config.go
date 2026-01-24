package config

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/richeek45/todo-tui/utils"
)

type FilterType string

const (
	Priority FilterType = "priority"
	Status   FilterType = "status"
	Category FilterType = "category"
)

type PageInfo struct {
	HasNextPage bool
	StartCursor string
	EndCursor   string
}

type SectionConfig struct {
	Title   string
	Filters string
	Limit   *int
	Type    *FilterType
}

type PreviewConfig struct {
	Open  bool
	Width int
}

type ColumnConfig struct {
	Width  *int
	Hidden *bool
}

type LayoutConfig struct {
	TaskId       ColumnConfig
	UpdatedAt    ColumnConfig
	CreatedAt    ColumnConfig
	Title        ColumnConfig
	Description  ColumnConfig
	ReviewStatus ColumnConfig
	State        ColumnConfig
	Ci           ColumnConfig
	Lines        ColumnConfig
}

type Defaults struct {
	Preview   PreviewConfig
	View      FilterType
	Layout    LayoutConfig
	TaskLimit int
}

type Config struct {
	Defaults     Defaults
	TaskSections []SectionConfig
}

var DefaultConfig = &Config{
	Defaults: Defaults{
		Preview: PreviewConfig{
			Open:  true,
			Width: 70,
		},
		TaskLimit: 10,
		View:      Status,
		Layout: LayoutConfig{
			UpdatedAt: ColumnConfig{
				Width: utils.IntPtr(lipgloss.Width("Updated at  ")),
			},
			CreatedAt: ColumnConfig{
				Width: utils.IntPtr(lipgloss.Width("Created at  ")),
			},
			TaskId: ColumnConfig{
				Width: utils.IntPtr(9),
			},
			Title: ColumnConfig{
				Width: utils.IntPtr(25),
			},
			Description: ColumnConfig{
				Width: utils.IntPtr(15),
			},
			ReviewStatus: ColumnConfig{
				Width: utils.IntPtr(15),
			},
			Lines: ColumnConfig{
				Width: utils.IntPtr(lipgloss.Width(" +31.4k -31.6k ")),
			},
		},
	},
	TaskSections: []SectionConfig{
		{
			Title:   "TODO",
			Filters: "status=todo",
		},
		{
			Title:   "IN PROGRESS",
			Filters: "status=inprogress",
		},
		{
			Title:   "DONE",
			Filters: "status=completed",
		},
	},
}

func ParseConfig() Config {
	return *DefaultConfig
}
