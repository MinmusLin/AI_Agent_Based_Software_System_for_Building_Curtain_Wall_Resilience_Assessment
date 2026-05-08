package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/core/biz"

	"icw_core_biz/configs"
	"icw_core_biz/repositories/mysql/model"
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
func JsonOrEmptyArray(value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "[]", err
	}
	if string(data) == "null" {
		return "[]", nil
	}
	return string(data), nil
}

// JsonStringOrEmptyObject 将空 JSON 对象字符串兜底为 "{}"
func JsonStringOrEmptyObject(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "{}"
	}
	return value
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

// NormalizeDetectionTaskCodes 标准化图像检测子任务代码
func NormalizeDetectionTaskCodes(taskCodes []string) ([]string, error) {
	normalized := make([]string, 0, len(taskCodes))
	seen := map[activitypb.DetectionTaskCode_Value]bool{}
	for _, taskCode := range taskCodes {
		if strings.TrimSpace(taskCode) == "" {
			continue
		}
		code := enum.ParseDetectionTaskCode(taskCode)
		if code == activitypb.DetectionTaskCode_Unknown {
			return nil, model.ErrUnsupportedDetectionTaskCode
		}
		if seen[code] {
			continue
		}
		normalized = append(normalized, enum.DetectionTaskCodeString(code))
		seen[code] = true
	}
	return normalized, nil
}

// ClassificationNodeStatus 根据图像检测主任务状态推导图像检测分类任务状态
func ClassificationNodeStatus(task *model.ProjectDetectionTaskRecord) bizpb.ProjectDetectionSubTaskStatus_Value {
	switch task.Status {
	case bizpb.ProjectDetectionTaskStatus_Classifying:
		return bizpb.ProjectDetectionSubTaskStatus_Pending
	case bizpb.ProjectDetectionTaskStatus_Detecting,
		bizpb.ProjectDetectionTaskStatus_Summarizing,
		bizpb.ProjectDetectionTaskStatus_Succeeded:
		return bizpb.ProjectDetectionSubTaskStatus_Succeeded
	case bizpb.ProjectDetectionTaskStatus_Failed:
		if !task.CorrosionShouldExecute &&
			!task.CrackShouldExecute &&
			!task.StainShouldExecute &&
			!task.FlatnessShouldExecute &&
			!task.SpallingShouldExecute &&
			!task.SummaryShouldExecute {
			return bizpb.ProjectDetectionSubTaskStatus_Failed
		}
		return bizpb.ProjectDetectionSubTaskStatus_Succeeded
	default:
		return bizpb.ProjectDetectionSubTaskStatus_Unknown
	}
}
