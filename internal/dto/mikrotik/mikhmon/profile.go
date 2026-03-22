package mikhmon

import (
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

type ProfileConfig struct {
	Name          string `json:"name" binding:"required"`
	AddressPool   string `json:"address_pool,omitempty"`
	RateLimit     string `json:"rate_limit,omitempty"`
	SharedUsers   int    `json:"shared_users,omitempty"`
	ParentQueue   string `json:"parent_queue,omitempty"`
	Price         int64  `json:"price,omitempty"`
	SellingPrice  int64  `json:"selling_price,omitempty"`
	Validity      string `json:"validity,omitempty"`
	ExpireMode    string `json:"expire_mode,omitempty"`
	LockUser      bool   `json:"lock_user,omitempty"`
	LockServer    bool   `json:"lock_server,omitempty"`
	OnLoginScript string `json:"on_login_script,omitempty"`
}

type CreateProfileRequest struct {
	Name        string        `json:"name" binding:"required"`
	AddressPool string        `json:"address_pool,omitempty"`
	RateLimit   string        `json:"rate_limit,omitempty"`
	SharedUsers int           `json:"shared_users,omitempty"`
	ParentQueue string        `json:"parent_queue,omitempty"`
	Config      ProfileConfig `json:"config"`
}

type UpdateProfileRequest struct {
	Name         string `json:"name,omitempty"`
	AddressPool  string `json:"address_pool,omitempty"`
	RateLimit    string `json:"rate_limit,omitempty"`
	SharedUsers  int    `json:"shared_users,omitempty"`
	ParentQueue  string `json:"parent_queue,omitempty"`
	Price        int64  `json:"price,omitempty"`
	SellingPrice int64  `json:"selling_price,omitempty"`
	Validity     string `json:"validity,omitempty"`
	ExpireMode   string `json:"expire_mode,omitempty"`
	LockUser     bool   `json:"lock_user,omitempty"`
	LockServer   bool   `json:"lock_server,omitempty"`
}

type GenerateScriptRequest struct {
	Mode         string `json:"mode" binding:"required"`
	Price        int64  `json:"price"`
	Validity     string `json:"validity"`
	SellingPrice int64  `json:"selling_price"`
	NoExp        bool   `json:"no_exp"`
	LockUser     string `json:"lock_user"`
	LockServer   string `json:"lock_server"`
	ProfileName  string `json:"profile_name" binding:"required"`
}

type OnLoginScriptResponse struct {
	Script string `json:"script"`
}

func (r *CreateProfileRequest) ToDomain() *mikhmonDomain.ProfileRequest {
	return &mikhmonDomain.ProfileRequest{
		Name:        r.Name,
		AddressPool: r.AddressPool,
		RateLimit:   r.RateLimit,
		SharedUsers: r.SharedUsers,
		ParentQueue: r.ParentQueue,
		Config: mikhmonDomain.ProfileConfig{
			Name:          r.Config.Name,
			AddressPool:   r.Config.AddressPool,
			RateLimit:     r.Config.RateLimit,
			SharedUsers:   r.Config.SharedUsers,
			ParentQueue:   r.Config.ParentQueue,
			Price:         r.Config.Price,
			SellingPrice:  r.Config.SellingPrice,
			Validity:      r.Config.Validity,
			ExpireMode:    r.Config.ExpireMode,
			LockUser:      r.Config.LockUser,
			LockServer:    r.Config.LockServer,
			OnLoginScript: r.Config.OnLoginScript,
		},
	}
}

func (r *UpdateProfileRequest) ToDomain(profileID string) *mikhmonDomain.ProfileRequest {
	req := &mikhmonDomain.ProfileRequest{
		Name: profileID,
		Config: mikhmonDomain.ProfileConfig{
			Name: profileID,
		},
	}
	if r.Name != "" {
		req.Name = r.Name
		req.Config.Name = r.Name
	}
	if r.AddressPool != "" {
		req.AddressPool = r.AddressPool
		req.Config.AddressPool = r.AddressPool
	}
	if r.RateLimit != "" {
		req.RateLimit = r.RateLimit
		req.Config.RateLimit = r.RateLimit
	}
	if r.SharedUsers > 0 {
		req.SharedUsers = r.SharedUsers
		req.Config.SharedUsers = r.SharedUsers
	}
	if r.ParentQueue != "" {
		req.ParentQueue = r.ParentQueue
		req.Config.ParentQueue = r.ParentQueue
	}
	req.Config.Price = r.Price
	req.Config.SellingPrice = r.SellingPrice
	if r.Validity != "" {
		req.Config.Validity = r.Validity
	}
	if r.ExpireMode != "" {
		req.Config.ExpireMode = r.ExpireMode
	}
	req.Config.LockUser = r.LockUser
	req.Config.LockServer = r.LockServer
	return req
}

func (r *GenerateScriptRequest) ToDomain() *mikhmonDomain.OnLoginScriptData {
	return &mikhmonDomain.OnLoginScriptData{
		Mode:         r.Mode,
		Price:        r.Price,
		Validity:     r.Validity,
		SellingPrice: r.SellingPrice,
		NoExp:        r.NoExp,
		LockUser:     r.LockUser,
		LockServer:   r.LockServer,
		ProfileName:  r.ProfileName,
	}
}
