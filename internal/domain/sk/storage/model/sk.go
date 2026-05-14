package model

import "time"

// Sk json нужен для сохранения истории
type Sk struct {
	Id        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
