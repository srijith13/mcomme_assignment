package main

import (
	"log"
	"monk-commerce/app/config"
	"monk-commerce/app/controllers"
	"monk-commerce/app/db"
	"monk-commerce/app/routers"
	"monk-commerce/app/services"
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
	db, err := db.InitDb()
	if err != nil {
		log.Printf("Error Initializing DB: %v \n", err)
	}

	couponService := services.NewCouponService(db)
	couponController := controllers.NewCouponController(couponService)

	routers.InitRoutes(couponController)
}
