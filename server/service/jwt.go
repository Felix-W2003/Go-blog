package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/utils"
	"time"

	"github.com/gofrs/uuid"
)

// JwtService 提供与JWT相关的服务
type JwtService struct {
}

// blacklistPrefix 黑名单 key 前缀：集中定义，避免多处硬编码写歪导致黑名单静默失效
const blacklistPrefix = "jwt:BlackList:"

// blacklistKey 由 jti 生成黑名单 key（存 SHA256 摘要，避免明文凭据落在 Redis 里）
func blacklistKey(jti string) string {
	return blacklistPrefix + utils.JtiSHA256(jti)
}

// SetRedisJWT 将JWT存储到Redis中
func (jwtService *JwtService) SetRedisJWT(jti string, uuid uuid.UUID) error {
	// 解析配置中的JWT过期时间
	dr, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	if err != nil {
		return err
	}
	if dr <= 0 {
		return errors.New("RefreshTokenExpiryTime 非法，拒绝写入以免会话 key 永不过期")
	}
	// //把jti转字节数组，做SHA256哈希，返回固定32字节数组
	// jtiHash := utils.JtiSHA256(jti)
	// 设置JWT在Redis中的过期时间
	return global.Redis.Set(uuid.String(), jti, dr).Err()
}

// GetRedisJWT 从Redis中获取JWT
func (jwtService *JwtService) GetRedisJWT(uuid uuid.UUID) (string, error) {
	// 从Redis获取指定uuid对应的JWT
	res, err := global.Redis.Get(uuid.String()).Result()
	return res, err
}

// JoinInBlacklist 将jti加入黑名单
// expTime <= 0 表示上游没有传有效期（例如冻结用户），此时用配置里的 Refresh Token 有效期兜底
func (jwtService *JwtService) JoinInBlacklist(jti string, expTime time.Duration) error {
	if expTime <= 0 {
		exp, err := utils.ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
		if err != nil {
			// 不能吞掉：err 被吞掉后 exp 会留在 0，而 Redis.Set 的过期时间为 0 表示「永不过期」
			return fmt.Errorf("解析 RefreshTokenExpiryTime 失败，无法确定黑名单有效期: %w", err)
		}
		expTime = exp
	}

	// 兜底：ParseDuration("0s") / ("0d") 会返回 (0, nil)，解析成功但值为 0；
	// 负数同理 —— go-redis 只在 expiration > 0 时才带 EX/PX，否则就是永久 key
	if expTime <= 0 {
		return errors.New("黑名单有效期非法，拒绝写入以避免 key 永不过期")
	}

	return global.Redis.Set(blacklistKey(jti), 1, expTime).Err()
}

// IsInBlacklist 检查JWT是否在黑名单中
func (jwtService *JwtService) IsInBlacklist(jti string) (int64, error) {
	// 从黑名单缓存中检查JWT是否存在
	// _, ok := global.BlackCache.Get(jti)
	return global.Redis.Exists(blacklistKey(jti)).Result()
}

// LoadAll 从数据库加载所有的JWT黑名单并加入缓存
// func LoadAll() {
// 	var data []string
// 	// 从数据库中获取所有的黑名单JWT
// 	if err := global.DB.Model(&database.JwtBlacklist{}).Pluck("jwt", &data).Error; err != nil {
// 		// 如果获取失败，记录错误日志
// 		global.Log.Error("Failed to load JWT blacklist from the database", zap.Error(err))
// 		return
// 	}
// 	// 将所有JWT添加到BlackCache缓存中
// 	for i := 0; i < len(data); i++ {
// 		global.BlackCache.SetDefault(data[i], struct{}{})
// 		fmt.Println(data[i])
// 	}

// }
