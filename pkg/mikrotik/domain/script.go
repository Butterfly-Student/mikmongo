package domain

// ProfileRequest represents a request for script generation with Mikhmon config.
// Name is the profile name; Validity is a RouterOS time string (e.g., "1d", "1h").
type ProfileRequest struct {
	Name        string `json:"name" validate:"required"`
	Validity    string `json:"validity,omitempty"`
	RateLimit   string `json:"rateLimit,omitempty"`
	SharedUsers int    `json:"sharedUsers,omitempty"`
	Price       int64  `json:"price,omitempty"`
}
