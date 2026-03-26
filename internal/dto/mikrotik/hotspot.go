package mikrotik

// AddHotspotProfileRequest is the request body for creating a hotspot profile.
type AddHotspotProfileRequest struct {
	Name              string `json:"name" binding:"required"`
	SharedUsers       string `json:"shared_users,omitempty"`
	RateLimit         string `json:"rate_limit,omitempty"`
	ExpiredMode       string `json:"expired_mode,omitempty"`
	ValidityTime      string `json:"validity_time,omitempty"`
	KeepAliveTimeout  string `json:"keepalive_timeout,omitempty"`
	StatusAutorefresh string `json:"status_autorefresh,omitempty"`
	OnLogin           string `json:"on_login,omitempty"`
	OnLogout          string `json:"on_logout,omitempty"`
}

// AddHotspotUserRequest is the request body for creating a hotspot user.
type AddHotspotUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password,omitempty"`
	Profile   string `json:"profile,omitempty"`
	Server    string `json:"server,omitempty"`
	LimitUptime string `json:"limit_uptime,omitempty"`
	LimitBytesTotal string `json:"limit_bytes_total,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Disabled  bool   `json:"disabled,omitempty"`
}
