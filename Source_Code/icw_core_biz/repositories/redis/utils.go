package redis

import (
	"fmt"
)

// genEmailCodeKey 生成邮箱验证码哈希 Key
func genEmailCodeKey(scene, email string) string {
	return fmt.Sprintf("auth:email_code:%s:%s", scene, email)
}

// genLoginFailureKey 生成指定登录方式的登录失败次数 Key
func genLoginFailureKey(loginScene, email string) string {
	return fmt.Sprintf("auth:%s_login_fail:%s", loginScene, email)
}

// genRefreshReuseLockKey 生成 Refresh Token 轮换并发锁 Key
func genRefreshReuseLockKey(tokenId string) string {
	return fmt.Sprintf("auth:refresh_reuse_lock:%s", tokenId)
}

// genAccessBlacklistKey 生成 Access Token 黑名单 Key
func genAccessBlacklistKey(tokenId string) string {
	return fmt.Sprintf("auth:access_blacklist:%s", tokenId)
}
