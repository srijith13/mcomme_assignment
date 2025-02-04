package model

// For Coupon CRUD  Operations
type CouponsRequest struct {
	CouponId       int64         `json:"coupon_id"`
	Type           string        `json:"type"`
	Details        CouponDetails `json:"details"`
	ExpirationDate string        `json:"expiration_date"`
	IsActive       bool          `json:"is_active"`
}

type CouponDetails struct {
	ProductId     int64                 `json:"product_id"`
	Threshold     int8                  `json:"threshold"`
	Discount      int8                  `json:"discount"`
	BuyProducts   []BxgyProductQuantity `json:"buy_products"`
	GetProducts   []BxgyProductQuantity `json:"get_products"`
	RepitionLimit int8                  `json:"repition_limit"`
}

type BxgyProductQuantity struct {
	ProductId int64 `json:"product_id"`
	Quantity  int8  `json:"quantity"`
}

// For Coupon Application
type CartRequest struct {
	Cart     CartItems `json:"cart"`
	CouponId int64     `json:"coupon_id"`
}

type UpdatedCartRequest struct {
	Cart          CartItems `json:"cart"`
	CouponId      int64     `json:"coupon_id"`
	TotalPrice    int64     `json:"total_price"`
	TotalDiscount int64     `json:"total_discount"`
	FinalPrice    int64     `json:"final_price"`
}

type CartItems struct {
	Items []ProdDetails `json:"items"`
}

type ProdDetails struct {
	ProductId     int64 `json:"product_id"`
	Quantity      int64 `json:"quantity"`
	Price         int64 `json:"price"`
	TotalDiscount int64 `json:"total_discount"`
}

// For  ApplicableCoupon
type ApplicableCoupon struct {
	ApplicableCoupon []ApplicableCouponDetails `json:"applicable_coupons"`
}

type ApplicableCouponDetails struct {
	CouponId int64  `json:"coupon_id"`
	Type     string `json:"type"`
	Discount int64  `json:"discount"`
}
