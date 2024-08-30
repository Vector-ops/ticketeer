package models

import (
	"context"
	"time"
)

type Ticket struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	EventID   uint      `json:"eventId"`
	UserID    uint      `json:"userId" gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Event     Event     `json:"event" gorm:"foreignkey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Entered   bool      `json:"entered" default:"false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TicketRepository interface {
	GetMany(ctx context.Context, userId uint) ([]*Ticket, error)
	GetOne(ctx context.Context, ticketId uint, userId uint) (*Ticket, error)
	CreateOne(ctx context.Context, ticket *Ticket, userId uint) (*Ticket, error)
	UpdateOne(ctx context.Context, ticketId uint, updateData map[string]interface{}, userId uint) (*Ticket, error)
	DeleteOne(ctx context.Context, ticketId uint, userId uint) error
}

type ValidateTicket struct {
	TicketId uint `json:"ticketId"`
	UserId   uint `json:"userId"`
}
