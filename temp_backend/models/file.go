package models

import (
	"time"
)

type File struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	OrderID   *uint     `json:"order_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`