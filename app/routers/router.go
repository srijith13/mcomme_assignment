package routers

import (
	"fmt"
	"monk-commerce/app/config"
	"monk-commerce/app/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(controller *controllers.CouponController) {
	router := gin.Default()
	discounts := router.Group("/discounts/v1") // Grouping the routes

	// notes.Use(middleware.AppAutherization) If authentication is added

	discounts.POST("/coupons", controller.CreateCoupons)
	discounts.GET("/coupons", controller.GetCoupons)
	discounts.GET("/coupons/:id", controller.GetCoupons)
	discounts.PUT("/coupons", controller.UpdateCoupons)
	discounts.DELETE("/coupons/:id", controller.DeleteCoupons)

	discounts.POST("/applicable-coupons", controller.GetApplicableCoupons)
	discounts.POST("/apply-coupon/:id", controller.ApplyCoupons)

	router.Run(fmt.Sprintf(`:%s`, config.AppPort))
}
