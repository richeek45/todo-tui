package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Cursor struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	OrderBy   string    `json:"order_by"`
	OrderDir  string    `json:"order_dir"`
}

func EncodeCursor(cursor Cursor) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

func DecodeCursor(encoded string) (Cursor, error) {
	var cursor Cursor

	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return cursor, err
	}

	err = json.Unmarshal(data, &cursor)
	return cursor, err
}

func GenerateOrderClause(orderBy, orderDir string) string {
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
