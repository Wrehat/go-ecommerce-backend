package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppEnv    string `mapstructure:"APP_ENV"`
	AppPort   string `mapstructure:"PORT"`
	DBUri     string `mapstructure:"DB_URI"`
	RedisURL  string `mapstructure:"REDIS_URL"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() *AppConfig {
	// Todo: masukkan nama file yang akan dibaca
	viper.SetConfigFile(".env")

	// Todo: baca env variable dari sistem opearsi
	viper.AutomaticEnv()

	// Todo: baca dari file .env
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️ File .env tidak ditemukan, menggunakan nilai dari sistem OS: %v\n", err)
	}

	// Todo: pindah isi .env ke struct
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("❌ Gagal memparsing konfigurasi: %v\n", err)
	}

	// Todo: beri nilai default jika .env tidak ada
	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	if config.AppEnv == "" {
		config.AppEnv = "development"
	}

	return &config
}
