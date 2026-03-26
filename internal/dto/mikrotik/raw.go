package mikrotik

// RawCommandRequest is the request body for executing a raw RouterOS command.
type RawCommandRequest struct {
	Args []string `json:"args" binding:"required,min=1,max=20"`
}

// RawListenRequest is the WebSocket initial message for starting a raw listen.
type RawListenRequest struct {
	Args []string `json:"args" binding:"required,min=1,max=20"`
}
