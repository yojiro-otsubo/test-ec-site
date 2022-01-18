package models

import (
	"database/sql"
	"log"
	"main/config"
)

func CheckTransferGroup(transferGroup string) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var transfergroup string
	err = DbConnection.QueryRow("SELECT transfer_group FROM payment_history WHERE transfer_group = $1", transferGroup).Scan(&transfergroup)
	if err != nil {
		return true
	} else {
		return false
	}
}
func GetTransferGroup(productid string) string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var transfergroup string
	err = DbConnection.QueryRow("SELECT transfer_group FROM payment_history WHERE product_id = $1", productid).Scan(&transfergroup)
	if err != nil {
		log.Println(err)
	}

	return transfergroup
}
func AddTransferGroup(user_id interface{}, product_id, transferGroup string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	var Transfer_Group string
	err = DbConnection.QueryRow("INSERT INTO payment_history(user_id, product_id, transfer_group) VALUES($1, $2, $3) RETURNING transfer_group", user_id, product_id, transferGroup).Scan(&Transfer_Group)
	if err != nil {
		log.Println(err)
	}
	log.Println(Transfer_Group)

}

type ProductId struct {
	productId string
}

func GetProductIdWithTg(transferGroup string) []string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	rows, err := DbConnection.Query("SELECT product_id FROM payment_history WHERE transfer_group = $1", transferGroup)
	if err != nil {
		log.Panicln("GetProductIdWithTg = ")
		log.Println(err)
	}
	defer rows.Close()

	var pid []ProductId
	for rows.Next() {
		var p ProductId
		err := rows.Scan(&p.productId)
		if err != nil {
			log.Println(err)
		}
		pid = append(pid, p)
	}

	var product_id []string
	for _, pp := range pid {
		product_id = append(product_id, pp.productId)
	}
	return product_id
}
