package database

import (
	"server/global"
	"time"
)

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	global.MODEL
	Jwt       string    `json:"jwt" gorm:"type:varchar(512);uniqueIndex;not null;comment:JWT黑名单"` // Jwt
	ExpiresAt time.Time `json:"expiresAt" gorm:"index;comment:过期时间(用于清理)"`
}
