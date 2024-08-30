package repository

import (
	"context"

	"github.com/vector-ops/ticketeer/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) models.TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) GetMany(ctx context.Context, userId uint) ([]*models.Ticket, error) {
	tickets := []*models.Ticket{}

	res := r.db.Model(&models.Ticket{}).Where("user_id = ?", userId).Preload("Event").Order("updated_at desc").Find(&tickets)

	if res.Error != nil {
		return nil, res.Error
	}

	return tickets, nil
}

func (r *TicketRepository) GetOne(ctx context.Context, ticketId uint, userId uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}

	res := r.db.Model(ticket).Where("id = ?", ticketId).Where("user_id = ?", userId).Preload("Event").First(ticket)
	if res.Error != nil {
		return nil, res.Error
	}

	return ticket, nil
}

func (r *TicketRepository) CreateOne(ctx context.Context, ticket *models.Ticket, userId uint) (*models.Ticket, error) {
	ticket.UserID = userId

	res := r.db.Model(ticket).Create(ticket)
	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetOne(ctx, ticket.ID, userId)
}

func (r *TicketRepository) UpdateOne(ctx context.Context, ticketId uint, updatedData map[string]interface{}, userId uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}

	updateRes := r.db.Model(ticket).Where("id = ?", ticketId).Where("user_id = ?", userId).Updates(updatedData)
	if updateRes.Error != nil {
		return nil, updateRes.Error
	}

	return r.GetOne(ctx, ticketId, userId)

}

func (r *TicketRepository) DeleteOne(ctx context.Context, ticketId uint, userId uint) error {
	res := r.db.Delete(models.Ticket{}, ticketId).Where("user_id = ?", userId)
	return res.Error
}
