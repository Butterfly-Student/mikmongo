package hotspot

import (
	"context"
	"fmt"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// activeRepository implements ActiveRepository interface
type activeRepository struct {
	client *client.Client
}

// NewActiveRepository creates a new active repository
func NewActiveRepository(c *client.Client) ActiveRepository {
	return &activeRepository{client: c}
}

func parseHotspotActive(m map[string]string) *domain.HotspotActive {
	return &domain.HotspotActive{
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
		BytesIn:          utils.ParseInt(m["bytes-in"]),
		BytesOut:         utils.ParseInt(m["bytes-out"]),
		LimitBytesIn:     utils.ParseInt(m["limit-bytes-in"]),
		LimitBytesOut:    utils.ParseInt(m["limit-bytes-out"]),
		LimitBytesTotal:  utils.ParseInt(m["limit-bytes-total"]),
	}
}

func (r *activeRepository) GetActive(ctx context.Context) ([]*domain.HotspotActive, error) {
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

func (r *activeRepository) GetActiveCount(ctx context.Context) (int, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/active/print", "=count-only=")
	if err != nil {
		return 0, err
	}
	if len(reply.Re) > 0 {
		return int(utils.ParseInt(reply.Re[0].Map["ret"])), nil
	}
	return 0, nil
}

func (r *activeRepository) RemoveActive(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/active/remove", "=.id="+id)
	return err
}

func (r *activeRepository) RemoveActives(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveActive(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *activeRepository) ListenActive(ctx context.Context, resultChan chan<- []*domain.HotspotActive) (func() error, error) {
	args := []string{"/ip/hotspot/active/print", "=follow="}
	// Add proplist - menggunakan default untuk Active
	pl := utils.BuildProplist(ProplistHotspotActiveDefault)
	args = utils.AppendProplist(args, pl)

	listenReply, err := r.client.ListenArgsContext(ctx, args)
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

func (r *activeRepository) ListenInactive(ctx context.Context, resultChan chan<- []*domain.HotspotUser) (func() error, error) {
	// Listen to users with proplist default
	usersArgs := []string{"/ip/hotspot/user/print", "=follow="}
	usersPl := utils.BuildProplist(ProplistHotspotUserDefault)
	usersArgs = utils.AppendProplist(usersArgs, usersPl)
	usersListen, err := r.client.ListenArgsContext(ctx, usersArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to start hotspot users listen: %w", err)
	}

	// Listen to active with proplist default
	activeArgs := []string{"/ip/hotspot/active/print", "=follow="}
	activePl := utils.BuildProplist(ProplistHotspotUserDefault)
	activeArgs = utils.AppendProplist(activeArgs, activePl)
	activeListen, err := r.client.ListenArgsContext(ctx, activeArgs)
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
