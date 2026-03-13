package domain

// IPAddress represents an IP address entry from /ip/address
type IPAddress struct {
	ID              string `json:"id,omitempty"`
	Address         string `json:"address"`
	Network         string `json:"network,omitempty"`
	Interface       string `json:"interface"`
	Disabled        bool   `json:"disabled,omitempty"`
	Comment         string `json:"comment,omitempty"`
}

// IPPool represents an IP pool from /ip/pool
type IPPool struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Ranges   string `json:"ranges"`
	NextPool string `json:"nextPool,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// IPPoolUsed represents an allocated entry from /ip/pool/used
type IPPoolUsed struct {
	Pool    string `json:"pool"`
	Address string `json:"address"`
	Owner   string `json:"owner,omitempty"`
	Info    string `json:"info,omitempty"`
}
