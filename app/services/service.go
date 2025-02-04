package services

import (
	"database/sql"
	"fmt"
	"log"
	"monk-commerce/app/db"
	"monk-commerce/app/helper"
	"monk-commerce/app/model"
	"strings"
	"time"
)

var dbCon *sql.DB = db.CreateDbConPool()

func CreateCoupons(request *model.CouponsRequest) (interface{}, error) {

	var query string

	// Could be designed better works best with NoSQL DB given the format or accept Json format in SQL column. This is done for ease of use later.
	// This SQL query is prone to SQL injection due to time constrain in designing assuming that wont happen

	switch request.Type {
	case "cart-wise":
		query = fmt.Sprintf(`INSERT INTO coupons (coupon_type,threshold,discount,expiration_date,is_active)
			VALUES ('%s', %d, %d, '%s', %t )`, request.Type, request.Details.Threshold, request.Details.Discount, request.ExpirationDate, true)
	case "product-wise":
		query = fmt.Sprintf(`INSERT INTO coupons (coupon_type,discount,buy_product_id,expiration_date,is_active)
			VALUES ('%s', %d,'{%d}','%s', %t )`, request.Type, request.Details.Discount, request.Details.ProductId, request.ExpirationDate, true)
	case "bxgy":
		buyProdIds, buyPrdQuant := helper.ProdeuctDetails(request.Details.BuyProducts) // Formating the request for buy and get productIds and quantity
		getProdIds, getProdQuant := helper.ProdeuctDetails(request.Details.GetProducts)

		query = fmt.Sprintf(`INSERT INTO coupons (coupon_type,buy_product_id,buy_product_quantity,get_product_id,get_product_quantity,repition_limit,expiration_date,is_active)
		VALUES ('%s', '%v', '%v', '%v', '%v', %d, '%s', %t )`, request.Type, buyProdIds, buyPrdQuant, getProdIds, getProdQuant, request.Details.RepitionLimit, request.ExpirationDate, true)
	}
	_, err := dbCon.Exec(query)

	if err != nil {
		log.Println("Error Executing query:", err)
		return "Coupon Creation Failed", err
	}
	return "Create Coupons Successful", nil
}

func GetCoupons(couponId int64) (interface{}, error) {
	var coupons []model.Coupons

	var query string
	// conditionally changing query based on whether id is passed or not
	if couponId != -1 {
		query = fmt.Sprintf("select * from coupons where coupon_id = %d and is_active = true", couponId)
	} else {
		query = "select * from coupons where is_active = true"

	}
	rows, err := dbCon.Query(query)
	if err != nil {
		log.Println("Error Executing query:", err)
	}
	helper.CouponsDtoMapper(rows, &coupons)

	if len(coupons) == 0 {
		return "Coupon not found", nil
	}
	return coupons, nil

}

// Function to dynamically update database
func UpdateCoupons(request *model.CouponsRequest) (interface{}, error) {
	// CompareData to check which are the keys that have new data to be updated for the given Id
	colVal, colName := helper.CompareData(request)

	// Construct the SQL update query dynamically
	var updateParts []string

	for key, value := range colName {
		updateParts = append(updateParts, fmt.Sprintf("%s = %v", value, colVal[key]))
	}

	query := fmt.Sprintf("update coupons set %s where coupon_id = %d", strings.Join(updateParts, ", "), request.CouponId)

	_, err := dbCon.Exec(query)
	if err != nil {
		log.Println("Error Executing query:", err)
		return "Error", err
	}

	return "success", nil
}

func DeleteCoupons(couponId int64) error {

	// For soft delete
	// query := "update coupons SET is_active = false WHERE coupon_id = $1"

	query := "delete from coupons where coupon_id = $1"

	_, err := dbCon.Exec(query, couponId)
	if err != nil {
		log.Println("Error Executing query:", err)
	}
	return err
}

