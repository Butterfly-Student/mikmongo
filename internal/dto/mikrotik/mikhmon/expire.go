package mikhmon

type ExpireMonitorStatusResponse struct {
	Enabled   bool   `json:"enabled"`
	LastRun   string `json:"last_run,omitempty"`
	NextRun   string `json:"next_run,omitempty"`
	UserCount int    `json:"user_count"`
}

type ScriptResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
