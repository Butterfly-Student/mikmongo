package mikhmon

import (
	"context"
	"fmt"
	"strings"

	"github.com/Butterfly-Student/go-ros/domain"
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/Butterfly-Student/go-ros/repository/hotspot"
)

// profileRepository implements ProfileRepository interface
type profileRepository struct {
	hotspotRepo hotspot.Repository
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(hr hotspot.Repository) ProfileRepository {
	return &profileRepository{hotspotRepo: hr}
}

// CreateProfile creates a new hotspot profile with Mikhmon on-login script
func (r *profileRepository) CreateProfile(ctx context.Context, req *mikhmonDomain.ProfileRequest) error {
	// Generate on-login script
	onLoginData := &mikhmonDomain.OnLoginScriptData{
		Mode:         req.Config.ExpireMode,
		Price:        req.Config.Price,
		Validity:     req.Config.Validity,
		SellingPrice: req.Config.SellingPrice,
		NoExp:        req.Config.ExpireMode == mikhmonDomain.ExpireModeNoExpire,
		LockUser:     boolToString(req.Config.LockUser),
		LockServer:   boolToString(req.Config.LockServer),
		ProfileName:  req.Name, // Profile name for recording script
	}

	onLoginScript := r.GenerateOnLoginScript(onLoginData)

	// Create profile using hotspot repository
	profile := &domain.UserProfile{
		Name:              req.Name,
		AddressPool:       req.AddressPool,
		RateLimit:         req.RateLimit,
		SharedUsers:       req.SharedUsers,
		ParentQueue:       req.ParentQueue,
		StatusAutorefresh: "1m",
		OnLogin:           onLoginScript,
	}

	_, err := r.hotspotRepo.Profile().AddProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

// UpdateProfile updates an existing hotspot profile
func (r *profileRepository) UpdateProfile(ctx context.Context, id string, req *mikhmonDomain.ProfileRequest) error {
	// Generate on-login script
	onLoginData := &mikhmonDomain.OnLoginScriptData{
		Mode:         req.Config.ExpireMode,
		Price:        req.Config.Price,
		Validity:     req.Config.Validity,
		SellingPrice: req.Config.SellingPrice,
		NoExp:        req.Config.ExpireMode == mikhmonDomain.ExpireModeNoExpire,
		LockUser:     boolToString(req.Config.LockUser),
		LockServer:   boolToString(req.Config.LockServer),
		ProfileName:  req.Name, // Profile name for recording script
	}

	onLoginScript := r.GenerateOnLoginScript(onLoginData)

	// Update profile using hotspot repository
	profile := &domain.UserProfile{
		Name:              req.Name,
		AddressPool:       req.AddressPool,
		RateLimit:         req.RateLimit,
		SharedUsers:       req.SharedUsers,
		ParentQueue:       req.ParentQueue,
		StatusAutorefresh: "1m",
		OnLogin:           onLoginScript,
	}

	err := r.hotspotRepo.Profile().UpdateProfile(ctx, id, profile)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// GenerateOnLoginScript generates the RouterOS on-login script
func (r *profileRepository) GenerateOnLoginScript(data *mikhmonDomain.OnLoginScriptData) string {
	var sb strings.Builder

	// Debug output with profile configuration
	sb.WriteString(fmt.Sprintf(":put (\",%s,%d,%s,%d,", data.Mode, data.Price, data.Validity, data.SellingPrice))
	if data.NoExp {
		sb.WriteString("noexp,")
	} else {
		sb.WriteString(",")
	}
	sb.WriteString(fmt.Sprintf("%s,%s,\");\n", data.LockUser, data.LockServer))

	// Set mode variable
	mode := "X" // Default to Remove
	if data.Mode == mikhmonDomain.ExpireModeNotify || data.Mode == mikhmonDomain.ExpireModeNotifyRecord {
		mode = "N" // Notify
	}

	sb.WriteString(fmt.Sprintf(":local mode \"%s\";\n", mode))
	sb.WriteString("\n")

	// Main script block
	sb.WriteString("{\n")
	sb.WriteString("    # Get current date\n")
	sb.WriteString("    :local date [ /system clock get date ];\n")
	sb.WriteString("    :local year [ :pick $date 7 11 ];\n")
	sb.WriteString("    :local month [ :pick $date 0 3 ];\n")
	sb.WriteString("    \n")
	sb.WriteString("    # Get user comment\n")
	sb.WriteString("    :local comment [ /ip hotspot user get [/ip hotspot user find where name=\"$user\"] comment];\n")
	sb.WriteString("    :local ucode [:pic $comment 0 2];\n")
	sb.WriteString("    \n")
	sb.WriteString("    # Check if user has code prefix (vc- or up-)\n")
	sb.WriteString("    :if ($ucode = \"vc\" or $ucode = \"up\" or $comment = \"\") do={\n")
	sb.WriteString("        # Create temporary scheduler to calculate expire date\n")
	sb.WriteString(fmt.Sprintf("        /sys sch add name=\"$user\" disable=no start-date=$date interval=\"%s\";\n", data.Validity))
	sb.WriteString("        :delay 2s;\n")
	sb.WriteString("        \n")
	sb.WriteString("        # Get next-run (expire date)\n")
	sb.WriteString("        :local exp [ /sys sch get [ /sys sch find where name=\"$user\" ] next-run];\n")
	sb.WriteString("        :local getxp [len $exp];\n")
	sb.WriteString("        \n")
	sb.WriteString("        # Format expire date based on length\n")
	sb.WriteString("        :if ($getxp = 15) do={\n")
	sb.WriteString("            # Format: jan/01/2024 12:00:00\n")
	sb.WriteString("            :local d [:pic $exp 0 6];\n")
	sb.WriteString("            :local t [:pic $exp 7 16];\n")
	sb.WriteString("            :local s (\"/\");\n")
	sb.WriteString("            :local exp (\"$d$s$year $t\");\n")
	sb.WriteString("            /ip hotspot user set comment=\"$exp $mode\" [find where name=\"$user\"];\n")
	sb.WriteString("        };\n")
	sb.WriteString("        \n")
	sb.WriteString("        :if ($getxp = 8) do={\n")
	sb.WriteString("            # Format: 12:00:00\n")
	sb.WriteString("            /ip hotspot user set comment=\"$date $exp $mode\" [find where name=\"$user\"];\n")
	sb.WriteString("        };\n")
	sb.WriteString("        \n")
	sb.WriteString("        :if ($getxp > 15) do={\n")
	sb.WriteString("            # Other format\n")
	sb.WriteString("            /ip hotspot user set comment=\"$exp $mode\" [find where name=\"$user\"];\n")
	sb.WriteString("        };\n")
	sb.WriteString("        \n")
	sb.WriteString("        # Remove temporary scheduler\n")
	sb.WriteString("        /sys sch remove [find where name=\"$user\"];\n")

	// Add recording script if mode is remc or ntfc
	if data.Mode == mikhmonDomain.ExpireModeRemoveRecord || data.Mode == mikhmonDomain.ExpireModeNotifyRecord {
		sb.WriteString("        \n")
		sb.WriteString("        # Recording script for report\n")
		sb.WriteString("        :local mac \"$mac-address\";\n")
		sb.WriteString("        :local time [/system clock get time ];\n")
		sb.WriteString(fmt.Sprintf("        /system script add name=\"$date-|-$time-|-$user-|-%d-|-$address-|-$mac-|-%s-|-%s-|-$comment\" owner=\"$month$year\" source=$date comment=mikhmon\n",
			data.Price, data.Validity, data.ProfileName))
	}

	sb.WriteString("    };\n")
	sb.WriteString("};\n")
	sb.WriteString("\n")

	// MAC Address Locking
	if data.LockUser == "Enable" {
		sb.WriteString("# MAC Address Locking\n")
		sb.WriteString("[:local mac \"$mac-address\"; /ip hotspot user set mac-address=$mac [find where name=$user]];\n")
		sb.WriteString("\n")
	}

	// Server Locking
	if data.LockServer == "Enable" {
		sb.WriteString("# Server Locking\n")
		sb.WriteString("[:local mac \"$mac-address\"; :local srv [/ip hotspot host get [find where mac-address=\"$mac\"] server]; /ip hotspot user set server=$srv [find where name=$user]]\n")
	}

	return sb.String()
}

func boolToString(b bool) string {
	if b {
		return "Enable"
	}
	return "Disable"
}
