package structs

type HealthStatus int

const (
	Down HealthStatus = iota - 1
	Unhealthy
	Healthy
)

type Health struct {
	Modules []ModuleHealth `json:"modules"`
}
type ModuleHealth struct {
	ModuleName string       `json:"module_name"`
	Status     HealthStatus `json:"status"`
	StatusCode int          `json:"status_code"`
	LastUpdate string       `json:"last_update"`
}
