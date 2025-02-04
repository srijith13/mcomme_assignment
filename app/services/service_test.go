package services

import (
	"monk-commerce/app/model"
	// "monk-commerce/app/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestService_GetApplicableCoupons(t *testing.T) {
	type args struct {
		request *model.CartRequest
	}
	type want struct {
		response model.ApplicableCoupon
		err      error
	}

	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error Failed to create db tp mock: %v ", err)
	}
	defer db.Close()
	query := `select * from coupons where is_active = true and buy_product_id && ARRAY[%s]::bigint[]  or threshold <= %d`

	tests := []struct {
		name string
		args args
		want want
		mock func()
	}{
		{
			name: "Success All Applicable Coupons",
			args: args{
				request: &model.CartRequest{
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId: 1,
								Quantity:  6,
								Price:     50,
							},
							{
								ProductId: 2,
								Quantity:  3,
								Price:     30,
							},
							{
								ProductId: 3,
								Quantity:  2,
								Price:     25,
							},
						},
					},
				},
			},

			mock: func() {
				mockDb.ExpectQuery(query).WithArgs([]int64{1, 2, 3}, 440).WillReturnRows(
					sqlmock.NewRows([]string{"coupon_id", "coupon_type", "threshold", "discount", "buy_product_id", "buy_product_quantity", "get_product_id", "get_product_quantity", "repition_limit", "expiration_date", "is_active"}).
						AddRow(1, "cart-wise", 100, 50, "{}", "{}", "{}", "{}", 0, "2025-02-05", true).
						AddRow(2, "product-wise", 0, 20, "{1}", "{}", "{}", "{}", 0, "2025-02-05", true).
						AddRow(3, "bxgy", 0, 0, "{1,2}", "{2,2}", "{3}", "{1}", 2, "2025-02-05", true).
						AddRow(4, "bxgy", 0, 0, "{1,2}", "{2,2}", "{3}", "{1}", 2, "2025-02-01", true),
				)
			},
			want: want{
				response: model.ApplicableCoupon{
					ApplicableCoupon: []model.ApplicableCouponDetails{
						{
							CouponId: 1,
							Type:     "cart-wise",
							Discount: 220,
						},
						{
							CouponId: 2,
							Type:     "product-wise",
							Discount: 60,
						},
						{
							CouponId: 2,
							Type:     "bxgy",
							Discount: 100,
						},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := GetApplicableCoupons(tt.args.request)

			if tt.want.err != nil {
				assert.Equal(t, tt.want.err, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.response, result)
			}

		})
	}

}

func TestService_ApplyCoupons(t *testing.T) {
	type args struct {
		request *model.CartRequest
	}
	type want struct {
		response model.UpdatedCartRequest
		err      error
	}

	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error Failed to create db tp mock: %v ", err)
	}
	defer db.Close()
	query := `select * from coupons where is_active = true and coupon_id = %d `

	tests := []struct {
		name string
		args args
		want want
		mock func()
	}{
		{
			name: "Success Applied Coupons `cart-wise`",
			args: args{
				request: &model.CartRequest{
					CouponId: 1,
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId: 1,
								Quantity:  6,
								Price:     50,
							},
							{
								ProductId: 2,
								Quantity:  3,
								Price:     30,
							},
							{
								ProductId: 3,
								Quantity:  2,
								Price:     25,
							},
						},
					},
				},
			},

			mock: func() {
				mockDb.ExpectQuery(query).WithArgs(1).WillReturnRows(
					sqlmock.NewRows([]string{"coupon_id", "coupon_type", "threshold", "discount", "buy_product_id", "buy_product_quantity", "get_product_id", "get_product_quantity", "repition_limit", "expiration_date", "is_active"}).
						AddRow(1, "cart-wise", 100, 50, "{}", "{}", "{}", "{}", 0, "2025-02-05", true),
				)
			},
			want: want{
				response: model.UpdatedCartRequest{
					CouponId:      1,
					TotalPrice:    440,
					TotalDiscount: 220,
					FinalPrice:    220,
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId:     1,
								Quantity:      6,
								Price:         50,
								TotalDiscount: 150,
							},
							{
								ProductId:     2,
								Quantity:      3,
								Price:         30,
								TotalDiscount: 45,
							},
							{
								ProductId:     3,
								Quantity:      2,
								Price:         25,
								TotalDiscount: 25,
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "Success Applied Coupons `product-wise`",
			args: args{
				request: &model.CartRequest{
					CouponId: 2,
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId: 1,
								Quantity:  6,
								Price:     50,
							},
							{
								ProductId: 2,
								Quantity:  3,
								Price:     30,
							},
							{
								ProductId: 3,
								Quantity:  2,
								Price:     25,
							},
						},
					},
				},
			},

			mock: func() {
				mockDb.ExpectQuery(query).WithArgs(2).WillReturnRows(
					sqlmock.NewRows([]string{"coupon_id", "coupon_type", "threshold", "discount", "buy_product_id", "buy_product_quantity", "get_product_id", "get_product_quantity", "repition_limit", "expiration_date", "is_active"}).
						AddRow(2, "product-wise", 0, 20, "{1}", "{}", "{}", "{}", 0, "2025-02-05", true),
				)
			},
			want: want{
				response: model.UpdatedCartRequest{
					CouponId:      1,
					TotalPrice:    440,
					TotalDiscount: 60,
					FinalPrice:    380,
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId:     1,
								Quantity:      6,
								Price:         50,
								TotalDiscount: 60,
							},
							{
								ProductId:     2,
								Quantity:      3,
								Price:         30,
								TotalDiscount: 0,
							},
							{
								ProductId:     3,
								Quantity:      2,
								Price:         25,
								TotalDiscount: 0,
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "Success Applied Coupons `cart-wise`",
			args: args{
				request: &model.CartRequest{
					CouponId: 3,
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId: 1,
								Quantity:  6,
								Price:     50,
							},
							{
								ProductId: 2,
								Quantity:  3,
								Price:     30,
							},
							{
								ProductId: 3,
								Quantity:  2,
								Price:     25,
							},
						},
					},
				},
			},

			mock: func() {
				mockDb.ExpectQuery(query).WithArgs(3).WillReturnRows(
					sqlmock.NewRows([]string{"coupon_id", "coupon_type", "threshold", "discount", "buy_product_id", "buy_product_quantity", "get_product_id", "get_product_quantity", "repition_limit", "expiration_date", "is_active"}).
						AddRow(3, "bxgy", 0, 0, "{1,2}", "{2,2}", "{3}", "{1}", 2, "2025-02-05", true),
				)
			},
			want: want{
				response: model.UpdatedCartRequest{
					CouponId:      1,
					TotalPrice:    440,
					TotalDiscount: 50,
					FinalPrice:    390,
					Cart: model.CartItems{
						Items: []model.ProdDetails{
							{
								ProductId:     1,
								Quantity:      6,
								Price:         50,
								TotalDiscount: 150,
							},
							{
								ProductId:     2,
								Quantity:      3,
								Price:         30,
								TotalDiscount: 45,
							},
							{
								ProductId:     3,
								Quantity:      2,
								Price:         25,
								TotalDiscount: 50,
							},
						},
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := GetApplicableCoupons(tt.args.request)

			if tt.want.err != nil {
				assert.Equal(t, tt.want.err, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.response, result)
			}

		})
	}

}
