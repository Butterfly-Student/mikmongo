package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mikmongo/internal/model"
)

// MockBandwidthProfileRepository is a mock implementation of repository.BandwidthProfileRepository
type MockBandwidthProfileRepository struct {
	mock.Mock
}

func (m *MockBandwidthProfileRepository) Create(ctx context.Context, profile *model.BandwidthProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockBandwidthProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.BandwidthProfile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) GetByCode(ctx context.Context, code string) (*model.BandwidthProfile, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) GetByRouterAndCode(ctx context.Context, routerID uuid.UUID, code string) (*model.BandwidthProfile, error) {
	args := m.Called(ctx, routerID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) GetByRouterAndName(ctx context.Context, routerID uuid.UUID, name string) (*model.BandwidthProfile, error) {
	args := m.Called(ctx, routerID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) Update(ctx context.Context, profile *model.BandwidthProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockBandwidthProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBandwidthProfileRepository) List(ctx context.Context, limit, offset int) ([]model.BandwidthProfile, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.BandwidthProfile, error) {
	args := m.Called(ctx, routerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBandwidthProfileRepository) CountByRouterID(ctx context.Context, routerID uuid.UUID) (int64, error) {
	args := m.Called(ctx, routerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBandwidthProfileRepository) ListActive(ctx context.Context) ([]model.BandwidthProfile, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.BandwidthProfile), args.Error(1)
}

func (m *MockBandwidthProfileRepository) ListActiveByRouterID(ctx context.Context, routerID uuid.UUID) ([]model.BandwidthProfile, error) {
	args := m.Called(ctx, routerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.BandwidthProfile), args.Error(1)
}

// MockSequenceCounterRepository is a mock implementation of repository.SequenceCounterRepository
type MockSequenceCounterRepository struct {
	mock.Mock
}

func (m *MockSequenceCounterRepository) Create(ctx context.Context, counter *model.SequenceCounter) error {
	args := m.Called(ctx, counter)
	return args.Error(0)
}

func (m *MockSequenceCounterRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SequenceCounter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SequenceCounter), args.Error(1)
}

func (m *MockSequenceCounterRepository) GetByName(ctx context.Context, name string) (*model.SequenceCounter, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SequenceCounter), args.Error(1)
}

func (m *MockSequenceCounterRepository) Update(ctx context.Context, counter *model.SequenceCounter) error {
	args := m.Called(ctx, counter)
	return args.Error(0)
}

func (m *MockSequenceCounterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSequenceCounterRepository) List(ctx context.Context, limit, offset int) ([]model.SequenceCounter, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SequenceCounter), args.Error(1)
}

func (m *MockSequenceCounterRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSequenceCounterRepository) NextNumber(ctx context.Context, name string) (int, error) {
	args := m.Called(ctx, name)
	return args.Int(0), args.Error(1)
}

// MockSystemSettingRepository is a mock implementation of repository.SystemSettingRepository
type MockSystemSettingRepository struct {
	mock.Mock
}

func (m *MockSystemSettingRepository) Create(ctx context.Context, setting *model.SystemSetting) error {
	args := m.Called(ctx, setting)
	return args.Error(0)
}

func (m *MockSystemSettingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSetting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SystemSetting), args.Error(1)
}

func (m *MockSystemSettingRepository) GetByGroupAndKey(ctx context.Context, group, key string) (*model.SystemSetting, error) {
	args := m.Called(ctx, group, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SystemSetting), args.Error(1)
}

func (m *MockSystemSettingRepository) Update(ctx context.Context, setting *model.SystemSetting) error {
	args := m.Called(ctx, setting)
	return args.Error(0)
}

func (m *MockSystemSettingRepository) List(ctx context.Context, limit, offset int) ([]model.SystemSetting, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SystemSetting), args.Error(1)
}

func (m *MockSystemSettingRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSystemSettingRepository) ListByGroup(ctx context.Context, group string) ([]model.SystemSetting, error) {
	args := m.Called(ctx, group)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.SystemSetting), args.Error(1)
}

func (m *MockSystemSettingRepository) Upsert(ctx context.Context, setting *model.SystemSetting) error {
	args := m.Called(ctx, setting)
	return args.Error(0)
}

// MockMessageTemplateRepository is a mock implementation of repository.MessageTemplateRepository
type MockMessageTemplateRepository struct {
	mock.Mock
}

func (m *MockMessageTemplateRepository) Create(ctx context.Context, tmpl *model.MessageTemplate) error {
	args := m.Called(ctx, tmpl)
	return args.Error(0)
}

func (m *MockMessageTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MessageTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MessageTemplate), args.Error(1)
}

func (m *MockMessageTemplateRepository) GetByEventAndChannel(ctx context.Context, event, channel string) (*model.MessageTemplate, error) {
	args := m.Called(ctx, event, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MessageTemplate), args.Error(1)
}

func (m *MockMessageTemplateRepository) Update(ctx context.Context, tmpl *model.MessageTemplate) error {
	args := m.Called(ctx, tmpl)
	return args.Error(0)
}

func (m *MockMessageTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMessageTemplateRepository) List(ctx context.Context, limit, offset int) ([]model.MessageTemplate, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MessageTemplate), args.Error(1)
}

func (m *MockMessageTemplateRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageTemplateRepository) ListByEvent(ctx context.Context, event string) ([]model.MessageTemplate, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MessageTemplate), args.Error(1)
}

func (m *MockMessageTemplateRepository) ListActive(ctx context.Context) ([]model.MessageTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.MessageTemplate), args.Error(1)
}
