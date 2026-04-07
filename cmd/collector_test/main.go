// cmd/collector_test/main.go - Testing 3-Pipeline Collector dengan 2 router
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	"mikmongo/pkg/mikrotik/collector"
	"mikmongo/pkg/mikrotik/collector/pipeline/ondemand"
	"mikmongo/pkg/mikrotik/collector/writer"
)

func main() {
	log.Println("=== MikroTik 3-Pipeline Collector Test ===")

	// Config InfluxDB
	influxConfig := writer.InfluxConfig{
		URL:    "http://localhost:8086",
		Token:  os.Getenv("INFLUX_TOKEN"),
		Org:    os.Getenv("INFLUX_ORG"),
		Bucket: os.Getenv("INFLUX_BUCKET"),
	}

	// Default values jika env tidak set
	if influxConfig.Token == "" {
		influxConfig.Token = "test-token"
	}
	if influxConfig.Org == "" {
		influxConfig.Org = "mikrotik"
	}
	if influxConfig.Bucket == "" {
		influxConfig.Bucket = "miktik"
	}

	// Config Redis
	redisConfig := writer.RedisConfig{
		Addr:   "localhost:6379",
		DB:     0,
		Prefix: "mikrotik",
	}

	// Config Batch Writer
	batchConfig := writer.Config{
		BatchSize:     50,
		FlushInterval: 100 * time.Millisecond,
	}

	// Create Manager
	managerConfig := collector.ManagerConfig{
		InfluxConfig: influxConfig,
		RedisConfig:  redisConfig,
		BatchConfig:  batchConfig,
	}

	manager, err := collector.NewManager(managerConfig)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Router A Config
	routerAConfig := client.Config{
		Host:     "192.168.233.1",
		Port:     8728,
		Username: "admin",
		Password: "r00t",
	}

	// Router B Config
	routerBConfig := client.Config{
		Host:     "192.168.27.1",
		Port:     8728,
		Username: "admin",
		Password: "r00t",
	}

	// Define Pipeline A Specs (Time-series → InfluxDB)
	timeSeriesSpecs := []collector.CommandSpec{
		// Interface stats - Tier 1
		{
			Name:      "interface_stats",
			Args:      []string{"/interface/print", "stats", "interval=1s"},
			KeyField:  "name",
			Tier:      collector.Tier1,
			TTL:       0, // No TTL untuk time-series
			ToInflux:  true,
			RedisKey:  "",
			Interval:  0,
		},
		// System resources - Tier 1
		{
			Name:      "system_resource",
			Args:      []string{"/system/resource/print", "interval=5s"},
			KeyField:  "",
			Tier:      collector.Tier1,
			TTL:       0,
			ToInflux:  true,
			RedisKey:  "",
			Interval:  0,
		},
	}

	// Define Pipeline B Specs (Operational → Redis)
	operationalSpecs := []collector.CommandSpec{
		// Tier 2: follow=yes untuk real-time state
		{
			Name:      "ppp_active",
			Args:      []string{"/ppp/active/print", "follow=yes"},
			KeyField:  "name",
			Tier:      collector.Tier2,
			TTL:       0, // No TTL untuk real-time state
			ToInflux:  false,
			RedisKey:  "ppp:active",
			Interval:  0,
		},
		{
			Name:      "hotspot_active",
			Args:      []string{"/ip/hotspot/active/print", "follow=yes"},
			KeyField:  "mac-address",
			Tier:      collector.Tier2,
			TTL:       0,
			ToInflux:  false,
			RedisKey:  "hotspot:active",
			Interval:  0,
		},
		// Tier 3: ticker untuk static data
		{
			Name:      "ppp_secrets",
			Args:      []string{"/ppp/secret/print"},
			KeyField:  "name",
			Tier:      collector.Tier3,
			TTL:       10 * time.Minute,
			ToInflux:  false,
			RedisKey:  "ppp:secrets",
			Interval:  5 * time.Minute,
		},
		{
			Name:      "ppp_profiles",
			Args:      []string{"/ppp/profile/print"},
			KeyField:  "name",
			Tier:      collector.Tier3,
			TTL:       10 * time.Minute,
			ToInflux:  false,
			RedisKey:  "ppp:profiles",
			Interval:  5 * time.Minute,
		},
		{
			Name:      "ip_pools",
			Args:      []string{"/ip/pool/print"},
			KeyField:  "name",
			Tier:      collector.Tier3,
			TTL:       20 * time.Minute,
			ToInflux:  false,
			RedisKey:  "ip:pools",
			Interval:  10 * time.Minute,
		},
	}

	// Add Router A
	log.Println("\n[1] Adding Router A (192.168.233.1)...")
	if err := manager.AddRouter("router-a", routerAConfig, timeSeriesSpecs, operationalSpecs); err != nil {
		log.Printf("Failed to add Router A: %v", err)
	} else {
		log.Println("✓ Router A added")
	}

	// Add Router B
	log.Println("\n[2] Adding Router B (192.168.27.1)...")
	if err := manager.AddRouter("router-b", routerBConfig, timeSeriesSpecs, operationalSpecs); err != nil {
		log.Printf("Failed to add Router B: %v", err)
	} else {
		log.Println("✓ Router B added")
	}

	// Start all collectors
	log.Println("\n[3] Starting all collectors...")
	manager.StartAll()

	// Wait untuk collectors start
	time.Sleep(2 * time.Second)

	// Test Pipeline C: On-demand operations
	log.Println("\n[4] Testing Pipeline C (On-demand)...")
	testOnDemand(manager)

	// Print stats
	log.Println("\n[5] Collector Statistics:")
	printStats(manager)

	// Wait for data collection
	log.Println("\n[6] Collecting data for 30 seconds...")
	time.Sleep(30 * time.Second)

	// Print stats lagi
	log.Println("\n[7] Statistics after 30s:")
	printStats(manager)

	// Test Pipeline C lagi
	log.Println("\n[8] Testing Pipeline C (Read operations)...")
	testOnDemandRead(manager)

	// Setup graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("\n[9] Running... Press Ctrl+C to stop")
	<-sigCh

	log.Println("\n[10] Stopping all collectors...")
	manager.StopAll()

	log.Println("\n=== Test Complete ===")
}

