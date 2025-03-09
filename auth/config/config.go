package config

import "github.com/spf13/viper"

type Config struct {
	Port     string `mapstructure:"port"`
	Keycloak struct {
		Realm        string `mapstructure:"realm"`
		URL          string `mapstructure:"url"`
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"` // Not asking for redirect URIs since we are not using oauth2
		RedirectURI  string `mapstructure:"redirect_uri"`
	}
	AuthProviders AuthProviders `mapstructure:"auth_providers"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
	}
	EmailService EmailServiceConfig `mapstructure:"email_service"`
	Cookies CookiesConfig `mapstructure:"cookies"`
	JWTSecretKey string `mapstructure:"jwt_secret_key"`
}

type CookiesConfig struct {
	HTTPOnly bool `mapstructure:"http_only"`
	Secure   bool `mapstructure:"secure"`
	Domain   struct {
		Enabled bool   `mapstructure:"enabled"`
		Value   string `mapstructure:"value"`
	} `mapstructure:"domain"`
	Auth struct {
		Enabled bool `mapstructure:"enabled"`
		AccessToken string `mapstructure:"access_token"`
		RefreshToken string `mapstructure:"refresh_token"`
	}
}

type AuthProviders []struct {
	Name         string `mapstructure:"name"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
}

type EmailServiceConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		panic("Error reading config file : " + err.Error())
	}
}

func GetConfig() Config {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic("Error unmarshalling config : " + err.Error())
	}
	return config
}
