// cmd/full_test/main.go - Full 3-Pipeline Test untuk Router B
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/redis/go-redis/v9"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// TestConfig holds test configuration
type TestConfig struct {
	RouterConfig client.Config
	InfluxURL    string
	InfluxToken  string
	InfluxOrg    string
	InfluxBucket string
	RedisAddr    string
	Duration     time.Duration
}

func main() {
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("     3-Pipeline Collector Full Test - Router B")
	log.Println("═══════════════════════════════════════════════════════════")

	config := TestConfig{
		RouterConfig: client.Config{
			Host:     "192.168.27.1",
			Port:     8728,
			Username: "admin",
			Password: "r00t",
		},
		InfluxURL:    "http://localhost:8086",
		InfluxToken:  "my-super-secret-token",
		InfluxOrg:    "myorg",
		InfluxBucket: "miktik",
		RedisAddr:    "localhost:6379",
		Duration:     60 * time.Second,
	}

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("\n[Signal] Shutting down...")
		cancel()
	}()

	// Test 1: Basic Connection
	log.Println("\n▶ [Test 1] Basic Connection Test")
	c := testConnection(ctx, config.RouterConfig)
	if c == nil {
		log.Fatal("Failed to connect to router")
	}
	defer c.Close()

	// Test 2: Setup InfluxDB and Redis
	log.Println("\n▶ [Test 2] Setup InfluxDB and Redis Clients")
	influxClient, writeAPI := setupInfluxDB(config)
	defer influxClient.Close()

	redisClient := setupRedis(config)
	defer redisClient.Close()

	// Test 3: Pipeline A - Time-series to InfluxDB
	log.Println("\n▶ [Test 3] Pipeline A: Time-series → InfluxDB")
	go runPipelineA(ctx, c, writeAPI)

	// Test 4: Pipeline B Tier 2 - follow=yes to Redis
	log.Println("\n▶ [Test 4] Pipeline B Tier 2: follow=yes → Redis")
	go runPipelineBTier2(ctx, c, redisClient)

	// Test 5: Pipeline B Tier 3 - Ticker to Redis
	log.Println("\n▶ [Test 5] Pipeline B Tier 3: Ticker → Redis")
	go runPipelineBTier3(ctx, config.RouterConfig, redisClient)

	// Test 6: Pipeline C - On-demand operations
	log.Println("\n▶ [Test 6] Pipeline C: On-demand Operations")
	time.Sleep(5 * time.Second) // Wait a bit first
	go runPipelineC(ctx, config.RouterConfig, redisClient)

	// Wait for collection
	log.Printf("\n⏱️  Collecting data for %v...", config.Duration)
	time.Sleep(config.Duration)

	// Verify data
	log.Println("\n▶ [Test 7] Verifying Collected Data")
	verifyData(ctx, config, influxClient, redisClient)

	log.Println("\n═══════════════════════════════════════════════════════════")
	log.Println("                    Test Complete")
	log.Println("═══════════════════════════════════════════════════════════")
}

// testConnection tests basic connectivity
func testConnection(ctx context.Context, cfg client.Config) *client.Client {
	c := client.NewClient(cfg, nil)
	if err := c.Connect(ctx); err != nil {
		log.Printf("  ✗ Failed to connect: %v", err)
		return nil
	}

	// Test basic command
	result, err := c.RunContext(ctx, "/system/identity/print")
	if err != nil {
		log.Printf("  ✗ Failed to get identity: %v", err)
		c.Close()
		return nil
	}

	if len(result.Re) > 0 {
		log.Printf("  ✓ Connected. Router Identity: %s", result.Re[0].Map["name"])
	}
	return c
}

// setupInfluxDB initializes InfluxDB client
func setupInfluxDB(config TestConfig) (influxdb2.Client, api.WriteAPI) {
	client := influxdb2.NewClient(config.InfluxURL, config.InfluxToken)
	
	// Health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	health, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("  ✗ InfluxDB health check failed: %v", err)
	}
	log.Printf("  ✓ InfluxDB connected (status: %s)", health.Status)

	writeAPI := client.WriteAPI(config.InfluxOrg, config.InfluxBucket)
	return client, writeAPI
}

// setupRedis initializes Redis client
func setupRedis(config TestConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("  ✗ Redis connection failed: %v", err)
	}
	log.Println("  ✓ Redis connected")
	return client
}

