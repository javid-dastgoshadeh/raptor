package database

import "time"

// ID , Database ID
type ID uint

// ModelWithSoftDelete ...
type ModelWithSoftDelete struct {
	ID        ID         `json:"id" gorm:"primary_key" `
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time ` json:"deleted_at" sql:"index"`
}

// Model a database base model
type Model struct {
	ID        ID        `json:"id" gorm:"primary_key" `
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
