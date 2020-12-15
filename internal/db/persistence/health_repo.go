package persistence

import "github.com/jmoiron/sqlx"

type HealthRepo struct {
	db *sqlx.DB
}

func NewHealthRepo(appDb *sqlx.DB) *HealthRepo {
	return &HealthRepo{db: appDb}
}

func (h *HealthRepo) PingDB() bool {
	if err := h.db.Ping(); err != nil {
		return false
	}
	return true
}