func GetApplicableCoupons(request *model.CartRequest) (interface{}, error) {
	// Assuming all prices are Integer for ease of calculation
	var coupons []model.Coupons // update structure for finalCoupons
	var finalCoupons model.ApplicableCoupon
	// var ApplicableCoupon
	var prodIds []int64
	price := make(map[int64]int64)
	quantity := make(map[int64]int64)

	var totalPrice int64
	for _, cartItems := range request.Cart.Items {
		prodIds = append(prodIds, cartItems.ProductId)
		quantity[cartItems.ProductId] = cartItems.Quantity
		price[cartItems.ProductId] = cartItems.Price * cartItems.Quantity // total price for each product
		totalPrice += cartItems.Price * cartItems.Quantity
	}

	var proIds []string
	for _, ele := range prodIds {
		proIds = append(proIds, fmt.Sprintf(`%v`, ele))
	}

	query := fmt.Sprintf(`select * from coupons where is_active = true and buy_product_id && ARRAY[%s]::bigint[]  or threshold <= %d `, strings.Join(proIds, ","), totalPrice)

	rows, err := dbCon.Query(query) // get all the coupons based on productIds and check if it matches the threshold (threshold less than or equal to chart total) for chart based coupon
	if err != nil {
		log.Println("Error Executing query:", err)
	}
	helper.CouponsDtoMapper(rows, &coupons)

	if len(coupons) == 0 {
		return "Coupon not found", nil
	}

	log.Println("Starting Calculation of Discount based on Types")

	for _, couponDetails := range coupons {
		expiryDate, _ := time.Parse("2006-01-02", couponDetails.ExpirationDate)
		today, _ := time.Parse("2006-01-02", strings.Split(time.Now().String(), " ")[0])

		// check if coupon is expired if it is expired set it to inactive state (or can delete). Ideally should also be done in /getCoupon api
		if today.Equal(expiryDate) || today.After(expiryDate) {
			query := "UPDATE coupons SET is_active = false WHERE coupon_id = $1"
			_, err := dbCon.Exec(query, couponDetails.CouponId)
			if err != nil {
				log.Println("Error Executing query:", err)
			}
			continue
		}

		var coupon model.ApplicableCouponDetails
		switch couponDetails.CouponType {
		case "cart-wise":
			coupon.CouponId = couponDetails.CouponId
			coupon.Type = couponDetails.CouponType
			discount := (int64(couponDetails.Discount) * totalPrice) / 100 // discount based on total price
			coupon.Discount = discount
			finalCoupons.ApplicableCoupon = append(finalCoupons.ApplicableCoupon, coupon)
		case "product-wise":
			if len(couponDetails.BuyProductId) > 0 && len(couponDetails.GetProductId) == 0 {
				coupon.CouponId = couponDetails.CouponId
				coupon.Type = couponDetails.CouponType
				discount := (int64(couponDetails.Discount) * price[couponDetails.BuyProductId[0]]) / 100 // discount based on particular product total price
				coupon.Discount = discount
				finalCoupons.ApplicableCoupon = append(finalCoupons.ApplicableCoupon, coupon)
			}
		case "bxgy":
			coupon.CouponId = couponDetails.CouponId
			coupon.Type = couponDetails.CouponType

			var buyCount, requiredFreeItems int64
			for _, item := range request.Cart.Items {
				for _, buyItem := range couponDetails.BuyProductId {
					if item.ProductId == buyItem {
						buyCount += item.Quantity
					}
				}
			}

			for _, buyQuant := range couponDetails.BuyProductQuantity {
				for _, getQuant := range couponDetails.GetProductQuantity {
					requiredFreeItems = (buyCount / buyQuant) * getQuant
				}
			}
			if requiredFreeItems > couponDetails.RepitionLimit { //RepetitionLimit
				requiredFreeItems = couponDetails.RepitionLimit
			}

			for _, items := range request.Cart.Items {
				for _, getItem := range couponDetails.GetProductId {
					if items.ProductId == getItem {
						discount := (items.Price * requiredFreeItems) // calcuate product based on the free items
						coupon.Discount = discount
						// requiredFreeItems--
					}
				}
			}
			finalCoupons.ApplicableCoupon = append(finalCoupons.ApplicableCoupon, coupon)
		}
	}
	log.Println("Ending Calculation of Discount based on Types")

	if len(finalCoupons.ApplicableCoupon) == 0 {
		return "Conditions not met", nil
	}
	return finalCoupons, nil
}

