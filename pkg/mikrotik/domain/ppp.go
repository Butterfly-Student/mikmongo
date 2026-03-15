package domain

// PPPSecret represents a PPP secret (user) from MikroTik API
// Fields from /ppp/secret/print: name, service, caller-id, password, profile, caller-id, routes, local-address, remote-address, ipv6-routes, limit-bytes-in, limit-bytes-out, last-logged-out, comment, last-caller-id, last-disconnect-reason, disabled
type PPPSecret struct {
	ID                   string `json:".id,omitempty"`
	Name                 string `json:"name"`
	Service              string `json:"service,omitempty"`
	CallerID             string `json:"callerID,omitempty"`
	Password             string `json:"password,omitempty"`
	Profile              string `json:"profile,omitempty"`
	Routes               string `json:"routes,omitempty"`
	LocalAddress         string `json:"localAddress,omitempty"`
	RemoteAddress        string `json:"remoteAddress,omitempty"`
	IPv6Routes           string `json:"ipv6Routes,omitempty"`
	LimitBytesIn         int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut        int64  `json:"limitBytesOut,omitempty"`
	LastLoggedOut        string `json:"lastLoggedOut,omitempty"`
	Comment              string `json:"comment,omitempty"`
	LastCallerID         string `json:"lastCallerID,omitempty"`
	LastDisconnectReason string `json:"lastDisconnectReason,omitempty"`
	Disabled             bool   `json:"disabled,omitempty"`
}

// PPPProfile represents a PPP profile from MikroTik API
// Fields from /ppp/profile/print: address-list, bridge, bridge-horizon, bridge-learning, bridge-path-cost, bridge-port-priority, change-tcp-mss, comment, copy-from, dns-server, idle-timeout, incoming-filter, insert-queue-before, interface-list, local-address, name, on-down, on-up, only-one, outgoing-filter, parent-queue, queue-type, rate-limit, remote-address, session-timeout, use-compression, use-encryption, use-mpls, use-upnp, wins-server
type PPPProfile struct {
	ID                 string `json:".id,omitempty"`
	AddressList        string `json:"addressList,omitempty"`
	Bridge             string `json:"bridge,omitempty"`
	BridgeHorizon      string `json:"bridgeHorizon,omitempty"`
	BridgeLearning     bool   `json:"bridgeLearning,omitempty"`
	BridgePathCost     int    `json:"bridgePathCost,omitempty"`
	BridgePortPriority int    `json:"bridgePortPriority,omitempty"`
	ChangeTCPMSS       bool   `json:"changeTCPMSS,omitempty"`
	Comment            string `json:"comment,omitempty"`
	DNSServer          string `json:"dnsServer,omitempty"`
	IdleTimeout        string `json:"idleTimeout,omitempty"`
	IncomingFilter     string `json:"incomingFilter,omitempty"`
	InsertQueueBefore  string `json:"insertQueueBefore,omitempty"`
	InterfaceList      string `json:"interfaceList,omitempty"`
	LocalAddress       string `json:"localAddress,omitempty"`
	Name               string `json:"name"`
	OnDown             string `json:"onDown,omitempty"`
	OnUp               string `json:"onUp,omitempty"`
	OnlyOne            bool   `json:"onlyOne,omitempty"`
	OutgoingFilter     string `json:"outgoingFilter,omitempty"`
	ParentQueue        string `json:"parentQueue,omitempty"`
	QueueType          string `json:"queueType,omitempty"`
	RateLimit          string `json:"rateLimit,omitempty"`
	RemoteAddress      string `json:"remoteAddress,omitempty"`
	SessionTimeout     string `json:"sessionTimeout,omitempty"`
	UseCompression     bool   `json:"useCompression,omitempty"`
	UseEncryption      bool   `json:"useEncryption,omitempty"`
	UseMPLS            bool   `json:"useMPLS,omitempty"`
	UseUPnP            bool   `json:"useUPnP,omitempty"`
	WinsServer         string `json:"winsServer,omitempty"`
}

// PPPActive represents an active PPP session from MikroTik API
// Fields from /ppp/active/print: name, service, caller-id, encoding, address, uptime
type PPPActive struct {
	Name     string `json:"name,omitempty"`
	Service  string `json:"service,omitempty"`
	CallerID string `json:"callerID,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	Address  string `json:"address,omitempty"`
	Uptime   string `json:"uptime,omitempty"`
}

// PPPSecretRequest represents a request to add a PPP secret
// Fields from /ppp/secret/add: caller-id, comment, copy-from, disabled, ipv6-routes, limit-bytes-in, limit-bytes-out, local-address, name, password, profile, remote-address, routes, service
type PPPSecretRequest struct {
	CallerID      string `json:"callerID,omitempty" validate:"max=100"`
	Comment       string `json:"comment,omitempty" validate:"max=200"`
	CopyFrom      string `json:"copyFrom,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
	IPv6Routes    string `json:"ipv6Routes,omitempty"`
	LimitBytesIn  int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut int64  `json:"limitBytesOut,omitempty"`
	LocalAddress  string `json:"localAddress,omitempty" validate:"max=50"`
	Name          string `json:"name" validate:"required,max=50"`
	Password      string `json:"password,omitempty" validate:"max=50"`
	Profile       string `json:"profile,omitempty" validate:"max=50"`
	RemoteAddress string `json:"remoteAddress,omitempty" validate:"max=50"`
	Routes        string `json:"routes,omitempty"`
	Service       string `json:"service,omitempty"`
}

