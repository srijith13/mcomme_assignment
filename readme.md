# **README - Coupon Management System**

## **General Overview**

This API is designed to manage and apply various types of discount coupons for an e-commerce platform. The system supports:

- **Cart-wise coupons**: Apply discounts based on the total cart amount.
- **Product-wise coupons**: Apply discounts to specific products.
- **BxGy (Buy X, Get Y) coupons**: Offer promotions where the customer gets free items upon purchasing a certain number of items from a selected set.
- **BxQn (Buy X, of Quantity n) coupons**: Offer promotions where the customer gets a discount on the total price upon purchasing a certain number of items from a selected set. (Eg. Product x costs 100. If you purchase 2 quantities of product x then you get a discount of 10% else you won't get.). _\[Not Implemented for applicable_coupons, apply-coupon\]_
- **Based on Payment Method:** Subcategory of Cart-wise where based on payment method a discount on the entire cart is applicable. \[_Not Implemented_ \]

### Assumptions

- The coupon applies only once per checkout process.
- The Cart total is calculated before the application of any coupon then after the application of the coupon, the cart total is recalculated.
- **Coupon Expiration**: In the current implementation, the expiration is checked only based on date and not time. The format for _expiration_date_ is **yyyy-mm-dd**
- **Threshold Boundaries**: If the cart's total is exactly the threshold value, currently the _coupon is still applicable._
- **No dynamic thresholds**: The system assumes static threshold values for cart-wise coupons, dynamic thresholds based on cart contents are not supported.
- **For Product-Wise**: The coupon applies only if the product is present in the cart.
- **For Product-Wise Coupon**: Product ID is the unique identifier to that product and each product will have a separate coupon.
- **Multiple Items of Same Product**: If the cart has multiple items of the same product, then all the items the discounts will be implemented
- **Repetition limits for BxGy**: The repetition limit for BxGy coupons applies only to the "buy" items in the cart; if the user has already received the maximum number of free items, no more free items will be granted
- **For BxGy**: The coupon applies only if the cart contains both "buy" and "get" items.
- **Coupon validity**: All coupons are assumed to be valid and the check of expiration is done only during the _applicable_coupons api call_ and then the status is updated to inactive.
- Cannot update **Coupon Type**

### Edge Cases

- **Threshold Boundaries**: If the cart's total is exactly the threshold value, should the coupon apply?
- **Multiple Items of Same Product**: If the cart has multiple items of the same product, then all the items the discounts will be implemented
- **Insufficient Products**: If the cart doesn't meet the "buy" condition (e.g., fewer than 2 items from \[X, Y, Z\]), the coupon shouldn’t apply.
- **Insufficient Free Products**: If the cart doesn't meet the required "get" product quantity even if the "buy" condition is satisfied
- **Exceeding Free Products**: If the user buys more than the required "buy" products but doesn’t have enough products in the "get" list, the system applies only available free products until they reach the repetition limit, and the remaining will not be covered under any discount

### Limitations

- **No coupon stacking**: Cart-wise, product-wise, and BxGy coupons cannot be stacked in a single checkout and are assumed to be the same way in the process.
- **Insufficient Free Products For BxGy**: If the cart doesn't meet the required "get" product quantity even if the "buy" condition is satisfied _this is not implemented_

## **Current System Overview**

### Assumptions

- **Simplified calculation**: The discount calculations are performed in a simplified manner using integer type for all prices, and discounts.

### Limitations

- Since I use SQL based db it is limited by the dynamic changes in conditions and types of coupons that can be used. Solution: Use of NoSQL based db would ideally give much more flexibility.
- Condition to check if more than one coupon is applied (assumed only one can be applied)
- Can only add new coupons which fit in the current db set up implying only sub categories of the existing coupons will work best. For new type of coupon it is best to have an NoSQL based DB

### Improvements

- Complex methods are used because of the way the database table is structured this could be avoided by using a mix of SQL and NoSql dbs or fully NoSQL db.
- A cache system like Redis can be used to keep the cart details so that if the user changes the coupon multiple times we dont have to update the db (A db for storing the cart/bill details ) which will be updated to the db only after the payment is successful.
- **Coupon validity**: The check of expiration should be done also during the _get_coupons_ api call.
- Expiration of coupon should be based on both date and time there by giving more flexibility.
- Addition of Docker for easy deployment and have standardized units.
- Test cases are currently written for two functions that had the most business logic and only success flows are written. Need to write test cases for failures, null conditions and possible edge cases.

## **Database Schema:**
```sh
-- Table: public.coupons
-- DROP TABLE IF EXISTS public.coupons;
CREATE TABLE IF NOT EXISTS public.coupons
(
    coupon_id bigint NOT NULL DEFAULT nextval('coupons_coupon_id_seq'::regclass),
    coupon_type text COLLATE pg_catalog."default" NOT NULL,
    threshold bigint DEFAULT 0,
    discount bigint DEFAULT 0,
    buy_product_id bigint[] DEFAULT '{}'::bigint[],
    buy_product_quantity integer[] DEFAULT '{}'::integer[],
    get_product_id bigint[] DEFAULT '{}'::bigint[],
    get_product_quantity integer[] DEFAULT '{}'::integer[],
    repition_limit integer DEFAULT 0,
    expiration_date character varying(20) COLLATE pg_catalog."default",
    is_active boolean DEFAULT false,
    CONSTRAINT coupon_table_pkey PRIMARY KEY (coupon_id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.coupons
    OWNER to postgres;
```
## **API Structures:**

To start application ``` go run start.go ```
Please populate the ```.env``` file details before starting the application (db details)

END POINT  \[ <http://localhost:8080/discounts/v1/> \]

1. **POST** /coupons
```javascript
{
    // "type":  "cart-wise" ,
    // "details":{
    //     "threshold":100,
    //     "discount":50   
    // },
    // "type":  "product-wise" , 
    // "details": {
    //     "product_id": 1,
    //     "discount": 20
    // } ,
    "type": "bxgy",
    "details": {
        "buy_products": [
            {
                "product_id": 1,
                "quantity": 3
            },
            {
                "product_id": 2,
                "quantity": 3
            }
        ],
        "get_products": [
            {
                "product_id": 3,
                "quantity": 1
            }
        ],
        "repition_limit": 2
    },
    "expiration_date": "2025-02-05"
}
```
2. **GET** /coupons or /coupons /{id}
3. **PUT** /coupons /{id}
```javascript
{
    "coupon_id":3,
    //  "type":  "cart-wise" ,
    // "details":{
    //     "discount":25   
    // },
    // "type":  "product-wise" , 
    // "details": {
    //     "product_id": 1,
    //     "discount": 20
    // } ,
    "type": "bxgy",
    "details": {
        "buy_products": [
            {
                "product_id": 1,
                "quantity": 2
            },
            {
                "product_id": 2,
                "quantity": 2
            }
        ]
    //     "get_products": [
    //         {
    //             "product_id": 3,
    //             "quantity": 1
    //         }
    //     ],
        // "repition_limit": 2
    }
    // "expiration_date": "2025-02-05"
}
```
1. **DELETE** /coupons /{id}
2. **POST** /applicable-coupons
```javascript
{
    "cart": {
        "items": [
            {
                "product_id": 1,
                "quantity": 6,
                "price": 50
            }, // Product X 
            {
                "product_id": 2,
                "quantity": 3,
                "price": 30
            }, // Product Y 
            {
                "product_id": 3,
                "quantity": 2,
                "price": 25
            } // Product Z 
        ]
    }
}
```

1. **POST** /apply-coupon/{id}
```javascript
{
    "cart": {
        "items": [
            {
                "product_id": 1,
                "quantity": 6,
                "price": 50
            }, // Product X 
            {
                "product_id": 2,
                "quantity": 3,
                "price": 30
            }, // Product Y 
            {
                "product_id": 3,
                "quantity": 2,
                "price": 25
            } // Product Z 
        ]
    }
}
```
