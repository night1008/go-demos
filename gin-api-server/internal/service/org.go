package service

import (
	"context"
	"errors"
	"gin-api-server/internal/authority"
	"gin-api-server/internal/extend/gormx"
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/internal/extend/timex"
	"gin-api-server/internal/model"
	"gin-api-server/pkg/auth"
	"strings"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

const MaxOrgCreateCount = 100

type OrgCreateRequest struct {
	Name string `json:"name" binding:"required,max=255"`
}

type OrgUpdateRequest struct {
	ID   uint64
	Name string `json:"name" binding:"required,max=255"`
}

type OrgListAppUserRequest struct {
	OrgID uint64
	AppID uint64 `form:"app_id"`
}

type OrgQueryUserRequest struct {
	OrgID uint64
	Email string `json:"email" form:"email" binding:"required,email,max=255"`
}

type OrgUserRoleUpdateRequest struct {
	OrgID   uint64
	Email   string `json:"email" binding:"required,email,max=255"`
	Updates []struct {
		AppID  uint64 `json:"app_id" binding:"required"`
		RoleID uint64 `json:"role_id" binding:"required"`
	} `json:"updates"`
	Deletes []uint64 `json:"deletes"`
}

type OrgRestoreDeletedAppsRequest struct {
	AppIDs []uint64 `json:"app_ids" binding:"required"`
}

type OrgSvc struct {
	DB        *gorm.DB
	Authority *authority.Authority
}

func (s *OrgSvc) Create(ctx context.Context, userID uint64, req *OrgCreateRequest) (*model.OrgDTO, error) {
	var count int64
	if err := s.DB.Model(&model.Org{}).
		Where("creator_id = ?", userID).
		Count(&count).Error; err != nil {
		return nil, err
	} else if count >= MaxOrgCreateCount {
		return nil, svcerrorsx.ErrOrgOverCountLimit
	}

	var org *model.Org
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// TODO 发送邮件
		org = &model.Org{
			Name:      req.Name,
			CreatorID: userID,
		}
		result := s.DB.Create(&org)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return svcerrorsx.ErrResourceCreateFailed
		}
		return s.Authority.AssignRole(userID, org.ID, 0, model.RoleOrgOwner)
	})
	if err != nil {
		return nil, err
	}

	isOwner, err := s.Authority.CheckOrgRole(userID, org.ID, model.RoleOrgOwner)
	if err != nil {
		return nil, err
	}

	dtoOrg := org.ToDTO()
	dtoOrg.IsOwner = isOwner
	return dtoOrg, nil
}

func (s *OrgSvc) Update(ctx context.Context, userID uint64, req *OrgUpdateRequest) (*model.OrgDTO, error) {
	var org model.Org
	exist, err := gormx.FindOne(s.DB.Where("id = ?", req.ID), &org)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrResourceNotFound
	}

	if ok, err := s.Authority.CheckPermission(userID, org.ID, 0, model.PermissionUpdateOrg); err != nil {
		return nil, err
	} else if !ok {
		return nil, svcerrorsx.ErrUserInvalidPermisson
	}

	org.Name = req.Name

	if err := s.DB.Save(&org).Error; err != nil {
		return nil, err
	}

	isOwner, err := s.Authority.CheckOrgRole(userID, org.ID, model.RoleOrgOwner)
	if err != nil {
		return nil, err
	}

	dtoOrg := org.ToDTO()
	dtoOrg.IsOwner = isOwner
	return dtoOrg, nil
}

func (s *OrgSvc) Delete(ctx context.Context, userID, id uint64) error {
	var org model.Org
	exist, err := gormx.FindOne(s.DB.Where("id = ?", id), &org)
	if err != nil {
		return err
	} else if !exist {
		return svcerrorsx.ErrResourceNotFound
	}

	if ok, err := s.Authority.CheckPermission(userID, org.ID, 0, model.PermissionDeleteOrg); err != nil {
		return err
	} else if !ok {
		return svcerrorsx.ErrUserInvalidPermisson
	}

	var count int64
	if err := s.DB.Model(&model.App{}).
		Where("org_id = ?", org.ID).
		Count(&count).Error; err != nil {
		return err
	} else if count > 0 {
		return svcerrorsx.ErrOrgHasAssociatedApp
	}

	return s.DB.Delete(&model.Org{}, id).Error
}

