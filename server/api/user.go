package api

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type UserApi struct {
}

// Register 注册
func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}

	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}

	// 注册成功后，生成 token 并返回
	userApi.TokenNext(c, user)
}

// Login 登录接口，根据不同的登录方式调用不同的登录方法
func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	case "qq":
		userApi.QQLogin(c)
	default:
		userApi.EmailLogin(c)
	}
}

// EmailLogin 邮箱登录
func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 校验验证码
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := userService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage("Failed to login", c)
			return
		}

		// 登录成功后生成 token
		userApi.TokenNext(c, user)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

// QQLogin QQ登录
func (userApi *UserApi) QQLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("Code is required", c)
		return
	}

	// 获取访问令牌
	accessTokenResponse, err := qqService.GetAccessTokenByCode(code)
	if err != nil || accessTokenResponse.Openid == "" {
		global.Log.Error("Invalid code", zap.Error(err))
		response.FailWithMessage("Invalid code", c)
		return
	}

	// 根据访问令牌进行QQ登录
	user, err := userService.QQLogin(accessTokenResponse)
	if err != nil {
		global.Log.Error("Failed to login:", zap.Error(err))
		response.FailWithMessage("Failed to login", c)
		return
	}

	// 登录成功后生成 token
	userApi.TokenNext(c, user)
}

func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	// 检查用户是否被冻结
	if user.Freeze {
		response.FailWithMessage("The user is frozen, contact the administrator", c)
		return
	}

	baseClaims := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()
	jti := j.CreateJti() // 复用同一个 JWT 实例生成 jti

	// 创建访问令牌
	accessClaims := j.CreateAccessClaims(baseClaims, jti)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get accessToken:", zap.Error(err))
		response.FailWithMessage("Failed to get accessToken", c)
		return
	}

	// 创建刷新令牌
	refreshClaims := j.CreateRefreshClaims(baseClaims, jti)
	refreshToken, err := j.CreateRefreshToken(refreshClaims)
	if err != nil {
		global.Log.Error("Failed to get refreshToken:", zap.Error(err))
		response.FailWithMessage("Failed to get refreshToken", c)
		return
	}

	// 未开启多地点登录拦截，直接返回
	if !global.Config.System.UseMultipoint {
		userApi.respondWithLogin(c, user, accessToken, accessClaims, refreshToken, refreshClaims)
		return
	}

	// 查询该用户当前活跃会话的 jti
	oldJti, err := jwtService.GetRedisJWT(user.UUID)
	switch {
	case errors.Is(err, redis.Nil):
		// 首次登录（Redis 中不存在该用户的会话）
		if err := jwtService.SetRedisJWT(refreshClaims.ID, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}
		userApi.respondWithLogin(c, user, accessToken, accessClaims, refreshToken, refreshClaims)

	case err != nil:
		// Redis 查询出错
		global.Log.Error("Failed to get login status:", zap.Error(err))
		response.FailWithMessage("Failed to get login status", c)

	default:
		// 已存在活跃会话，踢掉旧会话，写入新会话
		// 旧会话的剩余寿命 = Redis 中 uuid 这个 key 的 TTL（它正是按旧的 Refresh 有效期写入的），
		// 必须在下面 SetRedisJWT 覆盖它之前读取，否则读到的就是新会话的寿命
		oldTTL, ttlErr := global.Redis.TTL(user.UUID.String()).Result()
		if ttlErr != nil {
			global.Log.Warn("读取旧会话剩余有效期失败，改用黑名单默认有效期", zap.Error(ttlErr))
			oldTTL = 0 // 交给 JoinInBlacklist 用配置里的 Refresh 有效期兜底
		}

		if err := jwtService.JoinInBlacklist(oldJti, oldTTL); err != nil {
			global.Log.Error("Failed to invalidate jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt", c)
			return
		}

		if err := jwtService.SetRedisJWT(refreshClaims.ID, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		userApi.respondWithLogin(c, user, accessToken, accessClaims, refreshToken, refreshClaims)
	}
}

// respondWithLogin 设置刷新令牌并返回登录成功响应，供 TokenNext 各分支复用
func (userApi *UserApi) respondWithLogin(
	c *gin.Context,
	user database.User,
	accessToken string,
	accessClaims request.JwtCustomClaims,
	refreshToken string,
	refreshClaims request.JwtCustomRefreshClaims,
) {
	ttl := int(refreshClaims.ExpiresAt.Unix() - time.Now().Unix())
	utils.SetRefreshToken(c, refreshToken, ttl)
	c.Set("user_id", user.ID)
	response.OkWithDetailed(response.Login{
		User:                 user,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
	}, "Successful login", c)
}

// ForgotPassword 找回密码
func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)

	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	err = userService.ForgotPassword(req)
	if err != nil {
		global.Log.Error("Failed to retrieve the password:", zap.Error(err))
		response.FailWithMessage("Failed to retrieve the password", c)
		return
	}
	response.OkWithMessage("Successfully retrieved", c)
}

// UserCard 获取用户卡片信息
func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to get card:", zap.Error(err))
		response.FailWithMessage("Failed to get card", c)
		return
	}
	response.OkWithData(userCard, c)
}

// Logout 登出
func (userApi *UserApi) Logout(c *gin.Context) {
	userService.Logout(c)
	response.OkWithMessage("Successful logout", c)
}

// UserResetPassword 修改密码
func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = userService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("Failed to modify:", zap.Error(err))
		response.FailWithMessage("Failed to modify, orginal password does not match the current account", c)
		return
	}

	userService.Logout(c)
	response.OkWithMessage("Successfully changed password, please log in again", c)
}

// UserInfo 获取个人信息
func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user information:", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}
	response.OkWithData(user, c)
}

// UserChangeInfo 修改个人信息
func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = userService.UserChangeInfo(req)
	if err != nil {
		global.Log.Error("Failed to change user information:", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}
	response.OkWithMessage("Successfully changed user information", c)
}

// UserWeather 获取天气
func (userApi *UserApi) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	weather, err := userService.UserWeather(ip)
	if err != nil {
		global.Log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", c)
		return
	}
	response.OkWithData(weather, c)
}

// UserChart 获取用户图表数据，登录和注册人数
func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to user chart", c)
		return
	}
	response.OkWithData(data, c)
}

// UserList 获取用户列表
func (userApi *UserApi) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// UserFreeze 冻结用户
func (userApi *UserApi) UserFreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = userService.UserFreeze(req)
	if err != nil {
		global.Log.Error("Failed to freeze user:", zap.Error(err))
		response.FailWithMessage("Failed to freeze user", c)
		return
	}
	response.OkWithMessage("Successfully freeze user", c)
}

// UserUnfreeze 解冻用户
func (userApi *UserApi) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.UserUnfreeze(req)
	if err != nil {
		global.Log.Error("Failed to unfreeze user:", zap.Error(err))
		response.FailWithMessage("Failed to unfreeze user", c)
		return
	}
	response.OkWithMessage("Successfully unfreeze user", c)
}

// UserLoginList 获取登录日志列表
func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user login list:", zap.Error(err))
		response.FailWithMessage("Failed to get user login list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
