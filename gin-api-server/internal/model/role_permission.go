package model

type RolePermission struct {
	RoleID       uint64 `gorm:"primaryKey"` // 角色ID
	PermissionID uint64 `gorm:"primaryKey"` // 权限ID
}
