package main

import (
	"log"
	"monk-commerce/app/config"
	"monk-commerce/app/routers"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	log.Println("Application Initiated")
	loc, _ := time.LoadLocation(config.GinTZ)
	time.Local = loc
	gin.SetMode(config.GinMode)
	routers.InitRoutes()
}
