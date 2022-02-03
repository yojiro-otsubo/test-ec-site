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
	defer DbConnection.Close()

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

type Inquiry struct {
	Id, Name, Email, InquiryMsg, Date string
}

func GetInquiryAll() []Inquiry {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT * FROM inquiry ORDER BY id DESC")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var inquiryResult []Inquiry
	for rows.Next() {
		var i Inquiry
		err := rows.Scan(&i.Id, &i.Name, &i.Email, &i.InquiryMsg, &i.Date)
		if err != nil {
			log.Println(err)
		}

		inquiryResult = append(inquiryResult, i)

	}

	return inquiryResult

}
