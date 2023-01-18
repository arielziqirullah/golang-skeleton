package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AppDebug() {
	ginMode := os.Getenv("APP_DEBUG")
	if ginMode == "false" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func AppLoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}