// runPipelineA collects time-series data to InfluxDB
func runPipelineA(ctx context.Context, c *client.Client, writeAPI api.WriteAPI) {
	log.Println("  Starting interface stats collection...")
	
	resultCh := make(chan map[string]string, 100)
	stopListen, err := c.ListenRaw(ctx, []string{
		"/interface/print", 
		"stats", 
		"=interval=1s",
		"=.proplist=name,rx-byte,tx-byte,rx-packet,tx-packet",
	}, resultCh)
	
	if err != nil {
		log.Printf("  ✗ Failed to start listener: %v", err)
		return
	}
	defer stopListen()

	count := 0
	for {
		select {
		case <-ctx.Done():
			log.Printf("  ✓ Pipeline A stopped. Total points: %d", count)
			return
		case data := <-resultCh:
			if data == nil {
				continue
			}
			
			// Write to InfluxDB
			point := influxdb2.NewPoint(
				"interface_stats",
				map[string]string{
					"router": "router-b",
					"name":   data["name"],
				},
				map[string]interface{}{
					"rx_byte":    parseInt(data["rx-byte"]),
					"tx_byte":    parseInt(data["tx-byte"]),
					"rx_packet":  parseInt(data["rx-packet"]),
					"tx_packet":  parseInt(data["tx-packet"]),
				},
				time.Now(),
			)
			writeAPI.WritePoint(point)
			count++
			
			if count%10 == 0 {
				log.Printf("    Written %d points to InfluxDB", count)
			}
		}
	}
}

// runPipelineBTier2 collects real-time state to Redis
func runPipelineBTier2(ctx context.Context, c *client.Client, rdb *redis.Client) {
	log.Println("  Starting PPP active monitoring...")

	resultCh := make(chan map[string]string, 100)
	stopListen, err := c.ListenRaw(ctx, []string{
		"/ppp/active/print",
		"follow=yes",
	}, resultCh)

	if err != nil {
		log.Printf("  ✗ Failed to start PPP listener: %v", err)
		return
	}
	defer stopListen()

	count := 0
	for {
		select {
		case <-ctx.Done():
			log.Printf("  ✓ Pipeline B Tier 2 stopped. Total updates: %d", count)
			return
		case data := <-resultCh:
			if data == nil {
				continue
			}

			// Write to Redis Hash
			name := data["name"]
			if name == "" {
				continue
			}

			key := "mikrotik:router-b:ppp:active"
			_, err := rdb.HSet(ctx, key, name, formatHash(data)).Result()
			if err != nil {
				log.Printf("    Redis HSet error: %v", err)
			} else {
				count++
				if count%5 == 0 {
					log.Printf("    Written %d PPP active updates to Redis", count)
				}
			}
		}
	}
}

// runPipelineBTier3 collects static data with ticker to Redis
func runPipelineBTier3(ctx context.Context, cfg client.Config, rdb *redis.Client) {
	log.Println("  Starting PPP secrets collection (ticker)...")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	count := 0
	// Initial fetch
	fetchAndCachePPPSecrets(ctx, cfg, rdb, &count)

	for {
		select {
		case <-ctx.Done():
			log.Printf("  ✓ Pipeline B Tier 3 stopped. Total fetches: %d", count)
			return
		case <-ticker.C:
			fetchAndCachePPPSecrets(ctx, cfg, rdb, &count)
		}
	}
}

func fetchAndCachePPPSecrets(ctx context.Context, cfg client.Config, rdb *redis.Client, count *int) {
	c := client.NewClient(cfg, nil)
	if err := c.Connect(ctx); err != nil {
		log.Printf("    Connection error: %v", err)
		return
	}
	defer c.Close()

	result, err := c.RunContext(ctx, "/ppp/secret/print")
	if err != nil {
		log.Printf("    Run error: %v", err)
		return
	}

	key := "mikrotik:router-b:ppp:secrets"
	pipe := rdb.Pipeline()

	for _, reply := range result.Re {
		name := reply.Map["name"]
		if name == "" {
			continue
		}
		pipe.HSet(ctx, key, name, formatHash(reply.Map))
	}

	pipe.Expire(ctx, key, 2*time.Minute)
	_, err = pipe.Exec(ctx)
	
	if err != nil {
		log.Printf("    Redis pipeline error: %v", err)
	} else {
		*count++
		log.Printf("    Cached %d PPP secrets (fetch #%d)", len(result.Re), *count)
	}
}

