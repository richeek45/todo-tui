package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/richeek45/todo-tui/models"
	paginationPkg "github.com/richeek45/todo-tui/pagination"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) CreateTodo(ctx context.Context, task *models.Task) error {
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

func (r *TodoRepository) UpdateTodo(ctx context.Context, task *models.Task) error {
	query := `
        UPDATE todos 
        SET title = ?, 
            description = ?, 
            status = ?, 
            priority = ?, 
            due_date = ?, 
            updated_at = ?
        WHERE id = ?
    `

	task.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.UpdatedAt,
		task.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, taskID string) error {
	query := `DELETE FROM todos WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %s not found", taskID)
	}

	return nil
}

func (r *TodoRepository) GetTodosWithCursor(
	ctx context.Context,
	filter models.TodoFilter,
	pagination models.CursorPagination,
	direction string,
) (*models.PaginatedTodos, error) {
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}
	if pagination.OrderBy == "" {
		pagination.OrderBy = "created_at"
	}
	if pagination.OrderDir == "" {
		pagination.OrderDir = "DESC"
	}

	whereClause, args := buildWhereClause(filter)

	orderClause := r.generateOrderClause(pagination.OrderBy, pagination.OrderDir)

	if direction == "prev" || filter.Search != "" {
		orderClause = r.reverseOrderClause(pagination.OrderBy, pagination.OrderDir)
	}

	cursorClause := ""
	var cursorArgs []interface{}

	if pagination.Cursor != "" {
		cursor, err := paginationPkg.DecodeCursor(pagination.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %v", err)
		}

		// NOTE: Search should reverse the orderBy=created_at to search all the newly created tasks
		if direction == "next" && filter.Search == "" {
			cursorClause, cursorArgs = r.buildCursorCondition(cursor, pagination.OrderBy, pagination.OrderDir, false)
		} else {
			cursorClause, cursorArgs = r.buildCursorCondition(cursor, pagination.OrderBy, pagination.OrderDir, true)
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

	args = append(args, cursorArgs...)
	args = append(args, pagination.Limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
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

	totalCount, err := r.getTotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	var nextCursor, prevCursor string
	hasNext := false
	hasPrev := false

	// For forward pagination
	if direction == "next" {
		if len(tasks) > pagination.Limit {
			hasNext = true
			tasks = tasks[:pagination.Limit]

			lastTodo := tasks[len(tasks)-1]
			c := paginationPkg.Cursor{
				Id:        lastTodo.Id,
				CreatedAt: lastTodo.CreatedAt,
				OrderBy:   pagination.OrderBy,
				OrderDir:  pagination.OrderDir,
			}
			nextCursor, _ = paginationPkg.EncodeCursor(c)
		}

		// Check if we have previous items
		if pagination.Cursor != "" {
			hasPrev = true
			// To get prev cursor, we need to go to the first item
			if len(tasks) > 0 {
				firstTodo := tasks[0]
				prevC := paginationPkg.Cursor{
					Id:        firstTodo.Id,
					CreatedAt: firstTodo.CreatedAt,
					OrderBy:   pagination.OrderBy,
					OrderDir:  pagination.OrderDir,
				}
				prevCursor, _ = paginationPkg.EncodeCursor(prevC)
			}
		}
	} else {
		// For backward pagination
		if len(tasks) > pagination.Limit {
			hasPrev = true
			tasks = tasks[:pagination.Limit]
		}

		// Reverse the list for natural order
		for i, j := 0, len(tasks)-1; i < j; i, j = i+1, j-1 {
			tasks[i], tasks[j] = tasks[j], tasks[i]
		}

		// Set next cursor
		if len(tasks) > 0 {
			lastTodo := tasks[len(tasks)-1]
			c := paginationPkg.Cursor{
				Id:        lastTodo.Id,
				CreatedAt: lastTodo.CreatedAt,
				OrderBy:   pagination.OrderBy,
				OrderDir:  pagination.OrderDir,
			}
			nextCursor, _ = paginationPkg.EncodeCursor(c)
		}

		// Set previous cursor if we have more items
		if len(tasks) == pagination.Limit && hasPrev {
			firstTodo := tasks[0]
			c := paginationPkg.Cursor{
				Id:        firstTodo.Id,
				CreatedAt: firstTodo.CreatedAt,
				OrderBy:   pagination.OrderBy,
				OrderDir:  pagination.OrderDir,
			}
			prevCursor, _ = paginationPkg.EncodeCursor(c)
		}
		hasNext = true // Always has next when going back
	}

	return &models.PaginatedTodos{
		Todos:      tasks,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalCount: totalCount,
	}, nil

}

func buildWhereClause(filter models.TodoFilter) (string, []interface{}) {
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

func (r *TodoRepository) getTotalCount(ctx context.Context, filter models.TodoFilter) (int, error) {
	whereClause, args := buildWhereClause(filter)

	query := fmt.Sprintf("SELECT COUNT(*) FROM todos WHERE 1=1 %s", whereClause)

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *TodoRepository) buildCursorCondition(c paginationPkg.Cursor, orderBy, orderDir string, isReverse bool) (string, []interface{}) {
	comparator := "<"
	reverseComparator := ">"

	if orderDir == "ASC" && !isReverse {
		comparator = ">"
		reverseComparator = "<"
	} else if isReverse {
		// Swap comparators for reverse direction
		comparator, reverseComparator = reverseComparator, comparator
	}

	if orderBy == "priority" {
		priorityValue := map[string]int{"high": 1, "medium": 2, "low": 3}
		val := priorityValue[c.OrderBy]
		return fmt.Sprintf(
			"AND ((CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END) %s ? OR "+
				"((CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END) = ? AND created_at %s ?) OR "+
				"((CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END) = ? AND created_at = ? AND id %s ?))",
			comparator, comparator, comparator,
		), []interface{}{val, val, c.CreatedAt, val, c.CreatedAt, c.Id}
	}

	if orderBy == "status" {
		statusValue := map[string]int{"completed": 1, "in_progress": 2, "pending": 3}
		val := statusValue[c.OrderBy]
		return fmt.Sprintf(
			"AND ((CASE status WHEN 'completed' THEN 1 WHEN 'in_progress' THEN 2 WHEN 'pending' THEN 3 END) %s ? OR "+
				"((CASE status WHEN 'completed' THEN 1 WHEN 'in_progress' THEN 2 WHEN 'pending' THEN 3 END) = ? AND created_at %s ?) OR "+
				"((CASE status WHEN 'completed' THEN 1 WHEN 'in_progress' THEN 2 WHEN 'pending' THEN 3 END) = ? AND created_at = ? AND id %s ?))",
			comparator, comparator, comparator,
		), []interface{}{val, val, c.CreatedAt, val, c.CreatedAt, c.Id}
	}

	return fmt.Sprintf(
		"AND (%s %s ? OR (%s = ? AND id %s ?))",
		c.OrderBy, comparator, c.OrderBy, comparator,
	), []interface{}{c.CreatedAt, c.CreatedAt, c.Id}
}

func (r *TodoRepository) reverseOrderClause(orderBy, orderDir string) string {
	reverseDir := "ASC"
	if orderDir == "ASC" {
		reverseDir = "DESC"
	}
	return r.generateOrderClause(orderBy, reverseDir)
}

func (r *TodoRepository) generateOrderClause(orderBy, orderDir string) string {
	if orderBy == "" {
		orderBy = "created_at"
	}

	if orderDir == "" {
		orderDir = "DESC"
	}

	validOrderBy := map[string]bool{
		"id": true, "created_at": true, "due_date": true,
		"completed_at": true, "priority": true,
	}

	validOrderDir := map[string]bool{"ASC": true, "DESC": true}

	if !validOrderBy[orderBy] {
		orderBy = "created_at"
	}

	if !validOrderDir[orderDir] {
		orderDir = "DESC"
	}

	if orderBy == "priority" {
		return fmt.Sprintf("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END %s, created_at %s", orderDir, orderDir)
	}

	return fmt.Sprintf("%s %s, id %s", orderBy, orderDir, orderDir)
}
