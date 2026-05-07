package model

import (
	"errors"
)

var (
	// ErrRefreshTokenNotReplaceable 旧 Refresh Token 不存在或已吊销，不能继续完成轮换
	ErrRefreshTokenNotReplaceable = errors.New("refresh token not replaceable")
	// ErrProjectGroupCannotDeleteLast 项目应至少存在一个图像组
	ErrProjectGroupCannotDeleteLast = errors.New("project must keep at least one group")
	// ErrUnsupportedDetectionTaskCode 不支持的原子检测能力代码
	ErrUnsupportedDetectionTaskCode = errors.New("unsupported detection task code")
)
