// Package main streams live data from a MikroTik router using ListenManyArgsContext.
//
// It connects to a real RouterOS device and concurrently listens to multiple
// streaming commands, printing every event to stdout. All listeners share a
// single async TCP connection via the library's tag-multiplexing mechanism.
// The program runs until the user presses Ctrl+C.
//
// Usage:
//
//	go run ./cmd/mikrotik \
//	    -host 192.168.1.1 \
//	    -port 8728 \
//	    -user admin \
//	    -pass ""
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
)

func main() {
	host := flag.String("host", "192.168.233.1", "RouterOS host IP")
	port := flag.Int("port", 8728, "RouterOS API port")
	user := flag.String("user", "admin", "RouterOS username")
	pass := flag.String("pass", "r00t", "RouterOS password")
	flag.Parse()

	cfg := client.Config{
		Host:     *host,
		Port:     *port,
		Username: *user,
		Password: *pass,
		Timeout:  10 * time.Second,
	}

	fmt.Printf("Connecting to %s:%d as %s …\n", cfg.Host, cfg.Port, cfg.Username)

	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	fmt.Printf("Connected (async=%v). Listening for %s …\n\n", c.IsAsync())

	// ─────────────────────────────────────────────────────────────────────────
	// Define the commands to listen concurrently.
	//
	// RouterOS streaming commands ("follow" style): the device keeps sending
	// !re sentences every configured interval until we cancel.
	//
	// ".proplist" limits which keys come back, keeping output manageable.
	// ─────────────────────────────────────────────────────────────────────────
	commands := [][]string{
		// 0 – interface traffic (follow) - every interface, all traffic counters
		{
			"/interface/print",
			"=.proplist=name,type,rx-byte,tx-byte,rx-packet,tx-packet,running",
			"=follow=",
		},
		// 1 – interface stats only (follow-only = suppresses initial snapshot)
		{
			"/interface/print",
			"=.proplist=name,rx-byte,tx-byte",
			"=follow-only=",
		},
		// 2 – simple queue stats (follow)
		{
			"/queue/simple/print",
			"=.proplist=name,bytes,packets,dropped",
			"=follow=",
		},
		// 3 – system resource monitoring (follow) - CPU, memory, uptime
		{
			"/system/resource/print",
			"=.proplist=uptime,cpu-load,free-memory,total-memory",
			"=follow=",
		},
		// 4 – active PPPoE / hotspot sessions (plain print, not a follow — will
		//     return all current rows once then close; demonstrates mixing
		//     one-shot and streaming commands in the same fan-out)
		{
			"/ppp/active/print",
			"=.proplist=name,address,uptime,caller-id",
		},
		// 5 – IP address table (plain print, one-shot)
		{
			"/ip/address/print",
			"=.proplist=address,interface,network,disabled",
		},
	}

	labels := []string{
		"[0] iface/follow",
		"[1] iface/follow-only",
		"[2] queue/follow",
		"[3] resource/follow",
		"[4] ppp/active (one-shot)",
		"[5] ip/address  (one-shot)",
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Press Ctrl+C to stop.")

	ch, err := c.ListenManyArgsContext(ctx, commands, 128)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListenManyArgsContext: %v\n", err)
		os.Exit(1)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tSOURCE\tKEY=VALUE …")
	fmt.Fprintln(tw, "────\t──────\t──────────")

	total := 0
	for ev := range ch {
		label := labels[ev.Index]
		if ev.Err != nil {
			fmt.Fprintf(tw, "%s\t%s\tERROR: %v\n",
				time.Now().Format("15:04:05.000"), label, ev.Err)
			tw.Flush()
			continue
		}

		// Format the map as sorted key=value pairs.
		pairs := make([]string, 0, len(ev.Map))
		for k, v := range ev.Map {
			pairs = append(pairs, k+"="+v)
		}
		sort.Strings(pairs)

		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			time.Now().Format("15:04:05.000"),
			label,
			strings.Join(pairs, "  "),
		)
		tw.Flush()
		total++
	}

	fmt.Printf("\n─── done: received %d events ───\n", total)
}
