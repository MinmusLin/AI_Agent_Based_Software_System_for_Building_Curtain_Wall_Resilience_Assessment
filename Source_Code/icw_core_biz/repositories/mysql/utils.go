package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"icw_core_biz/configs"
)

// MySQLDSN 根据拆分后的 MySQL 配置生成连接 DSN
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

// nullInt64ToUint16 将 sql.NullInt64 类型转换为 uint16
func nullInt64ToUint16(value sql.NullInt64) uint16 {
	if !value.Valid || value.Int64 <= 0 {
		return 0
	}
	return uint16(value.Int64)
}

// timeToString 将 time.Time 类型转换为 string
func timeToString(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.DateTime)
}

// lockProjectForUpdate 添加 MySQL 项目锁
func lockProjectForUpdate(ctx context.Context, tx *sql.Tx, userId, projectId uint64) (bool, error) {
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
