package ticketmodel

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
)

type TicketPriority string

const (
	TicketPriorityLow    TicketPriority = "low"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityHigh   TicketPriority = "high"
	TicketPriorityUrgent TicketPriority = "urgent"
)

type TicketCategory string

const (
	TicketCategoryBugReport      TicketCategory = "bug report"
	TicketCategoryFeatureRequest TicketCategory = "feature request"
	TicketCategoryAccountIssue   TicketCategory = "account issue"
	TicketCategoryGeneralQuery   TicketCategory = "general query"
)

type Ticket struct {
	Id          uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey"`
	Title       string          `json:"title" gorm:"size:256;not null"`
	Description string          `json:"description" gorm:"type:text;not null"`
	Status      TicketStatus    `json:"status" gorm:"type:varchar(32);not null;default:'Open'"`
	Priority    TicketPriority  `json:"priority" gorm:"type:varchar(32);not null;default:'Medium'"`
	Category    TicketCategory  `json:"category" gorm:"type:varchar(32);not null"`
	UserId      uuid.UUID       `json:"user_id" gorm:"not null"`
	AssignedTo  uuid.UUID       `json:"assigned_to" gorm:"type:uuid;default:null"`
	CreatedAt   time.Time       `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	Comments    []TicketComment `json:"comments" gorm:"foreignKey:TicketId"`
}
