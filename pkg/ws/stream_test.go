package ws

import "testing"

func TestValidateInterfaceName_Valid(t *testing.T) {
	valid := []string{
		"ether1",
		"wlan1",
		"bridge-local",
		"vlan100",
		"pppoe-out1",
		"ether1.100",
		"my_interface",
		"a",
	}
	for _, name := range valid {
		if !ValidateInterfaceName(name) {
			t.Errorf("expected %q to be valid", name)
		}
	}
}

func TestValidateInterfaceName_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"ether 1",                                                                    // space
		"ether;drop",                                                                  // semicolon
		"$(cmd)",                                                                      // injection
		"a/b",                                                                         // slash
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // >64 chars
	}
	for _, name := range invalid {
		if ValidateInterfaceName(name) {
			t.Errorf("expected %q to be invalid", name)
		}
	}
}