func (s *OrgSvc) ListOwners(ctx context.Context, userID, orgID uint64) (model.UserDTOList, error) {
	ok, err := s.Authority.CheckAnyOrgRole(userID, orgID)
	if err != nil {
		return nil, err
	} else if !ok {
		var list model.UserDTOList
		return list, nil
	}

	var users model.Users
	if err := s.DB.Raw(
		`SELECT users.* FROM users
	    INNER JOIN (
				SELECT user_id FROM user_roles
		 		 WHERE org_id = ? AND role_id IN (
					SELECT id FROM roles WHERE name = ?
		 		)
			) AS ur
		 	  ON users.id = ur.user_id
		 ORDER BY last_active_at DESC`,
		orgID, model.RoleOrgOwner).Scan(&users).Error; err != nil {
		return nil, err
	}

	return users.ToDTOs(), nil
}

func listApps(db *gorm.DB, orgID, appID uint64) (model.AppDTOList, error) {
	var apps model.Apps
	q := db.Model(&model.App{}).Where("org_id = ?", orgID)
	if appID > 0 {
		q = q.Where("id = ?", appID)
	}
	if err := q.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps.ToDTOs(), nil
}

func (s *OrgSvc) ListApps(ctx context.Context, userID, orgID uint64) (model.AppDTOList, error) {
	ok, err := s.Authority.CheckAnyOrgRole(userID, orgID)
	if err != nil {
		return nil, err
	} else if !ok {
		var list model.AppDTOList
		return list, nil
	}

	apps, err := listApps(s.DB, orgID, 0)
	appUserRoles, err := listAppUserRoles(s.DB, orgID, 0, userID)
	appUserRolesMap := make(map[uint64]uint64, len(appUserRoles))
	for _, role := range appUserRoles {
		appUserRolesMap[role.AppID] = role.RoleID
	}
	userApps := make([]*model.AppDTO, 0)
	for _, app := range apps {
		if roleID, ok := appUserRolesMap[app.ID]; ok {
			app.RoleID = roleID
			userApps = append(userApps, app)
		}
	}
	return userApps, nil
}

func (s *OrgSvc) ListAppMembers(ctx context.Context, userID uint64, req *OrgListAppUserRequest) (model.UserWithAppDTOList, error) {
	ok, err := s.Authority.CheckAnyOrgRole(userID, req.OrgID)
	if err != nil {
		return nil, err
	} else if !ok {
		var list model.UserWithAppDTOList
		return list, nil
	}

	var users model.Users
	subQuery := s.DB.Model(&model.UserRole{}).
		Select("user_id").
		Distinct("user_id").
		Where("org_id = ?", req.OrgID).
		Where("app_id != 0")
	if req.AppID != 0 {
		subQuery = subQuery.Where("app_id = ?", req.AppID)
	}
	if err := s.DB.Model(&model.User{}).
		Joins("INNER JOIN (?) AS user_roles ON users.id = user_roles.user_id", subQuery).Scan(&users).Error; err != nil {
		return nil, err
	}

	var apps1 model.Apps
	q := s.DB.Model(&model.App{}).Where("org_id = ?", req.OrgID)
	if req.AppID > 0 {
		q = q.Where("id = ?", req.AppID)
	}
	if err := q.Find(&apps1).Error; err != nil {
		return nil, err
	}

	apps := apps1.ToDTOs()
	var appUsers []*model.UserRole
	q = s.DB.Model(&model.UserRole{}).Where("org_id = ?", req.OrgID).Where("app_id != 0")
	if req.AppID > 0 {
		q = q.Where("app_id = ?", req.AppID)
	}
	if err := q.Find(&appUsers).Error; err != nil {
		return nil, err
	}
	appUserRoles := appUsers

	userAppRoles := make(map[uint64]map[uint64]uint64)
	for _, appUserRole := range appUserRoles {
		if userAppRole, ok := userAppRoles[appUserRole.UserID]; ok {
			userAppRole[appUserRole.AppID] = appUserRole.RoleID
		} else {
			userAppRoles[appUserRole.UserID] = make(map[uint64]uint64)
			userAppRoles[appUserRole.UserID][appUserRole.AppID] = appUserRole.RoleID
		}
	}

	usersWithApp := make([]*model.UserWithAppDTO, len(users))
	for i, user := range users {
		userApps := make([]*model.AppDTO, 0)
		if userAppRole, ok := userAppRoles[user.ID]; ok {
			for _, app := range apps {
				if roleID, ok := userAppRole[app.ID]; ok {
					var userApp model.AppDTO
					copier.Copy(&userApp, &app)
					userApp.RoleID = roleID
					userApps = append(userApps, &userApp)
				}
			}
		}

		usersWithApp[i] = &model.UserWithAppDTO{
			UserDTO: *user.ToDTO(),
			Apps:    userApps,
		}
	}
	return usersWithApp, nil
}

