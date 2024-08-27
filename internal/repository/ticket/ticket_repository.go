package ticketrepository

import (
	"errors"
	ticketmodel "github.com/drunkleen/rasta/internal/models/ticket"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

func (r *TicketRepository) Create(ticket *ticketmodel.Ticket) error {
	if ticket.Title == "" {
		return errors.New("title is required")
	}
	if err := r.DB.Create(ticket).Error; err != nil {
		log.Printf("failed to create ticket: %v", err)
		return errors.New("could not create ticket")
	}
	return nil
}

func (r *TicketRepository) Delete(ticketID uuid.UUID) error {
	err := r.DB.Where("id = ?", ticketID).Delete(&ticketmodel.Ticket{}).Error
	if err != nil {
		log.Printf("failed to delete ticket: %v", err)
		return errors.New("could not delete ticket")
	}
	return nil
}

func (r *TicketRepository) FindById(ticketID uuid.UUID) (*ticketmodel.Ticket, error) {
	var ticket ticketmodel.Ticket
	err := r.DB.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		log.Printf("ticket not found: %v", err)
		return nil, errors.New("ticket not found")
	}
	return &ticket, nil
}

func (r *TicketRepository) FindByUserId(userID uuid.UUID) ([]ticketmodel.Ticket, error) {
	var tickets []ticketmodel.Ticket
	err := r.DB.Where("user_id = ?", userID).Find(&tickets).Error
	if err != nil {
		log.Printf("no tickets found for user: %v", err)
		return nil, errors.New("no tickets found")
	}
	return tickets, nil
}

func (r *TicketRepository) UpdateStatus(ticketID uuid.UUID, status ticketmodel.TicketStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	err := r.DB.Model(&ticketmodel.Ticket{}).Where("id = ?", ticketID).Updates(updates).Error
	if err != nil {
		log.Printf("failed to update ticket status: %v", err)
		return errors.New("could not update ticket status")
	}
	return nil
}

func (r *TicketRepository) UpdatePriority(ticketID uuid.UUID, priority ticketmodel.TicketPriority) error {
	updates := map[string]interface{}{
		"priority":   priority,
		"updated_at": time.Now(),
	}
	err := r.DB.Model(&ticketmodel.Ticket{}).Where("id = ?", ticketID).Updates(updates).Error
	if err != nil {
		log.Printf("failed to update ticket priority: %v", err)
		return errors.New("could not update ticket priority")
	}
	return nil
}

func (r *TicketRepository) AddComment(comment *ticketmodel.TicketComment) error {
	if comment.Comment == "" {
		return errors.New("comment is required")
	}
	if err := r.DB.Create(comment).Error; err != nil {
		log.Printf("failed to add comment: %v", err)
		return errors.New("could not add comment")
	}
	return nil
}

func (r *TicketRepository) GetComments(ticketID uuid.UUID) ([]ticketmodel.TicketComment, error) {
	var comments []ticketmodel.TicketComment
	err := r.DB.Where("ticket_id = ?", ticketID).Order("created_at asc").Find(&comments).Error
	if err != nil {
		log.Printf("no comments found for ticket: %v", err)
		return nil, errors.New("no comments found")
	}
	return comments, nil
}

func (r *TicketRepository) FindAll() ([]ticketmodel.Ticket, error) {
	var tickets []ticketmodel.Ticket
	err := r.DB.Find(&tickets).Error
	if err != nil {
		log.Printf("no tickets found: %v", err)
		return nil, errors.New("no tickets found")
	}
	return tickets, nil
}
