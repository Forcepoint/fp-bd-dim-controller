package structs

type Stats struct {
	Total             int    `json:"num_sources" db:"total"`
	NumBlockedIps     int    `json:"num_blocked_ips" db:"ip"`
	NumBlockedDomains int    `json:"num_blocked_domains" db:"domain"`
	NumBlockedUrls    int    `json:"num_blocked_urls" db:"url"`
	LastUpdate        string `json:"last_update"`
}
