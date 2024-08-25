package newsletterservice

import (
	newslettermodel "github.com/drunkleen/rasta/internal/models/newsletter"
	"github.com/drunkleen/rasta/internal/repository/newsletter"
	emailPkg "github.com/drunkleen/rasta/pkg/email"
)

type NewsletterService struct {
	Repository *newsletterrepository.NewsletterRepository
}

func NewNewsletterService(repository *newsletterrepository.NewsletterRepository) *NewsletterService {
	return &NewsletterService{Repository: repository}
}

func (s *NewsletterService) Create(email *string) error {
	return s.Repository.Create(email)
}

func (s *NewsletterService) DeleteByEmail(email *string) error {
	return s.Repository.Delete(email)
}

func (s *NewsletterService) FindByEmail(email *string) (*newslettermodel.Newsletter, error) {
	return s.Repository.FindByEmail(email)
}

func (s *NewsletterService) UpdateActiveStatus(email *string, isActive bool) error {
	return s.Repository.UpdateActiveStatus(email, &isActive)
}

func (s *NewsletterService) FindAllActive() ([]newslettermodel.Newsletter, error) {
	return s.Repository.FindAll(true)
}

func (s *NewsletterService) FindAllInactive() ([]newslettermodel.Newsletter, error) {
	return s.Repository.FindAll(false)
}

func (s *NewsletterService) CountActiveSubscribers() (int64, error) {
	return s.Repository.CountSubscribers(true)
}

func (s *NewsletterService) CountInactiveSubscribers() (int64, error) {
	return s.Repository.CountSubscribers(false)
}

func (s *NewsletterService) SendNewslettersEmail(emailMessage *string, limit int) error {
	count64, err := s.Repository.CountSubscribers(true)
	if err != nil {
		return err
	}
	if count64 == 0 {
		return nil
	}
	count := int(count64)
	pages := count / limit
	if count%limit != 0 {
		pages++
	}
	for start := 0; start < pages; start++ {
		var newsletters *[]newslettermodel.Newsletter
		newsletters, err = s.Repository.GetLimited(start, limit)
		if err != nil {
			return err
		}
		err = emailPkg.SendNewsletter(newsletters, emailMessage)
		if err != nil {
			return err
		}
	}
	return nil
}
