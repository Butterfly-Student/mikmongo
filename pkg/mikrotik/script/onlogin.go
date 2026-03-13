package script

import (
	"fmt"
	"regexp"
	"strings"

	"mikmongo/pkg/mikrotik/domain"
)

// OnLoginGenerator generates and parses on-login scripts for MikroTik user profiles
type OnLoginGenerator struct{}

// NewOnLoginGenerator creates a new on-login script generator
func NewOnLoginGenerator() *OnLoginGenerator {
	return &OnLoginGenerator{}
}

// ScriptComponents holds all parts of the on-login script
type ScriptComponents struct {
	Header     string
	Expiration string
	Record     string
	LockUser   string
	LockServer string
	Footer     string
}

// Generate creates an on-login script for user profile
func (g *OnLoginGenerator) Generate(req *domain.ProfileRequest) string {
	components := g.buildScriptComponents(req)
	return g.assembleScript(components)
}

func (g *OnLoginGenerator) buildScriptComponents(req *domain.ProfileRequest) *ScriptComponents {
	comp := &ScriptComponents{}
	comp.Header = g.buildHeader(req)
	comp.Expiration = g.buildExpirationLogic(req)
	if req.ExpireMode == "remc" || req.ExpireMode == "ntfc" {
		comp.Record = g.buildRecordScript(req)
	}
	if req.LockUser == "Enable" {
		comp.LockUser = g.buildLockUserScript()
	}
	if req.LockServer == "Enable" {
		comp.LockServer = g.buildLockServerScript()
	}
	comp.Footer = g.buildFooter(req)
	return comp
}

func (g *OnLoginGenerator) buildHeader(req *domain.ProfileRequest) string {
	var mode string
	switch req.ExpireMode {
	case "ntf", "ntfc":
		mode = "N"
	case "rem", "remc":
		mode = "X"
	default:
		mode = ""
	}
	return fmt.Sprintf(
		`:put (",%s,%.0f,%s,%.0f,,%s,%s,"); :local mode "%s"; {`,
		req.ExpireMode,
		req.Price,
		req.Validity,
		req.SellingPrice,
		req.LockUser,
		req.LockServer,
		mode,
	)
}

func (g *OnLoginGenerator) buildExpirationLogic(req *domain.ProfileRequest) string {
	if req.ExpireMode == "0" || req.ExpireMode == "" {
		return ""
	}
	return fmt.Sprintf(`
    :local date [/system clock get date];
    :local year [:pick $date 7 11];
    :local month [:pick $date 0 3];
    :local comment [/ip hotspot user get [/ip hotspot user find where name="$user"] comment];
    :local ucode [:pic $comment 0 2];

    :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
        /sys sch add name="$user" disable=no start-date=$date interval="%s";
        :delay 2s;
        :local exp [/sys sch get [/sys sch find where name="$user"] next-run];
        :local getxp [len $exp];

        :if ($getxp = 16) do={
            :local d [:pic $exp 0 6];
            :local t [:pic $exp 7 15];
            :local s ("/");
            :local exp ("$d$s$year $t");
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };

        :if ($getxp = 8) do={
            /ip hotspot user set comment="$date $exp $mode" [find where name="$user"];
        };

        :if ($getxp > 16) do={
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };

        /sys sch remove [find where name="$user"];
    }
`, req.Validity)
}

func (g *OnLoginGenerator) buildRecordScript(req *domain.ProfileRequest) string {
	return fmt.Sprintf(`
    :local mac $"mac-address";
    :local time [/system clock get time];
    /system script add name="$date-|-$time-|-$user-|-%.0f-|-$address-|-$mac-|-%s-|-%s-|-$comment" owner="$month$year" source=$date comment=mikhmon;`,
		req.Price,
		req.Validity,
		req.Name,
	)
}

func (g *OnLoginGenerator) buildLockUserScript() string {
	return `; [:local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]`
}

func (g *OnLoginGenerator) buildLockServerScript() string {
	return `; [:local mac $"mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]]`
}

func (g *OnLoginGenerator) buildFooter(req *domain.ProfileRequest) string {
	switch req.ExpireMode {
	case "rem", "remc", "ntf", "ntfc":
		return "}"
	case "0":
		if req.LockUser == "Enable" || req.LockServer == "Enable" {
			return "}"
		}
		return ""
	default:
		return "}"
	}
}