// runPipelineC performs on-demand operations
func runPipelineC(ctx context.Context, cfg client.Config, rdb *redis.Client) {
	log.Println("  Performing on-demand operations...")

	c := client.NewClient(cfg, nil)
	if err := c.Connect(ctx); err != nil {
		log.Printf("    Connection error: %v", err)
		return
	}
	defer c.Close()

	// Operation 1: Get system resource
	result, err := c.RunContext(ctx, "/system/resource/print")
	if err != nil {
		log.Printf("    ✗ System resource read failed: %v", err)
	} else if len(result.Re) > 0 {
		r := result.Re[0].Map
		log.Printf("    ✓ System: CPU=%s%%, Memory=%s/%s",
			r["cpu-load"], r["free-memory"], r["total-memory"])
	}

	// Operation 2: Get interface list and cache to Redis
	result, err = c.RunContext(ctx, "/interface/print")
	if err != nil {
		log.Printf("    ✗ Interface read failed: %v", err)
	} else {
		key := "mikrotik:router-b:interface:list"
		pipe := rdb.Pipeline()
		for _, reply := range result.Re {
			name := reply.Map["name"]
			if name != "" {
				pipe.HSet(ctx, key, name, formatHash(reply.Map))
			}
		}
		pipe.Expire(ctx, key, 5*time.Minute)
		pipe.Exec(ctx)
		log.Printf("    ✓ Cached %d interfaces to Redis", len(result.Re))
	}

	// Operation 3: Get IP pools
	result, err = c.RunContext(ctx, "/ip/pool/print")
	if err != nil {
		log.Printf("    ✗ IP pool read failed: %v", err)
	} else {
		log.Printf("    ✓ Found %d IP pools", len(result.Re))
	}
}

// verifyData checks collected data in InfluxDB and Redis
func verifyData(ctx context.Context, config TestConfig, influxClient influxdb2.Client, rdb *redis.Client) {
	log.Println("\n  [InfluxDB Verification]")
	
	// Query InfluxDB for interface stats
	queryAPI := influxClient.QueryAPI(config.InfluxOrg)
	query := fmt.Sprintf(`
		from(bucket: "%s")
		|> range(start: -2m)
		|> filter(fn: (r) => r._measurement == "interface_stats")
		|> count()
	`, config.InfluxBucket)

	result, err := queryAPI.Query(ctx, query)
	totalPoints := 0
	if err != nil {
		log.Printf("    ✗ InfluxDB query failed: %v", err)
	} else {
		for result.Next() {
			if v, ok := result.Record().Value().(int64); ok {
				totalPoints = int(v)
			}
		}
		if totalPoints > 0 {
			log.Printf("    ✓ Found %d interface_stats points in InfluxDB", totalPoints)
		} else {
			log.Println("    ⚠ No interface_stats data found")
		}
	}

	log.Println("\n  [Redis Verification]")
	
	// Check PPP active
	pppKey := "mikrotik:router-b:ppp:active"
	pppData, err := rdb.HGetAll(ctx, pppKey).Result()
	if err != nil {
		log.Printf("    ✗ PPP active read failed: %v", err)
	} else {
		log.Printf("    ✓ PPP active: %d entries", len(pppData))
	}

	// Check PPP secrets
	secretsKey := "mikrotik:router-b:ppp:secrets"
	secretsData, err := rdb.HGetAll(ctx, secretsKey).Result()
	if err != nil {
		log.Printf("    ✗ PPP secrets read failed: %v", err)
	} else {
		log.Printf("    ✓ PPP secrets: %d entries", len(secretsData))
		if len(secretsData) > 0 {
			i := 0
			for name := range secretsData {
				log.Printf("      - %s", name)
				i++
				if i >= 3 {
					log.Println("      ...")
					break
				}
			}
		}
	}

	// Check interface list
	ifKey := "mikrotik:router-b:interface:list"
	ifData, err := rdb.HGetAll(ctx, ifKey).Result()
	if err != nil {
		log.Printf("    ✗ Interface list read failed: %v", err)
	} else {
		log.Printf("    ✓ Interface list: %d entries", len(ifData))
	}

	// Summary
	log.Println("\n  [Summary]")
	log.Printf("    Pipeline A (InfluxDB): %s", map[bool]string{true: "✓ Data collected", false: "✗ No data"}[totalPoints > 0])
	log.Printf("    Pipeline B Tier 2 (Redis PPP Active): %d entries", len(pppData))
	log.Printf("    Pipeline B Tier 3 (Redis PPP Secrets): %d entries", len(secretsData))
	log.Printf("    Pipeline C (On-demand Interfaces): %d entries", len(ifData))
}

// Helper functions
func parseInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func formatHash(data map[string]string) string {
	pairs := make([]string, 0, len(data))
	for k, v := range data {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, ",")
}
