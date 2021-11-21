package authority

import (
	"errors"
	"gin-api-server/internal/extend/gormx"
	"gin-api-server/internal/model"

	"gorm.io/gorm"
)

type Authority struct {
	DB *gorm.DB
}

var (
	ErrPermissionInUse     = errors.New("cannot delete assigned permission")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrRoleAlreadyAssigned = errors.New("this role is already assigned to the user")
	ErrRoleInUse           = errors.New("cannot delete assigned role")
	ErrRoleNotFound        = errors.New("role not found")
)

func (a *Authority) CreateRole(name string) (*model.Role, error) {
	role := &model.Role{
		Name: name,
	}
	if err := a.DB.FirstOrCreate(&role, &role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (a *Authority) CreatePermission(name string) (*model.Permission, error) {
	permission := &model.Permission{
		Name: name,
	}
	if err := a.DB.FirstOrCreate(&permission, &permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (a *Authority) AssignPermissions(roleName string, permissionNames []string) error {
	var role model.Role
	if err := a.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	var permissions []model.Permission
	if err := a.DB.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
		return err
	}
	if len(permissionNames) != len(permissions) {
		return ErrPermissionNotFound
	}

	rolePermissions := make([]model.RolePermission, len(permissions))
	for i, p := range permissions {
		rolePermissions[i] = model.RolePermission{
			RoleID:       role.ID,
			PermissionID: p.ID,
		}
	}
	return a.DB.Transaction(func(tx *gorm.DB) error {
		if err := a.DB.Delete(&model.RolePermission{}, "role_id = ?", role.ID).Error; err != nil {
			return err
		}
		if err := a.DB.CreateInBatches(rolePermissions, 100).Error; err != nil {
			return err
		}
		return nil
	})
}

func (a *Authority) CheckAnyOrgRole(userID, orgID uint64) (bool, error) {
	var count int64
	if err := a.DB.Raw(
		`SELECT COUNT(1) FROM user_roles
		  WHERE user_id = ?
		    AND org_id = ?`,
		userID, orgID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a *Authority) CheckOrgRole(userID, orgID uint64, roleName string) (bool, error) {
	var count int64
	if err := a.DB.Raw(
		`SELECT COUNT(1) FROM user_roles
		  WHERE user_id = ?
		   AND org_id = ?
			 AND role_id IN (
				SELECT id FROM roles WHERE name = ?
		)`,
		userID, orgID, roleName).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a *Authority) AssignRole(userID, orgID, appID uint64, roleName string) error {
	var role model.Role
	exist, err := gormx.FindOne(a.DB.Where("name = ?", roleName), &role)
	if err != nil {
		return err
	} else if !exist {
		return ErrRoleNotFound
	}
	userRole := model.UserRole{
		OrgID:  orgID,
		AppID:  appID,
		UserID: userID,
		RoleID: role.ID,
	}

	if err := a.DB.FirstOrCreate(&userRole, &userRole).Error; err != nil {
		return err
	}
	return nil
}

func (s *Authority) CheckPermission(userID, orgID, appID uint64, permissioName string) (bool, error) {
	var count int64
	if err := s.DB.Raw(`
		SELECT count(1) FROM (
			SELECT permission_id FROM role_permissions
			 WHERE role_id IN (
				SELECT role_id FROM user_roles
				 WHERE org_id = ? AND app_id IN (?) AND user_id = ?
			)
		) AS rp
		INNER JOIN (SELECT id FROM permissions WHERE name = ?) AS p
		ON rp.permission_id = p.id
	`, orgID, []uint64{0, appID}, userID, permissioName).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a *Authority) GetRoles() (model.Roles, error) {
	var roles model.Roles
	if err := a.DB.Model(&model.Role{}).Order("id").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (a *Authority) GetPermissions() (model.Permissions, error) {
	var permissions model.Permissions
	if err := a.DB.Model(&model.Permission{}).Order("id").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
