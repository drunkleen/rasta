package usermodel

import (
	"time"

	"github.com/google/uuid"
)

// RegionType represents a valid region value.
type RegionType string

// Constants representing various regions.
const (
	RegionTypeNorthernAmerica RegionType = "Northern America"
	RegionTypeCentralAmerica  RegionType = "Central America"
	RegionTypeCaribbean       RegionType = "Caribbean"

	NorthernSouthAmerica RegionType = "Northern South America"
	SouthernSouthAmerica RegionType = "Southern South America"
	WesternSouthAmerica  RegionType = "Western South America"
	EasternSouthAmerica  RegionType = "Eastern South America"

	RegionTypeScandinavia    RegionType = "Scandinavia"
	RegionTypeSouthernEurope RegionType = "Southern Europe"
	RegionTypeWesternEurope  RegionType = "Western Europe"
	RegionTypeEasternEurope  RegionType = "Eastern Europe"
	RegionTypeCentralEurope  RegionType = "Central Europe"

	RegionTypeMiddleEast       RegionType = "Middle East"
	RegionTypeCentralAsia      RegionType = "Central Asia"
	RegionTypeEasternAsia      RegionType = "Eastern Asia"
	RegionTypeSouthernAsia     RegionType = "Southern Asia"
	RegionTypeSoutheasternAsia RegionType = "Southeastern Asia"
	RegionTypeSiberia          RegionType = "Siberia"

	RegionTypeNorthernAfrica RegionType = "Northern Africa"
	RegionTypeWesternAfrica  RegionType = "Western Africa"
	RegionTypeCentralAfrica  RegionType = "Central Africa"
	RegionTypeHornOfAfrica   RegionType = "Horn of Africa"
	RegionTypeSouthernAfrica RegionType = "Southern Africa"

	AustraliaAndNewZealand RegionType = "Australia and New Zealand"
	Melanesia              RegionType = "Melanesia"
	Micronesia             RegionType = "Micronesia"
	Polynesia              RegionType = "Polynesia"
)

type AccountType string

const (
	AccountTypeNormal AccountType = "User"
	AccountTypeSeller AccountType = "Seller"
	AccountTypeAdmin  AccountType = "Admin"
)

type User struct {
	Id         uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	FirstName  string      `json:"first_name" gorm:"size:64;not null"`
	LastName   string      `json:"last_name" gorm:"size:64;not null"`
	Username   string      `json:"username" gorm:"size:64;unique;not null"`
	Email      string      `json:"email" gorm:"size:128;unique;not null"`
	Password   string      `json:"password" gorm:"size:256;not null"`
	IsVerified bool        `json:"is_verified" gorm:"default:false"`
	IsDisabled bool        `json:"is_disabled" gorm:"default:false"`
	Account    AccountType `json:"account" gorm:"size:32;not null;default:'User'"`
	Region     RegionType  `json:"region" gorm:"size:32;not null"`
	CreatedAt  time.Time   `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt  time.Time   `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`

	OAuth    OAuth    `gorm:"foreignKey:UserId"`
	OtpEmail OtpEmail `gorm:"foreignKey:UserId"`
	ResetPwd ResetPwd `gorm:"foreignKey:UserId"`
}

// TODO Ticketing System
// TODO Gift Cards
// TODO Product Listing
// TODO Orders
