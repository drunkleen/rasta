package database

import (
	"github.com/drunkleen/rasta/config"
	newslettermodel "github.com/drunkleen/rasta/internal/models/newsletter"
	"github.com/drunkleen/rasta/internal/models/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	DB *gorm.DB
)

// InitDB initializes the database connection using the database string
// obtained from the configuration.  It also creates the tables for the
// models defined in the `models` package.
func InitDB() {
	dbString := config.GetDBString()

	var err error
	DB, err = gorm.Open(postgres.Open(dbString), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to database!")
	}
	if err = createTables(); err != nil {
		log.Panic("could not create tables")
	}
}

// createTables creates the necessary tables in the database.
//
// It uses the AutoMigrate function from the GORM library to create tables for the User, OAuth, and OtpEmail models.
// error - Returns an error if any of the table creations fail.
func createTables() error {
	if err := DB.AutoMigrate(&usermodel.User{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&usermodel.OAuth{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&usermodel.OtpEmail{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&usermodel.ResetPwd{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&newslettermodel.Newsletter{}); err != nil {
		return err
	}

	return nil
}
