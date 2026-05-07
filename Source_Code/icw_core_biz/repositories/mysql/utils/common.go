package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"icw_core_biz/configs"
)

// MySQLDSN 根据配置生成 MySQL 连接 DSN
func MySQLDSN(cfg configs.Config) string {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.MySQLUsername
	mysqlConfig.Passwd = cfg.MySQLPassword
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = cfg.MySQLAddr
	mysqlConfig.DBName = cfg.MySQLDatabase
	mysqlConfig.ParseTime = true
	mysqlConfig.Loc = time.Local
	mysqlConfig.Params = map[string]string{
		"charset": "utf8mb4",
	}
	return mysqlConfig.FormatDSN()
}

// IsDuplicateEntryError 判断是否为 MySQL 唯一键冲突错误
func IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

// LockProjectForUpdate 添加 MySQL 项目锁
func LockProjectForUpdate(ctx context.Context, tx *sql.Tx, userId, projectId uint64) (bool, error) {
	var id uint64
	err := tx.QueryRowContext(ctx, `
		SELECT id
		FROM projects
		WHERE id = ? AND user_id = ?
		FOR UPDATE
	`, projectId, userId).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// JsonOrEmptyArray 将空 JSON 数组字段兜底为 "[]"
func JsonOrEmptyArray(raw json.RawMessage) string {
	if len(raw) == 0 {
		return "[]"
	}
	return string(raw)
}

// CheckRowsAffected 检查 SQL 执行结果是否至少影响一行
func CheckRowsAffected(result sql.Result, err error) error {
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
