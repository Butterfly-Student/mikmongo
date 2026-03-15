package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	"github.com/Butterfly-Student/go-ros/utils"
	"github.com/go-routeros/routeros/v3/proto"
)

// interfaceRepository implements InterfaceRepository interface
type interfaceRepository struct {
	client *client.Client
}

// NewInterfaceRepository creates a new interface repository
func NewInterfaceRepository(c *client.Client) InterfaceRepository {
	return &interfaceRepository{client: c}
}

func parseInterface(m map[string]string) *domain.Interface {
	return &domain.Interface{
		ID:         m[".id"],
		Name:       m["name"],
		Type:       m["type"],
		MTU:        int(utils.ParseInt(m["actual-mtu"])),
		MacAddress: m["mac-address"],
		Running:    utils.ParseBool(m["running"]),
		Disabled:   utils.ParseBool(m["disabled"]),
		Comment:    m["comment"],
	}
}

func (r *interfaceRepository) GetInterfaces(ctx context.Context) ([]*domain.Interface, error) {
	reply, err := r.client.RunContext(ctx, "/interface/print")
	if err != nil {
		return nil, err
	}
	interfaces := make([]*domain.Interface, 0, len(reply.Re))
	for _, re := range reply.Re {
		interfaces = append(interfaces, parseInterface(re.Map))
	}
	return interfaces, nil
}

func (r *interfaceRepository) StartTrafficMonitorListen(ctx context.Context, name string, resultChan chan<- domain.TrafficMonitorStats) (func() error, error) {
	if name == "" {
		return nil, fmt.Errorf("interface name is required")
	}

	listenReply, err := r.client.ListenArgsContext(ctx, []string{
		"/interface/monitor-traffic",
		fmt.Sprintf("=interface=%s", name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start traffic monitor listen: %w", err)
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
				result := parseTrafficMonitorSentence(sentence, name)
				result.Timestamp = time.Now()
				select {
				case resultChan <- result:
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

func parseTrafficMonitorSentence(sentence *proto.Sentence, name string) domain.TrafficMonitorStats {
	m := sentence.Map
	return domain.TrafficMonitorStats{
		Name:                  name,
		RxBitsPerSecond:       client.ParseRate(m["rx-bits-per-second"]),
		TxBitsPerSecond:       client.ParseRate(m["tx-bits-per-second"]),
		RxPacketsPerSecond:    utils.ParseInt(m["rx-packets-per-second"]),
		TxPacketsPerSecond:    utils.ParseInt(m["tx-packets-per-second"]),
		FpRxBitsPerSecond:     client.ParseRate(m["fp-rx-bits-per-second"]),
		FpTxBitsPerSecond:     client.ParseRate(m["fp-tx-bits-per-second"]),
		FpRxPacketsPerSecond:  utils.ParseInt(m["fp-rx-packets-per-second"]),
		FpTxPacketsPerSecond:  utils.ParseInt(m["fp-tx-packets-per-second"]),
		RxDropsPerSecond:      utils.ParseInt(m["rx-drops-per-second"]),
		TxDropsPerSecond:      utils.ParseInt(m["tx-drops-per-second"]),
		TxQueueDropsPerSecond: utils.ParseInt(m["tx-queue-drops-per-second"]),
		RxErrorsPerSecond:     utils.ParseInt(m["rx-errors-per-second"]),
		TxErrorsPerSecond:     utils.ParseInt(m["tx-errors-per-second"]),
	}
}
