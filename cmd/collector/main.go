// cmd/collector/main.go
// Test runner untuk MikroTik Collector dengan 4 Tier System
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/collector"

	"github.com/Butterfly-Student/go-ros/client"
)

func main() {
	var (
		host     = flag.String("host", "192.168.27.1", "MikroTik router address")
		port     = flag.Int("port", 8728, "MikroTik API port")
		username = flag.String("user", "admin", "MikroTik username")
		password = flag.String("pass", "r00t", "MikroTik password")
		redis    = flag.String("redis", "localhost:6379", "Redis address")
		duration = flag.Int("duration", 60, "Test duration in seconds")
	)
	flag.Parse()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  MikroTik 4-Tier Collector Test Runner")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("Router:    %s:%d\n", *host, *port)
	fmt.Printf("Redis:     %s\n", *redis)
	fmt.Printf("Duration:  %d seconds\n", *duration)
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// Step 1: Create MikroTik Client (existing code)
	// ═══════════════════════════════════════════════════════════
	cfg := client.Config{
		Host:     *host,
		Port:     *port,
		Username: *username,
		Password: *password,
	}

	log.Println("[1/4] Connecting to MikroTik...")
	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer c.Close()
	log.Println("      Connected!")

	// Wrap dengan existing Client facade
	mtClient := mikrotik.NewClientFromConnection(c)

	// ═══════════════════════════════════════════════════════════
	// Step 2: Create Collector
	// ═══════════════════════════════════════════════════════════
	log.Println("[2/4] Initializing Collector...")

	collCfg := collector.Config{
		RedisAddr: *redis,
	}

	coll, err := collector.NewCollector(collCfg)
	if err != nil {
		log.Fatalf("Failed to create collector: %v", err)
	}
	defer coll.Stop()

	// ═══════════════════════════════════════════════════════════
	// Step 3: Combine specs dan register router
	// ═══════════════════════════════════════════════════════════
	log.Println("[3/4] Registering router with specs...")

	// Combine Tier 1/2 specs dengan Tier 3 specs
	allSpecs := append(collector.DefaultISPSpecs(), collector.Tier3Specs()...)

	routerID := fmt.Sprintf("router-%s", *host)
	if err := coll.RegisterRouter(routerID, mtClient, allSpecs); err != nil {
		log.Fatalf("Failed to register router: %v", err)
	}

	// Print spec summary
	fmt.Println()
	fmt.Println("Specs registered:")
	tier1Count, tier2Count, tier3Count := 0, 0, 0
	for _, spec := range allSpecs {
		switch spec.Tier {
		case collector.Tier1:
			tier1Count++
		case collector.Tier2:
			tier2Count++
		case collector.Tier3:
			tier3Count++
		}
	}
	fmt.Printf("  - Tier 1 (High Freq):   %d specs\n", tier1Count)
	fmt.Printf("  - Tier 2 (Event):       %d specs\n", tier2Count)
	fmt.Printf("  - Tier 3 (Static):      %d specs\n", tier3Count)
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// Step 4: Start collecting
	// ═══════════════════════════════════════════════════════════
	log.Println("[4/4] Starting collection...")
	if err := coll.StartRouter(routerID); err != nil {
		log.Fatalf("Failed to start router: %v", err)
	}

	log.Println("      Collection started!")
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// Step 5: Monitor dan test
	// ═══════════════════════════════════════════════════════════
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Setup test ticker untuk print stats
	statsTicker := time.NewTicker(10 * time.Second)
	defer statsTicker.Stop()

	// Setup test duration timer
	testTimer := time.NewTimer(time.Duration(*duration) * time.Second)
	defer testTimer.Stop()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Monitoring started (Ctrl+C to stop)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Test: Read cached data setiap 10 detik
	for {
		select {
		case <-sigCh:
			fmt.Println("\n[Signal] Shutting down...")
			return

		case <-testTimer.C:
			fmt.Println("\n[Timer] Test duration completed")
			fmt.Println("\nFinal pool stats:")
			printStats(coll, routerID)
			return

		case <-statsTicker.C:
			printStats(coll, routerID)

			// Test: Coba baca cache
			testCacheRead(coll, routerID)
		}
	}
}

// printStats prints pool statistics
func printStats(coll *collector.Collector, routerID string) {
	stats, err := coll.GetPoolStats(routerID)
	if err != nil {
		log.Printf("Stats error: %v", err)
		return
	}

	fmt.Printf("[%s] Pool Usage: T1[%d/%d] T2[%d/%d] T3[%d/%d]\n",
		time.Now().Format("15:04:05"),
		stats.Tier1Used, stats.Tier1Total,
		stats.Tier2Used, stats.Tier2Total,
		stats.Tier3Used, stats.Tier3Total,
	)
}

// testCacheRead mencoba membaca data dari cache
func testCacheRead(coll *collector.Collector, routerID string) {
	// Test baca Tier 1 cache
	if data, err := coll.GetCachedDataAll(routerID, "interface:stats"); err == nil && len(data) > 0 {
		fmt.Printf("  [Cache] interface:stats: %d entries\n", len(data))
	}

	// Test baca Tier 2 cache
	if data, err := coll.GetCachedDataAll(routerID, "ppp:active"); err == nil && len(data) > 0 {
		fmt.Printf("  [Cache] ppp:active: %d entries\n", len(data))
	}

	// Test baca Tier 3 cache
	if data, err := coll.GetCachedDataAll(routerID, "ppp:secrets"); err == nil && len(data) > 0 {
		fmt.Printf("  [Cache] ppp:secrets: %d entries\n", len(data))
	}
}
