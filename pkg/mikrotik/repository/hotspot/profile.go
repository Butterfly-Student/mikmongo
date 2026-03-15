package hotspot

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

func parseUserProfile(m map[string]string) *domain.UserProfile {
	return &domain.UserProfile{
		ID:                 m[".id"],
		Name:               m["name"],
		AddressPool:        m["address-pool"],
		SharedUsers:        int(utils.ParseInt(m["shared-users"])),
		RateLimit:          m["rate-limit"],
		ParentQueue:        m["parent-queue"],
		QueueType:          m["queue-type"],
		StatusAutorefresh:  m["status-autorefresh"],
		OnLogin:            m["on-login"],
		OnLogout:           m["on-logout"],
		OpenStatusPage:     m["open-status-page"],
		TransparentProxy:   utils.ParseBool(m["transparent-proxy"]),
		Advertise:          utils.ParseBool(m["advertise"]),
		AdvertiseInterval:  m["advertise-interval"],
		AdvertiseTimeout:   m["advertise-timeout"],
		AdvertiseURL:       m["advertise-url"],
		IdleTimeout:        m["idle-timeout"],
		SessionTimeout:     m["session-timeout"],
		KeepaliveTimeout:   m["keepalive-timeout"],
		MacCookieTimeout:   m["mac-cookie-timeout"],
		AddMacCookie:       utils.ParseBool(m["add-mac-cookie"]),
		AddressList:        m["address-list"],
		IncomingFilter:     m["incoming-filter"],
		IncomingPacketMark: m["incoming-packet-mark"],
		OutgoingFilter:     m["outgoing-filter"],
		OutgoingPacketMark: m["outgoing-packet-mark"],
		InsertQueueBefore:  m["insert-queue-before"],
	}
}

func (r *profileRepository) GetProfiles(ctx context.Context, proplist ...string) ([]*domain.UserProfile, error) {
	args := []string{"/ip/hotspot/user/profile/print"}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotProfileDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}

	profiles := make([]*domain.UserProfile, 0, len(reply.Re))
	for _, re := range reply.Re {
		profiles = append(profiles, parseUserProfile(re.Map))
	}
	return profiles, nil
}

func (r *profileRepository) GetProfileByID(ctx context.Context, id string, proplist ...string) (*domain.UserProfile, error) {
	args := []string{"/ip/hotspot/user/profile/print", "?.id=" + id}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotProfileDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseUserProfile(reply.Re[0].Map), nil
}

func (r *profileRepository) GetProfileByName(ctx context.Context, name string, proplist ...string) (*domain.UserProfile, error) {
	args := []string{"/ip/hotspot/user/profile/print", "?name=" + name}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotProfileDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseUserProfile(reply.Re[0].Map), nil
}

func (r *profileRepository) AddProfile(ctx context.Context, profile *domain.UserProfile) (string, error) {
	args := []string{
		"/ip/hotspot/user/profile/add",
		"=name=" + profile.Name,
	}

	if profile.AddressPool != "" {
		args = append(args, "=address-pool="+profile.AddressPool)
	}
	if profile.SharedUsers > 0 {
		args = append(args, "=shared-users="+utils.FormatInt(int64(profile.SharedUsers)))
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if profile.QueueType != "" {
		args = append(args, "=queue-type="+profile.QueueType)
	}
	if profile.StatusAutorefresh != "" {
		args = append(args, "=status-autorefresh="+profile.StatusAutorefresh)
	}
	if profile.OnLogin != "" {
		args = append(args, "=on-login="+profile.OnLogin)
	}
	if profile.OnLogout != "" {
		args = append(args, "=on-logout="+profile.OnLogout)
	}
	if profile.OpenStatusPage != "" {
		args = append(args, "=open-status-page="+profile.OpenStatusPage)
	}
	if profile.TransparentProxy {
		args = append(args, "=transparent-proxy=yes")
	}
	if profile.Advertise {
		args = append(args, "=advertise=yes")
	}
	if profile.AdvertiseInterval != "" {
		args = append(args, "=advertise-interval="+profile.AdvertiseInterval)
	}
	if profile.AdvertiseTimeout != "" {
		args = append(args, "=advertise-timeout="+profile.AdvertiseTimeout)
	}
	if profile.AdvertiseURL != "" {
		args = append(args, "=advertise-url="+profile.AdvertiseURL)
	}
	if profile.IdleTimeout != "" {
		args = append(args, "=idle-timeout="+profile.IdleTimeout)
	}
	if profile.SessionTimeout != "" {
		args = append(args, "=session-timeout="+profile.SessionTimeout)
	}
	if profile.KeepaliveTimeout != "" {
		args = append(args, "=keepalive-timeout="+profile.KeepaliveTimeout)
	}
	if profile.MacCookieTimeout != "" {
		args = append(args, "=mac-cookie-timeout="+profile.MacCookieTimeout)
	}
	if profile.AddMacCookie {
		args = append(args, "=add-mac-cookie=yes")
	}
	if profile.AddressList != "" {
		args = append(args, "=address-list="+profile.AddressList)
	}
	if profile.IncomingFilter != "" {
		args = append(args, "=incoming-filter="+profile.IncomingFilter)
	}
	if profile.IncomingPacketMark != "" {
		args = append(args, "=incoming-packet-mark="+profile.IncomingPacketMark)
	}
	if profile.OutgoingFilter != "" {
		args = append(args, "=outgoing-filter="+profile.OutgoingFilter)
	}
	if profile.OutgoingPacketMark != "" {
		args = append(args, "=outgoing-packet-mark="+profile.OutgoingPacketMark)
	}
	if profile.InsertQueueBefore != "" {
		args = append(args, "=insert-queue-before="+profile.InsertQueueBefore)
	}

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return "", err
	}
	if len(reply.Re) > 0 {
		return reply.Re[0].Map["ret"], nil
	}
	return "", nil
}