func (g *OnLoginGenerator) assembleScript(comp *ScriptComponents) string {
	var parts []string
	parts = append(parts, comp.Header)
	if comp.Expiration != "" {
		parts = append(parts, comp.Expiration)
	}
	if comp.Record != "" {
		parts = append(parts, comp.Record)
	}
	if comp.LockUser != "" {
		parts = append(parts, comp.LockUser)
	}
	if comp.LockServer != "" {
		parts = append(parts, comp.LockServer)
	}
	if comp.Footer != "" {
		parts = append(parts, comp.Footer)
	}
	return strings.TrimSpace(strings.Join(parts, ""))
}

// Parse extracts Mikhmon metadata from an existing on-login script
func (g *OnLoginGenerator) Parse(scriptStr string) *domain.ProfileRequest {
	req := &domain.ProfileRequest{}
	if scriptStr == "" {
		return req
	}

	putPattern := regexp.MustCompile(`:put \(",([\w]*),([\d\.]*),([^,]*),([\d\.]*),,([^,]*),([^,]*),"\)`)
	matches := putPattern.FindStringSubmatch(scriptStr)

	if len(matches) >= 7 {
		req.ExpireMode = matches[1]
		req.Price = parseFloat(matches[2])
		req.Validity = matches[3]
		req.SellingPrice = parseFloat(matches[4])
		req.LockUser = matches[5]
		req.LockServer = matches[6]
	}

	if strings.Contains(scriptStr, "/system script add") {
		if req.ExpireMode == "rem" {
			req.ExpireMode = "remc"
		} else if req.ExpireMode == "ntf" {
			req.ExpireMode = "ntfc"
		}
	}

	return req
}

// GenerateExpiredAction generates the action script for when user expires
func (g *OnLoginGenerator) GenerateExpiredAction(expireMode string) string {
	switch expireMode {
	case "rem", "remc":
		return "/ip hotspot user remove [find name=$user]"
	case "ntf", "ntfc":
		return "/ip hotspot user set limit-uptime=1s [find name=$user]"
	default:
		return ""
	}
}

// GenerateExpireMonitorScript generates the global scheduler script used by
// "Mikhmon-Expire-Monitor" to enforce expired users handling.
func (g *OnLoginGenerator) GenerateExpireMonitorScript() string {
	return `:local dateint do={:local montharray ("jan","feb","mar","apr","may","jun","jul","aug","sep","oct","nov","dec"); :local days [:pick $d 4 6]; :local month [:pick $d 0 3]; :local year [:pick $d 7 11]; :local monthint ([:find $montharray $month]); :local month ($monthint + 1); :if ([len $month] = 1) do={:local zero ("0"); :return [:tonum ("$year$zero$month$days")];} else={:return [:tonum ("$year$month$days")];}}; :local timeint do={:local hours [:pick $t 0 2]; :local minutes [:pick $t 3 5]; :return ($hours * 60 + $minutes);}; :local date [/system clock get date]; :local time [/system clock get time]; :local today [$dateint d=$date]; :local curtime [$timeint t=$time]; :local tyear [:pick $date 7 11]; :local lyear ($tyear - 1); :foreach i in=[/ip hotspot user find where comment~"/$tyear" || comment~"/$lyear"] do={:local comment [/ip hotspot user get $i comment]; :local limit [/ip hotspot user get $i limit-uptime]; :local name [/ip hotspot user get $i name]; :local gettime [:pick $comment 12 20]; :if ([:pick $comment 3] = "/" and [:pick $comment 6] = "/") do={:local expd [$dateint d=$comment]; :local expt [$timeint t=$gettime]; :if ((($expd < $today and $expt < $curtime) or ($expd < $today and $expt > $curtime) or ($expd = $today and $expt < $curtime)) and $limit != "00:00:01") do={:if ([:pick $comment 21] = "N") do={/ip hotspot user set limit-uptime=1s $i; /ip hotspot active remove [find where user=$name];} else={/ip hotspot user remove $i; /ip hotspot active remove [find where user=$name];}}}}`
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}
