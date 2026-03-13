package model

// CasbinRule represents RBAC authorization rules for Casbin
type CasbinRule struct {
	ID    int64  `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"type:varchar(10);index"`
	V0    string `gorm:"type:varchar(256);index"`
	V1    string `gorm:"type:varchar(256);index"`
	V2    string `gorm:"type:varchar(256)"`
	V3    string `gorm:"type:varchar(256)"`
	V4    string `gorm:"type:varchar(256)"`
	V5    string `gorm:"type:varchar(256)"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}
