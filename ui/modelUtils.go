package ui

import (
	"github.com/richeek45/todo-tui/components/section"
	"github.com/richeek45/todo-tui/config"
	"github.com/richeek45/todo-tui/models"
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

func (m *Model) getCurrRowData() *models.Task {
	section := m.GetCurrSection()
	if section == nil {
		return nil
	}
	return section.GetCurrRow()
}

func (m *Model) getPreviousSectionId() int {
	return max(0, m.currSectionId-1)
}

func (m *Model) getNextSectionId() int {
	return min(len(m.getCurrentViewSections())-1, m.currSectionId+1)
}

func (m *Model) getSectionAt(id int) section.Section {
	sections := m.getCurrentViewSections()
	if len(sections) <= id {
		return nil
	}
	return sections[id]
}
