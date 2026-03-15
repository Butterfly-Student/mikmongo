package hotspot

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// userRepository implements UserRepository interface
type userRepository struct {
	client *client.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(c *client.Client) UserRepository {
	return &userRepository{client: c}
}

func parseHotspotUser(re map[string]string) *domain.HotspotUser {
	return &domain.HotspotUser{
		ID:              re[".id"],
		Server:          re["server"],
		Name:            re["name"],
		Password:        re["password"],
		Profile:         re["profile"],
		MACAddress:      re["mac-address"],
		IPAddress:       re["address"],
		Routes:          re["routes"],
		Uptime:          re["uptime"],
		BytesIn:         utils.ParseInt(re["bytes-in"]),
		BytesOut:        utils.ParseInt(re["bytes-out"]),
		PacketsIn:       utils.ParseInt(re["packets-in"]),
		PacketsOut:      utils.ParseInt(re["packets-out"]),
		LimitUptime:     re["limit-uptime"],
		LimitBytesIn:    utils.ParseInt(re["limit-bytes-in"]),
		LimitBytesOut:   utils.ParseInt(re["limit-bytes-out"]),
		LimitBytesTotal: utils.ParseInt(re["limit-bytes-total"]),
		Comment:         re["comment"],
		Disabled:        utils.ParseBool(re["disabled"]),
		Email:           re["email"],
	}
}

func (r *userRepository) GetUsers(ctx context.Context, profile string, proplist ...string) ([]*domain.HotspotUser, error) {
	args := []string{"/ip/hotspot/user/print"}
	if profile != "" {
		args = append(args, "?profile="+profile)
	}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotUserDefault, proplist...)
	args = utils.AppendProplist(args, pl)

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

func (r *userRepository) GetUsersByComment(ctx context.Context, comment string, proplist ...string) ([]*domain.HotspotUser, error) {
	args := []string{"/ip/hotspot/user/print"}
	if comment != "" {
		args = append(args, "?comment="+comment)
	}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotUserDefault, proplist...)
	args = utils.AppendProplist(args, pl)

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

func (r *userRepository) GetUserByID(ctx context.Context, id string, proplist ...string) (*domain.HotspotUser, error) {
	args := []string{"/ip/hotspot/user/print", "?.id=" + id}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotUserDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseHotspotUser(reply.Re[0].Map), nil
}

func (r *userRepository) GetUserByName(ctx context.Context, name string, proplist ...string) (*domain.HotspotUser, error) {
	args := []string{"/ip/hotspot/user/print", "?name=" + name}
	// Add proplist
	pl := utils.BuildProplist(ProplistHotspotUserDefault, proplist...)
	args = utils.AppendProplist(args, pl)

	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseHotspotUser(reply.Re[0].Map), nil
}

func (r *userRepository) AddUser(ctx context.Context, user *domain.HotspotUser) (string, error) {
	args := []string{
		"/ip/hotspot/user/add",
		"=name=" + user.Name,
		"=profile=" + user.Profile,
	}
	if user.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
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
	if user.IPAddress != "" {
		args = append(args, "=address="+user.IPAddress)
	}
	if user.LimitUptime != "" {
		args = append(args, "=limit-uptime="+user.LimitUptime)
	}
	if user.LimitBytesIn > 0 {
		args = append(args, "=limit-bytes-in="+utils.FormatInt(user.LimitBytesIn))
	}
	if user.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+utils.FormatInt(user.LimitBytesOut))
	}
	if user.LimitBytesTotal > 0 {
		args = append(args, "=limit-bytes-total="+utils.FormatInt(user.LimitBytesTotal))
	}
	if user.Comment != "" {
		args = append(args, "=comment="+user.Comment)
	}
	if user.Email != "" {
		args = append(args, "=email="+user.Email)
	}
	if user.Routes != "" {
		args = append(args, "=routes="+user.Routes)
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

func (r *userRepository) UpdateUser(ctx context.Context, id string, user *domain.HotspotUser) error {
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
	if user.IPAddress != "" {
		args = append(args, "=address="+user.IPAddress)
	}
	if user.LimitUptime != "" {
		args = append(args, "=limit-uptime="+user.LimitUptime)
	}
	if user.LimitBytesIn > 0 {
		args = append(args, "=limit-bytes-in="+utils.FormatInt(user.LimitBytesIn))
	}
	if user.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+utils.FormatInt(user.LimitBytesOut))
	}
	if user.LimitBytesTotal > 0 {
		args = append(args, "=limit-bytes-total="+utils.FormatInt(user.LimitBytesTotal))
	}
	if user.Comment != "" {
		args = append(args, "=comment="+user.Comment)
	}
	if user.Email != "" {
		args = append(args, "=email="+user.Email)
	}
	if user.Routes != "" {
		args = append(args, "=routes="+user.Routes)
	}
	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *userRepository) RemoveUser(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/remove", "=.id="+id)
	return err
}

func (r *userRepository) RemoveUsersByComment(ctx context.Context, comment string) error {
	users, err := r.GetUsersByComment(ctx, comment)
	if err != nil {
		return err
	}
	for _, user := range users {
		_ = r.RemoveUser(ctx, user.ID)
	}
	return nil
}

func (r *userRepository) RemoveUsers(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.RemoveUser(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) DisableUser(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/disable", "=.id="+id)
	return err
}

func (r *userRepository) EnableUser(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/enable", "=.id="+id)
	return err
}

func (r *userRepository) DisableUsers(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.DisableUser(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) EnableUsers(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.EnableUser(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) ResetUserCounters(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/ip/hotspot/user/reset-counters", "=.id="+id)
	return err
}

func (r *userRepository) ResetUserCountersMultiple(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.ResetUserCounters(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) GetUsersCount(ctx context.Context) (int, error) {
	reply, err := r.client.RunContext(ctx, "/ip/hotspot/user/print", "=count-only=")
	if err != nil {
		return 0, err
	}
	if len(reply.Re) > 0 {
		count := int(utils.ParseInt(reply.Re[0].Map["ret"])) - 1
		if count < 0 {
			count = 0
		}
		return count, nil
	}
	return 0, nil
}

// PrintUsersRaw executes /ip/hotspot/user/print with custom arguments and returns raw data
func (r *userRepository) PrintUsersRaw(ctx context.Context, args ...string) ([]map[string]string, error) {
	cmdArgs := append([]string{"/ip/hotspot/user/print"}, args...)
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

// ListenUsersRaw executes /ip/hotspot/user/print with follow and returns raw data stream
func (r *userRepository) ListenUsersRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error) {
	cmdArgs := append([]string{"/ip/hotspot/user/print"}, args...)
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
