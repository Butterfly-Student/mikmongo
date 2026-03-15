package monitor

import (
	"context"
	"fmt"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
)

// logRepository implements LogRepository interface
type logRepository struct {
	client *client.Client
}

// NewLogRepository creates a new log repository
func NewLogRepository(c *client.Client) LogRepository {
	return &logRepository{client: c}
}

func parseLogEntry(m map[string]string) *domain.LogEntry {
	return &domain.LogEntry{
		ID:      m[".id"],
		Time:    m["time"],
		Topics:  m["topics"],
		Message: m["message"],
	}
}

func (r *logRepository) GetLogs(ctx context.Context, topics string, limit int) ([]*domain.LogEntry, error) {
	args := []string{"/log/print"}
	if topics != "" {
		args = append(args, fmt.Sprintf("?topics=%s", topics))
	}
	reply, err := r.client.RunContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	logs := make([]*domain.LogEntry, 0, len(reply.Re))
	for i, re := range reply.Re {
		if limit > 0 && i >= limit {
			break
		}
		logs = append(logs, parseLogEntry(re.Map))
	}
	return logs, nil
}

func (r *logRepository) GetHotspotLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error) {
	_ = r.EnableHotspotLogging(ctx)
	return r.GetLogs(ctx, "hotspot,info,debug", limit)
}

func (r *logRepository) GetPPPLogs(ctx context.Context, limit int) ([]*domain.LogEntry, error) {
	_ = r.EnablePPPLogging(ctx)
	return r.GetLogs(ctx, "ppp,pppoe,info", limit)
}

func (r *logRepository) EnableHotspotLogging(ctx context.Context) error {
	reply, err := r.client.RunContext(ctx, "/system/logging/print", "?prefix=->")
	if err != nil {
		return err
	}
	if len(reply.Re) > 0 {
		return nil
	}
	_, err = r.client.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=->",
		"=topics=hotspot,info,debug",
	)
	return err
}

func (r *logRepository) EnablePPPLogging(ctx context.Context) error {
	reply, err := r.client.RunContext(ctx, "/system/logging/print", "?prefix=ppp->")
	if err != nil {
		return err
	}
	if len(reply.Re) > 0 {
		return nil
	}
	_, err = r.client.RunContext(ctx,
		"/system/logging/add",
		"=action=disk",
		"=prefix=ppp->",
		"=topics=pppoe",
	)
	return err
}

func (r *logRepository) ListenLogs(ctx context.Context, topics string, resultChan chan<- *domain.LogEntry) (func() error, error) {
	args := []string{"/log/print", "=follow-only="}
	if topics != "" {
		args = append(args, fmt.Sprintf("?topics=%s", topics))
	}
	listenReply, err := r.client.ListenArgsContext(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to start log listen: %w", err)
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
				case resultChan <- parseLogEntry(sentence.Map):
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

func (r *logRepository) ListenHotspotLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return r.ListenLogs(ctx, "hotspot,info", resultChan)
}

func (r *logRepository) ListenPPPLogs(ctx context.Context, resultChan chan<- *domain.LogEntry) (func() error, error) {
	return r.ListenLogs(ctx, "pppoe", resultChan)
}