// testOnDemand tests Pipeline C write operations
func testOnDemand(manager *collector.Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test dengan Router A
	runnerA := manager.GetRunner("router-a")
	if runnerA == nil {
		log.Println("✗ Router A runner not available")
		return
	}

	// Test Ping
	log.Println("  Testing Ping on Router A...")
	if err := runnerA.Ping(ctx); err != nil {
		log.Printf("  ✗ Ping failed: %v", err)
	} else {
		log.Println("  ✓ Ping success")
	}

	// Test Read: System identity
	log.Println("  Testing Read: /system/identity/print...")
	results, err := runnerA.Run(ctx, []string{"/system/identity/print"})
	if err != nil {
		log.Printf("  ✗ Read failed: %v", err)
	} else {
		log.Printf("  ✓ Read success: %d items", len(results))
		if len(results) > 0 {
			log.Printf("    Identity: %v", results[0])
		}
	}

	// Test Read: Interface list
	log.Println("  Testing Read: /interface/print...")
	results, err = runnerA.Run(ctx, []string{"/interface/print"})
	if err != nil {
		log.Printf("  ✗ Read failed: %v", err)
	} else {
		log.Printf("  ✓ Read success: %d interfaces", len(results))
	}

	// Test dengan Router B
	runnerB := manager.GetRunner("router-b")
	if runnerB == nil {
		log.Println("✗ Router B runner not available")
		return
	}

	log.Println("  Testing Ping on Router B...")
	if err := runnerB.Ping(ctx); err != nil {
		log.Printf("  ✗ Ping failed: %v", err)
	} else {
		log.Println("  ✓ Ping success")
	}
}

// testOnDemandRead tests read operations
func testOnDemandRead(manager *collector.Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	runnerA := manager.GetRunner("router-a")
	if runnerA == nil {
		return
	}

	// Read PPP secrets
	log.Println("  Reading PPP secrets...")
	results, err := runnerA.Run(ctx, []string{"/ppp/secret/print"})
	if err != nil {
		log.Printf("  ✗ Failed: %v", err)
	} else {
		log.Printf("  ✓ Found %d PPP secrets", len(results))
		for i, r := range results {
			if i >= 3 { // Show max 3
				log.Println("    ...")
				break
			}
			log.Printf("    - %s (profile: %s)", r["name"], r["profile"])
		}
	}

	// Read IP pools
	log.Println("  Reading IP pools...")
	results, err = runnerA.Run(ctx, []string{"/ip/pool/print"})
	if err != nil {
		log.Printf("  ✗ Failed: %v", err)
	} else {
		log.Printf("  ✓ Found %d IP pools", len(results))
		for i, r := range results {
			if i >= 3 {
				log.Println("    ...")
				break
			}
			log.Printf("    - %s: %s", r["name"], r["ranges"])
		}
	}
}

// printStats prints manager statistics
func printStats(manager *collector.Manager) {
	stats := manager.GetStats()
	fmt.Printf("  Routers: %d\n", stats["router_count"])

	if routers, ok := stats["routers"].(map[string]interface{}); ok {
		for id, r := range routers {
			fmt.Printf("  Router %s:\n", id)
			if rs, ok := r.(map[string]interface{}); ok {
				fmt.Printf("    Time-series specs: %v\n", rs["time_series_specs"])
				fmt.Printf("    Operational specs: %v\n", rs["operational_specs"])
			}
		}
	}
}
