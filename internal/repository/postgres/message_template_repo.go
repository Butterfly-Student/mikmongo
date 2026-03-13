package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type messageTemplateRepository struct {
	db *gorm.DB
}

func NewMessageTemplateRepository(db *gorm.DB) repository.MessageTemplateRepository {
	return &messageTemplateRepository{db: db}
}

func (r *messageTemplateRepository) Create(ctx context.Context, tmpl *model.MessageTemplate) error {
	return r.db.WithContext(ctx).Create(tmpl).Error
}

func (r *messageTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MessageTemplate, error) {
	var tmpl model.MessageTemplate
	err := r.db.WithContext(ctx).First(&tmpl, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *messageTemplateRepository) GetByEventAndChannel(ctx context.Context, event, channel string) (*model.MessageTemplate, error) {
	var tmpl model.MessageTemplate
	err := r.db.WithContext(ctx).First(&tmpl, "event = ? AND channel = ?", event, channel).Error
	if err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func (r *messageTemplateRepository) Update(ctx context.Context, tmpl *model.MessageTemplate) error {
	return r.db.WithContext(ctx).Save(tmpl).Error
}

func (r *messageTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.MessageTemplate{}, "id = ?", id).Error
}

func (r *messageTemplateRepository) List(ctx context.Context, limit, offset int) ([]model.MessageTemplate, error) {
	var tmpls []model.MessageTemplate
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&tmpls).Error
	return tmpls, err
}

func (r *messageTemplateRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.MessageTemplate{}).Count(&count).Error
	return count, err
}

func (r *messageTemplateRepository) ListByEvent(ctx context.Context, event string) ([]model.MessageTemplate, error) {
	var tmpls []model.MessageTemplate
	err := r.db.WithContext(ctx).Where("event = ?", event).Find(&tmpls).Error
	return tmpls, err
}

func (r *messageTemplateRepository) ListActive(ctx context.Context) ([]model.MessageTemplate, error) {
	var tmpls []model.MessageTemplate
	err := r.db.WithContext(ctx).Where("is_active = true").Find(&tmpls).Error
	return tmpls, err
}
