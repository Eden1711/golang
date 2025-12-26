package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	TokenSecret   string        `mapstructure:"TOKEN_SECRET"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // Đường dẫn chứa file config
	viper.SetConfigName("app") // Tên file (không cần đuôi .env)
	viper.SetConfigType("env") // Loại file (env, json, xml...)

	viper.AutomaticEnv() // Tự động ghi đè nếu có biến môi trường (Environment Variable)

	// Đọc file config
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Parse dữ liệu vào struct
	err = viper.Unmarshal(&config)
	return
}
