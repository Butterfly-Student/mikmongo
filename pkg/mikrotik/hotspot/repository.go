package hotspot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/mikrotik/script"
)

// Repository handles Hotspot data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new Hotspot repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseBool(s string) bool {
	return s == "true" || s == "yes"
}

func formatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}

// ─── parse helpers ────────────────────────────────────────────────────────────

func parseHotspotUser(re map[string]string) *domain.HotspotUser {
	return &domain.HotspotUser{
		ID:              re[".id"],
		Server:          re["server"],
		Name:            re["name"],
		Password:        re["password"],
		Profile:         re["profile"],
		MACAddress:      re["mac-address"],
		IPAddress:       re["address"],
		Uptime:          re["uptime"],
		BytesIn:         parseInt(re["bytes-in"]),
		BytesOut:        parseInt(re["bytes-out"]),
		PacketsIn:       parseInt(re["packets-in"]),
		PacketsOut:      parseInt(re["packets-out"]),
		LimitUptime:     re["limit-uptime"],
		LimitBytesIn:    parseInt(re["limit-bytes-in"]),
		LimitBytesOut:   parseInt(re["limit-bytes-out"]),
		LimitBytesTotal: parseInt(re["limit-bytes-total"]),
		Comment:         re["comment"],
		Disabled:        parseBool(re["disabled"]),
		Email:           re["email"],
	}
}

func parseHotspotActive(m map[string]string) *domain.HotspotActive {
	return &domain.HotspotActive{
		ID:               m[".id"],
		Server:           m["server"],
		User:             m["user"],
		Address:          m["address"],
		MACAddress:       m["mac-address"],
		LoginBy:          m["login-by"],
		Uptime:           m["uptime"],
		SessionTimeLeft:  m["session-time-left"],
		IdleTime:         m["idle-time"],
		IdleTimeout:      m["idle-timeout"],
		KeepaliveTimeout: m["keepalive-timeout"],
		BytesIn:          parseInt(m["bytes-in"]),
		BytesOut:         parseInt(m["bytes-out"]),
		PacketsIn:        parseInt(m["packets-in"]),
		PacketsOut:       parseInt(m["packets-out"]),
		Radius:           parseBool(m["radius"]),
	}
}

// ─── Hotspot Users ────────────────────────────────────────────────────────────

