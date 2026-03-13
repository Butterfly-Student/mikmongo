package domain

// SystemResource represents MikroTik system resource
type SystemResource struct {
	Uptime               string  `json:"uptime,omitempty"`
	Version              string  `json:"version,omitempty"`
	BuildTime            string  `json:"buildTime,omitempty"`
	FreeMemory           int64   `json:"freeMemory,omitempty"`
	TotalMemory          int64   `json:"totalMemory,omitempty"`
	FreeHddSpace         int64   `json:"freeHddSpace,omitempty"`
	TotalHddSpace        int64   `json:"totalHddSpace,omitempty"`
	WriteSectSinceReboot int64   `json:"writeSectSinceReboot,omitempty"`
	WriteSectTotal       int64   `json:"writeSectTotal,omitempty"`
	BadBlocks            float64 `json:"badBlocks,omitempty"`
	ArchitectureName     string  `json:"architectureName,omitempty"`
	BoardName            string  `json:"boardName,omitempty"`
	Platform             string  `json:"platform,omitempty"`
	Cpu                  string  `json:"cpu,omitempty"`
	CpuCount             int     `json:"cpuCount,omitempty"`
	CpuFrequency         int     `json:"cpuFrequency,omitempty"`
	CpuLoad              int     `json:"cpuLoad,omitempty"`
}

// SystemHealth represents MikroTik system health
type SystemHealth struct {
	Voltage     string `json:"voltage,omitempty"`
	Temperature string `json:"temperature,omitempty"`
	FanSpeed    string `json:"fanSpeed,omitempty"`
	FanSpeed2   string `json:"fanSpeed2,omitempty"`
	FanSpeed3   string `json:"fanSpeed3,omitempty"`
}

// SystemIdentity represents MikroTik system identity
type SystemIdentity struct {
	Name string `json:"name,omitempty"`
}

// SystemClock represents MikroTik system clock
type SystemClock struct {
	Time         string `json:"time,omitempty"`
	Date         string `json:"date,omitempty"`
	TimeZoneName string `json:"timeZoneName,omitempty"`
	TimeZoneAuto string `json:"timeZoneAuto,omitempty"`
	DSTActive    string `json:"dstActive,omitempty"`
}

// RouterBoardInfo represents routerboard information
type RouterBoardInfo struct {
	RouterBoard     string `json:"routerboard,omitempty"`
	Model           string `json:"model,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	FirmwareType    string `json:"firmwareType,omitempty"`
	FactoryFirmware string `json:"factoryFirmware,omitempty"`
	CurrentFirmware string `json:"currentFirmware,omitempty"`
	UpgradeFirmware string `json:"upgradeFirmware,omitempty"`
}

// SystemResourceMonitorStats holds all fields from /system/resource/print interval=1s.
type SystemResourceMonitorStats struct {
	Uptime               string  `json:"uptime"`
	Version              string  `json:"version"`
	BuildTime            string  `json:"buildTime"`
	FreeMemory           int64   `json:"freeMemory"`
	TotalMemory          int64   `json:"totalMemory"`
	CPU                  string  `json:"cpu"`
	CPUCount             int     `json:"cpuCount"`
	CPUFrequency         int     `json:"cpuFrequency"`
	CPULoad              int     `json:"cpuLoad"`
	FreeHddSpace         int64   `json:"freeHddSpace"`
	TotalHddSpace        int64   `json:"totalHddSpace"`
	WriteSectSinceReboot int64   `json:"writeSectSinceReboot"`
	WriteSectTotal       int64   `json:"writeSectTotal"`
	BadBlocks            float64 `json:"badBlocks"`
	ArchitectureName     string  `json:"architectureName"`
	BoardName            string  `json:"boardName"`
	Platform             string  `json:"platform"`
}

// LogEntry represents a log entry
type LogEntry struct {
	ID      string `json:"id,omitempty"`
	Time    string `json:"time,omitempty"`
	Topics  string `json:"topics,omitempty"`
	Message string `json:"message,omitempty"`
}

// HotspotStats represents hotspot statistics
type HotspotStats struct {
	TotalUsers  int `json:"totalUsers"`
	ActiveUsers int `json:"activeUsers"`
}

// DashboardData represents complete dashboard data
type DashboardData struct {
	RouterID        uint             `json:"routerId"`
	RouterName      string           `json:"routerName,omitempty"`
	SystemTime      *SystemClock     `json:"systemTime,omitempty"`
	Resource        *SystemResource  `json:"resource,omitempty"`
	Health          *SystemHealth    `json:"health,omitempty"`
	Identity        *SystemIdentity  `json:"identity,omitempty"`
	RouterBoard     *RouterBoardInfo `json:"routerBoard,omitempty"`
	Stats           *HotspotStats    `json:"stats,omitempty"`
	Interfaces      []*Interface     `json:"interfaces,omitempty"`
	HotspotLogs     []*LogEntry      `json:"hotspotLogs,omitempty"`
	ConnectionError string           `json:"connectionError,omitempty"`
}

// TrafficRequest represents a traffic monitoring request
type TrafficRequest struct {
	Interface string `json:"interface" validate:"required"`
}

// ResourceUsage is kept for backward compatibility with old stub.
type ResourceUsage = SystemResource
