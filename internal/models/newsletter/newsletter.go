package newslettermodel

import (
	"time"
)

type Newsletter struct {
	Id        uint      `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Email     string    `json:"email" gorm:"not null;unique"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}
