package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Repository Redis 非关系型数据库服务
type Repository struct {
	redis *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{
		redis: rdb,
	}
}

// SaveEmailCode 保存邮箱验证码哈希
func (r *Repository) SaveEmailCode(ctx context.Context, scene, email, codeHash string, ttl time.Duration) error {
	return r.redis.Set(ctx, genEmailCodeKey(scene, email), codeHash, ttl).Err()
}

// GetEmailCode 获取邮箱验证码哈希
func (r *Repository) GetEmailCode(ctx context.Context, scene, email string) (string, error) {
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
func (r *Repository) ClearEmailCode(ctx context.Context, scene, email string) error {
	return r.redis.Del(ctx, genEmailCodeKey(scene, email)).Err()
}

// EmailCodeExists 判断邮箱验证码是否在有效期内
func (r *Repository) EmailCodeExists(ctx context.Context, scene, email string) (bool, error) {
	exists, err := r.redis.Exists(ctx, genEmailCodeKey(scene, email)).Result()
	return exists > 0, err
}

// RecordLoginFailure 记录一次登录失败
func (r *Repository) RecordLoginFailure(ctx context.Context, loginScene, email string, ttl time.Duration) error {
	return redis.NewScript(`
	local count = redis.call("INCR", KEYS[1])
	if count == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end
	return count
	`).Run(ctx, r.redis, []string{genLoginFailureKey(loginScene, email)}, int64(ttl.Seconds())).Err()
}

// ClearLoginFailure 清除登录失败计数
func (r *Repository) ClearLoginFailure(ctx context.Context, loginScene, email string) error {
	return r.redis.Del(ctx, genLoginFailureKey(loginScene, email)).Err()
}

// IsLoginLocked 判断指定登录方式是否达到登录失败次数上限
func (r *Repository) IsLoginLocked(ctx context.Context, loginScene, email string, limit int) (bool, time.Duration, error) {
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
		return false, 0, nil
	}
	return true, ttl, nil
}

// SetRefreshReuseLock 设置 Refresh Token 轮换并发锁
func (r *Repository) SetRefreshReuseLock(ctx context.Context, tokenId string, ttl time.Duration) (bool, error) {
	return r.redis.SetNX(ctx, genRefreshReuseLockKey(tokenId), "lock", ttl).Result()
}

// ClearRefreshReuseLock 清除 Refresh Token 轮换并发锁
func (r *Repository) ClearRefreshReuseLock(ctx context.Context, tokenId string) error {
	return r.redis.Del(ctx, genRefreshReuseLockKey(tokenId)).Err()
}

// BlacklistAccessToken 将 Access Token 的 JWT ID 加入黑名单
func (r *Repository) BlacklistAccessToken(ctx context.Context, tokenId string, ttl time.Duration) error {
	return r.redis.Set(ctx, genAccessBlacklistKey(tokenId), "blacklisted", ttl).Err()
}

// AccessTokenBlacklisted 判断 Access Token 的 JWT ID 是否在黑名单中
func (r *Repository) AccessTokenBlacklisted(ctx context.Context, tokenId string) (bool, error) {
	blacklisted, err := r.redis.Exists(ctx, genAccessBlacklistKey(tokenId)).Result()
	if err != nil {
		return false, err
	}
	return blacklisted > 0, nil
}

// NextProjectGroupSequence 分配新图像组的下一个序号
func (r *Repository) NextProjectGroupSequence(ctx context.Context, projectId uint64) (int64, error) {
	return r.redis.Incr(ctx, genProjectGroupSequenceKey(projectId)).Result()
}

// GetPresignURL 获取预签名 URL 缓存
func (r *Repository) GetPresignURL(ctx context.Context, key string) (string, error) {
	presignURL, err := r.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return presignURL, nil
}

// SavePresignURL 保存预签名 URL 缓存
func (r *Repository) SavePresignURL(ctx context.Context, key, presignURL string, ttl time.Duration) error {
	if key == "" || presignURL == "" || ttl <= 0 {
		return nil
	}
	return r.redis.Set(ctx, key, presignURL, ttl).Err()
}

// ClearPresignURL 清除预签名 URL 缓存
func (r *Repository) ClearPresignURL(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	return r.redis.Del(ctx, key).Err()
}
