package model

// PPPProfileConfig for PPPoE profile configuration
type PPPProfileConfig struct {
	LocalAddress   string `json:"local_address,omitempty"`
	RemoteAddress  string `json:"remote_address,omitempty"` // IP pool name
	DNSServer      string `json:"dns_server,omitempty"`
	SessionTimeout string `json:"session_timeout,omitempty"`
	IdleTimeout    string `json:"idle_timeout,omitempty"`
	RateLimit      string `json:"rate_limit,omitempty"` // "10M/10M"
	UseCompression bool   `json:"use_compression,omitempty"`
	UseEncryption  bool   `json:"use_encryption,omitempty"`
	OnlyOne        bool   `json:"only_one,omitempty"`
	ChangeTCPMSS   bool   `json:"change_tcp_mss,omitempty"`
	Bridge         string `json:"bridge,omitempty"`
	AddressList    string `json:"address_list,omitempty"`
}
