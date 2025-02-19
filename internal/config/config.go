package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Congig struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	JWTSecretKey        string
	DBConnectionRetries int
	DBConnectionDelay   int
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}
}

func New() *Congig {
	return &Congig{
		DBHost:              getEnvStr("DB_HOST", ""),
		DBPort:              getEnvStr("DB_PORT", ""),
		DBUser:              getEnvStr("DB_USER", ""),
		DBPassword:          getEnvStr("DB_PASSWORD", ""),
		DBName:              getEnvStr("DB_NAME", ""),
		JWTSecretKey:        getEnvStr("JWT_SECRET_KEY", ""),
		DBConnectionRetries: getEnvInt("DB_CONNECTION_RETRIES", 0),
		DBConnectionDelay:   getEnvInt("DB_CONNECTION_DELAY", 0),
	}
}

func getEnvStr(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Cannot convert %s to int, using default: %d\n", key, defaultVal)
	}
	return defaultVal
}
