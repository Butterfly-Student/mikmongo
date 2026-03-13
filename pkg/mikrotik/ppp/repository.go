package ppp

import (
	"context"
	"fmt"
	"strconv"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Repository handles PPP data access via RouterOS API
type Repository struct {
	client *client.Client
}

// NewRepository creates a new PPP repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

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

func parsePPPSecret(m map[string]string) *domain.PPPSecret {
	return &domain.PPPSecret{
		ID:                   m[".id"],
		Name:                 m["name"],
		Password:             m["password"],
		Profile:              m["profile"],
		Service:              m["service"],
		Disabled:             parseBool(m["disabled"]),
		CallerID:             m["caller-id"],
		LocalAddress:         m["local-address"],
		RemoteAddress:        m["remote-address"],
		Routes:               m["routes"],
		Comment:              m["comment"],
		LimitBytesIn:         parseInt(m["limit-bytes-in"]),
		LimitBytesOut:        parseInt(m["limit-bytes-out"]),
		LastLoggedOut:        m["last-logged-out"],
		LastCallerID:         m["last-caller-id"],
		LastDisconnectReason: m["last-disconnect-reason"],
	}
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
		OnlyOne:            parseBool(m["only-one"]),
		Comment:            m["comment"],
		RateLimit:          m["rate-limit"],
		ParentQueue:        m["parent-queue"],
		QueueType:          m["queue-type"],
		UseCompression:     parseBool(m["use-compression"]),
		UseEncryption:      parseBool(m["use-encryption"]),
		UseMPLS:            parseBool(m["use-mpls"]),
		UseUPnP:            parseBool(m["use-upnp"]),
		Bridge:             m["bridge"],
		AddressList:        m["address-list"],
		InterfaceList:      m["interface-list"],
		OnUp:               m["on-up"],
		OnDown:             m["on-down"],
		ChangeTCPMSS:       parseBool(m["change-tcp-mss"]),
		IncomingFilter:     m["incoming-filter"],
		OutgoingFilter:     m["outgoing-filter"],
		InsertQueueBefore:  m["insert-queue-before"],
		WinsServer:         m["wins-server"],
		BridgeHorizon:      m["bridge-horizon"],
		BridgeLearning:     parseBool(m["bridge-learning"]),
		BridgePathCost:     int(parseInt(m["bridge-path-cost"])),
		BridgePortPriority: int(parseInt(m["bridge-port-priority"])),
	}
}

func parsePPPActive(m map[string]string) *domain.PPPActive {
	bytesIn, bytesOut := client.SplitSlashValue(m["bytes"])
	packetsIn, packetsOut := client.SplitSlashValue(m["packets"])
	return &domain.PPPActive{
		ID:            m[".id"],
		Name:          m["name"],
		Service:       m["service"],
		CallerID:      m["caller-id"],
		Address:       m["address"],
		Uptime:        m["uptime"],
		SessionID:     m["session-id"],
		Encoding:      m["encoding"],
		BytesIn:       bytesIn,
		BytesOut:      bytesOut,
		PacketsIn:     packetsIn,
		PacketsOut:    packetsOut,
		LimitBytesIn:  parseInt(m["limit-bytes-in"]),
		LimitBytesOut: parseInt(m["limit-bytes-out"]),
	}
}

// ─── PPP Secrets ──────────────────────────────────────────────────────────────

