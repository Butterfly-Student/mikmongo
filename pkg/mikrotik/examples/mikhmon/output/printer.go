package output

import (
	"fmt"
	"strings"

	"github.com/Butterfly-Student/go-ros/domain"
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

// PrintVoucherBatch prints a batch of vouchers
func PrintVoucherBatch(batch *mikhmonDomain.VoucherBatch, title string) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Batch Code: %s\n", batch.Code)
	fmt.Printf("Quantity: %d\n", batch.Quantity)
	fmt.Printf("Profile: %s\n", batch.Profile)
	fmt.Printf("Server: %s\n", batch.Server)
	fmt.Printf("Time Limit: %s\n", batch.TimeLimit)
	fmt.Printf("Data Limit: %s\n", batch.DataLimit)
	fmt.Println(strings.Repeat("-", 60))

	for i, v := range batch.Vouchers {
		if batch.Vouchers[0].Mode == "vc" {
			fmt.Printf("%d. Username: %s (Password: %s)\n", i+1, v.Name, v.Password)
		} else {
			fmt.Printf("%d. Username: %s | Password: %s\n", i+1, v.Name, v.Password)
		}
	}
}

// PrintVoucherList prints a list of vouchers
func PrintVoucherList(vouchers []*mikhmonDomain.Voucher) {
	if len(vouchers) == 0 {
		fmt.Println("No vouchers found.")
		return
	}

	fmt.Println("\nVoucher List:")
	fmt.Println(strings.Repeat("-", 60))
	for i, v := range vouchers {
		fmt.Printf("%d. %s | Profile: %s | Comment: %s\n",
			i+1, v.Name, v.Profile, v.Comment)
	}
}

// PrintOnLoginScript prints the generated on-login script
func PrintOnLoginScript(script string) {
	fmt.Println("\nGenerated On-Login Script:")
	fmt.Println(strings.Repeat("-", 60))
	lines := strings.Split(script, "\n")
	for i, line := range lines {
		if i < 20 {
			fmt.Println(line)
		} else if i == 20 {
			fmt.Println("... (script continues)")
			break
		}
	}
	fmt.Printf("\nTotal lines: %d\n", len(lines))
}

// PrintProfileList prints a list of hotspot profiles
func PrintProfileList(profiles []*domain.UserProfile) {
	if len(profiles) == 0 {
		fmt.Println("No profiles found.")
		return
	}

	fmt.Println("\nHotspot Profiles:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-20s %-15s %-15s %-10s\n", "Name", "Address Pool", "Rate Limit", "Shared")
	fmt.Println(strings.Repeat("-", 80))

	for _, p := range profiles {
		fmt.Printf("%-20s %-15s %-15s %-10d\n",
			p.Name, p.AddressPool, p.RateLimit, p.SharedUsers)
	}
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Printf("✗ %s\n", message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Printf("ℹ %s\n", message)
}
