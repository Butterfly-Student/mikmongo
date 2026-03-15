package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
)

// RunMultiRouter demonstrates multiple router connection management
func RunMultiRouter(ctx context.Context) {
	fmt.Println("=====================================")
	fmt.Println("  Multi Router Test")
	fmt.Println("=====================================")
	fmt.Println()

	// Create manager
	manager := client.NewManager(nil)
	defer manager.CloseAll()

	// Register primary router
	fmt.Println("Registering router 'primary' (192.168.233.2)...")
	cfg1 := client.Config{
		Host:     "192.168.233.2",
		Port:     8728,
		Username: "admin",
		Password: "r00t",
		Timeout:  10 * time.Second,
	}

	if err := manager.Register(ctx, "primary", cfg1); err != nil {
		fmt.Printf("Failed to register primary router: %v\n", err)
		return
	}
	fmt.Println("Router 'primary' registered successfully")

	// Test connection
	fmt.Println("\nTesting connection...")
	c, err := manager.Get("primary")
	if err != nil {
		fmt.Printf("Failed to get client: %v\n", err)
		return
	}

	// Run simple command
	reply, err := c.Run("/system/identity/print")
	if err != nil {
		fmt.Printf("Failed to run command: %v\n", err)
		return
	}

	fmt.Println("Router identity:")
	for _, re := range reply.Re {
		if name, ok := re.Map["name"]; ok {
			fmt.Printf("  Name: %s\n", name)
		}
	}

	// List all registered routers
	fmt.Println("\nRegistered routers:")
	names := manager.Names()
	for _, name := range names {
		fmt.Printf("  - %s\n", name)
	}

	// Test GetOrConnect
	fmt.Println("\nTesting GetOrConnect...")
	c2, err := manager.GetOrConnect(ctx, "primary", cfg1)
	if err != nil {
		fmt.Printf("GetOrConnect failed: %v\n", err)
		return
	}

	if c == c2 {
		fmt.Println("GetOrConnect returned existing client (correct)")
	} else {
		fmt.Println("GetOrConnect created new client")
	}

	fmt.Println("\nMulti router test completed!")
}
