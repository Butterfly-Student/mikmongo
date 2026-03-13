package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

type sequenceCounterRepository struct {
	db *gorm.DB
}

func NewSequenceCounterRepository(db *gorm.DB) repository.SequenceCounterRepository {
	return &sequenceCounterRepository{db: db}
}

func (r *sequenceCounterRepository) Create(ctx context.Context, counter *model.SequenceCounter) error {
	return r.db.WithContext(ctx).Create(counter).Error
}

func (r *sequenceCounterRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SequenceCounter, error) {
	var counter model.SequenceCounter
	err := r.db.WithContext(ctx).First(&counter, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

func (r *sequenceCounterRepository) GetByName(ctx context.Context, name string) (*model.SequenceCounter, error) {
	var counter model.SequenceCounter
	err := r.db.WithContext(ctx).First(&counter, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &counter, nil
}

func (r *sequenceCounterRepository) Update(ctx context.Context, counter *model.SequenceCounter) error {
	return r.db.WithContext(ctx).Save(counter).Error
}

func (r *sequenceCounterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.SequenceCounter{}, "id = ?", id).Error
}

func (r *sequenceCounterRepository) List(ctx context.Context, limit, offset int) ([]model.SequenceCounter, error) {
	var counters []model.SequenceCounter
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&counters).Error
	return counters, err
}

func (r *sequenceCounterRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.SequenceCounter{}).Count(&count).Error
	return count, err
}

func (r *sequenceCounterRepository) NextNumber(ctx context.Context, name string) (int, error) {
	var counter model.SequenceCounter
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&counter, "name = ?", name).Error; err != nil {
			return err
		}
		counter.LastNumber++
		return tx.Save(&counter).Error
	})
	if err != nil {
		return 0, err
	}
	return counter.LastNumber, nil
}
