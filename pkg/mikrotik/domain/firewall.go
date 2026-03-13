package domain

// NATRule represents a firewall NAT rule
type NATRule struct {
	ID              string `json:"id,omitempty"`
	Chain           string `json:"chain,omitempty"`
	Action          string `json:"action,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	SrcAddress      string `json:"srcAddress,omitempty"`
	DstAddress      string `json:"dstAddress,omitempty"`
	SrcPort         string `json:"srcPort,omitempty"`
	DstPort         string `json:"dstPort,omitempty"`
	InInterface     string `json:"inInterface,omitempty"`
	OutInterface    string `json:"outInterface,omitempty"`
	ToAddresses     string `json:"toAddresses,omitempty"`
	ToPorts         string `json:"toPorts,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Dynamic         bool   `json:"dynamic,omitempty"`
	Invalid         bool   `json:"invalid,omitempty"`
	Bytes           int64  `json:"bytes,omitempty"`
	Packets         int64  `json:"packets,omitempty"`
	ConnectionBytes int64  `json:"connectionBytes,omitempty"`
}

// FirewallRule represents a firewall filter rule
type FirewallRule struct {
	ID           string `json:"id,omitempty"`
	Chain        string `json:"chain,omitempty"`
	Action       string `json:"action,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
	SrcAddress   string `json:"srcAddress,omitempty"`
	DstAddress   string `json:"dstAddress,omitempty"`
	SrcPort      string `json:"srcPort,omitempty"`
	DstPort      string `json:"dstPort,omitempty"`
	InInterface  string `json:"inInterface,omitempty"`
	OutInterface string `json:"outInterface,omitempty"`
	Comment      string `json:"comment,omitempty"`
	Disabled     bool   `json:"disabled,omitempty"`
}

// AddressList represents a firewall address list entry
type AddressList struct {
	ID       string `json:"id,omitempty"`
	List     string `json:"list,omitempty"`
	Address  string `json:"address,omitempty"`
	Timeout  string `json:"timeout,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}
