package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Repository MySQL 服务
type Repository struct {
	mysql *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		mysql: db,
	}
}

// CreateUser 创建用户
func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, name string) error {
	_, err := r.mysql.ExecContext(ctx, `
		INSERT INTO users(email, password_hash, name)
		VALUES (?, ?, ?)
	`, email, passwordHash, name)
	if err != nil {
		return err
	}
	return nil
}

// FindUserById 按用户 ID 查询用户
func (r *Repository) FindUserById(ctx context.Context, id uint64) (*UserRecord, error) {
	var user UserRecord
	err := r.mysql.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByEmail 按邮箱查询用户
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	var user UserRecord
	err := r.mysql.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePasswordByEmail 按邮箱更新用户密码
func (r *Repository) UpdatePasswordByEmail(ctx context.Context, email, passwordHash string) error {
	_, err := r.mysql.ExecContext(ctx, `
		UPDATE users
		SET password_hash = ?
		WHERE email = ?
	`, passwordHash, email)
	if err != nil {
		return err
	}
	return nil
}

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

	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET last_login_at = NOW(3)
		WHERE id = ?
	`, userId); err != nil {
		return err
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
	var token RefreshTokenRecord
	var user UserRecord
	err := r.mysql.QueryRowContext(ctx, `
		SELECT rt.token_id, rt.user_id, rt.expires_at, rt.revoked_at,
		       u.id, u.email, u.password_hash, u.name
		FROM refresh_tokens rt
		JOIN users u ON u.id = rt.user_id
		WHERE rt.token_id = ? AND rt.token_hash = ?
		LIMIT 1
	`, tokenId, tokenHash).Scan(
		&token.TokenId, &token.UserId, &token.ExpiresAt, &token.RevokedAt,
		&user.Id, &user.Email, &user.PasswordHash, &user.Name,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	return &token, &user, nil
}

// RevokeRefreshTokenByTokenId 按 Token Id 吊销 Refresh Token
func (r *Repository) RevokeRefreshTokenByTokenId(ctx context.Context, tokenId string) error {
	_, err := r.mysql.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW(3)
		WHERE token_id = ? AND revoked_at IS NULL
	`, tokenId)
	return err
}

// RevokeUserRefreshTokensByEmail 按邮箱吊销 Refresh Token
func (r *Repository) RevokeUserRefreshTokensByEmail(ctx context.Context, email string) error {
	_, err := r.mysql.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = NOW(3)
		WHERE user_id = (
			SELECT id
			FROM users
			WHERE email = ?
			LIMIT 1
		) AND revoked_at IS NULL
	`, email)
	return err
}

// CreateEmailSendLog 创建邮件发送记录
func (r *Repository) CreateEmailSendLog(ctx context.Context, receiverEmail string, senderEmail string, scene string, emailCode string, status EmailSendStatus, errorMessage string) error {
	_, err := r.mysql.ExecContext(ctx, `
		INSERT INTO email_send_logs(receiver_email, sender_email, scene, email_code, status, error_message)
		VALUES (?, ?, ?, ?, ?, ?)
	`, receiverEmail, senderEmail, scene, emailCode, status, errorMessage)
	if err != nil {
		return err
	}
	return nil
}
