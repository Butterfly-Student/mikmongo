// Package mikrotik provides a service-layer registry that wraps pkg/mikrotik repository
// implementations behind a per-router connection managed by RouterService.
package mikrotik

import (
	"context"

	"github.com/Butterfly-Student/go-ros/repository/firewall"
	goroshotspot "github.com/Butterfly-Student/go-ros/repository/hotspot"
	gorosipaddress "github.com/Butterfly-Student/go-ros/repository/ip-address"
	gorosmonitor "github.com/Butterfly-Student/go-ros/repository/monitor"
	gorosppp "github.com/Butterfly-Student/go-ros/repository/ppp"
	gorosqueue "github.com/Butterfly-Student/go-ros/repository/queue"
	"github.com/google/uuid"

	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	"mikmongo/internal/service"
	"mikmongo/internal/service/mikrotik/mikhmon"
)

// Registry exposes per-module Mikrotik service instances.
type Registry struct {
	PPP       *PPPService
	Hotspot   *HotspotService
	Queue     *QueueService
	Firewall  *FirewallService
	IPPool    *IPPoolService
	IPAddress *IPAddressService
	Monitor   *MonitorService
	Mikhmon   *mikhmon.Registry
}

// NewRegistry creates a Registry backed by routerSvc for connection management.
func NewRegistry(routerSvc *service.RouterService) *Registry {
	return &Registry{
		PPP:       &PPPService{routerSvc: routerSvc},
		Hotspot:   &HotspotService{routerSvc: routerSvc},
		Queue:     &QueueService{routerSvc: routerSvc},
		Firewall:  &FirewallService{routerSvc: routerSvc},
		IPPool:    &IPPoolService{routerSvc: routerSvc},
		IPAddress: &IPAddressService{routerSvc: routerSvc},
		Monitor:   &MonitorService{routerSvc: routerSvc},
		Mikhmon:   mikhmon.NewRegistry(routerSvc),
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// PPP Service
// ──────────────────────────────────────────────────────────────────────────────

// PPPService wraps go-ros PPP repositories.
type PPPService struct {
	routerSvc *service.RouterService
}

func (s *PPPService) pppRepo(ctx context.Context, routerID uuid.UUID) (gorosppp.Repository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	conn := c.Conn()
	return &pppRepo{
		secret:  gorosppp.NewSecretRepository(conn),
		profile: gorosppp.NewProfileRepository(conn),
		active:  gorosppp.NewActiveRepository(conn),
	}, nil
}

func (s *PPPService) GetProfiles(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.PPPProfile, error) {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Profile().GetProfiles(ctx)
}

func (s *PPPService) AddProfile(ctx context.Context, routerID uuid.UUID, p *mkdomain.PPPProfile) error {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.Profile().AddProfile(ctx, p)
}

func (s *PPPService) GetProfileByName(ctx context.Context, routerID uuid.UUID, name string) (*mkdomain.PPPProfile, error) {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Profile().GetProfileByName(ctx, name)
}

func (s *PPPService) RemoveProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.Profile().RemoveProfile(ctx, id)
}

func (s *PPPService) GetSecrets(ctx context.Context, routerID uuid.UUID, profile string) ([]*mkdomain.PPPSecret, error) {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Secret().GetSecrets(ctx, profile)
}

func (s *PPPService) AddSecret(ctx context.Context, routerID uuid.UUID, secret *mkdomain.PPPSecret) error {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.Secret().AddSecret(ctx, secret)
}

func (s *PPPService) GetSecretByName(ctx context.Context, routerID uuid.UUID, name string) (*mkdomain.PPPSecret, error) {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Secret().GetSecretByName(ctx, name)
}

func (s *PPPService) RemoveSecret(ctx context.Context, routerID uuid.UUID, id string) error {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.Secret().RemoveSecret(ctx, id)
}

func (s *PPPService) GetActiveUsers(ctx context.Context, routerID uuid.UUID, service string) ([]*mkdomain.PPPActive, error) {
	repo, err := s.pppRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Active().GetActive(ctx, service)
}

// pppRepo is a minimal in-package aggregator for go-ros PPP repositories.
type pppRepo struct {
	secret  gorosppp.SecretRepository
	profile gorosppp.ProfileRepository
	active  gorosppp.ActiveRepository
}

func (r *pppRepo) Secret() gorosppp.SecretRepository   { return r.secret }
func (r *pppRepo) Profile() gorosppp.ProfileRepository { return r.profile }
func (r *pppRepo) Active() gorosppp.ActiveRepository   { return r.active }

// ──────────────────────────────────────────────────────────────────────────────
// Hotspot Service
// ──────────────────────────────────────────────────────────────────────────────

// HotspotService wraps go-ros Hotspot repositories.
type HotspotService struct {
	routerSvc *service.RouterService
}

func (s *HotspotService) conn(ctx context.Context, routerID uuid.UUID) (goroshotspot.Repository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	raw := c.Conn()
	return &hotspotRepo{
		user:    goroshotspot.NewUserRepository(raw),
		profile: goroshotspot.NewProfileRepository(raw),
		active:  goroshotspot.NewActiveRepository(raw),
		host:    goroshotspot.NewHostRepository(raw),
		server:  goroshotspot.NewServerRepository(raw),
	}, nil
}

func (s *HotspotService) GetProfiles(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.UserProfile, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Profile().GetProfiles(ctx)
}

func (s *HotspotService) AddProfile(ctx context.Context, routerID uuid.UUID, p *mkdomain.UserProfile) (string, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return "", err
	}
	return repo.Profile().AddProfile(ctx, p)
}

