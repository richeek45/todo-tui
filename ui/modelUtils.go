package ui

import (
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/config"
)

func (m *Model) GetCurrSection() section.Section {
	sections := m.GetCurrViewSections()

	if len(sections) == 0 || m.currSectionId >= len(sections) {
		return nil
	}

	return sections[m.currSectionId]
}

func (m *Model) GetCurrViewSections() []section.Section {
	switch m.ctx.View {
	case config.Status:
		return m.statusTasks
	default:
		return m.priorityTasks
	}
}
