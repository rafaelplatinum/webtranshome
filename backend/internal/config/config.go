package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	Database struct {
		DataSource string
	}

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	PasswordReset struct {
		HMACSecret           string
		WhatsAppGatewayURL   string
		WhatsAppGatewayToken string
		EmailGatewayURL      string
		EmailGatewayToken    string
	}

	ImageStorage struct {
		Directory string
		MaxBytes  int64
	}
}