func listAppUserRoles(db *gorm.DB, orgID, appID, userID uint64) ([]*model.UserRole, error) {
	var userRoles []*model.UserRole
	q := db.Model(&model.UserRole{}).Where("org_id = ?", orgID)
	if appID > 0 {
		q = q.Where("app_id = ?", appID)
	}
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}

	if err := q.Find(&userRoles).Error; err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (s *OrgSvc) QueryAppUser(ctx context.Context, req *OrgQueryUserRequest) (*model.UserWithAppDTO, error) {
	var usersWithApp model.UserWithAppDTO
	var user model.User
	exist, err := gormx.FindOne(s.DB.Where("email = ?", req.Email), &user)
	if err != nil {
		return nil, err
	} else if !exist {
		usersWithApp = model.UserWithAppDTO{
			UserDTO: model.UserDTO{
				ID:              0,
				Name:            "",
				Logo:            "",
				Email:           req.Email,
				EmailVerifiedAt: 0,
			},
			IsOwner: false,
			Apps:    make([]*model.AppDTO, 0),
		}
		return &usersWithApp, nil
	}

	isOwner, err := s.Authority.CheckOrgRole(user.ID, req.OrgID, model.RoleOrgOwner)
	if err != nil {
		return nil, err
	}
	apps, err := listApps(s.DB, req.OrgID, 0)
	if err != nil {
		return nil, err
	}
	appUserRoles, err := listAppUserRoles(s.DB, req.OrgID, 0, user.ID)
	if err != nil {
		return nil, err
	}
	userAppRole := make(map[uint64]uint64)
	for _, appUserRole := range appUserRoles {
		userAppRole[appUserRole.AppID] = appUserRole.RoleID
	}

	userApps := make([]*model.AppDTO, 0)
	for _, app := range apps {
		if roleID, ok := userAppRole[app.ID]; ok {
			app.RoleID = roleID
			userApps = append(userApps, app)
		}
	}

	usersWithApp = model.UserWithAppDTO{
		UserDTO: *user.ToDTO(),
		IsOwner: isOwner,
		Apps:    userApps,
	}
	return &usersWithApp, nil
}

func createUserByEmail(db *gorm.DB, email string) (*model.UserDTO, bool, error) {
	emailParts := strings.SplitN(email, "@", 2)
	user := &model.User{
		Name:    emailParts[0],
		Email:   email,
		IsAdmin: false,
	}

	password, err := auth.HashPassword(user.Password, UserPasswordHashCost)
	if err != nil {
		return nil, false, err
	}
	user.Password = password

	result := db.Create(&user)
	if result.Error != nil {
		return nil, false, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, false, nil
	}
	return user.ToDTO(), true, nil
}

