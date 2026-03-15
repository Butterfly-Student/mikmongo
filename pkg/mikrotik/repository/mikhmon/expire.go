package mikhmon

import (
	"context"
	"fmt"
	"strings"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/repository/system"
)

// expireRepository implements ExpireRepository interface
type expireRepository struct {
	client     *client.Client
	systemRepo system.Repository
}

// NewExpireRepository creates a new expire repository
func NewExpireRepository(c *client.Client, sr system.Repository) ExpireRepository {
	return &expireRepository{
		client:     c,
		systemRepo: sr,
	}
}

// SetupExpireMonitor sets up the expire monitor scheduler
func (r *expireRepository) SetupExpireMonitor(ctx context.Context) error {
	script := r.GenerateExpireMonitorScript()

	// Check if scheduler already exists
	reply, err := r.client.RunContext(ctx, "/system/scheduler/print", "?name=Mikhmon-Expire-Monitor")
	if err != nil {
		return fmt.Errorf("failed to check existing scheduler: %w", err)
	}

	if len(reply.Re) > 0 {
		// Scheduler exists, update it
		id := reply.Re[0].Map[".id"]
		_, err = r.client.RunContext(ctx,
			"/system/scheduler/set",
			"=.id="+id,
			"=interval=00:01:00",
			"=on-event="+script,
			"=disabled=no",
		)
		if err != nil {
			return fmt.Errorf("failed to update expire monitor: %w", err)
		}
	} else {
		// Create new scheduler
		_, err = r.client.RunContext(ctx,
			"/system/scheduler/add",
			"=name=Mikhmon-Expire-Monitor",
			"=start-time=00:00:00",
			"=interval=00:01:00",
			"=on-event="+script,
			"=disabled=no",
			"=comment=Mikhmon Expire Monitor",
		)
		if err != nil {
			return fmt.Errorf("failed to create expire monitor: %w", err)
		}
	}

	return nil
}

// DisableExpireMonitor disables the expire monitor scheduler
func (r *expireRepository) DisableExpireMonitor(ctx context.Context) error {
	reply, err := r.client.RunContext(ctx, "/system/scheduler/print", "?name=Mikhmon-Expire-Monitor")
	if err != nil {
		return fmt.Errorf("failed to find scheduler: %w", err)
	}

	if len(reply.Re) > 0 {
		id := reply.Re[0].Map[".id"]
		_, err = r.client.RunContext(ctx,
			"/system/scheduler/set",
			"=.id="+id,
			"=disabled=yes",
		)
		if err != nil {
			return fmt.Errorf("failed to disable expire monitor: %w", err)
		}
	}

	return nil
}

// IsExpireMonitorEnabled checks if expire monitor is enabled
func (r *expireRepository) IsExpireMonitorEnabled(ctx context.Context) (bool, error) {
	reply, err := r.client.RunContext(ctx, "/system/scheduler/print", "?name=Mikhmon-Expire-Monitor", "?disabled=false")
	if err != nil {
		return false, fmt.Errorf("failed to check scheduler status: %w", err)
	}

	return len(reply.Re) > 0, nil
}

// GenerateExpireMonitorScript generates the RouterOS expire monitor script
func (r *expireRepository) GenerateExpireMonitorScript() string {
	var sb strings.Builder

	sb.WriteString("# Function to convert date to integer format YYYYMMDD\n")
	sb.WriteString(":local dateint do={\n")
	sb.WriteString("    :local montharray ( \"jan\",\"feb\",\"mar\",\"apr\",\"may\",\"jun\",\"jul\",\"aug\",\"sep\",\"oct\",\"nov\",\"dec\" );\n")
	sb.WriteString("    :local days [ :pick $d 4 6 ];\n")
	sb.WriteString("    :local month [ :pick $d 0 3 ];\n")
	sb.WriteString("    :local year [ :pick $d 7 11 ];\n")
	sb.WriteString("    :local monthint ([ :find $montharray $month]);\n")
	sb.WriteString("    :local month ($monthint + 1);\n")
	sb.WriteString("    :if ( [len $month] = 1) do={\n")
	sb.WriteString("        :local zero (\"0\");\n")
	sb.WriteString("        :return [:tonum (\"$year$zero$month$days\")];\n")
	sb.WriteString("    } else={\n")
	sb.WriteString("        :return [:tonum (\"$year$month$days\")];\n")
	sb.WriteString("    }\n")
	sb.WriteString("};\n")
	sb.WriteString("\n")

	sb.WriteString("# Function to convert time to minutes\n")
	sb.WriteString(":local timeint do={\n")
	sb.WriteString("    :local hours [ :pick $t 0 2 ];\n")
	sb.WriteString("    :local minutes [ :pick $t 3 5 ];\n")
	sb.WriteString("    :return ($hours * 60 + $minutes);\n")
	sb.WriteString("};\n")
	sb.WriteString("\n")

	sb.WriteString("# Get current date and time\n")
	sb.WriteString(":local date [ /system clock get date ];\n")
	sb.WriteString(":local time [ /system clock get time ];\n")
	sb.WriteString(":local today [$dateint d=$date];\n")
	sb.WriteString(":local curtime [$timeint t=$time];\n")
	sb.WriteString("\n")

	sb.WriteString("# Get current year and last year\n")
	sb.WriteString(":local tyear [ :pick $date 7 11 ];\n")
	sb.WriteString(":local lyear ($tyear-1);\n")
	sb.WriteString("\n")

	sb.WriteString("# Loop all users with comment containing year\n")
	sb.WriteString(":foreach i in [ /ip hotspot user find where comment~\"/\" . $tyear || comment~\"/\" . $lyear ] do={\n")
	sb.WriteString("    :local comment [ /ip hotspot user get $i comment];\n")
	sb.WriteString("    :local limit [ /ip hotspot user get $i limit-uptime];\n")
	sb.WriteString("    :local name [ /ip hotspot user get $i name];\n")
	sb.WriteString("    :local gettime [:pic $comment 12 20];\n")
	sb.WriteString("    \n")
	sb.WriteString("    # Check comment format (must have / at position 3 and 6)\n")
	sb.WriteString("    :if ([:pic $comment 3] = \"/\" and [:pic $comment 6] = \"/\") do={\n")
	sb.WriteString("        :local expd [$dateint d=$comment];\n")
	sb.WriteString("        :local expt [$timeint t=$gettime];\n")
	sb.WriteString("        \n")
	sb.WriteString("        # Check expired condition\n")
	sb.WriteString("        :if (($expd < $today and $expt < $curtime) or \n")
	sb.WriteString("              ($expd < $today and $expt > $curtime) or \n")
	sb.WriteString("              ($expd = $today and $expt < $curtime) and $limit != \"00:00:01\") do={\n")
	sb.WriteString("            \n")
	sb.WriteString("            # Mode N = Notify (disable user)\n")
	sb.WriteString("            :if ([:pic $comment 21] = \"N\") do={\n")
	sb.WriteString("                [ /ip hotspot user set limit-uptime=1s $i ];\n")
	sb.WriteString("                [ /ip hotspot active remove [find where user=$name] ];\n")
	sb.WriteString("            } else={\n")
	sb.WriteString("                # Mode X = Remove (remove user)\n")
	sb.WriteString("                [ /ip hotspot user remove $i ];\n")
	sb.WriteString("                [ /ip hotspot active remove [find where user=$name] ];\n")
	sb.WriteString("            }\n")
	sb.WriteString("        }\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	return sb.String()
}
