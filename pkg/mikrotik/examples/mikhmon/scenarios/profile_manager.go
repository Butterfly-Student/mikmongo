package scenarios

import (
	"context"
	"fmt"

	"github.com/Butterfly-Student/go-ros/client"
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/Butterfly-Student/go-ros/examples/mikhmon/output"
	"github.com/Butterfly-Student/go-ros/repository/hotspot"
	mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
)

// RunProfileManager demonstrates profile creation with Mikhmon on-login script
func RunProfileManager(ctx context.Context, c *client.Client) {
	fmt.Println("=====================================")
	fmt.Println("  Profile Manager")
	fmt.Println("=====================================")
	fmt.Println()

	// Init repositories
	hotspotRepo := hotspot.NewRepository(c)
	profileRepo := mikhmonRepo.NewProfileRepository(hotspotRepo)

	// Create profile with Mikhmon config
	fmt.Println("Creating profile 'Test-1H' with Mikhmon on-login script...")
	req := &mikhmonDomain.ProfileRequest{
		Name:        "Test-1H",
		AddressPool: "hs-pool",
		RateLimit:   "1M/2M",
		SharedUsers: 1,
		Config: mikhmonDomain.ProfileConfig{
			Name:         "Test-1H",
			Price:        5000,
			SellingPrice: 7000,
			Validity:     "1h",
			ExpireMode:   mikhmonDomain.ExpireModeRemove,
			LockUser:     false,
			LockServer:   false,
		},
	}

	if err := profileRepo.CreateProfile(ctx, req); err != nil {
		fmt.Printf("Error creating profile: %v\n", err)
		return
	}

	fmt.Println("Profile 'Test-1H' created successfully!")

	// Show generated on-login script
	scriptData := &mikhmonDomain.OnLoginScriptData{
		Mode:         mikhmonDomain.ExpireModeRemove,
		Price:        5000,
		Validity:     "1h",
		SellingPrice: 7000,
		NoExp:        false,
		LockUser:     "Disable",
		LockServer:   "Disable",
	}
	script := profileRepo.GenerateOnLoginScript(scriptData)
	output.PrintOnLoginScript(script)

	// List existing profiles
	fmt.Println("\nListing existing hotspot profiles...")
	profiles, err := hotspotRepo.Profile().GetProfiles(ctx)
	if err != nil {
		fmt.Printf("Error listing profiles: %v\n", err)
		return
	}

	output.PrintProfileList(profiles)

	// Cleanup - remove test profile
	fmt.Println("\nCleaning up test profile...")
	for _, p := range profiles {
		if p.Name == "Test-1H" {
			if err := hotspotRepo.Profile().RemoveProfile(ctx, p.ID); err != nil {
				fmt.Printf("Error removing profile: %v\n", err)
			} else {
				fmt.Println("Profile 'Test-1H' removed successfully")
			}
			break
		}
	}

	fmt.Println("\nProfile manager test completed!")
}