// PPPSecretUpdateRequest represents a request to update a PPP secret
type PPPSecretUpdateRequest struct {
	CallerID      string `json:"callerID,omitempty" validate:"max=100"`
	Comment       string `json:"comment,omitempty" validate:"max=200"`
	Disabled      bool   `json:"disabled,omitempty"`
	IPv6Routes    string `json:"ipv6Routes,omitempty"`
	LimitBytesIn  int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut int64  `json:"limitBytesOut,omitempty"`
	LocalAddress  string `json:"localAddress,omitempty" validate:"max=50"`
	Name          string `json:"name,omitempty" validate:"max=50"`
	Password      string `json:"password,omitempty" validate:"max=50"`
	Profile       string `json:"profile,omitempty" validate:"max=50"`
	RemoteAddress string `json:"remoteAddress,omitempty" validate:"max=50"`
	Routes        string `json:"routes,omitempty"`
	Service       string `json:"service,omitempty"`
}

// PPPProfileRequest represents a request to add/update a PPP profile
// Fields from /ppp/profile/add: address-list, bridge, bridge-horizon, bridge-learning, bridge-path-cost, bridge-port-priority, change-tcp-mss, comment, copy-from, dns-server, idle-timeout, incoming-filter, insert-queue-before, interface-list, local-address, name, on-down, on-up, only-one, outgoing-filter, parent-queue, queue-type, rate-limit, remote-address, session-timeout, use-compression, use-encryption, use-mpls, use-upnp, wins-server
type PPPProfileRequest struct {
	AddressList        string `json:"addressList,omitempty" validate:"max=50"`
	Bridge             string `json:"bridge,omitempty" validate:"max=50"`
	BridgeHorizon      string `json:"bridgeHorizon,omitempty"`
	BridgeLearning     bool   `json:"bridgeLearning,omitempty"`
	BridgePathCost     int    `json:"bridgePathCost,omitempty"`
	BridgePortPriority int    `json:"bridgePortPriority,omitempty"`
	ChangeTCPMSS       bool   `json:"changeTCPMSS,omitempty"`
	Comment            string `json:"comment,omitempty" validate:"max=200"`
	CopyFrom           string `json:"copyFrom,omitempty"`
	DNSServer          string `json:"dnsServer,omitempty" validate:"max=100"`
	IdleTimeout        string `json:"idleTimeout,omitempty" validate:"max=20"`
	IncomingFilter     string `json:"incomingFilter,omitempty" validate:"max=50"`
	InsertQueueBefore  string `json:"insertQueueBefore,omitempty" validate:"max=50"`
	InterfaceList      string `json:"interfaceList,omitempty" validate:"max=50"`
	LocalAddress       string `json:"localAddress,omitempty" validate:"max=50"`
	Name               string `json:"name" validate:"required,max=50"`
	OnDown             string `json:"onDown,omitempty"`
	OnUp               string `json:"onUp,omitempty"`
	OnlyOne            bool   `json:"onlyOne,omitempty"`
	OutgoingFilter     string `json:"outgoingFilter,omitempty" validate:"max=50"`
	ParentQueue        string `json:"parentQueue,omitempty" validate:"max=50"`
	QueueType          string `json:"queueType,omitempty" validate:"max=50"`
	RateLimit          string `json:"rateLimit,omitempty" validate:"max=100"`
	RemoteAddress      string `json:"remoteAddress,omitempty" validate:"max=50"`
	SessionTimeout     string `json:"sessionTimeout,omitempty" validate:"max=20"`
	UseCompression     bool   `json:"useCompression,omitempty"`
	UseEncryption      bool   `json:"useEncryption,omitempty"`
	UseMPLS            bool   `json:"useMPLS,omitempty"`
	UseUPnP            bool   `json:"useUPnP,omitempty"`
	WinsServer         string `json:"winsServer,omitempty" validate:"max=100"`
}

// PPPoEUser is an alias kept for backward compatibility (old pkg stub).
// Use PPPActive for new code.
type PPPoEUser = PPPActive

// PPPoESecret is an alias kept for backward compatibility (old pkg stub).
// Use PPPSecret for new code.
type PPPoESecret = PPPSecret
