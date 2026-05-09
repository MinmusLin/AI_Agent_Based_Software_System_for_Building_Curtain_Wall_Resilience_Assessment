package model

import (
	"database/sql"
	"time"

	"icw_common/gen/core/common"
)

// RefreshTokenRecord Refresh Token 记录
type RefreshTokenRecord struct {
	Id                uint64
	TokenId           string
	UserId            uint64
	TokenHash         string
	ExpiresAt         time.Time
	RevokedAt         sql.NullTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ReplacedByTokenId sql.NullString
}

// EmailSendLogRecord 邮件发送记录
type EmailSendLogRecord struct {
	Id            uint64
	ReceiverEmail string
	SenderEmail   string
	Scene         string
	EmailCode     string
	Status        commonpb.EmailSendStatus_Value
	ErrorMessage  sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
