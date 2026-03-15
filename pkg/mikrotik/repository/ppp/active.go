package ppp

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

func parsePPPActive(m map[string]string) *domain.PPPActive {
	return &domain.PPPActive{
		Name:     m["name"],
		Service:  m["service"],
		CallerID: m["caller-id"],
		Encoding: m["encoding"],
		Address:  m["address"],
		Uptime:   m["uptime"],
	}
}

func (r *activeRepository) GetActive(ctx context.Context, service string) ([]*domain.PPPActive, error) {
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

func (r *activeRepository) GetActiveByID(ctx context.Context, id string) (*domain.PPPActive, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/active/print", "?.id="+id)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPActive(reply.Re[0].Map), nil
}

func (r *activeRepository) GetActiveByName(ctx context.Context, name string) (*domain.PPPActive, error) {
	reply, err := r.client.RunContext(ctx, "/ppp/active/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPActive(reply.Re[0].Map), nil
}

func (r *activeRepository) RemoveActive(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/active/remove", "=.id="+id)
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

func (r *activeRepository) ListenActive(ctx context.Context, resultChan chan<- []*domain.PPPActive) (func() error, error) {
	args := []string{"/ppp/active/print", "=follow="}
	listenReply, err := r.client.ListenArgsContext(ctx, args)
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

func (r *activeRepository) ListenInactive(ctx context.Context, resultChan chan<- []*domain.PPPSecret) (func() error, error) {
	// Listen to secrets with proplist default
	secretsArgs := []string{"/ppp/secret/print", "=follow="}
	secretsPl := utils.BuildProplist(ProplistPPPSecretDefault)
	secretsArgs = utils.AppendProplist(secretsArgs, secretsPl)
	secretsListen, err := r.client.ListenArgsContext(ctx, secretsArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to start ppp secrets listen: %w", err)
	}

	// Listen to active
	activeArgs := []string{"/ppp/active/print", "=follow="}
	activeListen, err := r.client.ListenArgsContext(ctx, activeArgs)
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
