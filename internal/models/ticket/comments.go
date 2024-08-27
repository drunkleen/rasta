package ticketmodel

import (
	"time"

	"github.com/google/uuid"
)

type TicketComment struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	TicketId  uuid.UUID `json:"ticket_id" gorm:"not null"`
	UserId    uuid.UUID `json:"user_id" gorm:"not null"`
	Comment   string    `json:"comment" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}
