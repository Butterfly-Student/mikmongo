package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/domain/customer"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
)

// CustomerService handles customer business logic
type CustomerService struct {
	repo            repository.CustomerRepository
	seqRepo         repository.SequenceCounterRepository
	profileRepo     repository.BandwidthProfileRepository
	customerDomain  *customer.Domain
	subscriptionSvc *SubscriptionService
	routerSvc       *RouterService
}

// NewCustomerService creates a new customer service
func NewCustomerService(
	repo repository.CustomerRepository,
	seqRepo repository.SequenceCounterRepository,
	profileRepo repository.BandwidthProfileRepository,
	customerDomain *customer.Domain,
	routerSvc *RouterService,
) *CustomerService {
	return &CustomerService{
		repo:           repo,
		seqRepo:        seqRepo,
		profileRepo:    profileRepo,
		customerDomain: customerDomain,
		routerSvc:      routerSvc,
	}
}

// SetSubscriptionService injects subscription service (avoids circular dep)
func (s *CustomerService) SetSubscriptionService(sub *SubscriptionService) {
	s.subscriptionSvc = sub
}

// generateCustomerCode generates a new customer code like CST00001
func (s *CustomerService) generateCustomerCode(ctx context.Context) (string, error) {
	n, err := s.seqRepo.NextNumber(ctx, "customer_code")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("CST%05d", n), nil
}

// Create creates a new customer with auto-generated CustomerCode
func (s *CustomerService) Create(ctx context.Context, c *model.Customer) error {
	if err := s.customerDomain.ValidateCustomer(c); err != nil {
		return err
	}
	if c.CustomerCode == "" {
		code, err := s.generateCustomerCode(ctx)
		if err != nil {
			return fmt.Errorf("failed to generate customer code: %w", err)
		}
		c.CustomerCode = code
	}
	if c.Username == nil || *c.Username == "" {
		u := generateUsernameFromFullName(c.FullName)
		c.Username = &u
	}
	c.IsActive = true
	return s.repo.Create(ctx, c)
}

// CreateWithSubscription creates a new customer and subscription atomically
func (s *CustomerService) CreateWithSubscription(ctx context.Context, customer *model.Customer, subscription *model.Subscription) (*model.Customer, *model.Subscription, error) {
	// Validate customer
	if err := s.customerDomain.ValidateCustomer(customer); err != nil {
		return nil, nil, fmt.Errorf("invalid customer data: %w", err)
	}

	// Generate customer code
	if customer.CustomerCode == "" {
		code, err := s.generateCustomerCode(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate customer code: %w", err)
		}
		customer.CustomerCode = code
	}
	if customer.Username == nil || *customer.Username == "" {
		u := generateUsernameFromFullName(customer.FullName)
		customer.Username = &u
	}
	customer.IsActive = true

	// Validate plan exists and belongs to router
	planID, err := uuid.Parse(subscription.PlanID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid plan ID: %w", err)
	}
	profile, err := s.profileRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, nil, fmt.Errorf("plan not found: %w", err)
	}
	if profile.RouterID != subscription.RouterID {
		return nil, nil, fmt.Errorf("profile does not belong to the specified router")
	}

	// Validate credentials
	if subscription.Username == "" {
		return nil, nil, fmt.Errorf("username is required")
	}
	if subscription.Password == "" {
		return nil, nil, fmt.Errorf("password is required")
	}
	if len(subscription.Username) < 3 || len(subscription.Username) > 100 {
		return nil, nil, fmt.Errorf("username must be between 3 and 100 characters")
	}
	if len(subscription.Password) < 6 {
		return nil, nil, fmt.Errorf("password must be at least 6 characters")
	}

	// Create customer first
	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Set customer ID to subscription
	subscription.CustomerID = customer.ID
	subscription.Status = "pending"

	// Create subscription
	if s.subscriptionSvc != nil {
		if err := s.subscriptionSvc.Create(ctx, subscription, nil); err != nil {
			// Rollback: delete customer
			customerID, _ := uuid.Parse(customer.ID)
			_ = s.repo.Delete(ctx, customerID)
			return nil, nil, fmt.Errorf("failed to create subscription: %w", err)
		}
	} else {
		// Rollback: delete customer
		customerID, _ := uuid.Parse(customer.ID)
		_ = s.repo.Delete(ctx, customerID)
		return nil, nil, fmt.Errorf("subscription service not available")
	}

	return customer, subscription, nil
}

