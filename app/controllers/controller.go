package controllers

import (
	"log"
	"monk-commerce/app/model"
	"monk-commerce/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCoupons(c *gin.Context) {
	var request model.CouponsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error Binding Json", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err.Error()))
		return
	}
	// // Custom validator to check the request body matches the required format and data structure to prevent injections through json

	// if errors := helper.ValidateCustomBody(request); len(errors) != 0 {
	// 	c.JSON(http.StatusBadRequest, model.BuildErrorResponse("Bad Request", errors, nil))
	// 	return
	// }

	result, err := services.CreateCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func GetCoupons(c *gin.Context) {
	id := c.Param("id")
	var couponId int64 = -1
	var err error
	if id != "" {
		couponId, err = strconv.ParseInt(id, 10, 64)
	}
	if err != nil {
		log.Println("Error Parsing int to string", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err))
	}

	result, err := services.GetCoupons(couponId)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func UpdateCoupons(c *gin.Context) {
	var request model.CouponsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error Binding Json", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err.Error()))
		return
	}

	result, err := services.UpdateCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func DeleteCoupons(c *gin.Context) {
	id := c.Param("id")
	couponId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("Error Parsing int to string", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err))
	}

	err = services.DeleteCoupons(couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", "Success", err))
	}
}

func GetApplicableCoupons(c *gin.Context) {
	var request model.CartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error Binding Json", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err.Error()))
		return
	}

	result, err := services.GetApplicableCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func ApplyCoupons(c *gin.Context) {
	id := c.Param("id")
	couponId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("Error Parsing int to string", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err))
	}
	var request model.UpdatedCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err.Error()))
		return
	}
	request.CouponId = couponId
	result, err := services.ApplyCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}
