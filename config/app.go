package config

import "github.com/spf13/viper"

type AppConfig struct {
	AppGatewayPort     string
	ServiceAuthPort    string
	ServiceOrderPort   string
	ServiceProductPort string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		AppGatewayPort:     viper.GetString("APP_GATEWAY_PORT"),
		ServiceAuthPort:    viper.GetString("SERVICE_AUTH_PORT"),
		ServiceOrderPort:   viper.GetString("SERVICE_ORDER_PORT"),
		ServiceProductPort: viper.GetString("SERVICE_PRODUCT_PORT"),
	}
}
