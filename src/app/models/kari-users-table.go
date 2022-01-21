package models

import (
	"database/sql"
	"log"
	"main/config"
)

//user登録
func KariUserRegistration(username, email, hashpassword string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}
	//userにINSERTする

	var id string
	err = DbConnection.QueryRow("INSERT INTO kari_users(username, password, email) VALUES($1, $2, $3) RETURNING id", username, hashpassword, email).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	return id
}

func KariUserCheck(pk string) (string, string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var id, uname, email, pass string
	err = DbConnection.QueryRow("SELECT * FROM kari_users WHERE id = $1", pk).Scan(&id, &uname, &pass, &email)
	if err != nil {
		log.Println(err)
	}
	return uname, email
}
func GetKariUserALL(pk string) (string, string, string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var id, uname, email, pass string
	err = DbConnection.QueryRow("SELECT * FROM kari_users WHERE id = $1", pk).Scan(&id, &uname, &pass, &email)
	if err != nil {
		log.Println(err)
	}
	return uname, email, pass
}

func DeleteKariUser(pk string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("DELETE FROM kari_users WHERE id = $1")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(pk)
	if err != nil {
		log.Println(err)
	}
}
