package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// user:password@tcp(host:port)/dbname?parseTime=true
func (c Config) DatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func InitConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatalf("DB_HOST is required")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("DB_PORT must be a number")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatalf("DB_USER is required")
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("DB_NAME is required")
	}

	cfg := Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
	}

	return cfg
}
