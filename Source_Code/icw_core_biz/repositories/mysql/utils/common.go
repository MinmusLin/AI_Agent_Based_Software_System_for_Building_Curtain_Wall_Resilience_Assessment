package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
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

// JsonStringOrEmptyObject 将空 JSON 对象字符串兜底为 "{}"
func JsonStringOrEmptyObject(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "{}"
	}
	return value
}

// UnmarshalRegions 将数据库 JSON 字符串解析为检测区域列表
func UnmarshalRegions[T any](regions sql.NullString, target *[]*T) error {
	if target == nil {
		return nil
	}
	*target = make([]*T, 0)
	if !regions.Valid || strings.TrimSpace(regions.String) == "" {
		return nil
	}
	return json.Unmarshal([]byte(regions.String), target)
}

// NullBool 将 sql.NullBool 转换为 bool
func NullBool(value sql.NullBool) bool {
	return value.Valid && value.Bool
}

// NullFloat64 将 sql.NullFloat64 转换为 float64
func NullFloat64(value sql.NullFloat64) float64 {
	if !value.Valid {
		return 0
	}
	return value.Float64
}

// NullString 将 sql.NullString 转换为 string
func NullString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

// NullUint32 将 sql.NullInt64 转换为 uint32
func NullUint32(value sql.NullInt64) uint32 {
	if !value.Valid || value.Int64 <= 0 {
		return 0
	}
	return uint32(value.Int64)
}

// NullUint64 将 sql.NullInt64 转换为 uint64
func NullUint64(value sql.NullInt64) uint64 {
	if !value.Valid || value.Int64 <= 0 {
		return 0
	}
	return uint64(value.Int64)
}

// NullTimeString 将 sql.NullTime 转换为时间字符串
func NullTimeString(value sql.NullTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.DateTime)
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
