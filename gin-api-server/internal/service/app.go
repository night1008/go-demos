package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"gin-api-server/internal/authority"
	"gin-api-server/internal/extend/gormx"
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/internal/model"
	"hash/fnv"
	"time"

	"gorm.io/gorm"
)

const MaxSingleOrgAppCreateCount = 100

type AppCreateRequest struct {
	Name       string `json:"name" binding:"required,max=255"`
	Identifier string `json:"identifier" binding:"required,max=255"`
	OrgID      uint64 `json:"org_id" binding:"required"`
}

type AppUpdateRequest struct {
	ID   uint64
	Name string `json:"name" binding:"required,max=255"`
}

type AppSvc struct {
	DB        *gorm.DB
	Authority *authority.Authority
}

func (s *AppSvc) Get(ctx context.Context, userID, appID uint64) (*model.AppDTO, error) {
	var app model.App
	exist, err := gormx.FindOne(s.DB.Where("id = ?", appID), &app)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrResourceNotFound
	}

	if ok, err := s.Authority.CheckPermission(userID, app.OrgID, app.ID, model.PermissionGetApp); err != nil {
		return nil, err
	} else if !ok {
		return nil, svcerrorsx.ErrUserInvalidPermisson
	}

	return app.ToDTO(), err
}

func (s *AppSvc) Create(ctx context.Context, userID uint64, req *AppCreateRequest) (*model.AppDTO, error) {
	if ok, err := s.Authority.CheckPermission(userID, req.OrgID, 0, model.PermissionCreateApp); err != nil {
		return nil, err
	} else if !ok {
		return nil, svcerrorsx.ErrUserInvalidPermisson
	}

	var count int64
	if err := s.DB.Model(&model.App{}).
		Where("org_id = ?", req.OrgID).
		Count(&count).Error; err != nil {
		return nil, err
	} else if count >= MaxSingleOrgAppCreateCount {
		return nil, svcerrorsx.ErrOrgAppOverCountLimit
	}

	// TODO 创建成功后同步到 consul
	accessKeyID := createAccessKeyID(req.Identifier)
	accessKeySecret := createAccessKeySecret(req.Identifier)
	app := &model.App{
		Name:            req.Name,
		Identifier:      req.Identifier,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		DeployStatus:    model.AppDeployStatusWaiting,
		OrgID:           req.OrgID,
		CreatorID:       userID,
	}

	result := s.DB.Create(&app)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, svcerrorsx.ErrResourceCreateFailed
	}
	return app.ToDTO(), nil
}

func (s *AppSvc) Update(ctx context.Context, userID uint64, req *AppUpdateRequest) (*model.AppDTO, error) {
	var app model.App
	exist, err := gormx.FindOne(s.DB.Where("id = ?", req.ID), &app)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, svcerrorsx.ErrResourceNotFound
	}

	if ok, err := s.Authority.CheckPermission(userID, app.OrgID, app.ID, model.PermissionUpdateApp); err != nil {
		return nil, err
	} else if !ok {
		return nil, svcerrorsx.ErrUserInvalidPermisson
	}

	app.Name = req.Name
	if err := s.DB.Save(&app).Error; err != nil {
		return nil, err
	}
	return app.ToDTO(), nil
}

func (s *AppSvc) Delete(ctx context.Context, userID, appID uint64) error {
	var app model.App
	exist, err := gormx.FindOne(s.DB.Where("id = ?", appID), &app)
	if err != nil {
		return err
	} else if !exist {
		return svcerrorsx.ErrResourceNotFound
	}

	if ok, err := s.Authority.CheckPermission(userID, app.OrgID, app.ID, model.PermissionDeleteApp); err != nil {
		return err
	} else if !ok {
		return svcerrorsx.ErrUserInvalidPermisson
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.App{}).
			Where("id = ?", appID).
			Update("deletor_id", userID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.App{}, appID).Error; err != nil {
			return err
		}
		return nil
	})
}

func createAccessKeyID(identifier string) string {
	h := fnv.New64()
	h.Write([]byte(fmt.Sprintf("%s-%d", identifier, time.Now().UnixNano())))
	id := hex.EncodeToString(h.Sum(nil))
	return id
}

func createAccessKeySecret(identifier string) string {
	h := fnv.New128()
	h.Write([]byte(fmt.Sprintf("%s-%d", identifier, time.Now().UnixNano())))
	secret := hex.EncodeToString(h.Sum(nil))
	return secret
}
