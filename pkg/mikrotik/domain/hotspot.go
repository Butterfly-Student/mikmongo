package domain

// HotspotUser represents a hotspot user from MikroTik API
// Fields from /ip/hotspot/user/print: .id, address, comment, copy-from, disabled, email, limit-bytes-in, limit-bytes-out, limit-bytes-total, limit-uptime, mac-address, name, password, profile, routes, server, uptime, bytes-in, bytes-out, packets-in, packets-out
type HotspotUser struct {
	ID              string `json:".id,omitempty"`
	IPAddress       string `json:"address,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	Email           string `json:"email,omitempty"`
	LimitBytesIn    int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut   int64  `json:"limitBytesOut,omitempty"`
	LimitBytesTotal int64  `json:"limitBytesTotal,omitempty"`
	LimitUptime     string `json:"limitUptime,omitempty"`
	MACAddress      string `json:"macAddress,omitempty"`
	Name            string `json:"name"`
	Password        string `json:"password,omitempty"`
	Profile         string `json:"profile,omitempty"`
	Routes          string `json:"routes,omitempty"`
	Server          string `json:"server,omitempty"`
	OtpSecret       string `json:"otpSecret,omitempty"`

	//Read-Only fields
	Uptime     string `json:"uptime,omitempty"`
	BytesIn    int64  `json:"bytesIn,omitempty"`
	BytesOut   int64  `json:"bytesOut,omitempty"`
	PacketsIn  int64  `json:"packetsIn,omitempty"`
	PacketsOut int64  `json:"packetsOut,omitempty"`
}

// HotspotActive represents an active hotspot session
// Fields from /ip/hotspot/active/print: server, user, domain, address, uptime, idle-time, session-time-left, rx-rate, tx-rate
type HotspotActive struct {
	Server           string `json:"server,omitempty"`
	User             string `json:"user,omitempty"`
	Domain           string `json:"domain,omitempty"`
	Address          string `json:"address,omitempty"`
	MACAddress       string `json:"macAddress,omitempty"`
	LoginBy          string `json:"loginBy,omitempty"`
	Uptime           string `json:"uptime,omitempty"`
	IdleTime         string `json:"idleTime,omitempty"`
	SessionTimeLeft  string `json:"sessionTimeLeft,omitempty"`
	IdleTimeout      string `json:"idleTimeout,omitempty"`
	KeepaliveTimeout string `json:"keepaliveTimeout,omitempty"`
	BytesIn          int64  `json:"bytesIn,omitempty"`
	BytesOut         int64  `json:"bytesOut,omitempty"`
	LimitBytesIn     int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut    int64  `json:"limitBytesOut,omitempty"`
	LimitBytesTotal  int64  `json:"limitBytesTotal,omitempty"`
}

// HotspotHost represents a hotspot host
// Fields from /ip/hotspot/hosts/print: mac-address, address, to-address, server, idle-time, rx-rate, tx-rate
type HotspotHost struct {
	ID               string `json:".id,omitempty"`
	MACAddress       string `json:"macAddress,omitempty"`
	Address          string `json:"address,omitempty"`
	ToAddress        string `json:"toAddress,omitempty"`
	Server           string `json:"server,omitempty"`
	BridgePort       string `json:"bridgePort,omitempty"`
	Uptime           string `json:"uptime,omitempty"`
	IdleTime         string `json:"idleTime,omitempty"`
	IdleTimeout      string `json:"idleTimeout,omitempty"`
	KeepaliveTimeout string `json:"keepaliveTimeout,omitempty"`
	BytesIn          int64  `json:"bytesIn,omitempty"`
	BytesOut         int64  `json:"bytesOut,omitempty"`
	PacketsIn        int64  `json:"packetsIn,omitempty"`
	PacketsOut       int64  `json:"packetsOut,omitempty"`
}

// UserProfile represents a hotspot user profile
// Fields from /ip/hotspot/user/profile/add: .id, add-mac-cookie, address-list, address-pool, advertise, advertise-interval, advertise-timeout, advertise-url, copy-from, idle-timeout, incoming-filter, incoming-packet-mark, insert-queue-before, keepalive-timeout, mac-cookie-timeout, name, on-login, on-logout, open-status-page, outgoing-filter, outgoing-packet-mark, parent-queue, queue-type, rate-limit, session-timeout, shared-users, status-autorefresh, transparent-proxy
type UserProfile struct {
	ID                 string `json:".id,omitempty"`
	Name               string `json:"name"`
	AddressPool        string `json:"addressPool,omitempty"`
	SharedUsers        int    `json:"sharedUsers,omitempty"`
	RateLimit          string `json:"rateLimit,omitempty"`
	ParentQueue        string `json:"parentQueue,omitempty"`
	QueueType          string `json:"queueType,omitempty"`
	StatusAutorefresh  string `json:"statusAutorefresh,omitempty"`
	OnLogin            string `json:"onLogin,omitempty"`
	OnLogout           string `json:"onLogout,omitempty"`
	OpenStatusPage     string `json:"openStatusPage,omitempty"`
	TransparentProxy   bool   `json:"transparentProxy,omitempty"`
	Advertise          bool   `json:"advertise,omitempty"`
	AdvertiseInterval  string `json:"advertiseInterval,omitempty"`
	AdvertiseTimeout   string `json:"advertiseTimeout,omitempty"`
	AdvertiseURL       string `json:"advertiseURL,omitempty"`
	IdleTimeout        string `json:"idleTimeout,omitempty"`
	SessionTimeout     string `json:"sessionTimeout,omitempty"`
	KeepaliveTimeout   string `json:"keepaliveTimeout,omitempty"`
	MacCookieTimeout   string `json:"macCookieTimeout,omitempty"`
	AddMacCookie       bool   `json:"addMacCookie,omitempty"`
	AddressList        string `json:"addressList,omitempty"`
	IncomingFilter     string `json:"incomingFilter,omitempty"`
	IncomingPacketMark string `json:"incomingPacketMark,omitempty"`
	OutgoingFilter     string `json:"outgoingFilter,omitempty"`
	OutgoingPacketMark string `json:"outgoingPacketMark,omitempty"`
	InsertQueueBefore  string `json:"insertQueueBefore,omitempty"`
}

// HotspotIPBinding represents a /ip/hotspot/ip-binding entry
// Fields from /ip/hotspot/ip-binding/print: address, mac-address, to-address, server, type, comment, disabled
type HotspotIPBinding struct {
	ID         string `json:".id,omitempty"`
	Address    string `json:"address,omitempty"`
	MACAddress string `json:"macAddress,omitempty"`
	ToAddress  string `json:"toAddress,omitempty"`
	Server     string `json:"server,omitempty"`
	Type       string `json:"type,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Disabled   bool   `json:"disabled,omitempty"`
}

// HotspotCookie represents a hotspot cookie entry
// Fields from /ip/hotspot/cookie/print: user, domain, mac-address, expires-in
type HotspotCookie struct {
	User       string `json:"user,omitempty"`
	Domain     string `json:"domain,omitempty"`
	MACAddress string `json:"macAddress,omitempty"`
	ExpiresIn  string `json:"expiresIn,omitempty"`
}

// RemoveUserRequest represents a request to remove users
type RemoveUserRequest struct {
	IDs     []string `json:"ids,omitempty"`
	Comment string   `json:"comment,omitempty"`
}

// GetUsersRequest represents a request to get users
type GetUsersRequest struct {
	Profile string `json:"profile,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// UserFilter represents filter options for getting users
type UserFilter struct {
	Profile string
	Comment string
}

// BatchIDsRequest represents a request to perform batch operations on multiple IDs
type BatchIDsRequest struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

// HotspotProfile is an alias for UserProfile kept for backward compatibility.
type HotspotProfile = UserProfile
