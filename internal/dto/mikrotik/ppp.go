package mikrotik

// AddPPPProfileRequest is the request body for creating a PPP profile.
type AddPPPProfileRequest struct {
	Name          string `json:"name" binding:"required"`
	LocalAddress  string `json:"local_address,omitempty"`
	RemoteAddress string `json:"remote_address,omitempty"`
	RateLimit     string `json:"rate_limit,omitempty"`
	OnlyOne       string `json:"only_one,omitempty"`
	Comment       string `json:"comment,omitempty"`
}

// AddPPPSecretRequest is the request body for creating a PPP secret.
type AddPPPSecretRequest struct {
	Name          string `json:"name" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Profile       string `json:"profile,omitempty"`
	Service       string `json:"service,omitempty"`
	CallerID      string `json:"caller_id,omitempty"`
	LocalAddress  string `json:"local_address,omitempty"`
	RemoteAddress string `json:"remote_address,omitempty"`
	Routes        string `json:"routes,omitempty"`
	Comment       string `json:"comment,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
}
