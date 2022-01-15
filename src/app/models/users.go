package models

import (
	"database/sql"
	"log"
	"main/config"

	"golang.org/x/crypto/bcrypt"
)

//username存在チェック
func UserCheck(username string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var uname string
	err = DbConnection.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&uname)
	if err != nil {
		return true
	} else {
		return false
	}

}

//emailの存在チェック
func EmailCheck(email string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var mail string
	err = DbConnection.QueryRow("SELECT email FROM users WHERE email = $1", email).Scan(&mail)
	if err != nil {
		return true
	} else {
		return false
	}

}

//ログインチェック
func LoginCheck(username, password string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var uname string
	var hpass []byte
	//usernameの存在
	err = DbConnection.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&uname)
	if err != nil {
		return false
	} else {
		//passwordの存在チェック
		err = DbConnection.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hpass)
		if err != nil {
			return false
		}
		//passwordチェック
		err := bcrypt.CompareHashAndPassword(hpass, []byte(password))
		if err != nil {
			return false
		}
		return true
	}

}

//user登録
func UserRegistration(username, email, hashpassword string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	//userにINSERTする
	cmd, err := DbConnection.Prepare("INSERT INTO users(username, password, email) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(username, hashpassword, email)
	if err != nil {
		log.Println(err)
	}
}

//usernameでid取得
func GetUserID(username interface{}) int {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	var userid int
	err = DbConnection.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userid)
	if err != nil {
		log.Println(err)
	}

	return userid

}

//useridでusernameを取得
func GetUserName(user_id string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var username string
	err = DbConnection.QueryRow("SELECT username FROM users WHERE id = $1", user_id).Scan(&username)
	if err != nil {
		log.Println(err)
	}

	return username
}
