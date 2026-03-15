package ppp

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// profileRepository implements ProfileRepository interface
type profileRepository struct {
	client *client.Client
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(c *client.Client) ProfileRepository {
	return &profileRepository{client: c}
}

func parsePPPProfile(m map[string]string) *domain.PPPProfile {
	return &domain.PPPProfile{
		ID:                 m[".id"],
		Name:               m["name"],
		LocalAddress:       m["local-address"],
		RemoteAddress:      m["remote-address"],
		DNSServer:          m["dns-server"],
		SessionTimeout:     m["session-timeout"],
		IdleTimeout:        m["idle-timeout"],
		OnlyOne:            utils.ParseBool(m["only-one"]),
		Comment:            m["comment"],
		RateLimit:          m["rate-limit"],
		ParentQueue:        m["parent-queue"],
		QueueType:          m["queue-type"],
		UseCompression:     utils.ParseBool(m["use-compression"]),
		UseEncryption:      utils.ParseBool(m["use-encryption"]),
		UseMPLS:            utils.ParseBool(m["use-mpls"]),
		UseUPnP:            utils.ParseBool(m["use-upnp"]),
		Bridge:             m["bridge"],
		AddressList:        m["address-list"],
		InterfaceList:      m["interface-list"],
		OnUp:               m["on-up"],
		OnDown:             m["on-down"],
		ChangeTCPMSS:       utils.ParseBool(m["change-tcp-mss"]),
		IncomingFilter:     m["incoming-filter"],
		OutgoingFilter:     m["outgoing-filter"],
		InsertQueueBefore:  m["insert-queue-before"],
		WinsServer:         m["wins-server"],
		BridgeHorizon:      m["bridge-horizon"],
		BridgeLearning:     utils.ParseBool(m["bridge-learning"]),
		BridgePathCost:     int(utils.ParseInt(m["bridge-path-cost"])),
		BridgePortPriority: int(utils.ParseInt(m["bridge-port-priority"])),
	}
}

func appendPPPProfileArgs(args []string, p *domain.PPPProfile) []string {
	if p.LocalAddress != "" {
		args = append(args, "=local-address="+p.LocalAddress)
	}
	if p.RemoteAddress != "" {
		args = append(args, "=remote-address="+p.RemoteAddress)
	}
	if p.DNSServer != "" {
		args = append(args, "=dns-server="+p.DNSServer)
	}
	if p.SessionTimeout != "" {
		args = append(args, "=session-timeout="+p.SessionTimeout)
	}
	if p.IdleTimeout != "" {
		args = append(args, "=idle-timeout="+p.IdleTimeout)
	}
	if p.OnlyOne {
		args = append(args, "=only-one=yes")
	}
	if p.RateLimit != "" {
		args = append(args, "=rate-limit="+p.RateLimit)
	}
	if p.ParentQueue != "" {
		args = append(args, "=parent-queue="+p.ParentQueue)
	}
	if p.QueueType != "" {
		args = append(args, "=queue-type="+p.QueueType)
	}
	if p.UseCompression {
		args = append(args, "=use-compression=yes")
	}
	if p.UseEncryption {
		args = append(args, "=use-encryption=yes")
	}
	if p.UseMPLS {
		args = append(args, "=use-mpls=yes")
	}
	if p.UseUPnP {
		args = append(args, "=use-upnp=yes")
	}
	if p.Bridge != "" {
		args = append(args, "=bridge="+p.Bridge)
	}
	if p.BridgeHorizon != "" {
		args = append(args, "=bridge-horizon="+p.BridgeHorizon)
	}
	if p.BridgeLearning {
		args = append(args, "=bridge-learning=yes")
	}
	if p.BridgePathCost > 0 {
		args = append(args, "=bridge-path-cost="+utils.FormatInt(int64(p.BridgePathCost)))
	}
	if p.BridgePortPriority > 0 {
		args = append(args, "=bridge-port-priority="+utils.FormatInt(int64(p.BridgePortPriority)))
	}
	if p.ChangeTCPMSS {
		args = append(args, "=change-tcp-mss=yes")
	}
	if p.AddressList != "" {
		args = append(args, "=address-list="+p.AddressList)
	}
	if p.InterfaceList != "" {
		args = append(args, "=interface-list="+p.InterfaceList)
	}
	if p.OnUp != "" {
		args = append(args, "=on-up="+p.OnUp)
	}
	if p.OnDown != "" {
		args = append(args, "=on-down="+p.OnDown)
	}
	if p.IncomingFilter != "" {
		args = append(args, "=incoming-filter="+p.IncomingFilter)
	}
	if p.OutgoingFilter != "" {
		args = append(args, "=outgoing-filter="+p.OutgoingFilter)
	}
	if p.InsertQueueBefore != "" {
		args = append(args, "=insert-queue-before="+p.InsertQueueBefore)
	}
	if p.WinsServer != "" {
		args = append(args, "=wins-server="+p.WinsServer)
	}
	if p.Comment != "" {
		args = append(args, "=comment="+p.Comment)
	}
	return args
}

func (r *profileRepository) GetProfiles(ctx context.Context, proplist ...string) ([]*domain.PPPProfile, error) {
	args := []string{"/ppp/profile/print"}
	// Add proplist
	pl := utils.BuildProplist(ProplistPPPProfileDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	profiles := make([]*domain.PPPProfile, 0, len(reply.Re))
	for _, re := range reply.Re {
		profiles = append(profiles, parsePPPProfile(re.Map))
	}
	return profiles, nil
}

func (r *profileRepository) GetProfileByID(ctx context.Context, id string, proplist ...string) (*domain.PPPProfile, error) {
	args := []string{"/ppp/profile/print", "?.id=" + id}
	// Add proplist
	pl := utils.BuildProplist(ProplistPPPProfileDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPProfile(reply.Re[0].Map), nil
}

func (r *profileRepository) GetProfileByName(ctx context.Context, name string, proplist ...string) (*domain.PPPProfile, error) {
	args := []string{"/ppp/profile/print", "?name=" + name}
	// Add proplist
	pl := utils.BuildProplist(ProplistPPPProfileDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPProfile(reply.Re[0].Map), nil
}

func (r *profileRepository) AddProfile(ctx context.Context, profile *domain.PPPProfile) error {
	args := []string{"/ppp/profile/add", "=name=" + profile.Name}
	args = appendPPPProfileArgs(args, profile)
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *profileRepository) UpdateProfile(ctx context.Context, id string, profile *domain.PPPProfile) error {
	args := []string{"/ppp/profile/set", "=.id=" + id}
	if profile.Name != "" {
		args = append(args, "=name="+profile.Name)
	}
	args = appendPPPProfileArgs(args, profile)
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *profileRepository) RemoveProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/profile/remove", "=.id="+id)
	return err
}

func (r *profileRepository) DisableProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/profile/disable", "=.id="+id)
	return err
}

func (r *profileRepository) EnableProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/profile/enable", "=.id="+id)
	return err
}

func (r *profileRepository) RemoveProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *profileRepository) DisableProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisableProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *profileRepository) EnableProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnableProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}
