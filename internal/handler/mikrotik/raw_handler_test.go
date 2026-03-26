package mikrotik

import "testing"

func TestValidateRawArgs_EmptyArgs(t *testing.T) {
	if msg := validateRawArgs(nil); msg != "args is required" {
		t.Errorf("expected 'args is required', got %q", msg)
	}
	if msg := validateRawArgs([]string{}); msg != "args is required" {
		t.Errorf("expected 'args is required', got %q", msg)
	}
}

func TestValidateRawArgs_TooManyArgs(t *testing.T) {
	args := make([]string, 21)
	for i := range args {
		args[i] = "arg"
	}
	if msg := validateRawArgs(args); msg != "too many args (max 20)" {
		t.Errorf("expected 'too many args', got %q", msg)
	}
}

func TestValidateRawArgs_AllowedCommands(t *testing.T) {
	allowed := [][]string{
		{"/interface/print"},
		{"/ip/address/print"},
		{"/ppp/active/print"},
		{"/ip/hotspot/active/print"},
		{"/system/resource/print"},
		{"/queue/simple/print"},
		{"/log/print", "=follow="},
	}
	for _, args := range allowed {
		if msg := validateRawArgs(args); msg != "" {
			t.Errorf("expected allowed for %v, got %q", args, msg)
		}
	}
}

func TestValidateRawArgs_BlockedCommands(t *testing.T) {
	blocked := [][]string{
		{"/system/reboot"},
		{"/system/shutdown"},
		{"/system/reset-configuration"},
		{"/system/backup"},
		{"/system/backup/save"},
		{"/system/export"},
		{"/user/add", "=name=backdoor"},
		{"/user/set", "=.id=admin"},
		{"/user/remove", "=.id=admin"},
		{"/password"},
		{"/certificate"},
	}
	for _, args := range blocked {
		if msg := validateRawArgs(args); msg == "" {
			t.Errorf("expected blocked for %v", args)
		}
	}
}

func TestValidateRawArgs_BlockedVerbs(t *testing.T) {
	blocked := [][]string{
		{"/ip/address/remove", "=.id=*1"},
		{"/interface/set", "=.id=ether1"},
		{"/ip/firewall/filter/add"},
		{"/interface/disable", "=.id=ether1"},
		{"/interface/enable", "=.id=ether1"},
	}
	for _, args := range blocked {
		if msg := validateRawArgs(args); msg == "" {
			t.Errorf("expected blocked verb for %v", args)
		}
	}
}

func TestValidateRawArgs_CaseInsensitive(t *testing.T) {
	if msg := validateRawArgs([]string{"/SYSTEM/REBOOT"}); msg == "" {
		t.Error("expected blocked for uppercase /SYSTEM/REBOOT")
	}
	if msg := validateRawArgs([]string{"/User/Add"}); msg == "" {
		t.Error("expected blocked for mixed case /User/Add")
	}
}
