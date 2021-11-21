package service

import (
	"context"
	"gin-api-server/internal/authority"
	"gin-api-server/internal/extend/gormx"
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/internal/extend/timex"
	"gin-api-server/internal/model"
	"gin-api-server/pkg/auth"

	"gorm.io/gorm"
)

const (
	UserPasswordHashCost int = 10 // 用于用户密码加密
)

type UserLoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email,max=255"`
	Password string `json:"password" form:"password" binding:"required,max=255"`
}

type UserRegisterTokenRequest struct {
	Token string `json:"token" form:"token" binding:"required,max=255"`
}

type UserUpdatePasswordRequest struct {
	OldPassword        string `json:"old_password" form:"old_password" binding:"required,max=255"`
	NewPassword        string `json:"new_password" form:"new_password" binding:"required,min=6,max=255"`
	NewPasswordConfirm string `json:"new_password_confirm" form:"new_password_confirm" binding:"required,eqfield=NewPassword,min=6,max=255"`
}

type UserUpdateProfileRequest struct {
	Name string `json:"name" form:"name"  binding:"max=255"`
}

type UserRegisterEmailRequest struct {
	Email string `json:"email" form:"email" binding:"required,email,max=255"`
}

type UserRegisterProfileRequest struct {
	Token           string `json:"token" form:"token" binding:"required,max=255"`
	Email           string `json:"email" form:"email" binding:"required,email,max=255"`
	Name            string `json:"name" form:"name" binding:"required,max=255"`
	Password        string `json:"password" form:"password" binding:"required,min=6,max=255"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm" binding:"required,eqfield=Password,min=6,max=255"`
}

type UserSvc struct {
	Auther    auth.Auther
	DB        *gorm.DB
	Authority *authority.Authority
}

func (s *UserSvc) Verify(ctx context.Context, email, password string) (*model.UserDTO, error) {
	var user model.User
	exist, err := gormx.FindOne(s.DB.Where("email = ?", email), &user)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrUserWrongEmailOrPassword
	}

	if ok := auth.CheckPasswordHash(password, user.Password); !ok {
		return nil, svcerrorsx.ErrUserWrongEmailOrPassword
	} else {
		return user.ToDTO(), nil
	}
}

func (s *UserSvc) Get(ctx context.Context, id uint64) (*model.UserDTO, error) {
	var user model.User
	exist, err := gormx.FindOne(s.DB.Where("id = ?", id), &user)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrResourceNotFound
	}
	return user.ToDTO(), err
}

func (s *UserSvc) UpdatePassword(ctx context.Context, id uint64, req *UserUpdatePasswordRequest) error {
	var user model.User
	exist, err := gormx.FindOne(s.DB.Where("id = ?", id), &user)
	if err != nil {
		return err
	} else if !exist {
		return svcerrorsx.ErrResourceNotFound
	}

	if checked := auth.CheckPasswordHash(req.OldPassword, user.Password); !checked {
		return svcerrorsx.ErrUserWrongOldPassword
	}

	hashPassword, err := auth.HashPassword(req.NewPassword, UserPasswordHashCost)
	if err != nil {
		return err
	}

	return s.DB.Model(&model.User{}).Where("id = ?", id).Update("password", hashPassword).Error
}

func (s *UserSvc) UpdateProfile(ctx context.Context, id uint64, req *UserUpdateProfileRequest) (*model.UserDTO, error) {
	var user model.User
	exist, err := gormx.FindOne(s.DB.Where("id = ?", id), &user)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrResourceNotFound
	}

	user.Name = req.Name
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}

func (s *UserSvc) ListUserOrgs(ctx context.Context, userID uint64) (model.UserOrgDTOList, error) {
	var orgs model.Orgs
	if err := s.DB.Raw(
		"SELECT * FROM orgs WHERE id IN (SELECT org_id FROM user_roles WHERE user_id = ?) ORDER BY id",
		userID, userID).Scan(&orgs).Error; err != nil {
		return nil, err
	}
	// TOOO 增加返回 app count
	dtoOrgs := orgs.ToDTOs()
	for i := range dtoOrgs {
		isOwner, err := s.Authority.CheckOrgRole(userID, orgs[i].ID, model.RoleOrgOwner)
		if err != nil {
			return nil, err
		}
		dtoOrgs[i].IsOwner = isOwner
	}
	return dtoOrgs, nil
}

func (s *UserSvc) Create(ctx context.Context, req *UserRegisterProfileRequest) (*model.UserDTO, error) {
	if req.Password != req.PasswordConfirm {
		return nil, svcerrorsx.ErrUserConfirmPasswordInconsistency
	}

	var u model.User
	if exist, err := gormx.FindOne(s.DB.Where("email = ?", req.Email), &u); err != nil {
		return nil, err
	} else if exist {
		return nil, svcerrorsx.ErrUserEmailAlreadyExists
	}

	var token model.Token
	exist, err := gormx.FindOne(s.DB.Where("id = ?", req.Token), &token)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrTokenIDNotExist
	}
	if token.ExpiredAt < timex.GetNowUTCMilli() {
		return nil, svcerrorsx.ErrTokenIDExpired
	}

	user := &model.User{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		IsAdmin:         false,
		EmailVerifiedAt: timex.GetNowUTCMilli(),
	}
	password, err := auth.HashPassword(user.Password, UserPasswordHashCost)
	if err != nil {
		return nil, err
	}
	user.Password = password

	result := s.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return user.ToDTO(), nil
}

func (s *UserSvc) GetRegisterToken(ctx context.Context, id string) (*model.TokenDTO, error) {
	var token model.Token
	exist, err := gormx.FindOne(s.DB.Where("id = ?", id), &token)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrTokenIDNotExist
	}

	dtoToken := token.ToDTO()
	var user model.User
	exist, err = gormx.FindOne(s.DB.Where("email = ?", token.Email), &user)
	if err != nil {
		return nil, err
	} else if exist {
		dtoToken.EmailVerifiedAt = user.EmailVerifiedAt
	}

	return dtoToken, nil
}

func (s *UserSvc) CreateRegisterToken(ctx context.Context, email string) (*model.TokenDTO, error) {
	var user model.User
	if exist, err := gormx.FindOne(s.DB.Where("email = ?", email), &user); err != nil {
		return nil, err
	} else if exist {
		return nil, svcerrorsx.ErrUserEmailAlreadyExists
	}

	expiredAt := timex.GetNowUTCMilli() + 7*24*3600*1e3
	token := model.Token{
		Email:     email,
		Type:      model.TokenTypeRegister,
		ExpiredAt: expiredAt,
	}
	result := s.DB.Create(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return token.ToDTO(), nil
}
