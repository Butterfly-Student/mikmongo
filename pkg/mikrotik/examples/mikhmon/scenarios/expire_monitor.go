package scenarios

import (
	"context"
	"fmt"

	"github.com/Butterfly-Student/go-ros/client"
	mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
	"github.com/Butterfly-Student/go-ros/repository/system"
)

// RunExpireMonitor demonstrates expire monitor setup and management
func RunExpireMonitor(ctx context.Context, c *client.Client) {
	fmt.Println("=====================================")
	fmt.Println("  Expire Monitor Manager")
	fmt.Println("=====================================")
	fmt.Println()

	// Init repositories
	sysRepo := system.NewRepository(c)
	expireRepo := mikhmonRepo.NewExpireRepository(c, sysRepo)

	// Check current status
	fmt.Println("Checking expire monitor status...")
	enabled, err := expireRepo.IsExpireMonitorEnabled(ctx)
	if err != nil {
		fmt.Printf("Error checking status: %v\n", err)
		return
	}

	if enabled {
		fmt.Println("✓ Expire monitor is currently ENABLED")
	} else {
		fmt.Println("✗ Expire monitor is currently DISABLED")
	}

	// Show generated script
	fmt.Println("\nGenerated Expire Monitor Script:")
	fmt.Println("-----------------------------------")
	script := expireRepo.GenerateExpireMonitorScript()
	lines := splitLines(script)
	for i, line := range lines {
		if i < 30 {
			fmt.Println(line)
		} else if i == 30 {
			fmt.Println("... (script continues)")
			break
		}
	}
	fmt.Printf("\nTotal script lines: %d\n", len(lines))

	// Setup expire monitor
	fmt.Println("\nSetting up expire monitor...")
	if err := expireRepo.SetupExpireMonitor(ctx); err != nil {
		fmt.Printf("Error setting up expire monitor: %v\n", err)
		return
	}
	fmt.Println("✓ Expire monitor setup successfully!")

	// Verify status after setup
	enabled, err = expireRepo.IsExpireMonitorEnabled(ctx)
	if err != nil {
		fmt.Printf("Error verifying status: %v\n", err)
		return
	}

	if enabled {
		fmt.Println("✓ Expire monitor is now ENABLED and running")
	} else {
		fmt.Println("✗ Expire monitor is still DISABLED")
	}

	// Show scheduler info
	fmt.Println("\nFetching scheduler info...")
	schedulers, err := sysRepo.Scheduler().GetSchedulers(ctx)
	if err != nil {
		fmt.Printf("Error fetching schedulers: %v\n", err)
		return
	}

	fmt.Println("\nAll Schedulers:")
	for _, s := range schedulers {
		status := "enabled"
		if s.Disabled {
			status = "disabled"
		}
		fmt.Printf("  - %s (interval: %s, status: %s)\n", s.Name, s.Interval, status)
	}

	fmt.Println("\nExpire monitor test completed!")
	fmt.Println("\nNote: The expire monitor runs every 1 minute to check for expired users.")
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
