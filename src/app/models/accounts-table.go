package models

import (
	"database/sql"
	"log"
	"main/config"
)

//accountsテーブルにuser_idとstripe_accountをINSERT
func AccountRegist(userid int, stripeid interface{}) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	cmd, err := DbConnection.Prepare("INSERT INTO accounts(user_id, stripe_account) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(userid, stripeid)
	if err != nil {
		log.Println(err)
	}

}

//ussr_idでstripe_account取得
func GetStripeAccountId(userid int) (string, bool) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	var stripeid string
	err = DbConnection.QueryRow("SELECT stripe_account FROM accounts WHERE user_id = $1", userid).Scan(&stripeid)
	if err != nil {
		log.Println(err)
		return stripeid, false
	}

	return stripeid, true
}

//accountsテーブルのuserid存在チェック
func UserIdCheck(userid int) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

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

func CheckStripeAccountId(stripe_account_id interface{}) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var stripeAccount string
	err = DbConnection.QueryRow("SELECT stripe_account FROM accounts WHERE stripe_account = $1", stripe_account_id).Scan(&stripeAccount)
	if err != nil {
		return false
	} else {
		return true
	}
}
