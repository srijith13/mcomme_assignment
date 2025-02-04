package helper

import (
	"database/sql"
	"log"
	"monk-commerce/app/model"
	"reflect"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

// Creating string format for the product details
func ProdeuctDetails(products []model.BxgyProductQuantity) (prodIds string, prodQuant string) {
	prodIds = "{"
	prodQuant = "{"
	var tempPrdId, tempProdQuant []string
	for _, val := range products {
		tempPrdId = append(tempPrdId, strconv.Itoa(int(val.ProductId)))
		tempProdQuant = append(tempProdQuant, strconv.Itoa(int(val.Quantity)))
	}
	prodIds += strings.Join(tempPrdId, ",")
	prodQuant += strings.Join(tempProdQuant, ",")
	prodIds += "}"
	prodQuant += "}"

	return prodIds, prodQuant
}

// Function to Compare data to check which fields need to be updated
func CompareData(request *model.CouponsRequest) ([]interface{}, []string) {

	var buyProdIds, buyQuantity, getProdIds, getQuantity []int64
	if request.Type == "product-wise" && request.Details.ProductId > 0 {
		buyProdIds = append(buyProdIds, request.Details.ProductId)
	} else if request.Type == "bxgy" {
		for _, cartItems := range request.Details.BuyProducts {
			buyProdIds = append(buyProdIds, cartItems.ProductId)
			buyQuantity = append(buyProdIds, int64(cartItems.Quantity))
		}
		for _, cartItems := range request.Details.GetProducts {
			getProdIds = append(buyProdIds, cartItems.ProductId)
			getQuantity = append(buyProdIds, int64(cartItems.Quantity))
		}
	}

	var updateDetails model.Coupons
	// updateDetails.CouponType = request.Type
	updateDetails.Threshold = int64(request.Details.Threshold)
	updateDetails.Discount = int64(request.Details.Discount)
	updateDetails.BuyProductId = buyProdIds
	updateDetails.BuyProductQuantity = buyQuantity
	updateDetails.GetProductId = getProdIds
	updateDetails.GetProductQuantity = getQuantity
	updateDetails.RepitionLimit = int64(request.Details.RepitionLimit)
	updateDetails.ExpirationDate = request.ExpirationDate
	updateDetails.IsActive = request.IsActive

	val := reflect.ValueOf(updateDetails)

	if val.Kind() != reflect.Struct {
		log.Println("Provided value is not a struct")
	}

	typ := val.Type()
	var columnName []string
	var columnValue []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			if fieldType.Type.Kind() == reflect.String {
				columnValue = append(columnValue, "'"+field.Interface().(string)+"'")
			} else if fieldType.Type.Kind() == reflect.Slice && fieldType.Type == reflect.TypeOf([]int64{}) {
				columnValue = append(columnValue, IntArrayConverter(field.Interface().([]int64)))
			} else {
				columnValue = append(columnValue, field.Interface())
			}
			columnName = append(columnName, jsonTag)
		}
	}
	return columnValue, columnName
}

func IntArrayConverter(array []int64) string {
	arrayString := "'{"
	var tempArray []string
	for _, val := range array {
		tempArray = append(tempArray, strconv.Itoa(int(val)))
	}
	arrayString += strings.Join(tempArray, ",")
	arrayString += "}'"

	return arrayString
}

// Mapping the db queried details
func CouponsDtoMapper(rows *sql.Rows, coupons *[]model.Coupons) {
	log.Println("Scanning SQl rows to map")
	for rows.Next() {
		couponsDto := model.Coupons{}
		err := rows.Scan(&couponsDto.CouponId,
			&couponsDto.CouponType,
			&couponsDto.Threshold,
			&couponsDto.Discount,
			pq.Array(&couponsDto.BuyProductId),
			pq.Array(&couponsDto.BuyProductQuantity),
			pq.Array(&couponsDto.GetProductId),
			pq.Array(&couponsDto.GetProductQuantity),
			&couponsDto.RepitionLimit,
			&couponsDto.ExpirationDate,
			&couponsDto.IsActive)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		*coupons = append(*coupons, couponsDto)
	}
}

func CouponDtoMapper(row *sql.Row, coupons *model.Coupons) {
	log.Println("Scanning SQl row to map")
	err := row.Scan(&coupons.CouponId,
		&coupons.CouponType,
		&coupons.Threshold,
		&coupons.Discount,
		pq.Array(&coupons.BuyProductId),
		pq.Array(&coupons.BuyProductQuantity),
		pq.Array(&coupons.GetProductId),
		pq.Array(&coupons.GetProductQuantity),
		&coupons.RepitionLimit,
		&coupons.ExpirationDate,
		&coupons.IsActive)
	if err != nil {
		log.Println("Error scanning row:", err)

	}
}

// For validating the request json for correct formats
// func ValidateRequestBody(request dto.RouteRequest, params string) []string {
// }
