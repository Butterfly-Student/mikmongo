package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/examples/mikhmon/scenarios"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("  Mikhmon Real Usage Examples")
	fmt.Println("  Router: 192.168.233.1:8728")
	fmt.Println("=====================================")
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create manager and connect using GetOrConnect
	manager := client.NewManager(nil)
	defer manager.CloseAll()

	cfg := client.Config{
		Host:     "192.168.233.1",
		Port:     8728,
		Username: "admin",
		Password: "r00t",
		Timeout:  10 * time.Second,
	}

	fmt.Println("Connecting to MikroTik using Manager.GetOrConnect...")
	c, err := manager.GetOrConnect(ctx, "primary", cfg)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected successfully!")
	fmt.Printf("Registered routers: %v\n", manager.Names())
	fmt.Println()

	// Show menu
	for {
		showMenu()
		choice := getInput("Select scenario (1-6): ")

		switch choice {
		case "1":
			scenarios.RunVoucherGenerator(ctx, c)
		case "2":
			scenarios.RunProfileManager(ctx, c)
		case "3":
			scenarios.RunMultiRouter(ctx)
		case "4":
			scenarios.RunReportViewer(ctx, c)
		case "5":
			scenarios.RunExpireMonitor(ctx, c)
		case "6":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		fmt.Println()
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("Available Scenarios:")
	fmt.Println("  1. Voucher Generator - Generate hotspot vouchers")
	fmt.Println("  2. Profile Manager - Create profiles with on-login script")
	fmt.Println("  3. Multi Router - Test multiple router connections")
	fmt.Println("  4. Report Viewer - View sales reports")
	fmt.Println("  5. Expire Monitor - Setup expire monitoring")
	fmt.Println("  6. Exit")
	fmt.Println()
}

func getInput(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}