func ApplyCoupons(request *model.UpdatedCartRequest) (interface{}, error) {
	// Assuming all prices are Integer for ease of calculation
	var coupons model.Coupons
	var finalBillDetails model.UpdatedCartRequest = *request
	var prodIds []int64
	price := make(map[int64]int64)
	quantity := make(map[int64]int64)
	var totalPrice int64 //finalPrice
	for _, cartItems := range request.Cart.Items {
		prodIds = append(prodIds, cartItems.ProductId)
		quantity[cartItems.ProductId] = cartItems.Quantity
		price[cartItems.ProductId] = cartItems.Price * cartItems.Quantity // total price for each product
		totalPrice += cartItems.Price * cartItems.Quantity
	}

	// assuming only one coupon can be applied

	query := fmt.Sprintf(`select * from coupons where is_active = true and coupon_id = %d `, request.CouponId)

	row := dbCon.QueryRow(query)
	helper.CouponDtoMapper(row, &coupons)

	log.Println("Starting  Calculation of Discount, Final Total Price and discount on each product based on Types")

	switch coupons.CouponType {
	case "cart-wise":
		finalBillDetails.CouponId = coupons.CouponId
		discount := (int64(coupons.Discount) * totalPrice) / 100
		finalBillDetails.TotalPrice = totalPrice
		finalBillDetails.TotalDiscount = discount
		finalBillDetails.FinalPrice = totalPrice - discount
		for val, cartItems := range finalBillDetails.Cart.Items {
			finalBillDetails.Cart.Items[val].TotalDiscount = price[cartItems.ProductId] - ((price[cartItems.ProductId] * discount) / totalPrice) // calculation for dividing the total discount amount among each product
		}
	case "product-wise":
		if len(coupons.BuyProductId) > 0 && len(coupons.GetProductId) == 0 {
			finalBillDetails.CouponId = coupons.CouponId
			discount := (int64(coupons.Discount) * price[coupons.BuyProductId[0]]) / 100
			finalBillDetails.TotalPrice = totalPrice
			finalBillDetails.TotalDiscount = discount
			finalBillDetails.FinalPrice = totalPrice - discount
			finalBillDetails.Cart.Items[0].TotalDiscount = discount
			for val, _ := range finalBillDetails.Cart.Items {
				if finalBillDetails.Cart.Items[val].ProductId == coupons.BuyProductId[0] {
					finalBillDetails.Cart.Items[val].TotalDiscount = discount // calculations for getting total discount for the selected product
				}
			}
		}
	case "bxgy":
		finalBillDetails.CouponId = coupons.CouponId

		var discount, buyCount, requiredFreeItems int64
		for _, item := range request.Cart.Items {
			for _, buyItem := range coupons.BuyProductId {
				if item.ProductId == buyItem {
					buyCount += item.Quantity
				}
			}
		}
		for _, buyQuant := range coupons.BuyProductQuantity {
			for _, getQuant := range coupons.GetProductQuantity {
				requiredFreeItems = (buyCount / buyQuant) * getQuant
			}
		}
		if requiredFreeItems > coupons.RepitionLimit { //RepetitionLimit
			requiredFreeItems = coupons.RepitionLimit
		}

		for val, items := range request.Cart.Items {
			for _, getItem := range coupons.GetProductId {
				if items.ProductId == getItem {
					discount = (items.Price * requiredFreeItems)
					finalBillDetails.Cart.Items[val].TotalDiscount = discount // total discount amount
				}
			}
		}

		finalBillDetails.TotalPrice = totalPrice
		finalBillDetails.TotalDiscount = discount
		finalBillDetails.FinalPrice = totalPrice - discount
	}
	log.Println("Ending  Calculation of Discount, Final Total Price and discount on each product based on Types")

	return finalBillDetails, nil
}
