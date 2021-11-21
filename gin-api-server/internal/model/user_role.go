package model

type UserRole struct {
	OrgID     uint64 `gorm:"primaryKey"`       // 组织ID
	AppID     uint64 `gorm:"primaryKey"`       // 应用ID
	UserID    uint64 `gorm:"primaryKey;index"` // 用户ID
	RoleID    uint64 // 角色ID
	CreatedAt int64  `gorm:"autoCreateTime:milli"` // 记录创建时间，单位 milli
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"` // 记录更新时间，单位 milli
}
