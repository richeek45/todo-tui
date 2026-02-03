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
		{
			Id:          "11",
			Title:       "Refactor authentication module",
			Description: "Clean up and optimize the OAuth2 implementation for better security",
			Status:      models.InProgress,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(3 * 24 * time.Hour)),
		},
		{
			Id:          "12",
			Title:       "Write API documentation",
			Description: "Create comprehensive Swagger/OpenAPI docs for all endpoints",
			Status:      models.NotStarted,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(7 * 24 * time.Hour)),
		},
		{
			Id:          "13",
			Title:       "Setup monitoring alerts",
			Description: "Configure Prometheus alerts for critical services and error rates",
			Status:      models.Completed,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(-5 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now().Add(-4 * 24 * time.Hour)),
		},
		{
			Id:          "14",
			Title:       "Code review backlog",
			Description: "Clear pending PR reviews from the last sprint",
			Status:      models.InProgress,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(2 * 24 * time.Hour)),
		},
		{
			Id:          "15",
			Title:       "Update dependencies",
			Description: "Upgrade Node.js packages to latest LTS versions with security patches",
			Status:      models.NotStarted,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(10 * 24 * time.Hour)),
		},
		{
			Id:          "16",
			Title:       "Implement caching layer",
			Description: "Add Redis caching for frequently accessed user data",
			Status:      models.InProgress,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(5 * 24 * time.Hour)),
		},
		{
			Id:          "17",
			Title:       "User acceptance testing",
			Description: "Coordinate UAT session with product team for new checkout flow",
			Status:      models.NotStarted,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(4 * 24 * time.Hour)),
		},
		{
			Id:          "18",
			Title:       "Fix payment gateway integration",
			Description: "Resolve intermittent failures with Stripe webhook processing",
			Status:      models.InProgress,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(1 * 24 * time.Hour)),
		},
		{
			Id:          "19",
			Title:       "Optimize database indexes",
			Description: "Analyze query performance and add missing indexes to speed up reports",
			Status:      models.Completed,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(-2 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now().Add(-1 * 24 * time.Hour)),
		},
		{
			Id:          "20",
			Title:       "Prepare quarterly report",
			Description: "Compile metrics and achievements for stakeholder presentation",
			Status:      models.NotStarted,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(14 * 24 * time.Hour)),
		},
		{
			Id:          "21",
			Title:       "Setup CI/CD pipeline",
			Description: "Configure GitHub Actions for automated testing and deployment",
			Status:      models.InProgress,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(6 * 24 * time.Hour)),
		},
		{
			Id:          "22",
			Title:       "Implement dark mode",
			Description: "Add theme switching functionality to frontend application",
			Status:      models.NotStarted,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(21 * 24 * time.Hour)),
		},
		{
			Id:          "23",
			Title:       "Security audit",
			Description: "Conduct penetration testing and vulnerability assessment",
			Status:      models.InProgress,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(3 * 24 * time.Hour)),
		},
		{
			Id:          "24",
			Title:       "Migrate legacy data",
			Description: "Transfer user records from old system to new database schema",
			Status:      models.NotStarted,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(12 * 24 * time.Hour)),
		},
		{
			Id:          "25",
			Title:       "Create onboarding tutorial",
			Description: "Build interactive guide for new users",
			Status:      models.Completed,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(-7 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now().Add(-5 * 24 * time.Hour)),
		},
		{
			Id:          "26",
			Title:       "Add export functionality",
			Description: "Implement CSV/Excel export for reports dashboard",
			Status:      models.InProgress,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(8 * 24 * time.Hour)),
		},
		{
			Id:          "27",
			Title:       "Monitor server logs",
			Description: "Setup ELK stack for centralized logging and analysis",
			Status:      models.NotStarted,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(5 * 24 * time.Hour)),
		},
		{
			Id:          "28",
			Title:       "Refactor CSS architecture",
			Description: "Convert to CSS-in-JS and improve component styling system",
			Status:      models.InProgress,
			Priority:    models.PriorityLow,
			DueDate:     ptrTime(time.Now().Add(15 * 24 * time.Hour)),
		},
		{
			Id:          "29",
			Title:       "Implement WebSocket",
			Description: "Add real-time notifications for user activities",
			Status:      models.NotStarted,
			Priority:    models.PriorityMedium,
			DueDate:     ptrTime(time.Now().Add(9 * 24 * time.Hour)),
		},
		{
			Id:          "30",
			Title:       "Performance testing",
			Description: "Run load tests on API endpoints before production release",
			Status:      models.Completed,
			Priority:    models.PriorityHigh,
			DueDate:     ptrTime(time.Now().Add(-1 * 24 * time.Hour)),
			CompletedAt: ptrTime(time.Now()),
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
