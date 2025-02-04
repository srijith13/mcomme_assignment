package model

type Coupons struct {
	CouponId           int64   `json:"coupon_id"`
	CouponType         string  `json:"coupon_type"`
	Threshold          int64   `json:"threshold"`
	Discount           int64   `json:"discount"`
	BuyProductId       []int64 `json:"buy_product_id"`
	BuyProductQuantity []int64 `json:"buy_product_quantity"`
	GetProductId       []int64 `json:"get_product_id"`
	GetProductQuantity []int64 `json:"get_product_quantity"`
	RepitionLimit      int64   `json:"repition_limit"`
	ExpirationDate     string  `json:"expiration_date"`
	IsActive           bool    `json:"is_active"`
}
