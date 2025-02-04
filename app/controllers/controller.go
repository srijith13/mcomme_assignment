package controllers

import (
	"log"
	"monk-commerce/app/model"
	"monk-commerce/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CouponController struct {
	couponService services.IService
}

func NewCouponController(couponService services.IService) *CouponController {
	return &CouponController{couponService: couponService}
}

func (h *CouponController) CreateCoupons(c *gin.Context) {
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

	result, err := h.couponService.CreateCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func (h *CouponController) GetCoupons(c *gin.Context) {
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

	result, err := h.couponService.GetCoupons(couponId)

	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func (h *CouponController) UpdateCoupons(c *gin.Context) {
	var request model.CouponsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error Binding Json", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err.Error()))
		return
	}

	result, err := h.couponService.UpdateCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func (h *CouponController) DeleteCoupons(c *gin.Context) {
	id := c.Param("id")
	couponId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("Error Parsing int to string", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err))
	}

	err = h.couponService.DeleteCoupons(couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", "Success", err))
	}
}

func (h *CouponController) GetApplicableCoupons(c *gin.Context) {
	var request model.CartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error Binding Json", err)
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", nil, err.Error()))
		return
	}

	result, err := h.couponService.GetApplicableCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}

func (h *CouponController) ApplyCoupons(c *gin.Context) {
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
	result, err := h.couponService.ApplyCoupons(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildResponse("Bad Request", result, err))
	} else {
		c.JSON(http.StatusOK, model.BuildResponse("Message", result, err))
	}
}
