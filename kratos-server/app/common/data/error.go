package data

import "github.com/go-kratos/kratos/v2/errors"

var (
	SqlError = func(msg string) *errors.Error {
		return errors.New(500, "查询异常", msg)
	}

	ServerError = func(msg string) *errors.Error {
		return errors.New(500, "服务异常", msg)
	}

	PlanExceeds = func(msg string) *errors.Error {
		return errors.New(400, "计划数超额", msg)
	}

	ParamError = func(msg string) *errors.Error {
		return errors.New(400, "参数错误", msg)
	}

	RecordNotFoundError = func(msg string) *errors.Error {
		return errors.New(404, "记录不存在", msg)
	}

	UserNotLoginError = func(msg string) *errors.Error {
		return errors.New(401, "用户未登录", msg)
	}

	ResourceWithoutPermissionError = func(msg string) *errors.Error {
		return errors.New(403, "资源无权限", msg)
	}
)
