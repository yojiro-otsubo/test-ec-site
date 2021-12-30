package models

import (
	"database/sql"
	"fmt"
	"log"
	"main/config"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

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
	cmd2 := "CREATE TABLE IF NOT EXISTS products (id serial PRIMARY KEY, user_id INT, stripe_product_id VARCHAR(255), stripe_price_id VARCHAR(255));"
	_, err = DbConnection.Exec(cmd2)
	if err != nil {
		log.Fatalln(err)
	}

	//settlementテーブル作成
	cmd3 := "CREATE TABLE IF NOT EXISTS settlement (id serial PRIMARY KEY, user_id INT, product_id INT);"
	_, err = DbConnection.Exec(cmd3)
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

//動作テスト
func TestDb() {
	DbConnection, err := sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	var id int

	err = DbConnection.QueryRow("INSERT INTO test_db(test) VALUES($1) RETURNING id", "hoge").Scan(&id)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("PRIMARY KEY = ", id)
}

//userテーブルのusername存在チェック
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

//userテーブルのemailの存在チェック
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

//accountsテーブルにuser_idとstripe_accountをINSERT
func AccountRegist(userid int, stripeid interface{}) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("INSERT INTO accounts(user_id, stripe_account) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(userid, stripeid)
	if err != nil {
		log.Println(err)
	}

}

//usernameでusersテーブルのid取得
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

//ussr_idでstripe_account取得
func GetStripeAccountId(userid int) (string, bool) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	var stripeid string
	err = DbConnection.QueryRow("SELECT stripe_account FROM accounts WHERE user_id = $1", userid).Scan(&stripeid)
	if err != nil {
		log.Println(err)
		return stripeid, false
	}

	return stripeid, true
}

func RegistProduct(pkid int, productid, priceid string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	var id int
	var user_id int
	var stripe_product_id string
	var stripe_price_id string
	err = DbConnection.QueryRow("UPDATE products SET stripe_product_id = $2, stripe_price_id = $3 WHERE id = $1 RETURNING id, user_id, stripe_product_id, stripe_price_id").Scan(&id, &user_id, &stripe_product_id, &stripe_price_id)
	if err != nil {
		log.Println(err)
	}
	log.Println("id = ", id, "\nuser_id = ", user_id, "\nstripe_product_id = ", stripe_product_id, "\nstripe_price_id = ", stripe_price_id)

}

func RegistUserIdAndGetProductId(userid int) int {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	var id int

	err = DbConnection.QueryRow("INSERT INTO products(user_id) VALUES($1) RETURNING id", userid).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	log.Println("PRIMARY KEY = ", id)

	return id

}

//accountsテーブルのuserid存在チェック
func UserIdCheck(userid int) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var uid string
	err = DbConnection.QueryRow("SELECT user_id FROM accounts WHERE user_id = $1", userid).Scan(&uid)
	if err != nil {
		return true
	} else {
		return false
	}

}
