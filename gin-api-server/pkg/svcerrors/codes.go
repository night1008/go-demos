package svcerrors

import (
	"net/http"
)

// 公共业务错误码
var (
	ErrServerError          = NewError("InternalError", http.StatusInternalServerError, "服务器发生错误")
	ErrServiceUnavailable   = NewError("ServiceUnavailable", http.StatusServiceUnavailable, "系统繁忙，请稍后重试")
	ErrAuthLoginFailed      = NewError("Auth.LoginFailed", http.StatusInternalServerError, "用户登录失败")
	ErrAuthLogoutFailed     = NewError("Auth.LogoutFailed", http.StatusInternalServerError, "用户登出失败")
	ErrResourceCreateFailed = NewError("Resource.CreateFailed", http.StatusInternalServerError, "资源创建失败")

	ErrInvalidParameter                 = NewError("InvalidParameter", http.StatusBadRequest, "输入参数错误")
	ErrResourceNotFound                 = NewError("Resource.NotFound", http.StatusNotFound, "资源不存在")
	ErrResourceUnavailable              = NewError("Resource.Unavailable", http.StatusForbidden, "资源不可用")
	ErrMethodNotAllow                   = NewError("MethodNotAllow", http.StatusMethodNotAllowed, "方法不被允许")
	ErrTooManyRequests                  = NewError("TooManyRequests", http.StatusTooManyRequests, "请求过于频繁")
	ErrUserWrongEmailOrPassword         = NewError("User.WrongEmailOrPassword", http.StatusBadRequest, "邮箱或密码不正确")
	ErrUserUnauthorized                 = NewError("User.Unauthorized", http.StatusUnauthorized, "用户未登录")
	ErrUserWrongOldPassword             = NewError("User.WrongOldPassword", http.StatusBadRequest, "用户旧密码错误")
	ErrUserWrongPassword                = NewError("User.WrongPassword", http.StatusUnauthorized, "用户密码错误")
	ErrUserAlreadyExists                = NewError("User.AlreadyExists", http.StatusBadRequest, "用户已存在")
	ErrUserEmailAlreadyExists           = NewError("User.EmailAlreadyExists", http.StatusBadRequest, "邮箱已注册")
	ErrUserNotFound                     = NewError("User.NotFound", http.StatusNotFound, "用户不存在")
	ErrUserInvalidPassword              = NewError("User.InvalidPassword", http.StatusBadRequest, "用户密码不符合规范")
	ErrUserConfirmPasswordInconsistency = NewError("User.ConfirmPasswordInconsistency", http.StatusBadRequest, "用户确认密码不一致")
	ErrUserAlreadyDisabled              = NewError("User.AlreadyDisabled", http.StatusForbidden, "用户已被禁用")
	ErrUserInvalidPermisson             = NewError("User.InvalidPermisson", http.StatusForbidden, "用户权限验证失败")
	ErrUserInvalidToken                 = NewError("User.InvalidToken", http.StatusForbidden, "用户 token 已失效")
	ErrUserTokenExpired                 = NewError("User.TokenExpired", http.StatusForbidden, "用户 token 已过期")
)
