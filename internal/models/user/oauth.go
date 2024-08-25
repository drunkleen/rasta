package usermodel

import "github.com/google/uuid"

type OAuth struct {
	UserId  uuid.UUID `json:"user_id,omitempty" gorm:"not null;unique"`
	Enabled bool      `json:"oauth_enabled,omitempty" gorm:"default:false"`
	Secret  string    `json:"oauth_secret,omitempty" gorm:"size:512"`
}
