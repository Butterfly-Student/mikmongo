package ppp

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// secretRepository implements SecretRepository interface
type secretRepository struct {
	client *client.Client
}

// NewSecretRepository creates a new secret repository
func NewSecretRepository(c *client.Client) SecretRepository {
	return &secretRepository{client: c}
}

func parsePPPSecret(m map[string]string) *domain.PPPSecret {
	return &domain.PPPSecret{
		ID:                   m[".id"],
		Name:                 m["name"],
		Password:             m["password"],
		Profile:              m["profile"],
		Service:              m["service"],
		Disabled:             utils.ParseBool(m["disabled"]),
		CallerID:             m["caller-id"],
		LocalAddress:         m["local-address"],
		RemoteAddress:        m["remote-address"],
		Routes:               m["routes"],
		Comment:              m["comment"],
		LimitBytesIn:         utils.ParseInt(m["limit-bytes-in"]),
		LimitBytesOut:        utils.ParseInt(m["limit-bytes-out"]),
		LastLoggedOut:        m["last-logged-out"],
		LastCallerID:         m["last-caller-id"],
		LastDisconnectReason: m["last-disconnect-reason"],
	}
}

func (r *secretRepository) GetSecrets(ctx context.Context, profile string, proplist ...string) ([]*domain.PPPSecret, error) {
	args := []string{"/ppp/secret/print"}
	if profile != "" {
		args = append(args, "?profile="+profile)
	}
	// Add proplist
	pl := utils.BuildProplist(ProplistPPPSecretDefault, proplist...)
	args = utils.AppendProplist(args, pl)

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

func (r *secretRepository) GetSecretByID(ctx context.Context, id string, proplist ...string) (*domain.PPPSecret, error) {
	args := []string{"/ppp/secret/print", "?.id=" + id}
	// Add proplist
	pl := utils.BuildProplist(ProplistPPPSecretDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPSecret(reply.Re[0].Map), nil
}

func (r *secretRepository) GetSecretByName(ctx context.Context, name string, proplist ...string) (*domain.PPPSecret, error) {
	args := []string{"/ppp/secret/print", "?name=" + name}
	// Add proplist
	pl := utils.BuildProplist(ProplistPPPSecretDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parsePPPSecret(reply.Re[0].Map), nil
}

func (r *secretRepository) AddSecret(ctx context.Context, secret *domain.PPPSecret) error {
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
		args = append(args, "=limit-bytes-in="+utils.FormatInt(secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+utils.FormatInt(secret.LimitBytesOut))
	}
	if secret.IPv6Routes != "" {
		args = append(args, "=ipv6-routes="+secret.IPv6Routes)
	}
	if secret.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *secretRepository) UpdateSecret(ctx context.Context, id string, secret *domain.PPPSecret) error {
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
		args = append(args, "=limit-bytes-in="+utils.FormatInt(secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+utils.FormatInt(secret.LimitBytesOut))
	}
	if secret.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *secretRepository) RemoveSecret(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/secret/remove", "=.id="+id)
	return err
}

func (r *secretRepository) DisableSecret(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/secret/disable", "=.id="+id)
	return err
}

func (r *secretRepository) EnableSecret(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ppp/secret/enable", "=.id="+id)
	return err
}

func (r *secretRepository) RemoveSecrets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveSecret(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *secretRepository) DisableSecrets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisableSecret(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *secretRepository) EnableSecrets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnableSecret(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// PrintSecretsRaw executes /ppp/secret/print with custom arguments and returns raw data
func (r *secretRepository) PrintSecretsRaw(ctx context.Context, args ...string) ([]map[string]string, error) {
	cmdArgs := append([]string{"/ppp/secret/print"}, args...)
	reply, err := r.client.RunArgsContext(ctx, cmdArgs)
	if err != nil {
		return nil, err
	}
	results := make([]map[string]string, 0, len(reply.Re))
	for _, re := range reply.Re {
		results = append(results, re.Map)
	}
	return results, nil
}

// ListenSecretsRaw executes /ppp/secret/print with follow and returns raw data stream
func (r *secretRepository) ListenSecretsRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error) {
	cmdArgs := append([]string{"/ppp/secret/print"}, args...)
	listenReply, err := r.client.ListenArgsContext(ctx, cmdArgs)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(resultChan)
		for {
			select {
			case <-ctx.Done():
				listenReply.Cancel()
				return
			case sentence, ok := <-listenReply.Chan():
				if !ok {
					return
				}
				select {
				case resultChan <- sentence.Map:
				case <-ctx.Done():
					listenReply.Cancel()
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
