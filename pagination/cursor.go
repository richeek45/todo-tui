package pagination

import (
	"encoding/base64"
	"encoding/json"
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
