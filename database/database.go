package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/richeek45/todo-tui/models"
)

// db, err := sql.Open("sqlite3", "./test.db")
func NewDB(datasourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", datasourceName)

	if err != nil {
		return nil, err
	}
	// defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		priority TEXT NOT NULL DEFAULT 'medium',
		due_date DATETIME,
		completed_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CHECK (status IN ('pending', 'in_progress', 'completed')),
		CHECK (priority IN ('low', 'medium', 'high'))
	);
	
	CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
	CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
	CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
	CREATE INDEX IF NOT EXISTS idx_todos_completed_at ON todos(completed_at);
	CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);
	`

	_, err := db.Exec(migrationSQL)
	return err

}

func SeedSampleData(db *sql.DB) {
	repo := NewTodoRepository(db)

	sampleTodos := []*models.Task{
		{
			Id:          "1",
			Title:       "Complete project proposal",
			Description: "Write and submit the Q2 project proposal",
			Status:      models.Completed,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(-2 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now().Add(-1 * 24 * time.Hour)),
		},
		{
			Id:          "2",
			Title:       "Review team PRs",
			Description: "Review open pull requests from the team",
			Status:      models.InProgress,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(1 * 24 * time.Hour)),
		},
		{
			Id:          "3",
			Title:       "Update documentation",
			Description: "Update API documentation for new features",
			Status:      models.NotStarted,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(7 * 24 * time.Hour)),
		},
		// Add more sample todos...

		// task := models.Task{
		// 	Id:          "1",
		// 	Title:       "First Task",
		// 	Priority:    models.PriorityHigh,
		// 	Description: "This is going to be a long title for the thing that I am going to talk about. Nothing can change that",
		// 	Status:      models.NotStarted,
		// 	Error:       nil,
		// }
		// task2 := models.Task{
		// 	Id:          "2",
		// 	Title:       "Second Task",
		// 	Priority:    models.PriorityMedium,
		// 	Description: "How life has changed since the time I have first started doing something that means some other thing. Well, no",
		// 	Status:      models.Completed,
		// 	Error:       nil,
		// }
		// task3 := models.Task{
		// 	Id:          "3",
		// 	Title:       "Third Taskj",
		// 	Priority:    models.PriorityMedium,
		// 	Description: "Wow, still this is working. I cannot believe it. The brain can spew nonsense if we keep prompting it to provide something",
		// 	Status:      models.InProgress,
		// 	Error:       nil,
		// }
	}

	for _, todo := range sampleTodos {
		repo.CreateTodo(context.Background(), todo)
	}

	fmt.Println("Seeded sample data")
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
