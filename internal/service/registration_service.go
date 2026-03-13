package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// RegistrationService handles customer self-registration
type RegistrationService struct {
	regRepo         repository.CustomerRegistrationRepository
	customerSvc     *CustomerService
	subscriptionSvc *SubscriptionService
	notificationSvc *NotificationService
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(
	regRepo repository.CustomerRegistrationRepository,
	customerSvc *CustomerService,
	subscriptionSvc *SubscriptionService,
) *RegistrationService {
	return &RegistrationService{
		regRepo:         regRepo,
		customerSvc:     customerSvc,
		subscriptionSvc: subscriptionSvc,
	}
}

// SetNotificationService injects notification service
func (s *RegistrationService) SetNotificationService(n *NotificationService) {
	s.notificationSvc = n
}

// Create creates a new registration request
func (s *RegistrationService) Create(ctx context.Context, reg *model.CustomerRegistration) error {
	reg.Status = "pending"
	return s.regRepo.Create(ctx, reg)
}

// GetByID gets registration by ID
func (s *RegistrationService) GetByID(ctx context.Context, id uuid.UUID) (*model.CustomerRegistration, error) {
	return s.regRepo.GetByID(ctx, id)
}

// List lists registrations with pagination
func (s *RegistrationService) List(ctx context.Context, limit, offset int) ([]model.CustomerRegistration, int64, error) {
	regs, err := s.regRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.regRepo.Count(ctx)
	return regs, count, err
}

// ListByStatus lists registrations by status
func (s *RegistrationService) ListByStatus(ctx context.Context, status string) ([]model.CustomerRegistration, error) {
	return s.regRepo.ListByStatus(ctx, status)
}

// Approve approves a registration and creates customer + subscription
func (s *RegistrationService) Approve(ctx context.Context, regID uuid.UUID, approvedByID, routerID string, profileID *string) error {
	reg, err := s.regRepo.GetByID(ctx, regID)
	if err != nil {
		return err
	}
	if reg.Status != "pending" {
		return fmt.Errorf("registration is not pending")
	}

	// Create customer
	customer := &model.Customer{
		FullName:  reg.FullName,
		Email:     reg.Email,
		Phone:     reg.Phone,
		Address:   reg.Address,
		Latitude:  reg.Latitude,
		Longitude: reg.Longitude,
		Notes:     reg.Notes,
		CreatedBy: &approvedByID,
	}
	if err := s.customerSvc.Create(ctx, customer); err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	// Create subscription if profile and router are provided
	var username, password string
	if profileID != nil && routerID != "" {
		sub := &model.Subscription{
			CustomerID: customer.ID,
			PlanID:     *profileID,
			RouterID:   routerID,
			Username:   reg.Phone,
			Status:     "pending",
		}
		if err := s.subscriptionSvc.Create(ctx, sub); err == nil {
			username = sub.Username
			password = sub.Password
		}
	}

	// Update registration
	now := time.Now()
	customerIDStr := customer.ID
	reg.Status = "approved"
	reg.ApprovedBy = &approvedByID
	reg.ApprovedAt = &now
	reg.CustomerID = &customerIDStr
	if err := s.regRepo.Update(ctx, reg); err != nil {
		return err
	}

	// Send notification
	if s.notificationSvc != nil && username != "" {
		_ = s.notificationSvc.SendRegistrationApproved(ctx, customer, username, password)
	}

	return nil
}

// Reject rejects a registration
func (s *RegistrationService) Reject(ctx context.Context, regID uuid.UUID, reason, rejectedByID string) error {
	reg, err := s.regRepo.GetByID(ctx, regID)
	if err != nil {
		return err
	}
	if reg.Status != "pending" {
		return fmt.Errorf("registration is not pending")
	}

	if err := s.regRepo.UpdateStatus(ctx, regID, "rejected", reason, &rejectedByID); err != nil {
		return err
	}

	if s.notificationSvc != nil {
		_ = s.notificationSvc.SendRegistrationRejected(ctx, reg.Phone, reg.FullName, reason)
	}

	return nil
}
