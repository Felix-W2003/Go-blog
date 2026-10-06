package service

import (
	"errors"
	"server/config"
	"server/global"
	"server/model/appTypes"
	"server/utils"

	"gorm.io/gorm"
)

type ConfigService struct {
}

func (configService *ConfigService) UpdateWebsite(website config.Website) error {
	oldArray := []string{
		global.Config.Website.Logo,
		global.Config.Website.FullLogo,
		global.Config.Website.QQImage,
		global.Config.Website.WechatImage,
	}

	newArray := []string{
		website.Logo,
		website.FullLogo,
		website.QQImage,
		website.WechatImage,
	}

	added, removed := utils.DiffArrays(oldArray, newArray)

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.InitImagesCategory(global.DB, removed); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(global.DB, added, appTypes.System); err != nil {
			return err
		}
		global.Config.Website = website
		if err := utils.SaveYAML(); err != nil {
			return err
		}
		return nil
	})
}

func (configService *ConfigService) UpdateSystem(system config.System) error {
	global.Config.System.UseMultipoint = system.UseMultipoint
	global.Config.System.SessionsSecret = system.SessionsSecret
	global.Config.System.OssType = system.OssType
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateEmail(email config.Email) error {
	global.Config.Email = email
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateQQ(qq config.QQ) error {
	global.Config.QQ = qq
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateQiniu(qiniu config.Qiniu) error {
	global.Config.Qiniu = qiniu
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateJwt(jwt config.Jwt) error {
	refreshExp, err := utils.ParseDuration(jwt.RefreshTokenExpiryTime)
	if err != nil || refreshExp <= 0 {
		return errors.New("refresh_token_expiry_time 格式非法（示例：7d、2h、30m）")
	}
	accessExp, err := utils.ParseDuration(jwt.AccessTokenExpiryTime)
	if err != nil || accessExp <= 0 {
		return errors.New("access_token_expiry_time 格式非法（示例：7d、2h、30m）")
	}
	// 空密钥 = []byte("") 是合法的 HMAC key，等于谁都能伪造任意身份的 Token
	if jwt.AccessTokenSecret == "" || jwt.RefreshTokenSecret == "" {
		return errors.New("令牌密钥不能为空")
	}
	global.Config.Jwt = jwt
	return utils.SaveYAML()
}

func (configService *ConfigService) UpdateGaode(gaode config.Gaode) error {
	global.Config.Gaode = gaode
	return utils.SaveYAML()
}
