package redis

import (
	"fmt"
	"strings"

	"icw_common/utils"
)

// genEmailCodeKey 生成邮箱验证码哈希 Key
func genEmailCodeKey(scene, emailHash string) string {
	return fmt.Sprintf("auth:email_code:%s:%s", scene, emailHash)
}

// genLoginFailureKey 生成指定登录方式的登录失败次数 Key
func genLoginFailureKey(scene, emailHash string) string {
	return fmt.Sprintf("auth:login_fail:%s:%s", scene, emailHash)
}

// genRefreshReuseLockKey 生成 Refresh Token 轮换并发锁 Key
func genRefreshReuseLockKey(tokenId string) string {
	return fmt.Sprintf("auth:refresh_reuse_lock:%s", tokenId)
}

// genAccessBlacklistKey 生成 Access Token 黑名单 Key
func genAccessBlacklistKey(tokenId string) string {
	return fmt.Sprintf("auth:access_blacklist:%s", tokenId)
}

// genProjectGroupSequenceKey 生成新图像组的下一个序号 Key
func genProjectGroupSequenceKey(projectId uint64) string {
	return fmt.Sprintf("project:assets:group_seq:%s", utils.Encode(projectId))
}

// genSocketTicketKey 生成 WebSocket 连接票据上下文 Key
func genSocketTicketKey(ticketHash string) string {
	return fmt.Sprintf("socket:ticket:%s", ticketHash)
}

// genPresignURLKey 按 MinIO 对象 Key 生成预签名 URL 缓存 Key
func genPresignURLKey(objectKey string) string {
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return ""
	}
	return fmt.Sprintf("presign:get:%s", objectKey)
}
