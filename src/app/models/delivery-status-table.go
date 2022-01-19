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

func InsertArrives(productid string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("INSERT INTO delivery_status(product_id, arrives) VALUES($1, $2)")
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

func CheckSipping(productid string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var arrives string
	err = DbConnection.QueryRow("SELECT shipping FROM delivery_status WHERE product_id = $1", productid).Scan(&arrives)
	if err != nil {
		log.Println("CheckArrives", err)
		return "0"
	}
	return arrives
}

func CheckArrives(productid string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var arrives string
	err = DbConnection.QueryRow("SELECT arrives FROM delivery_status WHERE product_id = $1", productid).Scan(&arrives)
	if err != nil {
		log.Println("CheckArrives", err)
		return "0"
	}
	return arrives
}

func UpdateArrives(productid string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	var arrives int
	err = DbConnection.QueryRow("UPDATE delivery_status SET arrives = $2 WHERE product_id = $1 RETURNING arrives", productid, "1").Scan(&arrives)
	if err != nil {
		log.Println(err)
	}

	log.Println("arrives = ", arrives)
}
