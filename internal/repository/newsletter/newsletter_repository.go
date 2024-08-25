package newsletterrepository

import (
	"errors"
	newslettermodel "github.com/drunkleen/rasta/internal/models/newsletter"
	"gorm.io/gorm"
	"log"
	"time"
)

type NewsletterRepository struct {
	DB *gorm.DB
}

// NewNewsletterRepository creates a new NewsletterRepository.
//
// It takes a pointer to a gorm.DB as a parameter.
// Returns a pointer to a NewsletterRepository.
func NewNewsletterRepository(db *gorm.DB) *NewsletterRepository {
	return &NewsletterRepository{DB: db}
}

// Create creates a new newsletter in the database.
//
// It takes an email address as a parameter.
// Returns an error if the newsletter could not be created.
// Create creates a new newsletter in the database.
//
// It takes an email address as a parameter.
// Returns an error if the newsletter could not be created.
func (r *NewsletterRepository) Create(email *string) error {
	if *email == "" {
		return errors.New("email is required")
	}
	now := time.Now()
	newsletter := &newslettermodel.Newsletter{
		Email:     *email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.DB.Create(&newsletter).Error; err != nil {
		log.Printf("failed to create newsletter: %v", err)
		return errors.New("could not create newsletter")
	}
	return nil
}

// Delete deletes a newsletter from the database.
//
// It takes an email address as argument and returns an error.
//
// If the email address is empty, it returns an error.
//
// It returns an error if the newsletter could not be deleted.
func (r *NewsletterRepository) Delete(email *string) error {
	if *email == "" {
		return errors.New("email is required")
	}
	err := r.DB.Where("email = ?", *email).Delete(&newslettermodel.Newsletter{}).Error
	if err != nil {
		log.Printf("failed to delete newsletter: %v", err)
		return errors.New("could not delete newsletter")
	}
	return nil
}

// FindByEmail retrieves a newsletter from the database based on the provided email.
//
// The email parameter specifies the email address of the newsletter to retrieve.
// Returns a pointer to a newslettermodel.Newsletter and an error.
// FindByEmail retrieves a newsletter from the database based on the provided email.
//
// The email parameter specifies the email address of the newsletter to retrieve.
// Returns a pointer to a newslettermodel.Newsletter and an error.
func (r *NewsletterRepository) FindByEmail(email *string) (*newslettermodel.Newsletter, error) {
	var newsletter newslettermodel.Newsletter
	err := r.DB.Where("email = ?", *email).First(&newsletter).Error
	if err != nil {
		log.Printf("newsletter not found: %v", err)
		return nil, errors.New("newsletter not found")
	}
	return &newsletter, nil
}

// FindAll retrieves all newsletters from the database based on the provided status.
//
// The status parameter specifies whether to retrieve active or inactive newsletters.
// Returns a slice of newslettermodel.Newsletter and an error.
func (r *NewsletterRepository) FindAll(status bool) ([]newslettermodel.Newsletter, error) {
	var newsletters []newslettermodel.Newsletter
	err := r.DB.Find(&newsletters).Where("is_active = ?", status).Error
	if err != nil {
		log.Printf("no newsletters found: %v", err)
		return nil, errors.New("no newsletters found")
	}
	return newsletters, nil
}

// UpdateActiveStatus updates the newsletter's active status.
//
// It takes an email address as argument and whether the newsletter should be
// active or not.
//
// If the email address is empty, it returns an error.
//
// It returns an error if the newsletter could not be updated.
func (r *NewsletterRepository) UpdateActiveStatus(email *string, isActive *bool) error {
	if *email == "" {
		return errors.New("email is required")
	}
	updates := map[string]interface{}{
		"is_active":  *isActive,
		"updated_at": time.Now(),
	}
	err := r.DB.Model(&newslettermodel.Newsletter{}).Where("email = ?", *email).Updates(updates).Error
	if err != nil {
		log.Printf("failed to update newsletter: %v", err)
		return errors.New("could not update newsletter")
	}
	return nil
}

// CountSubscribers returns the number of subscribers based on their status.
//
// Parameter status: a boolean indicating whether to count active or inactive subscribers.
// Returns int64: the number of subscribers, and error: any error that occurred during the operation.
func (r *NewsletterRepository) CountSubscribers(status bool) (int64, error) {
	var count int64
	err := r.DB.Model(&newslettermodel.Newsletter{}).
		Where("is_active = ?", status).
		Count(&count).Error

	if err != nil {
		log.Printf("failed to count subscribers: %v", err)
		return 0, errors.New("could not count subscribers")
	}
	return count, nil
}

// GetLimited retrieves a limited number of newsletters from the database.
//
// The index parameter specifies the starting point for the query, and the limit parameter specifies the maximum number of newsletters to retrieve.
// Returns a pointer to a slice of newslettermodel.Newsletter and an error.
func (r *NewsletterRepository) GetLimited(index, limit int) (*[]newslettermodel.Newsletter, error) {
	var newsletters []newslettermodel.Newsletter
	err := r.DB.Offset(index).Limit(limit).Find(&newsletters).Error
	if err != nil {
		log.Printf("no newsletters found: %v", err)
		return nil, errors.New("no newsletters found")
	}
	return &newsletters, nil
}
