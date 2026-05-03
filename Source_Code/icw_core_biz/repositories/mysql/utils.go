package mysql

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
)

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
