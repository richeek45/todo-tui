package taskrow

import (
	"github.com/richeek45/todo-tui/components/table"
	"github.com/richeek45/todo-tui/context"
)

type TaskRow struct {
	Ctx     *context.ProgramContext
	Data    context.Task
	Columns []table.Column
}

// get data from t.Data for each component in taskRow and create separate functions that returns string
func (t *TaskRow) ToTableRow() table.Row {

	return table.Row{
		t.Data.StartText,
		t.Data.FinishedText,
		string(t.Data.Status),
		"Now This",
	}

}
