package svcerrorsx

import (
	"gin-api-server/pkg/svcerrors"
	"net/http"
)

var (
	ErrServerError          = svcerrors.ErrServerError
	ErrServiceUnavailable   = svcerrors.ErrServiceUnavailable
	ErrAuthLoginFailed      = svcerrors.ErrAuthLoginFailed
	ErrAuthLogoutFailed     = svcerrors.ErrAuthLogoutFailed
	ErrResourceCreateFailed = svcerrors.ErrResourceCreateFailed

	ErrInvalidParameter                 = svcerrors.ErrInvalidParameter
	ErrResourceNotFound                 = svcerrors.ErrResourceNotFound
	ErrResourceUnavailable              = svcerrors.ErrResourceUnavailable
	ErrMethodNotAllow                   = svcerrors.ErrMethodNotAllow
	ErrTooManyRequests                  = svcerrors.ErrMethodNotAllow
	ErrUserWrongEmailOrPassword         = svcerrors.ErrUserWrongEmailOrPassword
	ErrUserUnauthorized                 = svcerrors.ErrUserUnauthorized
	ErrUserWrongOldPassword             = svcerrors.ErrUserWrongOldPassword
	ErrUserWrongPassword                = svcerrors.ErrUserWrongPassword
	ErrUserAlreadyExists                = svcerrors.ErrUserAlreadyExists
	ErrUserEmailAlreadyExists           = svcerrors.ErrUserEmailAlreadyExists
	ErrUserNotFound                     = svcerrors.ErrUserNotFound
	ErrUserInvalidPassword              = svcerrors.ErrUserInvalidPassword
	ErrUserConfirmPasswordInconsistency = svcerrors.ErrUserConfirmPasswordInconsistency
	ErrUserAlreadyDisabled              = svcerrors.ErrUserAlreadyDisabled
	ErrUserInvalidPermisson             = svcerrors.ErrUserInvalidPermisson
	ErrUserInvalidToken                 = svcerrors.ErrUserInvalidPermisson
	ErrUserTokenExpired                 = svcerrors.ErrUserInvalidPermisson

	ErrOrgNotOwner            = svcerrors.NewError("Org.NotOwner", http.StatusForbidden, "当前用户不是组织拥有者")
	ErrOrgNotAppAdmin         = svcerrors.NewError("Org.NotAppAdmin", http.StatusForbidden, "当前用户不是应用管理者")
	ErrOrgNotMember           = svcerrors.NewError("Org.NotMember", http.StatusForbidden, "当前用户不是组织成员")
	ErrOrgNameAlreadyExisted  = svcerrors.NewError("Org.NameAlreadyExisted", http.StatusBadRequest, "组织名称已存在")
	ErrOrgAppCountCheckFailed = svcerrors.NewError("Org.AppCountCheckFailed", http.StatusBadRequest, "组织应用检查失败")
	ErrOrgOverCountLimit      = svcerrors.NewError("Org.OverCountLimit", http.StatusBadRequest, "组织超过数量限制")
	ErrOrgHasAssociatedApp    = svcerrors.NewError("Org.HasAssociatedApp", http.StatusBadRequest, "该组织下存在应用")
	ErrOrgAppOverCountLimit   = svcerrors.NewError("Org.AppOverCountLimit", http.StatusBadRequest, "该组织下应用超过数量限制")
	ErrTokenIDAlreadyExisted  = svcerrors.NewError("Token.IDAlreadyExisted", http.StatusBadRequest, "Token ID 已存在")
	ErrTokenIDNotExist        = svcerrors.NewError("Token.IDNotExist", http.StatusBadRequest, "Token ID 不存在")
	ErrTokenIDExpired         = svcerrors.NewError("Token.IDExpired", http.StatusBadRequest, "Token ID 已过期")
)
