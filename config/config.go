package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`

	DBHost         string        `mapstructure:"DB_HOST"`
	DBUser         string        `mapstructure:"DB_USER"`
	DBPassword     string        `mapstructure:"DB_PASSWORD"`
	DBName         string        `mapstructure:"DB_NAME"`
	DBPort         string        `mapstructure:"DB_PORT"`
	DBMaxOpenConns int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxLifeTime  time.Duration `mapstructure:"DB_MAX_LIFE_TIME"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	AdminSecretKey    string        `mapstructure:"ADMIN_SECRET_KEY"`
	JWTSecret         string        `mapstructure:"JWT_SECRET"`
	AccessTokenExpiry time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRY"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("PORT")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_MAX_OPEN_CONNS")
	viper.BindEnv("DB_MAX_IDLE_CONNS")
	viper.BindEnv("DB_MAX_LIFE_TIME")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("REDIS_DB")
	viper.BindEnv("ADMIN_SECRET_KEY")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("ACCESS_TOKEN_EXPIRY")
	viper.BindEnv("REFRESH_TOKEN_EXPIRY")

	viper.SetDefault("ENVIRONMENT", "dev")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
	viper.SetDefault("DB_MAX_LIFE_TIME", 5*time.Minute)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("ACCESS_TOKEN_EXPIRY", 15*time.Minute)
	viper.SetDefault("REFRESH_TOKEN_EXPIRY", 168*time.Hour)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Viper unmarshal error: %v, falling back to manual mapping", err)
	}

	if val := viper.GetString("DB_HOST"); val != "" { config.DBHost = val }
	if val := viper.GetString("DB_USER"); val != "" { config.DBUser = val }
	if val := viper.GetString("DB_PASSWORD"); val != "" { config.DBPassword = val }
	if val := viper.GetString("DB_NAME"); val != "" { config.DBName = val }
	if val := viper.GetString("DB_PORT"); val != "" { config.DBPort = val }
	
	if val := viper.GetString("REDIS_HOST"); val != "" { config.RedisHost = val }
	if val := viper.GetString("REDIS_PORT"); val != "" { config.RedisPort = val }
	
	if val := viper.GetString("JWT_SECRET"); val != "" { config.JWTSecret = val }
	if val := viper.GetString("ADMIN_SECRET_KEY"); val != "" { config.AdminSecretKey = val }

	return &config
}
