package models

import (
	"database/sql"
	"log"
	"main/config"
)

func PersonalInsert(userid int, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	cmd, err := DbConnection.Prepare(`INSERT INTO personal_info(
		user_id, 
		kanji_f_name,
		kanji_l_name, 
		kana_f_name, 
		kana_l_name, 
		postal_code, 
		address_level1, 
		address_level2, 
		address_line1, 
		address_line2, 
		organization
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`)
	if err != nil {
		log.Println(err)
	}
	_, err = cmd.Exec(userid, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization)
	if err != nil {
		log.Println(err)
	}

}

func PersonalUpdate(userid int, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization string) {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())
	if err != nil {
		log.Fatalln(err)
	}

	var id string
	err = DbConnection.QueryRow(`UPDATE personal_info SET 
		kanji_f_name = $2,
		kanji_l_name = $3, 
		kana_f_name = $4, 
		kana_l_name = $5, 
		postal_code = $6, 
		address_level1 = $7, 
		address_level2 = $8, 
		address_line1 = $9, 
		address_line2 = $10, 
		organization = $11
		WHERE user_id = $1 RETURNING id`, userid, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization).Scan(&id)
	if err != nil {
		log.Println(err)
	}

}

func PersonalUserIdCheck(userid int) bool {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}

	var id string
	err = DbConnection.QueryRow("SELECT user_id FROM personal_info WHERE user_id = $1", userid).Scan(&id)
	if err != nil {
		return true
	} else {
		return false
	}

}
func GetPersonal(userid int) [12]string {
	var err error
	DbConnection, err = sql.Open(config.Config.DBdriver, ConnectionInfo())

	if err != nil {
		log.Fatalln(err)
	}
	var Id, UserId, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization string

	err = DbConnection.QueryRow("SELECT * FROM personal_info WHERE user_id = $1", userid).Scan(&Id, &UserId, &kanji_f_name, &kanji_l_name, &kana_f_name, &kana_l_name, &postal_code, &address_level1, &address_level2, &address_line1, &address_line2, &organization)
	if err != nil {
		log.Println(err)
	}

	arr := [...]string{Id, UserId, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization}
	return arr
}
