package mysql

import (
	"context"
	"database/sql"
	"errors"
)

// CreateUser 创建用户
func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, name string) error {
	_, err := r.mysql.ExecContext(ctx, `
		INSERT INTO users(email, password_hash, name)
		VALUES (?, ?, ?)
	`, email, passwordHash, name)
	return err
}

// FindUserById 按用户 ID 查询用户
func (r *Repository) FindUserById(ctx context.Context, id uint64) (*UserRecord, error) {
	user := &UserRecord{}

	err := r.mysql.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name, last_login_at, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByEmail 按邮箱查询用户
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	user := &UserRecord{}

	err := r.mysql.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name, last_login_at, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Name, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePasswordByEmail 按邮箱更新用户密码
func (r *Repository) UpdatePasswordByEmail(ctx context.Context, email, passwordHash string) error {
	result, err := r.mysql.ExecContext(ctx, `
		UPDATE users
		SET password_hash = ?
		WHERE email = ?
	`, passwordHash, email)
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

	return nil
}