func (r *Repository) GetHotspotUsers(ctx context.Context, profile string) ([]*domain.HotspotUser, error) {
	args := []string{"/ip/hotspot/user/print"}
	if profile != "" {
		args = append(args, "?profile="+profile)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	users := make([]*domain.HotspotUser, 0, len(reply.Re))
	for _, re := range reply.Re {
		users = append(users, parseHotspotUser(re.Map))
	}
	return users, nil
}

func (r *Repository) GetHotspotUsersByComment(ctx context.Context, comment string) ([]*domain.HotspotUser, error) {
	args := []string{"/ip/hotspot/user/print"}
	if comment != "" {
		args = append(args, "?comment="+comment)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	users := make([]*domain.HotspotUser, 0, len(reply.Re))
	for _, re := range reply.Re {
		users = append(users, parseHotspotUser(re.Map))
	}
	return users, nil
}

func (r *Repository) GetHotspotUserByID(ctx context.Context, id string) (*domain.HotspotUser, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseHotspotUser(reply.Re[0].Map), nil
}

func (r *Repository) GetHotspotUserByName(ctx context.Context, name string) (*domain.HotspotUser, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseHotspotUser(reply.Re[0].Map), nil
}

func (r *Repository) AddHotspotUser(ctx context.Context, user *domain.HotspotUser) (string, error) {
	args := []string{
		"/ip/hotspot/user/add",
		"=name=" + user.Name,
		"=profile=" + user.Profile,
		"=disabled=no",
	}
	if user.Server != "" && user.Server != "all" {
		args = append(args, "=server="+user.Server)
	}
	if user.Password != "" {
		args = append(args, "=password="+user.Password)
	}
	if user.MACAddress != "" {
		args = append(args, "=mac-address="+user.MACAddress)
	}
	if user.LimitUptime != "" {
		args = append(args, "=limit-uptime="+user.LimitUptime)
	}
	if user.LimitBytesTotal > 0 {
		args = append(args, "=limit-bytes-total="+formatInt(user.LimitBytesTotal))
	}
	if user.Comment != "" {
		args = append(args, "=comment="+user.Comment)
	}
	if user.Email != "" {
		args = append(args, "=email="+user.Email)
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

func (r *Repository) UpdateHotspotUser(ctx context.Context, id string, user *domain.HotspotUser) error {
	args := []string{"/ip/hotspot/user/set", "=.id=" + id}
	if user.Name != "" {
		args = append(args, "=name="+user.Name)
	}
	if user.Password != "" {
		args = append(args, "=password="+user.Password)
	}
	if user.Profile != "" {
		args = append(args, "=profile="+user.Profile)
	}
	if user.Server != "" {
		args = append(args, "=server="+user.Server)
	}
	if user.MACAddress != "" {
		args = append(args, "=mac-address="+user.MACAddress)
	}
	if user.LimitUptime != "" {
		args = append(args, "=limit-uptime="+user.LimitUptime)
	}
	if user.LimitBytesTotal > 0 {
		args = append(args, "=limit-bytes-total="+formatInt(user.LimitBytesTotal))
	}
	if user.Comment != "" {
		args = append(args, "=comment="+user.Comment)
	}
	if user.Email != "" {
		args = append(args, "=email="+user.Email)
	}
	if user.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) RemoveHotspotUser(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/remove", "=.id="+id)
	return err
}

func (r *Repository) RemoveHotspotUsersByComment(ctx context.Context, comment string) error {
	users, err := r.GetHotspotUsersByComment(ctx, comment)
	if err != nil {
		return err
	}
	for _, user := range users {
		_ = r.RemoveHotspotUser(ctx, user.ID)
	}
	return nil
}

func (r *Repository) ResetHotspotUserCounters(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/reset-counters", "=.id="+id)
	return err
}

func (r *Repository) GetHotspotUsersCount(ctx context.Context) (int, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/print", "=count-only=")
	if err != nil {
		return 0, err
	}
	if len(reply.Re) > 0 {
		count := int(parseInt(reply.Re[0].Map["ret"])) - 1
		if count < 0 {
			count = 0
		}
		return count, nil
	}
	return 0, nil
}

func (r *Repository) RemoveHotspotUsers(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveHotspotUser(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DisableHotspotUser(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/disable", "=.id="+id)
	return err
}

func (r *Repository) EnableHotspotUser(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/enable", "=.id="+id)
	return err
}

func (r *Repository) DisableHotspotUsers(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisableHotspotUser(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) EnableHotspotUsers(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnableHotspotUser(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ResetHotspotUserCountersMultiple(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.ResetHotspotUserCounters(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ─── Hotspot Profiles ─────────────────────────────────────────────────────────

func (r *Repository) GetUserProfiles(ctx context.Context) ([]*domain.UserProfile, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/print")
	if err != nil {
		return nil, err
	}

	generator := script.NewOnLoginGenerator()
	profiles := make([]*domain.UserProfile, 0, len(reply.Re))

	for _, re := range reply.Re {
		profile := &domain.UserProfile{
			ID:                re.Map[".id"],
			Name:              re.Map["name"],
			AddressPool:       re.Map["address-pool"],
			SharedUsers:       int(parseInt(re.Map["shared-users"])),
			RateLimit:         re.Map["rate-limit"],
			ParentQueue:       re.Map["parent-queue"],
			StatusAutorefresh: re.Map["status-autorefresh"],
			OnLogin:           re.Map["on-login"],
			OnLogout:          re.Map["on-logout"],
			OnUp:              re.Map["on-up"],
			OnDown:            re.Map["on-down"],
			TransparentProxy:  parseBool(re.Map["transparent-proxy"]),
			OpenStatusPage:    re.Map["open-status-page"],
			Advertise:         parseBool(re.Map["advertise"]),
			AdvertiseInterval: re.Map["advertise-interval"],
			AdvertiseTimeout:  re.Map["advertise-timeout"],
			AdvertiseURL:      re.Map["advertise-url"],
		}

		if profile.OnLogin != "" {
			parsed := generator.Parse(profile.OnLogin)
			profile.ExpireMode = parsed.ExpireMode
			profile.Validity = parsed.Validity
			profile.Price = parsed.Price
			profile.SellingPrice = parsed.SellingPrice
			profile.LockUser = parsed.LockUser
			profile.LockServer = parsed.LockServer
		}

		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (r *Repository) GetUserProfileByID(ctx context.Context, id string) (*domain.UserProfile, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	re := reply.Re[0]
	generator := script.NewOnLoginGenerator()
	profile := &domain.UserProfile{
		ID:                re.Map[".id"],
		Name:              re.Map["name"],
		AddressPool:       re.Map["address-pool"],
		SharedUsers:       int(parseInt(re.Map["shared-users"])),
		RateLimit:         re.Map["rate-limit"],
		ParentQueue:       re.Map["parent-queue"],
		StatusAutorefresh: re.Map["status-autorefresh"],
		OnLogin:           re.Map["on-login"],
		OnLogout:          re.Map["on-logout"],
		OnUp:              re.Map["on-up"],
		OnDown:            re.Map["on-down"],
		TransparentProxy:  parseBool(re.Map["transparent-proxy"]),
		OpenStatusPage:    re.Map["open-status-page"],
	}
	if profile.OnLogin != "" {
		parsed := generator.Parse(profile.OnLogin)
		profile.ExpireMode = parsed.ExpireMode
		profile.Validity = parsed.Validity
		profile.Price = parsed.Price
		profile.SellingPrice = parsed.SellingPrice
		profile.LockUser = parsed.LockUser
		profile.LockServer = parsed.LockServer
	}
	return profile, nil
}

func (r *Repository) GetUserProfileByName(ctx context.Context, name string) (*domain.UserProfile, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	re := reply.Re[0]
	generator := script.NewOnLoginGenerator()
	profile := &domain.UserProfile{
		ID:          re.Map[".id"],
		Name:        re.Map["name"],
		AddressPool: re.Map["address-pool"],
		SharedUsers: int(parseInt(re.Map["shared-users"])),
		RateLimit:   re.Map["rate-limit"],
		ParentQueue: re.Map["parent-queue"],
		OnLogin:     re.Map["on-login"],
	}
	if profile.OnLogin != "" {
		parsed := generator.Parse(profile.OnLogin)
		profile.ExpireMode = parsed.ExpireMode
		profile.Validity = parsed.Validity
		profile.Price = parsed.Price
		profile.SellingPrice = parsed.SellingPrice
		profile.LockUser = parsed.LockUser
		profile.LockServer = parsed.LockServer
	}
	return profile, nil
}

func (r *Repository) AddUserProfile(ctx context.Context, profile *domain.UserProfile) (string, error) {
	var onLoginScript string
	if profile.ExpireMode != "" {
		generator := script.NewOnLoginGenerator()
		req := &domain.ProfileRequest{
			Name:         profile.Name,
			ExpireMode:   profile.ExpireMode,
			Validity:     profile.Validity,
			Price:        profile.Price,
			SellingPrice: profile.SellingPrice,
			LockUser:     profile.LockUser,
			LockServer:   profile.LockServer,
		}
		onLoginScript = generator.Generate(req)
	}

	args := []string{
		"/ip/hotspot/user/profile/add",
		"=name=" + profile.Name,
		"=status-autorefresh=1m",
	}
	if profile.AddressPool != "" && profile.AddressPool != "none" {
		args = append(args, "=address-pool="+profile.AddressPool)
	}
	if profile.SharedUsers > 0 {
		args = append(args, "=shared-users="+formatInt(int64(profile.SharedUsers)))
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" && profile.ParentQueue != "none" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if onLoginScript != "" {
		scriptStr := strings.ReplaceAll(onLoginScript, "; ", "\n")
		args = append(args, "=on-login="+scriptStr)
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

func (r *Repository) UpdateUserProfile(ctx context.Context, id string, profile *domain.UserProfile) error {
	var onLoginScript string
	if profile.ExpireMode != "" {
		generator := script.NewOnLoginGenerator()
		req := &domain.ProfileRequest{
			Name:         profile.Name,
			ExpireMode:   profile.ExpireMode,
			Validity:     profile.Validity,
			Price:        profile.Price,
			SellingPrice: profile.SellingPrice,
			LockUser:     profile.LockUser,
			LockServer:   profile.LockServer,
		}
		onLoginScript = generator.Generate(req)
	}

	args := []string{"/ip/hotspot/user/profile/set", "=.id=" + id}
	if profile.Name != "" {
		args = append(args, "=name="+profile.Name)
	}
	if profile.AddressPool != "" {
		args = append(args, "=address-pool="+profile.AddressPool)
	}
	if profile.SharedUsers > 0 {
		args = append(args, "=shared-users="+formatInt(int64(profile.SharedUsers)))
	}
	if profile.RateLimit != "" {
		args = append(args, "=rate-limit="+profile.RateLimit)
	}
	if profile.ParentQueue != "" {
		args = append(args, "=parent-queue="+profile.ParentQueue)
	}
	if onLoginScript != "" {
		scriptStr := strings.ReplaceAll(onLoginScript, "; ", "\n")
		args = append(args, "=on-login="+scriptStr)
	}

	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) RemoveUserProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/remove", "=.id="+id)
	return err
}

func (r *Repository) DisableUserProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/disable", "=.id="+id)
	return err
}

func (r *Repository) EnableUserProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/profile/enable", "=.id="+id)
	return err
}

func (r *Repository) RemoveUserProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveUserProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DisableUserProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisableUserProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) EnableUserProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnableUserProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ─── Hotspot Active ───────────────────────────────────────────────────────────

func (r *Repository) GetHotspotActive(ctx context.Context) ([]*domain.HotspotActive, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/active/print")
	if err != nil {
		return nil, err
	}
	active := make([]*domain.HotspotActive, 0, len(reply.Re))
	for _, re := range reply.Re {
		active = append(active, parseHotspotActive(re.Map))
	}
	return active, nil
}

func (r *Repository) GetHotspotActiveCount(ctx context.Context) (int, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/active/print", "=count-only=")
	if err != nil {
		return 0, err
	}
	if len(reply.Re) > 0 {
		return int(parseInt(reply.Re[0].Map["ret"])), nil
	}
	return 0, nil
}

func (r *Repository) RemoveHotspotActive(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/active/remove", "=.id="+id)
	return err
}

func (r *Repository) RemoveHotspotActives(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveHotspotActive(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ─── Hotspot Hosts ────────────────────────────────────────────────────────────

func (r *Repository) GetHotspotHosts(ctx context.Context) ([]*domain.HotspotHost, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/host/print")
	if err != nil {
		return nil, err
	}
	hosts := make([]*domain.HotspotHost, 0, len(reply.Re))
	for _, re := range reply.Re {
		hosts = append(hosts, &domain.HotspotHost{
			ID:           re.Map[".id"],
			MACAddress:   re.Map["mac-address"],
			Address:      re.Map["address"],
			ToAddress:    re.Map["to-address"],
			Server:       re.Map["server"],
			Authorized:   parseBool(re.Map["authorized"]),
			Bypassed:     parseBool(re.Map["bypassed"]),
			Blocked:      parseBool(re.Map["blocked"]),
			FoundBy:      re.Map["found-by"],
			HostDeadTime: re.Map["host-dead-time"],
			Comment:      re.Map["comment"],
		})
	}
	return hosts, nil
}

func (r *Repository) RemoveHotspotHost(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/host/remove", "=.id="+id)
	return err
}

// ─── Hotspot Servers ──────────────────────────────────────────────────────────

func (r *Repository) GetHotspotServers(ctx context.Context) ([]string, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/print")
	if err != nil {
		return nil, err
	}
	servers := []string{"all"}
	for _, re := range reply.Re {
		if name := re.Map["name"]; name != "" {
			servers = append(servers, name)
		}
	}
	return servers, nil
}

// ─── Streaming ────────────────────────────────────────────────────────────────

// ListenHotspotActive streams active hotspot sessions.
func (r *Repository) ListenHotspotActive(
	ctx context.Context,
	resultChan chan<- []*domain.HotspotActive,
) (func() error, error) {
	listenReply, err := r.client.ListenArgsContext(ctx, []string{"/ip/hotspot/active/print", "=follow="})
	if err != nil {
		return nil, fmt.Errorf("failed to start hotspot active listen: %w", err)
	}

	go func() {
		defer close(resultChan)
		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel() //nolint:errcheck
				return
			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}
				select {
				case resultChan <- []*domain.HotspotActive{parseHotspotActive(sentence.Map)}:
				case <-ctx.Done():
					listenReply.Cancel() //nolint:errcheck
					return
				}
			}
		}
	}()

	return func() error {
		_, err := listenReply.Cancel()
		return err
	}, nil
}

// ListenHotspotInactive streams inactive hotspot users.
func (r *Repository) ListenHotspotInactive(
	ctx context.Context,
	resultChan chan<- []*domain.HotspotUser,
) (func() error, error) {
	usersListen, err := r.client.ListenArgsContext(ctx, []string{"/ip/hotspot/user/print", "=follow="})
	if err != nil {
		return nil, fmt.Errorf("failed to start hotspot users listen: %w", err)
	}

	activeListen, err := r.client.ListenArgsContext(ctx, []string{"/ip/hotspot/active/print", "=follow="})
	if err != nil {
		usersListen.Cancel() //nolint:errcheck
		return nil, fmt.Errorf("failed to start hotspot active listen for inactive: %w", err)
	}

	usersBatches := client.ListenBatches(ctx, usersListen.Chan(), client.BatchDebounce)
	activeBatches := client.ListenBatches(ctx, activeListen.Chan(), client.BatchDebounce)

	go func() {
		defer close(resultChan)

		var latestUsers []*domain.HotspotUser
		var latestActive []*domain.HotspotActive

		sendDiff := func() {
			if latestUsers == nil {
				return
			}
			activeSet := make(map[string]struct{}, len(latestActive))
			for _, a := range latestActive {
				activeSet[a.User] = struct{}{}
			}
			inactive := make([]*domain.HotspotUser, 0)
			for _, u := range latestUsers {
				if _, ok := activeSet[u.Name]; !ok {
					inactive = append(inactive, u)
				}
			}
			select {
			case resultChan <- inactive:
			case <-ctx.Done():
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case batch, ok := <-usersBatches:
				if !ok {
					return
				}
				latestUsers = make([]*domain.HotspotUser, 0, len(batch))
				for _, s := range batch {
					latestUsers = append(latestUsers, parseHotspotUser(s.Map))
				}
				sendDiff()
			case batch, ok := <-activeBatches:
				if !ok {
					return
				}
				latestActive = make([]*domain.HotspotActive, 0, len(batch))
				for _, s := range batch {
					latestActive = append(latestActive, parseHotspotActive(s.Map))
				}
				sendDiff()
			}
		}
	}()

	return func() error {
		_, err1 := usersListen.Cancel()
		_, err2 := activeListen.Cancel()
		if err1 != nil {
			return err1
		}
		return err2
	}, nil
}

// ─── IP Binding ───────────────────────────────────────────────────────────────

func parseIPBinding(m map[string]string) *domain.HotspotIPBinding {
	return &domain.HotspotIPBinding{
		ID:         m[".id"],
		MACAddress: m["mac-address"],
		Address:    m["address"],
		Server:     m["server"],
		Type:       m["type"],
		Comment:    m["comment"],
		Disabled:   parseBool(m["disabled"]),
	}
}

func (r *Repository) GetIPBindings(ctx context.Context) ([]*domain.HotspotIPBinding, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/print")
	if err != nil {
		return nil, err
	}
	bindings := make([]*domain.HotspotIPBinding, 0, len(reply.Re))
	for _, re := range reply.Re {
		bindings = append(bindings, parseIPBinding(re.Map))
	}
	return bindings, nil
}

func (r *Repository) AddIPBinding(ctx context.Context, b *domain.HotspotIPBinding) (string, error) {
	args := []string{
		"/ip/hotspot/ip-binding/add",
		"=mac-address=" + b.MACAddress,
		"=type=regular",
	}
	if b.Address != "" {
		args = append(args, "=address="+b.Address)
	}
	if b.Server != "" {
		args = append(args, "=server="+b.Server)
	}
	if b.Type != "" {
		args[2] = "=type=" + b.Type
	}
	if b.Comment != "" {
		args = append(args, "=comment="+b.Comment)
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

func (r *Repository) RemoveIPBinding(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/remove", "=.id="+id)
	return err
}

func (r *Repository) EnableIPBinding(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/enable", "=.id="+id)
	return err
}

func (r *Repository) DisableIPBinding(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/ip-binding/disable", "=.id="+id)
	return err
}

// ─── Stubs for backward compatibility ────────────────────────────────────────

func (r *Repository) GetUsers() ([]domain.HotspotUser, error) {
	return nil, nil
}

func (r *Repository) GetProfiles() ([]domain.HotspotProfile, error) {
	return nil, nil
}

func (r *Repository) AddUser(user *domain.HotspotUser) error {
	return nil
}
