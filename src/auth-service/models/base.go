package models

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	CreatedAt time.Time       `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index column:deleted_at; default:null" swaggertype:"string"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = time.Now()
	return
}
