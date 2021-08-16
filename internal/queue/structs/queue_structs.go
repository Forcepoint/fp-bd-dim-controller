package structs

import (
	"time"
)

type Status string
type UpdateType string
type ElementType string

const (
	SUCCESS    Status = "success"
	PENDING    Status = "pending"
	FAILED     Status = "failed"
	INCOMPLETE Status = "incomplete"
	NEW        Status = "new"

	ADD    UpdateType = "add"
	DELETE UpdateType = "delete"

	IP     ElementType = "IP"
	DOMAIN ElementType = "DOMAIN"
	URL    ElementType = "URL"
	RANGE  ElementType = "RANGE"
	SNORT  ElementType = "SNORT"
)

func (s Status) String() string {
	return string(s)
}

type ProcessedItems struct {
	UpdateType UpdateType    `json:"update_type"`
	SafeList   bool          `json:"safe_list"`
	Items      []ListElement `json:"items"`
	Item       ListElement   `json:"item"`
	BatchId    int64         `json:"batch_id"`
}

type PaginatedElements struct {
	Elements       []ListElement `json:"elements"`
	TotalPageCount int           `json:"total_page_count"`
	PageNumber     int           `json:"page_number"`
}

type ElementBatch struct {
	ID             int64          `json:"id" db:"id"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at" db:"deleted_at"`
	ProcessedItems []ListElement  `json:"processed_items"`
	UpdateStatus   []UpdateStatus `json:"update_status"`
}

type ListElement struct {
	ID            int64       `json:"id" db:"id"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time  `json:"deleted_at" db:"deleted_at"`
	Source        string      `json:"source" db:"source"`
	ServiceName   string      `json:"service_name" db:"service_name"`
	Type          ElementType `json:"type"`
	Value         string      `json:"value"`
	Safe          bool        `json:"safe"`
	UpdateBatchId int64       `json:"batch_number" db:"update_batch_id"`
}

type UpdateStatus struct {
	ID               int64      `json:"id" db:"id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
	ServiceName      string     `json:"service_name" db:"service_name"`
	Status           Status     `json:"status"`
	UpdateType       UpdateType `json:"update_type" db:"update_type"`
	UpdateBatchId    int64      `json:"update_batch_id" db:"update_batch_id"`
	ModuleMetadataId int64      `json:"module_metadata_id" db:"module_metadata_id"`
}

type PaginatedStatus struct {
	Statuses       []UpdateStatus `json:"items"`
	TotalPageCount int            `json:"total_page_count"`
	PageNumber     int            `json:"page_number"`
}
