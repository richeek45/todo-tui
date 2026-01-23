package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	context1 "github.com/richeek45/todo-tui/context"
	"github.com/richeek45/todo-tui/pagination"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) CreateTodo(ctx context.Context, task context1.Task) error {
	query := `
		INSERT INTO todos (id, title, description, status, priority, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		task.Id,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

func (r *TodoRepository) GetTodosWithCursor(
	ctx context.Context,
	filter context1.TodoFilter,
	paginationVar context1.CursorPagination,
) (*context1.PaginatedTodos, error) {
	if paginationVar.Limit == 0 {
		paginationVar.Limit = 10
	}
	if paginationVar.OrderBy == "" {
		paginationVar.OrderBy = "created_at"
	}
	if paginationVar.OrderDir == "" {
		paginationVar.OrderDir = "DESC"
	}

	whereClause, args := buildWhereClause(filter)

	orderClause := pagination.GenerateOrderClause(paginationVar.OrderBy, paginationVar.OrderDir)

	cursorClause := ""
	if paginationVar.Cursor != "" {
		cursor, err := pagination.DecodeCursor(paginationVar.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %v", err)
		}

		comparator := "<"
		if paginationVar.OrderDir == "ASC" {
			comparator = ">"
		}

		cursorClause = fmt.Sprintf("AND (%s %s ? OR (%s = ? AND id %s ?))", cursor.OrderBy, comparator, cursor.OrderBy, comparator)

		if cursor.OrderBy == "priority" {
			priorityValue := map[string]int{"high": 1, "medium": 2, "low": 3}
			args = append(args, priorityValue[cursor.OrderBy], priorityValue[cursor.OrderBy], cursor.Id)
		} else {
			args = append(args, cursor.CreatedAt, cursor.CreatedAt, cursor.Id)
		}
	}

	query := fmt.Sprintf(`
		SELECT id, title, description, status, priority, due_date, completed_at, created_at, updated_at
		FROM todos
		WHERE 1=1 %s %s
		ORDER BY %s
		LIMIT ?`,
		whereClause, cursorClause, orderClause,
	)

	args = append(args, paginationVar.Limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []context1.Task

	for rows.Next() {
		var task context1.Task
		var dueDate, completedAt sql.NullTime

		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Priority,
			&dueDate, &completedAt, &task.CreatedAt, &task.UpdatedAt)

		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, task)
	}

	hasNext := false
	var nextCursor string

	if len(tasks) > paginationVar.Limit {
		hasNext = true
		tasks = tasks[:paginationVar.Limit]

		// Create cursor from last item
		lastTodo := tasks[len(tasks)-1]
		cursor := pagination.Cursor{
			Id:        lastTodo.Id,
			CreatedAt: lastTodo.CreatedAt,
			OrderBy:   paginationVar.OrderBy,
			OrderDir:  paginationVar.OrderDir,
		}

		nextCursor, err = pagination.EncodeCursor(cursor)
		if err != nil {
			return nil, err
		}
	}

	return &context1.PaginatedTodos{
		Todos:      tasks,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	}, nil

}

func buildWhereClause(filter context1.TodoFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Priority != "" {
		conditions = append(conditions, "priority = ?")
		args = append(args, filter.Priority)
	}

	if filter.FromDate != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, filter.FromDate)
	}

	if filter.ToDate != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, filter.ToDate)
	}

	if filter.CompletedFrom != nil {
		conditions = append(conditions, "completed_at >= ?")
		args = append(args, filter.CompletedFrom)
	}

	if filter.CompletedTo != nil {
		conditions = append(conditions, "completed_at <= ?")
		args = append(args, filter.CompletedTo)
	}

	if filter.Search != "" {
		conditions = append(conditions, "(title LIKE ? OR description LIKE ?)")
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "AND " + strings.Join(conditions, " AND "), args
}
