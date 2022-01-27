package models

import (
	"database/sql"
	"log"
	"main/config"
	"math"
	"strconv"
)

func RegistProduct(pkid int, productid, priceid string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	var id, user_id, amount int
	var stripe_product_id, stripe_price_id, item_name, description string
	err = DbConnection.QueryRow("UPDATE products SET stripe_product_id = $2, stripe_price_id = $3 WHERE id = $1 RETURNING id, user_id, stripe_product_id, stripe_price_id, item_name, description, amount", pkid, productid, priceid).Scan(&id, &user_id, &stripe_product_id, &stripe_price_id, &item_name, &description, &amount)
	if err != nil {
		log.Println(err)
	}
	log.Println("id = ", id, "\nuser_id = ", user_id, "\nstripe_product_id = ", stripe_product_id, "\nstripe_price_id = ", stripe_price_id)
	log.Println("item_name = ", item_name, "\ndescription = ", description, "\namount = ", amount)

}

func RegistUserIdAndGetProductId(userid, amount int, item_name, description string) int {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	var id int
	//temporary := "Temporary"

	err = DbConnection.QueryRow("INSERT INTO products(user_id, item_name, description, amount, sold_out) VALUES($1, $2, $3, $4, $5) RETURNING id", userid, item_name, description, amount, "0").Scan(&id)
	if err != nil {
		log.Println(err)
	}
	log.Println("PRIMARY KEY = ", id)

	return id

}

type Product struct {
	Id, UserId, StripeProductId, StripePriceId, ItemName, Description, Amount, SoldOut string
}

func GetTheProductOfUserId(userid int) []Product {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE user_id = $1", userid)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" && p.SoldOut == "0" {
			productResult = append(productResult, p)
		}
	}
	//log.Println(productResult)

	return productResult

}

func GetAllProductOfUserId(userid string) []Product {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE user_id = $1", userid)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" {
			productResult = append(productResult, p)
		}
	}
	//log.Println(productResult)

	return productResult

}

func GetSoldOutProductOfUserId(userid int) []Product {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE user_id = $1", userid)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" && p.SoldOut == "1" {
			if CheckDeliveryStatusProductId(p.Id) == "なし" {
				productResult = append(productResult, p)
			}
		}
	}
	//log.Println(productResult)

	return productResult

}

func GetSippingOkProductOfUserId(userid int) []ProductArrival {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE user_id = $1", userid)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []ProductArrival
	for rows.Next() {
		var p ProductArrival
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" && p.SoldOut == "1" {
			if CheckDeliveryStatusProductId(p.Id) == "あり" {
				if CheckArrives(p.Id) == "1" {
					p.Arrival = "1"
					productResult = append(productResult, p)
				} else {
					productResult = append(productResult, p)
				}
			}
		}
	}
	//log.Println(productResult)

	return productResult

}

func GetProductTop() []Product {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products LIMIT 21")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" {
			productResult = append(productResult, p)
		}
	}
	//log.Println(productResult)

	return productResult

}

func GetProduct(product_id string) [8]string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}
	var Id, UserId, StripeProductId, StripePriceId, ItemName, Description, Amount, SoldOut string

	err = DbConnection.QueryRow("SELECT * FROM products WHERE id = $1", product_id).Scan(&Id, &UserId, &StripeProductId, &StripePriceId, &ItemName, &Description, &Amount, &SoldOut)
	if err != nil {
		log.Println(err)
	}

	arr := [...]string{Id, UserId, StripeProductId, StripePriceId, ItemName, Description, Amount, SoldOut}
	return arr
}

//useridでcartテーブルから取得したproductidでproductテーブルの全て取得
func GetProductFromCartDB(user_id interface{}, tax float64) []Product {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE id IN (SELECT product_id FROM cart WHERE user_id = $1)", user_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" && p.SoldOut == "0" {
			f, err := strconv.ParseFloat(p.Amount, 64)
			if err != nil {
				log.Println(err)
			}
			f = f * tax
			taxamount := int(math.Round(f))
			p.Amount = strconv.Itoa(taxamount)
			productResult = append(productResult, p)
		}
	}
	//log.Println(productResult)

	return productResult
}

type ProductArrival struct {
	Id, UserId, StripeProductId, StripePriceId, ItemName, Description, Amount, SoldOut, Arrival string
}

func GetProductIdFromPaymentHistory(user_id interface{}) []ProductArrival {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE id IN (SELECT product_id FROM payment_history WHERE user_id = $1)", user_id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []ProductArrival
	for rows.Next() {
		var p ProductArrival
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" && p.SoldOut == "1" {
			if CheckArrives(p.Id) == "1" {
				p.Arrival = "1"
				log.Println(p.Arrival)
				productResult = append(productResult, p)
			} else {
				productResult = append(productResult, p)
			}
		}
	}
	//log.Println(productResult)

	return productResult
}

func UpdataSoldOutValue(productid, soldout string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}
	var soldOutValue int
	err = DbConnection.QueryRow("UPDATE products SET sold_out = $2 WHERE id = $1 RETURNING sold_out", productid, soldout).Scan(&soldOutValue)
	if err != nil {
		log.Println(err)
	}

	log.Println("UpdataSoldOutValue = ", soldOutValue)
}

func GetProductPurchaseConfirmation(productid []int) []Product {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM products WHERE id IN ($1)", productid)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var productResult []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.UserId, &p.StripeProductId, &p.StripePriceId, &p.ItemName, &p.Description, &p.Amount, &p.SoldOut)
		if err != nil {
			log.Println(err)
		}
		if p.Id != "" && p.UserId != "" && p.StripeProductId != "" && p.StripePriceId != "" && p.ItemName != "" && p.Description != "" && p.Amount != "" && p.SoldOut == "0" {

			productResult = append(productResult, p)
		}
	}
	//log.Println(productResult)

	return productResult
}

func DeleteProduct(userid int, productid string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("DELETE FROM products WHERE user_id = $1 AND id = $2")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(userid, productid)
	if err != nil {
		log.Println(err)
	}
}

func CountProduct(user_id int) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var count string
	err = DbConnection.QueryRow("SELECT COUNT( * ) FROM products WHERE user_id = $1", user_id).Scan(&count)
	if err != nil {
		return "0"
	} else {
		log.Println("count = ", count)
		return count
	}
}
