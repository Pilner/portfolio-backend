package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Values struct {
	PortNumber                int    `mapstructure:"PORT"`
	DbConnectionUrl           string `mapstructure:"DB_CONNECTION_URL"`
	JwtSecretKey              string `mapstructure:"JWT_SECRET_KEY"`
	JwtTokenExpiryMinutes     int    `mapstructure:"JWT_TOKEN_EXPIRY_MINUTES"`
	RefreshTokenSecretKey     string `mapstructure:"REFRESH_TOKEN_SECRET_KEY"`
	RefreshTokenExpiryMinutes int    `mapstructure:"REFRESH_TOKEN_TOKEN_EXPIRY_MINUTES"`
	AppEnv                    string `mapstructure:"APP_ENV"`
	IsProd                    bool
}

func LoadConfig(env, path string) (c Values, err error) {
	viper.SetConfigName(".default")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	// Load the base .default.env file
	if err := viper.ReadInConfig(); err != nil {
		return Values{}, err
	}

	// If an environment is given
	if env != "" {
		env = strings.ToLower(env)
		viper.SetConfigName(fmt.Sprintf(".%s", env))
		if err = viper.MergeInConfig(); err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	// Try to merge in a .local.env for local, developer-override config
	viper.SetConfigName(".local")
	if err = viper.MergeInConfig(); err != nil {
		// Print error but continue if .local.env not found
		fmt.Printf("%v\n", err)
	}

	viper.AutomaticEnv()

	// Unmarshal config into the Values struct
	if err = viper.Unmarshal(&c); err != nil {
		return Values{}, nil
	}

	if c.AppEnv == "production" {
		c.IsProd = true
	}

	return c, nil
}
