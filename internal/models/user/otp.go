package usermodel

import (
	"github.com/google/uuid"
	"time"
)

type OtpEmail struct {
	UserId uuid.UUID `json:"user_id,omitempty" gorm:"not null;unique"`
	Code   string    `json:"otp_code,omitempty" gorm:"not null"`
	Expiry time.Time `json:"otp_expiry,omitempty" gorm:"type:timestamp with time zone"`
}
