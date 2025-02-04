package routers

import (
	"fmt"
	"monk-commerce/app/config"
	"monk-commerce/app/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()
	discounts := router.Group("/discounts/v1") // Grouping the routes

	// notes.Use(middleware.AppAutherization) If authentication is added

	discounts.POST("/coupons", controllers.CreateCoupons)
	discounts.GET("/coupons", controllers.GetCoupons)
	discounts.GET("/coupons/:id", controllers.GetCoupons)
	discounts.PUT("/coupons", controllers.UpdateCoupons)
	discounts.DELETE("/coupons/:id", controllers.DeleteCoupons)

	discounts.POST("/applicable-coupons", controllers.GetApplicableCoupons)
	discounts.POST("/apply-coupon/:id", controllers.ApplyCoupons)

	router.Run(fmt.Sprintf(`:%s`, config.AppPort))
}
