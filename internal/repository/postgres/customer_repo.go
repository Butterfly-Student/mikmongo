package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// customerRepository implements repository.CustomerRepository
type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *gorm.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).First(&customer, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) List(ctx context.Context, limit, offset int) ([]model.Customer, error) {
	var customers []model.Customer
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&customers).Error
	return customers, err
}

func (r *customerRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Customer{}).Count(&count).Error
	return count, err
}
