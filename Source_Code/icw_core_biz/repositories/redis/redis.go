package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRepository Redis 服务
type RedisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{
		redis: rdb,
	}
}

// SaveEmailCode 保存邮箱验证码哈希
func (r *RedisRepository) SaveEmailCode(ctx context.Context, scene, email string, codeHash string, ttl time.Duration) error {
	return r.redis.Set(ctx, genEmailCodeKey(scene, email), codeHash, ttl).Err()
}

// GetEmailCode 获取邮箱验证码哈希
func (r *RedisRepository) GetEmailCode(ctx context.Context, scene, email string) (string, error) {
	codeHash, err := r.redis.Get(ctx, genEmailCodeKey(scene, email)).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return codeHash, nil
}

// ClearEmailCode 清除邮箱验证码哈希
func (r *RedisRepository) ClearEmailCode(ctx context.Context, scene, email string) error {
	return r.redis.Del(ctx, genEmailCodeKey(scene, email)).Err()
}

// EmailCodeExists 判断邮箱验证码是否在有效期内
func (r *RedisRepository) EmailCodeExists(ctx context.Context, scene, email string) (bool, error) {
	exists, err := r.redis.Exists(ctx, genEmailCodeKey(scene, email)).Result()
	return exists > 0, err
}

// RecordLoginFailure 记录一次登录失败
func (r *RedisRepository) RecordLoginFailure(ctx context.Context, loginScene, email string, ttl time.Duration) error {
	return redis.NewScript(`
	local count = redis.call("INCR", KEYS[1])
	if count == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end
	return count
	`).Run(ctx, r.redis, []string{genLoginFailureKey(loginScene, email)}, int64(ttl.Seconds())).Err()
}

// ClearLoginFailure 清除登录失败计数
func (r *RedisRepository) ClearLoginFailure(ctx context.Context, loginScene, email string) error {
	return r.redis.Del(ctx, genLoginFailureKey(loginScene, email)).Err()
}

// IsLoginLocked 判断指定登录方式是否达到登录失败次数上限
func (r *RedisRepository) IsLoginLocked(ctx context.Context, loginScene, email string, limit int) (bool, time.Duration, error) {
	key := genLoginFailureKey(loginScene, email)
	count, err := r.redis.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	if count < limit {
		return false, 0, nil
	}
	ttl, err := r.redis.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}
	if ttl <= 0 {
		_ = r.redis.Del(ctx, key).Err()
		return false, 0, nil
	}
	return true, ttl, nil
}

// SetRefreshReuseLock 设置 Refresh Token 轮换并发锁
func (r *RedisRepository) SetRefreshReuseLock(ctx context.Context, tokenId string, ttl time.Duration) (bool, error) {
	return r.redis.SetNX(ctx, genRefreshReuseLockKey(tokenId), "lock", ttl).Result()
}

// ClearRefreshReuseLock 清除 Refresh Token 轮换并发锁
func (r *RedisRepository) ClearRefreshReuseLock(ctx context.Context, tokenId string) error {
	return r.redis.Del(ctx, genRefreshReuseLockKey(tokenId)).Err()
}

// BlacklistAccessToken 将 Access Token 的 JWT ID 加入黑名单
func (r *RedisRepository) BlacklistAccessToken(ctx context.Context, tokenId string, ttl time.Duration) error {
	return r.redis.Set(ctx, genAccessBlacklistKey(tokenId), "blacklisted", ttl).Err()
}

// AccessTokenBlacklisted 判断 Access Token 的 JWT ID 是否在黑名单中
func (r *RedisRepository) AccessTokenBlacklisted(ctx context.Context, tokenId string) (bool, error) {
	blacklisted, err := r.redis.Exists(ctx, genAccessBlacklistKey(tokenId)).Result()
	if err != nil {
		return false, err
	}
	return blacklisted > 0, nil
}

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