func (s *HotspotService) GetProfileByName(ctx context.Context, routerID uuid.UUID, name string) (*mkdomain.UserProfile, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Profile().GetProfileByName(ctx, name)
}

func (s *HotspotService) GetProfileByID(ctx context.Context, routerID uuid.UUID, id string) (*mkdomain.UserProfile, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Profile().GetProfileByID(ctx, id)
}

func (s *HotspotService) RemoveProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.Profile().RemoveProfile(ctx, id)
}

func (s *HotspotService) GetUsers(ctx context.Context, routerID uuid.UUID, profile string) ([]*mkdomain.HotspotUser, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.User().GetUsers(ctx, profile)
}

func (s *HotspotService) AddUser(ctx context.Context, routerID uuid.UUID, u *mkdomain.HotspotUser) (string, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return "", err
	}
	return repo.User().AddUser(ctx, u)
}

func (s *HotspotService) GetUserByName(ctx context.Context, routerID uuid.UUID, name string) (*mkdomain.HotspotUser, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.User().GetUserByName(ctx, name)
}

func (s *HotspotService) GetUserByID(ctx context.Context, routerID uuid.UUID, id string) (*mkdomain.HotspotUser, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.User().GetUserByID(ctx, id)
}

func (s *HotspotService) RemoveUser(ctx context.Context, routerID uuid.UUID, id string) error {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.User().RemoveUser(ctx, id)
}

func (s *HotspotService) GetActive(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.HotspotActive, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Active().GetActive(ctx)
}

func (s *HotspotService) GetHosts(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.HotspotHost, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Host().GetHosts(ctx)
}

