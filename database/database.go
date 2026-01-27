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
		{
			Id:          "4",
			Title:       "Prepare quarterly presentation",
			Description: "Create slides for Q2 results presentation to stakeholders",
			Status:      models.InProgress,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(3 * 24 * time.Hour)),
		},
		{
			Id:          "5",
			Title:       "Fix authentication bug",
			Description: "Resolve JWT token expiration issue in mobile app",
			Status:      models.NotStarted,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(2 * 24 * time.Hour)),
		},
		{
			Id:          "6",
			Title:       "Onboard new team member",
			Description: "Set up workstation and provide orientation for new developer",
			Status:      models.Completed,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(-1 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now().Add(-6 * time.Hour)),
		},
		{
			Id:          "7",
			Title:       "Research new framework",
			Description: "Evaluate React Native vs Flutter for upcoming mobile project",
			Status:      models.InProgress,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(14 * 24 * time.Hour)),
		},
		{
			Id:          "8",
			Title:       "Update SSL certificates",
			Description: "Renew and deploy SSL certificates for all production domains",
			Status:      models.NotStarted,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(5 * 24 * time.Hour)),
		},
		{
			Id:          "9",
			Title:       "Performance optimization",
			Description: "Improve database query performance for user analytics dashboard",
			Status:      models.Completed,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(-3 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now().Add(-2 * 24 * time.Hour)),
		},
		{
			Id:          "10",
			Title:       "Plan team offsite",
			Description: "Organize venue, agenda, and logistics for quarterly team offsite",
			Status:      models.NotStarted,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(30 * 24 * time.Hour)),
		},
	}

	for _, todo := range sampleTodos {
		repo.CreateTodo(context.Background(), todo)
	}

	fmt.Println("Seeded sample data")
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
