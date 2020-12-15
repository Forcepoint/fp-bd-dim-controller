package structs

import (
	"time"
)

type LogEntry struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
	ModuleName string     `json:"module_name" db:"module_name"`
	Level      string     `json:"level"`
	Message    string     `json:"message"`
	Caller     string     `json:"caller"`
	Time       time.Time  `json:"time"`
}

type LogEvents struct {
	Events         []LogEntry `json:"events"`
	TotalPageCount int        `json:"total_page_count"`
	PageNumber     int        `json:"page_number"`
}
