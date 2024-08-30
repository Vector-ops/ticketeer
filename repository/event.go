package repository

import (
	"context"

	"github.com/vector-ops/ticketeer/models"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) models.EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (er *EventRepository) GetMany(ctx context.Context) ([]*models.Event, error) {
	events := []*models.Event{}

	res := er.db.Model(models.Event{}).Order("updated_at desc").Find(&events)
	if res.Error != nil {
		return nil, res.Error
	}

	return events, nil
}

func (er *EventRepository) GetOne(ctx context.Context, eventId uint) (*models.Event, error) {
	event := &models.Event{}

	res := er.db.Model(event).Where("id = ?", eventId).First(event)
	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (er *EventRepository) CreateOne(ctx context.Context, event *models.Event) (*models.Event, error) {

	res := er.db.Model(event).Create(event)
	if res.Error != nil {
		return nil, res.Error
	}

	return er.GetOne(ctx, event.ID)
}

func (er *EventRepository) UpdateOne(ctx context.Context, eventId uint, updatedData map[string]interface{}) (*models.Event, error) {
	event := &models.Event{}

	updateRes := er.db.Model(event).Where("id = ?", eventId).Updates(updatedData)
	if updateRes.Error != nil {
		return nil, updateRes.Error
	}

	return er.GetOne(ctx, eventId)

}

func (er *EventRepository) DeleteOne(ctx context.Context, eventId uint) error {
	res := er.db.Delete(models.Event{}, eventId)
	return res.Error
}
