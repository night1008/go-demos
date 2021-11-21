package service

import (
	"context"
	"gin-api-server/internal/extend/gormx"
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/internal/model"

	"gorm.io/gorm"
)

type RoleListRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type RoleCreateRequest struct {
	Name string `json:"name" binding:"required,max=255"`
}

type RoleUpdateRequest struct {
	ID   uint64
	Name string `json:"name" binding:"required,max=255"`
}

type RoleSvc struct {
	DB *gorm.DB
}

func (s *RoleSvc) ListAll(ctx context.Context) (model.RoleDTOList, error) {
	var roles model.Roles
	if err := s.DB.Model(&model.Role{}).Order("id").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles.ToDTOs(), nil
}

func (s *RoleSvc) Create(ctx context.Context, userID uint64, req *RoleCreateRequest) (*model.RoleDTO, error) {
	role := &model.Role{
		Name: req.Name,
	}
	result := s.DB.Create(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, svcerrorsx.ErrResourceCreateFailed
	}
	return role.ToDTO(), nil
}

func (s *RoleSvc) Update(ctx context.Context, userID uint64, req *RoleUpdateRequest) (*model.RoleDTO, error) {
	var role model.Role
	exist, err := gormx.FindOne(s.DB.Where("id = ?", req.ID), &role)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrResourceCreateFailed
	}

	role.Name = req.Name

	if err := s.DB.Save(&role).Error; err != nil {
		return nil, err
	}
	return role.ToDTO(), nil
}

func (s *RoleSvc) Delete(ctx context.Context, userID, id uint64) error {
	return s.DB.Delete(&model.Role{}, id).Error
}
