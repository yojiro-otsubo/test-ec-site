package models

import (
	"database/sql"
	"log"
	"main/config"
)

func InsertInquiry(name, email, inquiry string, date string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}
	//userにINSERTする
	var id, returndate string
	err = DbConnection.QueryRow("INSERT INTO inquiry(name, email, inquiry, date) VALUES($1, $2, $3, $4) RETURNING id, date", name, email, inquiry, date).Scan(&id, &returndate)
	if err != nil {
		log.Println(err)
	}
	log.Println("returndate", returndate)
	return id
}
func GetInquiry(inquiry_id interface{}) (string, bool) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var date string
	var id string
	err = DbConnection.QueryRow("SELECT id, date FROM inquiry WHERE id = $1", inquiry_id).Scan(&id, &date)
	if err != nil {
		log.Println(false, id, date)
		return date, false
	}
	log.Println(true, id, date)
	return date, true

}
