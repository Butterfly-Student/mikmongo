// cmd/collector_test/simple_test.go - Simple test untuk 2 router
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
)

func main() {
	log.Println("=== Simple MikroTik Connection Test ===")

	// Router A Config
	routerA := client.Config{
		Host:     "192.168.233.1",
		Port:     8728,
		Username: "admin",
		Password: "r00t",
	}

	// Router B Config
	routerB := client.Config{
		Host:     "192.168.27.1",
		Port:     8728,
		Username: "admin",
		Password: "r00t",
	}

	// Test Router A
	log.Println("\n[1] Testing Router A (192.168.233.1)...")
	testRouter("Router-A", routerA)

	// Test Router B
	log.Println("\n[2] Testing Router B (192.168.27.1)...")
	testRouter("Router-B", routerB)

	log.Println("\n=== Test Complete ===")
}

func testRouter(name string, cfg client.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create client
	c := client.NewClient(cfg, nil)
	if err := c.Connect(ctx); err != nil {
		log.Printf("  ✗ Failed to connect: %v", err)
		return
	}
	defer c.Close()

	log.Printf("  ✓ Connected to %s", name)

	// Test 1: Get system identity
	log.Println("  Testing: /system/identity/print")
	results, err := c.RunContext(ctx, "/system/identity/print")
	if err != nil {
		log.Printf("    ✗ Failed: %v", err)
	} else {
		log.Printf("    ✓ Success: %d items", len(results.Re))
		for _, r := range results.Re {
			log.Printf("      Identity: %s", r.Map["name"])
		}
	}

	// Test 2: Get interfaces
	log.Println("  Testing: /interface/print")
	results, err = c.RunContext(ctx, "/interface/print")
	if err != nil {
		log.Printf("    ✗ Failed: %v", err)
	} else {
		log.Printf("    ✓ Success: %d interfaces", len(results.Re))
		for i, r := range results.Re {
			if i >= 3 {
				log.Println("      ...")
				break
			}
			log.Printf("      - %s (type: %s, running: %s)", 
				r.Map["name"], r.Map["type"], r.Map["running"])
		}
	}

	// Test 3: Get PPP active (jika ada)
	log.Println("  Testing: /ppp/active/print")
	results, err = c.RunContext(ctx, "/ppp/active/print")
	if err != nil {
		log.Printf("    ✗ Failed: %v", err)
	} else {
		log.Printf("    ✓ Success: %d active PPP sessions", len(results.Re))
		for i, r := range results.Re {
			if i >= 3 {
				log.Println("      ...")
				break
			}
			log.Printf("      - %s (address: %s)", 
				r.Map["name"], r.Map["address"])
		}
	}

	// Test 4: Get system resources
	log.Println("  Testing: /system/resource/print")
	results, err = c.RunContext(ctx, "/system/resource/print")
	if err != nil {
		log.Printf("    ✗ Failed: %v", err)
	} else {
		log.Printf("    ✓ Success")
		if len(results.Re) > 0 {
			r := results.Re[0]
			log.Printf("      CPU Load: %s%%", r.Map["cpu-load"])
			log.Printf("      Free Memory: %s", r.Map["free-memory"])
			log.Printf("      Total Memory: %s", r.Map["total-memory"])
			log.Printf("      Uptime: %s", r.Map["uptime"])
		}
	}

	// Test 5: Test streaming dengan follow=yes (dengan timeout pendek)
	log.Println("  Testing: Streaming /interface/print stats interval=1s (5s)...")
	testStreaming(ctx, c)
}

func testStreaming(ctx context.Context, c *client.Client) {
	// Create short timeout context untuk streaming test
	streamCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resultCh := make(chan map[string]string, 100)
	
	// Start listening
	stopListen, err := c.ListenRaw(streamCtx, []string{"/interface/print", "stats", "=interval=1s"}, resultCh)
	if err != nil {
		log.Printf("    ✗ Failed to start listener: %v", err)
		return
	}
	defer stopListen()

	count := 0
	for {
		select {
		case <-streamCtx.Done():
			log.Printf("    ✓ Streaming test complete: %d events received", count)
			return
		case data := <-resultCh:
			if data != nil {
				count++
				if count == 1 {
					log.Printf("      First event: interface=%s rx-byte=%s", 
						data["name"], data["rx-byte"])
				}
			}
		}
	}
}

// Helper untuk menampilkan map
func formatMap(m map[string]string) string {
	result := ""
	for k, v := range m {
		result += fmt.Sprintf("%s=%s ", k, v)
	}
	return result
}
