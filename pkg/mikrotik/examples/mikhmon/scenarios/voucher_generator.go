package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/Butterfly-Student/go-ros/examples/mikhmon/output"
	"github.com/Butterfly-Student/go-ros/repository/hotspot"
	mikhmonRepo "github.com/Butterfly-Student/go-ros/repository/mikhmon"
)

// RunVoucherGenerator demonstrates voucher generation
func RunVoucherGenerator(ctx context.Context, c *client.Client) {
	fmt.Println("=====================================")
	fmt.Println("  Voucher Generator")
	fmt.Println("=====================================")
	fmt.Println()

	// Init repositories
	hotspotRepo := hotspot.NewRepository(c)
	generatorRepo := mikhmonRepo.NewGeneratorRepository()
	voucherRepo := mikhmonRepo.NewVoucherRepository(c, hotspotRepo, generatorRepo)

	// Generate vouchers - Mode VC (Voucher Card)
	fmt.Println("Generating 5 Voucher Cards (username = password)...")
	vcReq := &mikhmonDomain.VoucherGenerateRequest{
		Quantity:   5,
		Profile:    "default",
		Mode:       mikhmonDomain.VoucherModeVoucher,
		NameLength: 6,
		CharSet:    mikhmonDomain.CharSetUpplow1,
		TimeLimit:  "1h",
		DataLimit:  "1G",
	}

	vcBatch, err := voucherRepo.GenerateBatch(ctx, vcReq)
	if err != nil {
		fmt.Printf("Error generating VC batch: %v\n", err)
		return
	}

	output.PrintVoucherBatch(vcBatch, "Voucher Cards (VC)")

	// Generate vouchers - Mode UP (User/Password)
	fmt.Println("\nGenerating 5 User/Password vouchers...")
	upReq := &mikhmonDomain.VoucherGenerateRequest{
		Quantity:   5,
		Profile:    "default",
		Mode:       mikhmonDomain.VoucherModeUserPassword,
		NameLength: 6,
		CharSet:    mikhmonDomain.CharSetUpplow1,
		TimeLimit:  "2h",
		DataLimit:  "2G",
	}

	upBatch, err := voucherRepo.GenerateBatch(ctx, upReq)
	if err != nil {
		fmt.Printf("Error generating UP batch: %v\n", err)
		return
	}

	output.PrintVoucherBatch(upBatch, "User/Password (UP)")

	// List vouchers by comment
	fmt.Println("\nListing vouchers by comment...")
	vouchers, err := voucherRepo.GetVouchersByComment(ctx, vcBatch.Vouchers[0].Comment)
	if err != nil {
		fmt.Printf("Error listing vouchers: %v\n", err)
		return
	}
	output.PrintVoucherList(vouchers)

	// Cleanup
	fmt.Println("\nCleaning up test vouchers...")
	time.Sleep(2 * time.Second)

	if err := voucherRepo.RemoveVoucherBatch(ctx, vcBatch.Vouchers[0].Comment); err != nil {
		fmt.Printf("Error removing VC batch: %v\n", err)
	} else {
		fmt.Println("VC batch removed successfully")
	}

	if err := voucherRepo.RemoveVoucherBatch(ctx, upBatch.Vouchers[0].Comment); err != nil {
		fmt.Printf("Error removing UP batch: %v\n", err)
	} else {
		fmt.Println("UP batch removed successfully")
	}

	fmt.Println("\nVoucher generator test completed!")
}