func (r *profileRepository) UpdateProfile(ctx context.Context, id string, profile *domain.UserProfile) error {
	args := []string{"/ip/hotspot/user/profile/set", "=.id=" + id}

	if profile.Name != "" {
		args = append(args, "=name="+profile.Name)
	}
	if profile.AddressPool != "" {
		args = append(args, "=address-pool="+profile.AddressPool)
	}
	if profile.SharedUsers > 0 {
		args = append(args, "=shared-users="+utils.FormatInt(int64(profile.SharedUsers)))
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if profile.QueueType != "" {
		args = append(args, "=queue-type="+profile.QueueType)
	}
	if profile.StatusAutorefresh != "" {
		args = append(args, "=status-autorefresh="+profile.StatusAutorefresh)
	}
	if profile.OnLogin != "" {
		args = append(args, "=on-login="+profile.OnLogin)
	}
	if profile.OnLogout != "" {
		args = append(args, "=on-logout="+profile.OnLogout)
	}
	if profile.OpenStatusPage != "" {
		args = append(args, "=open-status-page="+profile.OpenStatusPage)
	}
	if profile.TransparentProxy {
		args = append(args, "=transparent-proxy=yes")
	} else {
		args = append(args, "=transparent-proxy=no")
	}
	if profile.Advertise {
		args = append(args, "=advertise=yes")
	} else {
		args = append(args, "=advertise=no")
	}
	if profile.AdvertiseInterval != "" {
		args = append(args, "=advertise-interval="+profile.AdvertiseInterval)
	}
	if profile.AdvertiseTimeout != "" {
		args = append(args, "=advertise-timeout="+profile.AdvertiseTimeout)
	}
	if profile.AdvertiseURL != "" {
		args = append(args, "=advertise-url="+profile.AdvertiseURL)
	}
	if profile.IdleTimeout != "" {
		args = append(args, "=idle-timeout="+profile.IdleTimeout)
	}
	if profile.SessionTimeout != "" {
		args = append(args, "=session-timeout="+profile.SessionTimeout)
	}
	if profile.KeepaliveTimeout != "" {
		args = append(args, "=keepalive-timeout="+profile.KeepaliveTimeout)
	}
	if profile.MacCookieTimeout != "" {
		args = append(args, "=mac-cookie-timeout="+profile.MacCookieTimeout)
	}
	if profile.AddMacCookie {
		args = append(args, "=add-mac-cookie=yes")
	} else {
		args = append(args, "=add-mac-cookie=no")
	}
	if profile.AddressList != "" {
		args = append(args, "=address-list="+profile.AddressList)
	}
	if profile.IncomingFilter != "" {
		args = append(args, "=incoming-filter="+profile.IncomingFilter)
	}
	if profile.IncomingPacketMark != "" {
		args = append(args, "=incoming-packet-mark="+profile.IncomingPacketMark)
	}
	if profile.OutgoingFilter != "" {
		args = append(args, "=outgoing-filter="+profile.OutgoingFilter)
	}
	if profile.OutgoingPacketMark != "" {
		args = append(args, "=outgoing-packet-mark="+profile.OutgoingPacketMark)
	}
	if profile.InsertQueueBefore != "" {
		args = append(args, "=insert-queue-before="+profile.InsertQueueBefore)
	}

	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *profileRepository) RemoveProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/remove", "=.id="+id)
	return err
}

func (r *profileRepository) DisableProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/disable", "=.id="+id)
	return err
}

func (r *profileRepository) EnableProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/enable", "=.id="+id)
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
