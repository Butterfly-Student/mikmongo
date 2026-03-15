package scenarios

import (
	"context"
	"fmt"
	"strings"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/repository/system"
)

// RunReportViewer demonstrates viewing sales reports from /system/script
func RunReportViewer(ctx context.Context, c *client.Client) {
	fmt.Println("=====================================")
	fmt.Println("  Report Viewer")
	fmt.Println("=====================================")
	fmt.Println()

	// Init repository
	sysRepo := system.NewRepository(c)

	// Get all scripts
	fmt.Println("Fetching scripts from /system/script...")
	scripts, err := sysRepo.Scripts().GetScripts(ctx)
	if err != nil {
		fmt.Printf("Error fetching scripts: %v\n", err)
		return
	}

	// Filter mikhmon reports
	fmt.Println("\nMikhmon Sales Reports:")
	fmt.Println("----------------------")

	found := false
	for _, script := range scripts {
		if script.Comment == "mikhmon" || strings.Contains(script.Name, "mikhmon") {
			found = true
			fmt.Printf("\nScript: %s\n", script.Name)
			fmt.Printf("  Owner: %s\n", script.Owner)
			fmt.Printf("  Comment: %s\n", script.Comment)
			fmt.Printf("  Source: %s\n", script.Source)
		}
	}

	if !found {
		fmt.Println("No mikhmon reports found.")
		fmt.Println("Reports are created when users login with profiles that have recording enabled.")
	}

	// Show script count
	fmt.Printf("\nTotal scripts: %d\n", len(scripts))

	// Show non-mikhmon scripts (sample)
	fmt.Println("\nOther scripts (sample):")
	count := 0
	for _, script := range scripts {
		if script.Comment != "mikhmon" && !strings.Contains(script.Name, "mikhmon") {
			fmt.Printf("  - %s (owner: %s)\n", script.Name, script.Owner)
			count++
			if count >= 5 {
				fmt.Println("  ...")
				break
			}
		}
	}

	fmt.Println("\nReport viewer test completed!")
}
