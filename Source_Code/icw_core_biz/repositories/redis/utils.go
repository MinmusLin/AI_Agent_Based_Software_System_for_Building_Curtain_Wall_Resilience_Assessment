package redis

import (
	"fmt"

	"icw_core_biz/utils"
)

// genEmailCodeKey 生成邮箱验证码哈希 Key
func genEmailCodeKey(scene, email string) string {
	return fmt.Sprintf("auth:email_code:%s:%s", scene, email) // todo email hash
}

// genLoginFailureKey 生成指定登录方式的登录失败次数 Key
func genLoginFailureKey(scene, email string) string {
	return fmt.Sprintf("auth:login_fail:%s:%s", scene, email) // todo email hash
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

// GenDefaultAvatarPresignURLKey 生成用户默认头像预签名 URL 缓存 Key
func GenDefaultAvatarPresignURLKey(emailHash string) string {
	return fmt.Sprintf("presign:avatar:default:%s", emailHash)
}

// GenCustomAvatarPresignURLKey 生成用户自定义头像预签名 URL 缓存 Key
func GenCustomAvatarPresignURLKey(emailHash string) string {
	return fmt.Sprintf("presign:avatar:custom:%s", emailHash)
}

// GenProjectThumbnailPresignURLKey 生成项目缩略图预签名 URL 缓存 Key
func GenProjectThumbnailPresignURLKey(userId, projectId uint64) string {
	return fmt.Sprintf("presign:project:thumbnail:%s:%s", utils.Encode(userId), utils.Encode(projectId))
}

// GenProjectImageOriginalPresignURLKey 生成项目图像原图预签名 URL 缓存 Key
func GenProjectImageOriginalPresignURLKey(userId, projectId uint64, imageUuid string) string {
	return fmt.Sprintf("presign:image:original:%s:%s:%s", utils.Encode(userId), utils.Encode(projectId), imageUuid)
}

// GenProjectImageThumbnailPresignURLKey 生成项目图像缩略图预签名 URL 缓存 Key
func GenProjectImageThumbnailPresignURLKey(userId, projectId uint64, imageUuid string) string {
	return fmt.Sprintf("presign:image:thumbnail:%s:%s:%s", utils.Encode(userId), utils.Encode(projectId), imageUuid)
}
