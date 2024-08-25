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

func NewNewsletterRepository(db *gorm.DB) *NewsletterRepository {
	return &NewsletterRepository{DB: db}
}

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

func (r *NewsletterRepository) FindByEmail(email *string) (*newslettermodel.Newsletter, error) {
	var newsletter newslettermodel.Newsletter
	err := r.DB.Where("email = ?", *email).First(&newsletter).Error
	if err != nil {
		log.Printf("newsletter not found: %v", err)
		return nil, errors.New("newsletter not found")
	}
	return &newsletter, nil
}

func (r *NewsletterRepository) FindAll(status bool) ([]newslettermodel.Newsletter, error) {
	var newsletters []newslettermodel.Newsletter
	err := r.DB.Find(&newsletters).Where("is_active = ?", status).Error
	if err != nil {
		log.Printf("no newsletters found: %v", err)
		return nil, errors.New("no newsletters found")
	}
	return newsletters, nil
}

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

func (r *NewsletterRepository) GetLimited(index, limit int) (*[]newslettermodel.Newsletter, error) {
	var newsletters []newslettermodel.Newsletter
	err := r.DB.Offset(index).Limit(limit).Find(&newsletters).Error
	if err != nil {
		log.Printf("no newsletters found: %v", err)
		return nil, errors.New("no newsletters found")
	}
	return &newsletters, nil
}
