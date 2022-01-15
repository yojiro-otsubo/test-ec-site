package models

import (
	"database/sql"
	"fmt"
	"log"
	"main/config"

	_ "github.com/lib/pq"
)

// 初期設定

var DbConnection *sql.DB

func ConnectionInfo() string {
	info := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Config.DBhost, config.Config.DBuser, config.Config.DBpassword, config.Config.DBname)
	return info
}

func ConnectionDB() {
	DbConnection, err := sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	//usersテーブル作成
	cmd := "CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, username VARCHAR(50), password VARCHAR(255), email VARCHAR(255));"
	_, err = DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	//accountsテーブル作成
	cmd1 := "CREATE TABLE IF NOT EXISTS accounts (id serial PRIMARY KEY, user_id INT, stripe_account VARCHAR(255));"
	_, err = DbConnection.Exec(cmd1)
	if err != nil {
		log.Fatalln(err)
	}

	//productsテーブル作成
	cmd2 := "CREATE TABLE IF NOT EXISTS products (id serial PRIMARY KEY, user_id INT, stripe_product_id VARCHAR(255), stripe_price_id VARCHAR(255), item_name VARCHAR(255), description VARCHAR(1000), amount INT, sold_out INT);"
	_, err = DbConnection.Exec(cmd2)
	if err != nil {
		log.Fatalln(err)
	}

	//payment_history
	cmd3 := "CREATE TABLE IF NOT EXISTS payment_history (id serial PRIMARY KEY, user_id INT, product_id INT, transfer_group VARCHAR(50));"
	_, err = DbConnection.Exec(cmd3)
	if err != nil {
		log.Fatalln(err)
	}

	cmd4 := "CREATE TABLE IF NOT EXISTS shipping (id serial PRIMARY KEY, product_id INT, shipping INT, arrives INT);"
	_, err = DbConnection.Exec(cmd4)
	if err != nil {
		log.Fatalln(err)
	}

	cmd5 := "CREATE TABLE IF NOT EXISTS cart (id serial PRIMARY KEY, user_id INT, product_id INT);"
	_, err = DbConnection.Exec(cmd5)
	if err != nil {
		log.Fatalln(err)
	}

	//testテーブル作成
	cmdtest := "CREATE TABLE IF NOT EXISTS test_db (id serial PRIMARY KEY, test VARCHAR(50));"
	_, err = DbConnection.Exec(cmdtest)
	if err != nil {
		log.Fatalln(err)
	}

	defer DbConnection.Close()
}
