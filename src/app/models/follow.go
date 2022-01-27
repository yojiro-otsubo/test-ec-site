package models

import (
	"database/sql"
	"log"
	"main/config"
)

func AddToFollow(my_user_id, follow_user_id int) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}
	defer DbConnection.Close()

	cmd, err := DbConnection.Prepare("INSERT INTO follow(user_id, follow_user_id) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(my_user_id, follow_user_id)
	if err != nil {
		log.Println(err)
	}
}

//cartのデータ削除
func DeleteFollow(my_user_id, follow_user_id int) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare("DELETE FROM follow WHERE user_id = $1 AND follow_user_id = $2")
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(my_user_id, follow_user_id)
	if err != nil {
		log.Println(err)
	}
}

func CountFollower(user_id int) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var count string
	err = DbConnection.QueryRow("SELECT COUNT( * ) FROM follow WHERE follow_user_id = $1", user_id).Scan(&count)
	if err != nil {
		return "0"
	} else {
		log.Println("count = ", count)
		return count
	}
}

func CheckFollow(my_user_id, follow_user_id int) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	defer DbConnection.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var count string
	err = DbConnection.QueryRow("SELECT follow_user_id = $2 FROM follow WHERE user_id = $1", my_user_id, follow_user_id).Scan(&count)
	if err != nil {
		log.Println("なし")
		return "なし"
	} else {
		log.Println("あり")
		return "あり"
	}
}