// generateUsernameFromFullName generates username from full name
func generateUsernameFromFullName(fullName string) string {
	// Lowercase and replace spaces with -
	username := strings.ToLower(fullName)
	username = strings.ReplaceAll(username, " ", "-")

	// Remove special characters (keep only alphanumeric and -)
	var result strings.Builder
	for _, char := range username {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// GetByID gets customer by ID
func (s *CustomerService) GetByID(ctx context.Context, id uuid.UUID) (*model.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCode gets customer by customer code
func (s *CustomerService) GetByCode(ctx context.Context, code string) (*model.Customer, error) {
	customers, err := s.repo.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}
	for _, c := range customers {
		if c.CustomerCode == code {
			return &c, nil
		}
	}
	return nil, errors.New("customer not found")
}

// List lists customers with pagination
func (s *CustomerService) List(ctx context.Context, limit, offset int) ([]model.Customer, int64, error) {
	customers, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	return customers, count, err
}

// Update updates a customer
func (s *CustomerService) Update(ctx context.Context, c *model.Customer) error {
	return s.repo.Update(ctx, c)
}

// Delete soft-deletes a customer
func (s *CustomerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ActivateAccount activates a customer account (sets is_active = true)
func (s *CustomerService) ActivateAccount(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.customerDomain.CanActivate(c); err != nil {
		return err
	}
	c.IsActive = true
	return s.repo.Update(ctx, c)
}

// DeactivateAccount deactivates a customer account (sets is_active = false)
func (s *CustomerService) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.customerDomain.CanDeactivate(c); err != nil {
		return err
	}
	c.IsActive = false
	return s.repo.Update(ctx, c)
}

// SuspendAllSubscriptions suspends all active subscriptions for a customer
func (s *CustomerService) SuspendAllSubscriptions(ctx context.Context, id uuid.UUID, reason string) error {
	// Suspend all active subscriptions
	if s.subscriptionSvc != nil {
		subs, _ := s.subscriptionSvc.GetByCustomerID(ctx, id)
		for _, sub := range subs {
			if sub.Status == "active" || sub.Status == "isolated" {
				subID, err := uuid.Parse(sub.ID)
				if err != nil {
					continue
				}
				_ = s.subscriptionSvc.Suspend(ctx, subID, reason)
			}
		}
	}
	return nil
}

// IsolateAllSubscriptions isolates all active subscriptions for a customer
func (s *CustomerService) IsolateAllSubscriptions(ctx context.Context, id uuid.UUID) error {
	// Isolate all active subscriptions
	if s.subscriptionSvc != nil {
		subs, _ := s.subscriptionSvc.GetByCustomerID(ctx, id)
		for _, sub := range subs {
			if sub.Status == "active" {
				subID, err := uuid.Parse(sub.ID)
				if err != nil {
					continue
				}
				_ = s.subscriptionSvc.Isolate(ctx, subID, "invoice_overdue")
			}
		}
	}
	return nil
}

// RestoreAllSubscriptions restores all isolated subscriptions for a customer
func (s *CustomerService) RestoreAllSubscriptions(ctx context.Context, id uuid.UUID) error {
	// Restore all isolated subscriptions
	if s.subscriptionSvc != nil {
		subs, _ := s.subscriptionSvc.GetByCustomerID(ctx, id)
		for _, sub := range subs {
			if sub.Status == "isolated" {
				subID, err := uuid.Parse(sub.ID)
				if err != nil {
					continue
				}
				_ = s.subscriptionSvc.Restore(ctx, subID)
			}
		}
	}
	return nil
}

// TerminateAllSubscriptions terminates all subscriptions for a customer
func (s *CustomerService) TerminateAllSubscriptions(ctx context.Context, id uuid.UUID) error {
	if s.subscriptionSvc != nil {
		subs, _ := s.subscriptionSvc.GetByCustomerID(ctx, id)
		for _, sub := range subs {
			if sub.Status != "terminated" {
				subID, err := uuid.Parse(sub.ID)
				if err != nil {
					continue
				}
				_ = s.subscriptionSvc.Terminate(ctx, subID)
			}
		}
	}
	return nil
}

// SetPortalPassword sets customer portal password
func (s *CustomerService) SetPortalPassword(ctx context.Context, id uuid.UUID, password string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	h := string(hashed)
	c.PortalPasswordHash = &h
	return s.repo.Update(ctx, c)
}

// AuthPortal authenticates a customer for portal access using username or email
func (s *CustomerService) AuthPortal(ctx context.Context, identifier, password string) (*model.Customer, error) {
	var c *model.Customer
	var err error

	c, err = s.repo.GetByUsername(ctx, identifier)
	if err != nil {
		c, err = s.repo.GetByEmail(ctx, identifier)
	}
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if c.PortalPasswordHash == nil {
		return nil, errors.New("portal access not configured")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*c.PortalPasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	now := time.Now()
	c.PortalLastLogin = &now
	_ = s.repo.Update(ctx, c)
	return c, nil
}
