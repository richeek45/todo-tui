package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// db, err := sql.Open("sqlite3", "./test.db")
func NewDB(datasourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", datasourceName)

	if err != nil {
		return nil, err
	}
	defer db.Close()

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
