package mysql

import (
	"database/sql"
)

// Repository MySQL 关系型数据库服务
type Repository struct {
	mysql *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		mysql: db,
	}
}
