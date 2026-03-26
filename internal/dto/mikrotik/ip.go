package mikrotik

// AddIPPoolRequest is the request body for creating an IP pool.
type AddIPPoolRequest struct {
	Name    string `json:"name" binding:"required"`
	Ranges  string `json:"ranges" binding:"required"`
	Comment string `json:"comment,omitempty"`
}
