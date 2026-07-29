package config

type Jwt struct {
	AccessTokenSecret      string `json:"access_token_secret" yaml:"access_token_secret"`
	RefreshTokenSecret     string `json:"refresh_token_secret" yaml:"refresh_token_secret"`
	AccessTokenExpiryTime  string `json:"access_token_expiry_time" yaml:"access_token_expiry_time"`
	RefreshTokenExpiryTime string `json:"refresh_token_expiry_time" yaml:"refresh_token_expiry_time"`
	Issuer                 string `json:"issuer" yaml:"issuer"`
}
