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
	defer DbConnection.Close()

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

func SelfIntroductionCheck(user_id int) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var self_introduction string
	err = DbConnection.QueryRow("SELECT self_introduction FROM users WHERE id = $1", user_id).Scan(&self_introduction)
	if err != nil {
		return false
	} else {
		return true
	}

}

func GetSelfIntroduction(user_id int) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var self_introduction string
	err = DbConnection.QueryRow("SELECT self_introduction FROM users WHERE id = $1", user_id).Scan(&self_introduction)
	if err != nil {
		return ""
	} else {
		return self_introduction
	}

}

func SelfIntroductionRegistration(user_id int, self_introduction string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}
	var selfintroduction string

	err = DbConnection.QueryRow("UPDATE users SET self_introduction = $2 WHERE id = $1", user_id, self_introduction).Scan(&selfintroduction)
	if err != nil {
		log.Println(err)
	}

}

//ログインチェック
func LoginCheck(username, password string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

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
	defer DbConnection.Close()

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
	defer DbConnection.Close()

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
	defer DbConnection.Close()

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
func GetUserEmail(user_id int) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var email string
	err = DbConnection.QueryRow("SELECT email FROM users WHERE id = $1", user_id).Scan(&email)
	if err != nil {
		log.Println(err)
	}

	return email
}

func UpdateToken(user_id int, token string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}
	var returntoken string

	err = DbConnection.QueryRow("UPDATE users SET token = $2 WHERE id = $1", user_id, token).Scan(&returntoken)
	if err != nil {
		log.Println(err)
	}
}

func TokenCheck(token string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var booltoken string
	err = DbConnection.QueryRow("SELECT token FROM users WHERE token = $1", token).Scan(&booltoken)
	if err != nil {
		return true
	} else {
		return false
	}

}
func LoginTokenCheck(username, token interface{}) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var id string
	err = DbConnection.QueryRow("SELECT id FROM users WHERE username = $1 AND token = $2", username, token).Scan(&id)
	if err != nil {
		log.Println("login token check false")
		log.Println(err)
		return false
	} else {
		log.Println("login token check true")
		return true
	}

}
