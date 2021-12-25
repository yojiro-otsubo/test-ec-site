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

	cmd := "CREATE TABLE IF NOT EXISTS user (id serial PRIMARY KEY, username VARCHAR(50), password VARCHAR(255), email VARCHAR(255));"
	_, err = DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	cmd1 := "CREATE TABLE IF NOT EXISTS accouts (id serial PRIMARY KEY, user_id INT, stripe_account VARCHAR(100));"
	_, err = DbConnection.Exec(cmd1)
	if err != nil {
		log.Fatalln(err)
	}

	cmd2 := "CREATE TABLE IF NOT EXISTS products (id serial PRIMARY KEY, user_id INT, product_name VARCHAR(100), amount INT, quantity INT);"
	_, err = DbConnection.Exec(cmd2)
	if err != nil {
		log.Fatalln(err)
	}

	cmd3 := "CREATE TABLE IF NOT EXISTS settlement (id serial PRIMARY KEY, user_id INT, product_id INT);"
	_, err = DbConnection.Exec(cmd3)
	if err != nil {
		log.Fatalln(err)
	}

	cmdtest := "CREATE TABLE IF NOT EXISTS test_db (id serial PRIMARY KEY, test VARCHAR(50));"
	_, err = DbConnection.Exec(cmdtest)
	if err != nil {
		log.Fatalln(err)
	}

	defer DbConnection.Close()
}

func TestDb() {
	DbConnection, err := sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	InsertCmd, err := DbConnection.Prepare("INSERT INTO test_db(test) VALUES($1) RETURNING id")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = InsertCmd.Exec("hogehoge")
	if err != nil {
		log.Fatalln(err)
	}

	var test string

	err = DbConnection.QueryRow("SELECT test FROM test_db WHERE test = $1", "hogehoge").Scan(&test)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(test)
}

func UserCheck(username string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var uname string
	err = DbConnection.QueryRow("SELECT username FROM user_account WHERE username = $1", username).Scan(&uname)
	if err != nil {
		return true
	} else {
		return false
	}

}

func EmailCheck(email string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var mail string
	err = DbConnection.QueryRow("SELECT email FROM user_account WHERE email = $1", email).Scan(&mail)
	if err != nil {
		return true
	} else {
		return false
	}

}

func LoginCheck(username, password string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var uname string
	var hpass []byte
	err = DbConnection.QueryRow("SELECT username FROM user_account WHERE username = $1", username).Scan(&uname)
	if err != nil {
		return false
	} else {
		err = DbConnection.QueryRow("SELECT password FROM user_account WHERE username = $1", username).Scan(&hpass)
		if err != nil {
			log.Fatalln(err)
			return false
		}
		err := bcrypt.CompareHashAndPassword(hpass, []byte(password))
		if err != nil {
			return false
		}
		return true
	}

}

func UserRegistration(username, email, hashpassword string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("INSERT INTO user_account(username, password, email) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = cmd.Exec(username, hashpassword, email)
	if err != nil {
		log.Fatalln(err)
	}
}

func ItemRegistrationDB(username interface{}, item_name, item_description, price string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("INSERT INTO user_account(username, itemname, item_description) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = cmd.Exec(username, item_name, price)
	if err != nil {
		log.Fatalln(err)
	}
}
