package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	RedisAddr  string
	JWTSecret  string 
	AdminToken string 
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, продолжаем с переменными окружения")
	}

	return &Config{
		ServerAddr: getEnv("SERVER_ADDR"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBName:     getEnv("DB_NAME"),
		RedisAddr:  getEnv("REDIS_ADDR"),
		JWTSecret:  getEnv("JWT_SECRET"),
		AdminToken: getEnv("ADMIN_TOKEN"),
	}
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Переменная окружения %s не установлена", key)
	return ""
}