func (s *HotspotService) GetServers(ctx context.Context, routerID uuid.UUID) ([]string, error) {
	repo, err := s.conn(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Server().GetServers(ctx)
}

type hotspotRepo struct {
	user    goroshotspot.UserRepository
	profile goroshotspot.ProfileRepository
	active  goroshotspot.ActiveRepository
	host    goroshotspot.HostRepository
	server  goroshotspot.ServerRepository
}

func (r *hotspotRepo) User() goroshotspot.UserRepository       { return r.user }
func (r *hotspotRepo) Profile() goroshotspot.ProfileRepository { return r.profile }
func (r *hotspotRepo) Active() goroshotspot.ActiveRepository   { return r.active }
func (r *hotspotRepo) Host() goroshotspot.HostRepository       { return r.host }
func (r *hotspotRepo) Server() goroshotspot.ServerRepository   { return r.server }
func (r *hotspotRepo) IPBinding() goroshotspot.IPBindingRepository {
	return nil
}
func (r *hotspotRepo) Cookie() goroshotspot.CookieRepository { return nil }

// ──────────────────────────────────────────────────────────────────────────────
// Queue Service
// ──────────────────────────────────────────────────────────────────────────────

// QueueService wraps go-ros Queue repositories.
type QueueService struct {
	routerSvc *service.RouterService
}

func (s *QueueService) GetSimpleQueues(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.SimpleQueue, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	repo := gorosqueue.NewSimpleQueueRepository(c.Conn())
	return repo.GetSimpleQueues(ctx)
}

// ──────────────────────────────────────────────────────────────────────────────
// Firewall Service
// ──────────────────────────────────────────────────────────────────────────────

// FirewallService wraps go-ros Firewall repositories.
type FirewallService struct {
	routerSvc *service.RouterService
}

func (s *FirewallService) fwRepo(ctx context.Context, routerID uuid.UUID) (firewall.Repository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	conn := c.Conn()
	return &firewallRepo{
		nat:      firewall.NewNATRepository(conn),
		filter:   firewall.NewFilterRepository(conn),
		addrList: firewall.NewAddressListRepository(conn),
	}, nil
}

func (s *FirewallService) GetFilterRules(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.FirewallRule, error) {
	repo, err := s.fwRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Filter().GetRules(ctx)
}

func (s *FirewallService) GetNATRules(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.NATRule, error) {
	repo, err := s.fwRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.NAT().GetNATRules(ctx)
}

func (s *FirewallService) GetAddressLists(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.AddressList, error) {
	repo, err := s.fwRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.AddressList().GetAddressLists(ctx)
}

type firewallRepo struct {
	nat      firewall.NATRepository
	filter   firewall.FilterRepository
	addrList firewall.AddressListRepository
}

func (r *firewallRepo) NAT() firewall.NATRepository                 { return r.nat }
func (r *firewallRepo) Filter() firewall.FilterRepository           { return r.filter }
func (r *firewallRepo) AddressList() firewall.AddressListRepository { return r.addrList }

// ──────────────────────────────────────────────────────────────────────────────
// IP Pool Service
// ──────────────────────────────────────────────────────────────────────────────

// IPPoolService wraps go-ros IP Pool repositories.
type IPPoolService struct {
	routerSvc *service.RouterService
}

func (s *IPPoolService) ipRepo(ctx context.Context, routerID uuid.UUID) (gorosipaddress.Repository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return gorosipaddress.NewRepository(c.Conn()), nil
}

func (s *IPPoolService) GetPools(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.IPPool, error) {
	repo, err := s.ipRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Pool().GetPools(ctx)
}

func (s *IPPoolService) AddPool(ctx context.Context, routerID uuid.UUID, pool *mkdomain.IPPool) (string, error) {
	repo, err := s.ipRepo(ctx, routerID)
	if err != nil {
		return "", err
	}
	return repo.Pool().AddPool(ctx, pool)
}

func (s *IPPoolService) GetPoolByName(ctx context.Context, routerID uuid.UUID, name string) (*mkdomain.IPPool, error) {
	repo, err := s.ipRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Pool().GetPoolByName(ctx, name)
}

func (s *IPPoolService) GetPoolByID(ctx context.Context, routerID uuid.UUID, id string) (*mkdomain.IPPool, error) {
	repo, err := s.ipRepo(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return repo.Pool().GetPoolByID(ctx, id)
}

func (s *IPPoolService) RemovePool(ctx context.Context, routerID uuid.UUID, id string) error {
	repo, err := s.ipRepo(ctx, routerID)
	if err != nil {
		return err
	}
	return repo.Pool().RemovePool(ctx, id)
}

// ──────────────────────────────────────────────────────────────────────────────
// IP Address Service
// ──────────────────────────────────────────────────────────────────────────────

// IPAddressService wraps go-ros IP Address repositories.
type IPAddressService struct {
	routerSvc *service.RouterService
}

func (s *IPAddressService) GetAddresses(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.IPAddress, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	repo := gorosipaddress.NewRepository(c.Conn())
	return repo.Address().GetAddresses(ctx)
}

// ──────────────────────────────────────────────────────────────────────────────
// Monitor Service
// ──────────────────────────────────────────────────────────────────────────────

// MonitorService wraps go-ros Monitor repositories.
type MonitorService struct {
	routerSvc *service.RouterService
}

func (s *MonitorService) GetSystemResource(ctx context.Context, routerID uuid.UUID) (*mkdomain.SystemResource, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	repo := gorosmonitor.NewRepository(c.Conn())
	return repo.System().GetSystemResource(ctx)
}

func (s *MonitorService) GetInterfaces(ctx context.Context, routerID uuid.UUID) ([]*mkdomain.Interface, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	repo := gorosmonitor.NewRepository(c.Conn())
	return repo.Interface().GetInterfaces(ctx)
}