func createRegisterToken(db *gorm.DB, email string) (*model.TokenDTO, bool, error) {
	expiredAt := timex.GetNowUTCMilli() + 7*24*3600*1e3
	token := model.Token{
		Email:     email,
		Type:      model.TokenTypeRegister,
		ExpiredAt: expiredAt,
	}
	result := db.Create(&token)
	if result.Error != nil {
		return nil, false, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	return token.ToDTO(), true, nil
}

func (s *OrgSvc) UpdateUserRoles(ctx context.Context, userID uint64, req *OrgUserRoleUpdateRequest) (*model.TokenDTO, error) {
	if len(req.Updates) == 0 && len(req.Deletes) == 0 {
		return nil, nil
	}

	var dtoUser *model.UserDTO
	var token *model.TokenDTO
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var user model.User
		exist, err := gormx.FindOne(tx.Where("email = ?", req.Email), &user)
		if err != nil {
			return err
		} else if !exist {
			dtoUserTmp, ok, err := createUserByEmail(tx, req.Email)
			if err != nil {
				return err
			} else if !ok {
				return svcerrorsx.ErrResourceCreateFailed
			}
			dtoUser = dtoUserTmp
			tokenTmp, ok, err := createRegisterToken(tx, req.Email)
			if err != nil {
				return err
			} else if !ok {
				return svcerrorsx.ErrResourceCreateFailed
			}
			token = tokenTmp
			// TODO 发送邮件
		} else {
			dtoUser = user.ToDTO()
		}

		// MEMO 不能编辑自己的权限
		if userID == dtoUser.ID {
			return errors.New("can not update yourself roles")
		}

		for _, appID := range req.Deletes {
			if ok, err := s.Authority.CheckPermission(userID, req.OrgID, appID, model.PermissionUpdateAppMember); err != nil {
				return err
			} else if !ok {
				return svcerrorsx.ErrUserInvalidPermisson
			}
		}
		if err := tx.Where("org_id = ?", req.OrgID).
			Where("app_id IN ?", req.Deletes).
			Where("user_id = ?", dtoUser.ID).
			Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		for _, appRole := range req.Updates {
			if ok, err := s.Authority.CheckPermission(userID, req.OrgID, appRole.AppID, model.PermissionUpdateAppMember); err != nil {
				return err
			} else if !ok {
				return svcerrorsx.ErrUserInvalidPermisson
			}
			// todo 检查赋予权限
			userRole := &model.UserRole{
				OrgID:  req.OrgID,
				AppID:  appRole.AppID,
				UserID: dtoUser.ID,
			}
			if err := tx.FirstOrCreate(&userRole, &userRole).Error; err != nil {
				return err
			}
			if err := tx.Model(&userRole).Updates(model.UserRole{RoleID: appRole.RoleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func (s *OrgSvc) ListDeletedApps(ctx context.Context, userID, orgID uint64) (model.DeletedAppDTOList, error) {
	if ok, err := s.Authority.CheckPermission(userID, orgID, 0, model.PermissionListOrgDeletedApp); err != nil {
		return nil, err
	} else if !ok {
		return nil, svcerrorsx.ErrUserInvalidPermisson
	}

	var apps model.Apps
	if err := s.DB.Unscoped().Model(&model.App{}).
		Where("org_id = ?", orgID).
		Where("deleted_at > 0").
		Order("deleted_at DESC").
		Find(&apps).Error; err != nil {
		return nil, err
	}
	deletorIDs := make([]uint64, len(apps))
	for i, app := range apps {
		deletorIDs[i] = app.DeletorID
	}
	var users model.Users
	if err := s.DB.Model(&model.User{}).
		Where("id IN ?", deletorIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	usersMap := make(map[uint64]*model.User)
	for _, user := range users {
		usersMap[user.ID] = user
	}

	dtoApps := apps.ToDeletedAppDTOs()
	for _, app := range dtoApps {
		if user, ok := usersMap[app.DeletorID]; ok {
			app.DeletedBy = user.ToDTO()
		}
	}
	return dtoApps, nil
}

func (s *OrgSvc) RestoreDeletedApps(ctx context.Context, userID, orgID uint64, req *OrgRestoreDeletedAppsRequest) error {
	if ok, err := s.Authority.CheckPermission(userID, orgID, 0, model.PermissionRestoreOrgDeletedApp); err != nil {
		return err
	} else if !ok {
		return svcerrorsx.ErrUserInvalidPermisson
	}

	return s.DB.Unscoped().Model(&model.App{}).
		Where("org_id = ?", orgID).
		Where("id IN ?", req.AppIDs).
		Updates(map[string]interface{}{"deleted_at": 0, "deletor_id": 0}).Error
}
