package models

import (
	"database/sql"
	"log"
	"main/config"
)

//useridとproductidを追加
func AddToCart(user_id int, product_id string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("INSERT INTO cart(user_id, product_id) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(user_id, product_id)
	if err != nil {
		log.Println(err)
	}
}

//cartのデータ削除
func DeleteCartItem(user_id int, product_id string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("DELETE FROM cart WHERE user_id = $1 AND product_id = $2")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(user_id, product_id)
	if err != nil {
		log.Println(err)
	}
}
