package config

import (
	"os"

	"github.com/joho/godotenv"
)

var GinMode string
var GinTZ string
var AppPort string

// var StageApiLimit int
// var LogPath string
var DbHost string
var DbPort string
var DbUser string
var DbPass string

func init() {
	godotenv.Load()

	GinMode = os.Getenv("GIN_MODE") // Gin mode to run the application
	GinTZ = os.Getenv("GIN_TZ")     // Setting Time zone
	AppPort = os.Getenv("APP_PORT") // Port the application should run

	// StageApiLimit, _ = strconv.Atoi(os.Getenv("API_LIMIT"))
	// LogPath = os.Getenv("LOG_PATH")

	// The Database connection details
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPass = os.Getenv("DB_PASS")

}
