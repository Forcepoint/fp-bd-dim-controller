package structs

import (
	"fp-dynamic-elements-manager-controller/internal/health/structs"
	structs2 "fp-dynamic-elements-manager-controller/internal/queue/structs"
	"time"
)

type ModuleType string

const (
	EGRESS     ModuleType = "egress"
	INGRESS    ModuleType = "ingress"
	FUNCTIONAL ModuleType = "functional"
)

type ModuleMetadata struct {
	ID                   int64                `json:"id"`
	CreatedAt            time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time           `json:"deleted_at" db:"deleted_at"`
	ModuleServiceName    string               `json:"module_service_name" db:"module_service_name"`
	ModuleDisplayName    string               `json:"module_display_name" db:"module_display_name"`
	ModuleType           ModuleType           `json:"module_type" db:"module_type"`
	ModuleDescription    string               `json:"module_description" db:"module_description"`
	InboundRoute         string               `json:"inbound_route" db:"inbound_route"`
	InternalIP           string               `json:"internal_ip" db:"internal_ip"`
	InternalPort         string               `json:"internal_port" db:"internal_port"`
	IconURL              string               `json:"icon_url" db:"icon_url"`
	Configured           bool                 `json:"configured"`
	Configurable         bool                 `json:"configurable"`
	LastPing             time.Time            `json:"last_ping" db:"last_ping"`
	AcceptedElementTypes ElementTypesWrapper  `json:"accepted_element_types" db:"-"`
	ModuleEndpoints      []ModuleEndpoint     `json:"internal_endpoints" db:"-"`
	ModuleHealth         structs.ModuleHealth `json:"module_health" db:"-"`
}

type ModuleEndpoint struct {
	ID               int64      `json:"id" db:"id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at" db:"deleted_at"`
	Secure           bool       `json:"secure"`
	Endpoint         string     `json:"endpoint"`
	ModuleMetadataId uint       `json:"module_metadata_id" db:"module_metadata_id"`
}

type ElementTypesWrapper struct {
	ElementTypes []structs2.ElementType `json:"element_types"`
}

type ModuleElementType struct {
	ID          int64                `json:"id" db:"id"`
	CreatedAt   time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time           `json:"deleted_at" db:"deleted_at"`
	ElementType structs2.ElementType `json:"secure" db:"element_type"`
	ModuleId    int64                `json:"module_id" db:"module_id"`
}

type RunningModules struct {
	Containers []string `json:"containers"`
}
