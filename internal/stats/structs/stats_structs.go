package structs

type Stats struct {
	Total             int    `json:"num_sources" db:"total"`
	NumBlockedIps     int    `json:"num_blocked_ips" db:"ip"`
	NumBlockedDomains int    `json:"num_blocked_domains" db:"domain"`
	NumBlockedUrls    int    `json:"num_blocked_urls" db:"url"`
	NumBlockedRanges  int    `json:"num_blocked_ranges" db:"ip_range"`
	NumBlockedSnorts  int    `json:"num_blocked_snorts" db:"snort"`
	LastUpdate        string `json:"last_update"`
}
