package system

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
)

// scriptsRepository implements ScriptsRepository interface
type scriptsRepository struct {
	client *client.Client
}

// NewScriptsRepository creates a new scripts repository
func NewScriptsRepository(c *client.Client) ScriptsRepository {
	return &scriptsRepository{client: c}
}

func parseSystemScript(m map[string]string) *domain.SystemScript {
	return &domain.SystemScript{
		ID:                     m[".id"],
		Name:                   m["name"],
		Owner:                  m["owner"],
		Policy:                 m["policy"],
		Source:                 m["source"],
		Comment:                m["comment"],
		DontRequirePermissions: utils.ParseBool(m["dont-require-permissions"]),
		RunCount:               m["run-count"],
		LastTimeStarted:        m["last-time-started"],
	}
}

func (r *scriptsRepository) GetScripts(ctx context.Context) ([]*domain.SystemScript, error) {
	reply, err := r.client.RunContext(ctx, "/system/script/print")
	if err != nil {
		return nil, err
	}
	scripts := make([]*domain.SystemScript, 0, len(reply.Re))
	for _, re := range reply.Re {
		scripts = append(scripts, parseSystemScript(re.Map))
	}
	return scripts, nil
}

func (r *scriptsRepository) GetScriptByName(ctx context.Context, name string) (*domain.SystemScript, error) {
	reply, err := r.client.RunContext(ctx, "/system/script/print", "?name="+name)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return nil, nil
	}
	return parseSystemScript(reply.Re[0].Map), nil
}

func (r *scriptsRepository) AddScript(ctx context.Context, script *domain.SystemScript) (string, error) {
	args := []string{
		"/system/script/add",
		"=name=" + script.Name,
		"=source=" + script.Source,
	}

	if script.Owner != "" {
		args = append(args, "=owner="+script.Owner)
	}
	if script.Comment != "" {
		args = append(args, "=comment="+script.Comment)
	}
	if script.Policy != "" {
		args = append(args, "=policy="+script.Policy)
	}
	if script.DontRequirePermissions {
		args = append(args, "=dont-require-permissions=yes")
	} else {
		args = append(args, "=dont-require-permissions=no")
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

func (r *scriptsRepository) UpdateScript(ctx context.Context, id string, script *domain.SystemScript) error {
	args := []string{"/system/script/set", "=.id=" + id}

	if script.Name != "" {
		args = append(args, "=name="+script.Name)
	}
	if script.Source != "" {
		args = append(args, "=source="+script.Source)
	}
	if script.Owner != "" {
		args = append(args, "=owner="+script.Owner)
	}
	if script.Comment != "" {
		args = append(args, "=comment="+script.Comment)
	}
	if script.Policy != "" {
		args = append(args, "=policy="+script.Policy)
	}
	if script.DontRequirePermissions {
		args = append(args, "=dont-require-permissions=yes")
	} else {
		args = append(args, "=dont-require-permissions=no")
	}

	_, err := r.client.RunArgsContext(ctx, args)
	return err
}

func (r *scriptsRepository) RemoveScript(ctx context.Context, id string) error {
	_, err := r.client.RunContext(ctx, "/system/script/remove", "=.id="+id)
	return err
}

func (r *scriptsRepository) RunScript(ctx context.Context, name string) error {
	_, err := r.client.RunContext(ctx, "/system/script/run", "=number="+name)
	return err
}
