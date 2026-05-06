package mysql

import (
	"database/sql"
)

// Repository MySQL 关系型数据库服务
type Repository struct {
	mysql *sql.DB
}

// NewRepository 创建 MySQL 关系型数据库服务
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		mysql: db,
	}
}
