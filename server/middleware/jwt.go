package middleware

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var jwtService = service.ServiceGroupApp.JwtService

// JWTAuth 验证请求中的 Access Token 是否合法，并检查其 jti 是否在黑名单中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := utils.GetAccessToken(c)
		refreshToken := utils.GetRefreshToken(c)

		j := utils.NewJWT()

		// ========== 1. 先尝试解析 Access Token ==========
		claims, err := j.ParseAccessToken(accessToken)

		if err == nil {
			// Access Token 有效，直接用它的 jti 查黑名单
			res, bErr := jwtService.IsInBlacklist(claims.ID)
			if bErr != nil {
				utils.ClearRefreshToken(c)
				response.NoAuth("Service unavailable", c)
				c.Abort()
				return
			}
			if res > 0 {
				utils.ClearRefreshToken(c)
				response.NoAuth("Account logged in from another location or token is invalid", c)
				c.Abort()
				return
			}

			c.Set("claims", claims)
			c.Next()
			return
		}

		// ========== 2. Access Token 无效，判断是否走 Refresh 流程 ==========
		// 只有「Access 为空」或「Access 过期」才尝试用 Refresh 换新
		if accessToken != "" && !errors.Is(err, utils.TokenExpired) {
			// 其他错误（格式错误、签名错误等），直接拒绝
			utils.ClearRefreshToken(c)
			response.NoAuth("Invalid access token", c)
			c.Abort()
			return
		}

		// ========== 3. 解析 Refresh Token ==========
		refreshClaims, err := j.ParseRefreshToken(refreshToken)
		if err != nil {
			utils.ClearRefreshToken(c)
			response.NoAuth("Refresh token expired or invalid", c)
			c.Abort()
			return
		}

		// 用 Refresh Token 的 jti 查黑名单
		// 因为 Access/Refresh 共享 jti，这里查的是同一个 key
		res, bErr := jwtService.IsInBlacklist(refreshClaims.ID)
		if bErr != nil {
			utils.ClearRefreshToken(c)
			response.NoAuth("Service unavailable", c)
			c.Abort()
			return
		}
		if res > 0 {
			utils.ClearRefreshToken(c)
			response.NoAuth("Account logged in from another location or token is invalid", c)
			c.Abort()
			return
		}

		// ========== 4. 用 Refresh 换发新 Access ==========
		var user database.User
		if err := global.DB.Select("uuid", "role_id").Take(&user, refreshClaims.UserID).Error; err != nil {
			utils.ClearRefreshToken(c)
			response.NoAuth("The user does not exist", c)
			c.Abort()
			return
		}

		// 关键：复用 Refresh 的 jti，保证会话 ID 不变
		newAccessClaims := j.CreateAccessClaims(request.BaseClaims{
			UserID: refreshClaims.UserID,
			UUID:   user.UUID,
			RoleID: user.RoleID,
		}, refreshClaims.ID)

		newAccessToken, err := j.CreateAccessToken(newAccessClaims)
		if err != nil {
			utils.ClearRefreshToken(c)
			response.NoAuth("Failed to create new access token", c)
			c.Abort()
			return
		}

		c.Header("new-access-token", newAccessToken)
		c.Header("new-access-expires-at", strconv.FormatInt(newAccessClaims.ExpiresAt.Unix(), 10))

		c.Set("claims", &newAccessClaims)
		c.Next()
	}
}
