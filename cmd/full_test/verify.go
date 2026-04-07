// cmd/full_test/verify.go - Quick verification after test
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Setup clients
	influxClient := influxdb2.NewClient("http://localhost:8086", "my-super-secret-token")
	defer influxClient.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 0})
	defer redisClient.Close()

	log.Println("Verifying data in storage...")
	
	// Query InfluxDB
	queryAPI := influxClient.QueryAPI("myorg")
	query := `
		from(bucket: "miktik")
		|> range(start: -5m)
		|> filter(fn: (r) => r._measurement == "interface_stats")
		|> group(columns: ["name"])
		|> count()
		|> yield(name: "count")
	`
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		log.Printf("InfluxDB query error: %v", err)
	} else {
		log.Println("\n📊 InfluxDB - Interface Stats by Interface:")
		for result.Next() {
			record := result.Record()
			log.Printf("  - %s: %v points", record.ValueByKey("name"), record.Value())
		}
	}

	// Query Redis
	log.Println("\n📦 Redis - Keys:")
	keys, err := redisClient.Keys(ctx, "mikrotik:router-b:*").Result()
	if err != nil {
		log.Printf("  Error: %v", err)
	} else if len(keys) == 0 {
		log.Println("  (No keys found - TTL expired)")
	} else {
		for _, key := range keys {
			typ, _ := redisClient.Type(ctx, key).Result()
			if typ == "hash" {
				count, _ := redisClient.HLen(ctx, key).Result()
				log.Printf("  - %s: %d entries", key, count)
			} else {
				log.Printf("  - %s: %s", key, typ)
			}
		}
	}

	fmt.Println("\n✅ Verification complete")
}
