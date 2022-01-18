package models

import (
	"database/sql"
	"log"
	"main/config"
)

func InsertSipping(productid string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("INSERT INTO delivery_status(product_id, shipping) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(productid, "1")
	if err != nil {
		log.Println(err)
	}

}

func CheckDeliveryStatusProductId(productid string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var id string
	err = DbConnection.QueryRow("SELECT product_id FROM delivery_status WHERE product_id = $1", productid).Scan(&id)
	if err != nil {
		return "なし"
	} else {
		return "あり"
	}
}
