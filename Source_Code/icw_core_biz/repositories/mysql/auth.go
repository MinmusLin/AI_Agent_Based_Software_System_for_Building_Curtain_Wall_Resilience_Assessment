package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"icw_core_biz/internal/services/auth/consts"
)

// CreateLoginSession 登录时保存登录态 Refresh Token，并更新用户最近登录时间
func (r *Repository) CreateLoginSession(ctx context.Context, tokenId string, userId uint64, tokenHash string, expiresAt time.Time) error {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO refresh_tokens(token_id, user_id, token_hash, expires_at)
		VALUES (?, ?, ?, ?)
	`, tokenId, userId, tokenHash, expiresAt); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET last_login_at = NOW(3)
		WHERE id = ?
	`, userId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

// RotateRefreshToken 刷新时签发新的 Access Token 和 Refresh Token，并吊销旧 Refresh Token
func (r *Repository) RotateRefreshToken(ctx context.Context, oldTokenId, newTokenId string, userId uint64, newTokenHash string, newExpiresAt time.Time) error {
	tx, err := r.mysql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO refresh_tokens(token_id, user_id, token_hash, expires_at)
		VALUES (?, ?, ?, ?)
	`, newTokenId, userId, newTokenHash, newExpiresAt); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW(3), replaced_by_token_id = ?
		WHERE token_id = ? AND user_id = ? AND revoked_at IS NULL
	`, newTokenId, oldTokenId, userId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrRefreshTokenNotReplaceable
	}

	return tx.Commit()
}

// FindRefreshToken 按 Token Id 和 Token Hash 查询 Refresh Token 及所属用户
func (r *Repository) FindRefreshToken(ctx context.Context, tokenId, tokenHash string) (*RefreshTokenRecord, *UserRecord, error) {
	token := &RefreshTokenRecord{}
	user := &UserRecord{}

	err := r.mysql.QueryRowContext(ctx, `
		SELECT
			rt.id,
			rt.token_id,
			rt.user_id,
			rt.token_hash,
			rt.expires_at,
			rt.revoked_at,
			rt.created_at,
			rt.updated_at,
			rt.replaced_by_token_id,
			u.id,
			u.email,
			u.password_hash,
			u.name,
			u.last_login_at,
			u.created_at,
			u.updated_at
		FROM refresh_tokens rt
		JOIN users u ON u.id = rt.user_id
		WHERE rt.token_id = ? AND rt.token_hash = ?
		LIMIT 1
	`, tokenId, tokenHash).Scan(
		&token.Id,
		&token.TokenId,
		&token.UserId,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.ReplacedByTokenId,
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	return token, user, nil
}

// RevokeRefreshTokensByTokenId 按 Token Id 吊销 Refresh Token
func (r *Repository) RevokeRefreshTokensByTokenId(ctx context.Context, tokenId string) error {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW(3)
		WHERE token_id = ? AND revoked_at IS NULL
	`, tokenId)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

// RevokeRefreshTokensByEmail 按邮箱吊销 Refresh Token
func (r *Repository) RevokeRefreshTokensByEmail(ctx context.Context, email string) error {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW(3)
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE email = ?
			LIMIT 1
		) AND revoked_at IS NULL
	`, email)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	return err
}

// CreateEmailSendLog 创建邮件发送记录
func (r *Repository) CreateEmailSendLog(ctx context.Context, receiverEmail, senderEmail, scene, emailCode string, status consts.EmailSendStatus, errorMessage string) error {
	var nullErrorMessage sql.NullString
	if strings.TrimSpace(errorMessage) != "" {
		nullErrorMessage = sql.NullString{
			String: strings.TrimSpace(errorMessage),
			Valid:  true,
		}
	}

	_, err := r.mysql.ExecContext(ctx, `
		INSERT INTO email_send_logs(receiver_email, sender_email, scene, email_code, status, error_message)
		VALUES (?, ?, ?, ?, ?, ?)
	`, receiverEmail, senderEmail, scene, emailCode, status.String(), nullErrorMessage)

	return err
}