func (r *Repository) GetPPPSecrets(ctx context.Context, profile string) ([]*domain.PPPSecret, error) {
	args := []string{"/ppp/secret/print"}
	if profile != "" {
		args = append(args, "?profile="+profile)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	secrets := make([]*domain.PPPSecret, 0, len(reply.Re))
	for _, re := range reply.Re {
		secrets = append(secrets, parsePPPSecret(re.Map))
	}
	return secrets, nil
}

func (r *Repository) GetPPPSecretByID(ctx context.Context, id string) (*domain.PPPSecret, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/secret/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPSecret(reply.Re[0].Map), nil
}

func (r *Repository) GetPPPSecretByName(ctx context.Context, name string) (*domain.PPPSecret, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/secret/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPSecret(reply.Re[0].Map), nil
}

func (r *Repository) AddPPPSecret(ctx context.Context, secret *domain.PPPSecret) error {
	args := []string{"/ppp/secret/add", "=name=" + secret.Name}
	if secret.Password != "" {
		args = append(args, "=password="+secret.Password)
	}
	if secret.Profile != "" {
		args = append(args, "=profile="+secret.Profile)
	}
	if secret.Service != "" {
		args = append(args, "=service="+secret.Service)
	}
	if secret.CallerID != "" {
		args = append(args, "=caller-id="+secret.CallerID)
	}
	if secret.LocalAddress != "" {
		args = append(args, "=local-address="+secret.LocalAddress)
	}
	if secret.RemoteAddress != "" {
		args = append(args, "=remote-address="+secret.RemoteAddress)
	}
	if secret.Routes != "" {
		args = append(args, "=routes="+secret.Routes)
	}
	if secret.Comment != "" {
		args = append(args, "=comment="+secret.Comment)
	}
	if secret.LimitBytesIn > 0 {
		args = append(args, "=limit-bytes-in="+formatInt(secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+formatInt(secret.LimitBytesOut))
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) UpdatePPPSecret(ctx context.Context, id string, secret *domain.PPPSecret) error {
	args := []string{"/ppp/secret/set", "=.id=" + id}
	if secret.Name != "" {
		args = append(args, "=name="+secret.Name)
	}
	if secret.Password != "" {
		args = append(args, "=password="+secret.Password)
	}
	if secret.Profile != "" {
		args = append(args, "=profile="+secret.Profile)
	}
	if secret.Service != "" {
		args = append(args, "=service="+secret.Service)
	}
	if secret.CallerID != "" {
		args = append(args, "=caller-id="+secret.CallerID)
	}
	if secret.LocalAddress != "" {
		args = append(args, "=local-address="+secret.LocalAddress)
	}
	if secret.RemoteAddress != "" {
		args = append(args, "=remote-address="+secret.RemoteAddress)
	}
	if secret.Routes != "" {
		args = append(args, "=routes="+secret.Routes)
	}
	if secret.Comment != "" {
		args = append(args, "=comment="+secret.Comment)
	}
	if secret.LimitBytesIn > 0 {
		args = append(args, "=limit-bytes-in="+formatInt(secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+formatInt(secret.LimitBytesOut))
	}
	if secret.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) RemovePPPSecret(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/secret/remove", "=.id="+id)
	return err
}

func (r *Repository) DisablePPPSecret(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/secret/disable", "=.id="+id)
	return err
}

func (r *Repository) EnablePPPSecret(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/secret/enable", "=.id="+id)
	return err
}

func (r *Repository) RemovePPPSecrets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemovePPPSecret(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DisablePPPSecrets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisablePPPSecret(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) EnablePPPSecrets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnablePPPSecret(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ─── PPP Profiles ─────────────────────────────────────────────────────────────

func (r *Repository) GetPPPProfiles(ctx context.Context) ([]*domain.PPPProfile, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/profile/print")
	if err != nil {
		return nil, err
	}
	profiles := make([]*domain.PPPProfile, 0, len(reply.Re))
	for _, re := range reply.Re {
		profiles = append(profiles, parsePPPProfile(re.Map))
	}
	return profiles, nil
}

func (r *Repository) GetPPPProfileByID(ctx context.Context, id string) (*domain.PPPProfile, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/profile/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPProfile(reply.Re[0].Map), nil
}

func (r *Repository) GetPPPProfileByName(ctx context.Context, name string) (*domain.PPPProfile, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/profile/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPProfile(reply.Re[0].Map), nil
}

func (r *Repository) AddPPPProfile(ctx context.Context, profile *domain.PPPProfile) error {
	args := []string{"/ppp/profile/add", "=name=" + profile.Name}
	args = appendPPPProfileArgs(args, profile)
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *Repository) UpdatePPPProfile(ctx context.Context, id string, profile *domain.PPPProfile) error {
	args := []string{"/ppp/profile/set", "=.id=" + id}
	if profile.Name != "" {
		args = append(args, "=name="+profile.Name)
	}
	args = appendPPPProfileArgs(args, profile)
	_, err := r.client.RunArgsContext(ctx, args)
	return err
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
	if p.RateLimit != "" {
		args = append(args, "=rate-limit="+p.RateLimit)
	}
	if p.ParentQueue != "" {
		args = append(args, "=parent-queue="+p.ParentQueue)
	}
	if p.QueueType != "" {
		args = append(args, "=queue-type="+p.QueueType)
	}
	if p.Bridge != "" {
		args = append(args, "=bridge="+p.Bridge)
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

func (r *Repository) RemovePPPProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/profile/remove", "=.id="+id)
	return err
}

func (r *Repository) DisablePPPProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/profile/disable", "=.id="+id)
	return err
}

func (r *Repository) EnablePPPProfile(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/profile/enable", "=.id="+id)
	return err
}

func (r *Repository) RemovePPPProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemovePPPProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DisablePPPProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisablePPPProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) EnablePPPProfiles(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnablePPPProfile(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ─── PPP Active ───────────────────────────────────────────────────────────────

func (r *Repository) GetPPPActive(ctx context.Context, service string) ([]*domain.PPPActive, error) {
	args := []string{"/ppp/active/print"}
	if service != "" {
		args = append(args, "?service="+service)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	active := make([]*domain.PPPActive, 0, len(reply.Re))
	for _, re := range reply.Re {
		active = append(active, parsePPPActive(re.Map))
	}
	return active, nil
}

func (r *Repository) GetPPPActiveByID(ctx context.Context, id string) (*domain.PPPActive, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/active/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPActive(reply.Re[0].Map), nil
}

func (r *Repository) RemovePPPActive(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/active/remove", "=.id="+id)
	return err
}

func (r *Repository) RemovePPPActives(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemovePPPActive(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// ─── Streaming ────────────────────────────────────────────────────────────────

// ListenPPPActive streams active PPP sessions. resultChan is closed when ctx is done.
func (r *Repository) ListenPPPActive(
	ctx context.Context,
	resultChan chan<- []*domain.PPPActive,
) (func() error, error) {
	listenReply, err := r.client.ListenArgsContext(ctx, []string{"/ppp/active/print", "=follow="})
	if err != nil {
		return nil, fmt.Errorf("failed to start ppp active listen: %w", err)
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
				case resultChan <- []*domain.PPPActive{parsePPPActive(sentence.Map)}:
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

// ListenPPPInactive streams inactive PPP secrets (secrets not in active sessions).
func (r *Repository) ListenPPPInactive(
	ctx context.Context,
	resultChan chan<- []*domain.PPPSecret,
) (func() error, error) {
	secretsListen, err := r.client.ListenArgsContext(ctx, []string{"/ppp/secret/print", "=follow="})
	if err != nil {
		return nil, fmt.Errorf("failed to start ppp secrets listen: %w", err)
	}

	activeListen, err := r.client.ListenArgsContext(ctx, []string{"/ppp/active/print", "=follow="})
	if err != nil {
		secretsListen.Cancel() //nolint:errcheck
		return nil, fmt.Errorf("failed to start ppp active listen for inactive: %w", err)
	}

	secretsBatches := client.ListenBatches(ctx, secretsListen.Chan(), client.BatchDebounce)
	activeBatches := client.ListenBatches(ctx, activeListen.Chan(), client.BatchDebounce)

	go func() {
		defer close(resultChan)

		var latestSecrets []*domain.PPPSecret
		var latestActive []*domain.PPPActive

		sendDiff := func() {
			if latestSecrets == nil {
				return
			}
			activeSet := make(map[string]struct{}, len(latestActive))
			for _, a := range latestActive {
				activeSet[a.Name] = struct{}{}
			}
			inactive := make([]*domain.PPPSecret, 0)
			for _, s := range latestSecrets {
				if _, ok := activeSet[s.Name]; !ok {
					inactive = append(inactive, s)
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
			case batch, ok := <-secretsBatches:
				if !ok {
					return
				}
				latestSecrets = make([]*domain.PPPSecret, 0, len(batch))
				for _, s := range batch {
					latestSecrets = append(latestSecrets, parsePPPSecret(s.Map))
				}
				sendDiff()
			case batch, ok := <-activeBatches:
				if !ok {
					return
				}
				latestActive = make([]*domain.PPPActive, 0, len(batch))
				for _, s := range batch {
					latestActive = append(latestActive, parsePPPActive(s.Map))
				}
				sendDiff()
			}
		}
	}()

	return func() error {
		_, err1 := secretsListen.Cancel()
		_, err2 := activeListen.Cancel()
		if err1 != nil {
			return err1
		}
		return err2
	}, nil
}
