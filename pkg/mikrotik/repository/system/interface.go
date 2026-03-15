package system

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// IdentityRepository defines the interface for system identity data access
type IdentityRepository interface {
	GetIdentity(ctx context.Context) (*domain.SystemIdentity, error)
	SetIdentity(ctx context.Context, name string) error
}

// ResourcesRepository defines the interface for system resources data access
type ResourcesRepository interface {
	GetResources(ctx context.Context) (*domain.SystemResource, error)
}

// RouterBoardRepository defines the interface for routerboard data access
type RouterBoardRepository interface {
	GetRouterBoardInfo(ctx context.Context) (*domain.RouterBoardInfo, error)
	GetFirmware(ctx context.Context) (current, upgrade string, err error)
}

// SchedulerRepository defines the interface for scheduler data access
type SchedulerRepository interface {
	GetSchedulers(ctx context.Context) ([]*domain.Scheduler, error)
	GetSchedulerByName(ctx context.Context, name string) (*domain.Scheduler, error)
	AddScheduler(ctx context.Context, scheduler *domain.Scheduler) (string, error)
	UpdateScheduler(ctx context.Context, id string, scheduler *domain.Scheduler) error
	RemoveScheduler(ctx context.Context, id string) error
	EnableScheduler(ctx context.Context, id string) error
	DisableScheduler(ctx context.Context, id string) error
}

// ScriptsRepository defines the interface for system scripts data access
type ScriptsRepository interface {
	GetScripts(ctx context.Context) ([]*domain.SystemScript, error)
	GetScriptByName(ctx context.Context, name string) (*domain.SystemScript, error)
	AddScript(ctx context.Context, script *domain.SystemScript) (string, error)
	UpdateScript(ctx context.Context, id string, script *domain.SystemScript) error
	RemoveScript(ctx context.Context, id string) error
	RunScript(ctx context.Context, name string) error
}

// Repository is the aggregator interface for all system repositories
type Repository interface {
	Identity() IdentityRepository
	Resources() ResourcesRepository
	RouterBoard() RouterBoardRepository
	Scheduler() SchedulerRepository
	Scripts() ScriptsRepository
}

// repository implements Repository interface
type repository struct {
	identity    IdentityRepository
	resources   ResourcesRepository
	routerBoard RouterBoardRepository
	scheduler   SchedulerRepository
	scripts     ScriptsRepository
}

// NewRepository creates a new system repository aggregator
func NewRepository(c *client.Client) Repository {
	return &repository{
		identity:    NewIdentityRepository(c),
		resources:   NewResourcesRepository(c),
		routerBoard: NewRouterBoardRepository(c),
		scheduler:   NewSchedulerRepository(c),
		scripts:     NewScriptsRepository(c),
	}
}

func (r *repository) Identity() IdentityRepository {
	return r.identity
}

func (r *repository) Resources() ResourcesRepository {
	return r.resources
}

func (r *repository) RouterBoard() RouterBoardRepository {
	return r.routerBoard
}

func (r *repository) Scheduler() SchedulerRepository {
	return r.scheduler
}

func (r *repository) Scripts() ScriptsRepository {
	return r.scripts
}
