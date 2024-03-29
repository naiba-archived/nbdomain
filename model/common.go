package model

import "time"

// Common ..
type Common struct {
	ID        uint64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
