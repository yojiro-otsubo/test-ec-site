package models

import (
	"database/sql"
	"log"
	"main/config"
)

//動作テスト
func TestDb() {
	DbConnection, err := sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	var id int

	err = DbConnection.QueryRow("INSERT INTO test_db(test) VALUES($1) RETURNING id", "hoge").Scan(&id)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("PRIMARY KEY = ", id)
}
