package controllers

import (
	"main/app/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//マイページ
func PersonalInformation(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")

	if UserInfo.UserName != nil {
		userid := models.GetUserID(UserInfo.UserName)
		personal := models.GetPersonal(userid)
		r_personal := models.GetReturnPersonal(userid)

		c.HTML(200, "personalInformation", gin.H{
			"title":            "個人情報",
			"login":            true,
			"kanji_f_name":     personal[2],
			"kanji_l_name":     personal[3],
			"kana_f_name":      personal[4],
			"kana_l_name":      personal[5],
			"postal_code":      personal[6],
			"address_level1":   personal[7],
			"address_level2":   personal[8],
			"address_line1":    personal[9],
			"address_line2":    personal[10],
			"organization":     personal[11],
			"r_kanji_f_name":   r_personal[2],
			"r_kanji_l_name":   r_personal[3],
			"r_kana_f_name":    r_personal[4],
			"r_kana_l_name":    r_personal[5],
			"r_phone_number":   r_personal[6],
			"r_postal_code":    r_personal[7],
			"r_address_level1": r_personal[8],
			"r_address_level2": r_personal[9],
			"r_address_line1":  r_personal[10],
			"r_address_line2":  r_personal[11],
			"r_organization":   r_personal[12],
			"username":         session.Get("UserName"),
			"csrfToken":        csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

func PersonalInformationInput(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName != nil {
		userid := models.GetUserID(UserInfo.UserName)
		personal := models.GetPersonal(userid)
		c.HTML(200, "PersonalInformationInput", gin.H{
			"title":          "お届け先入力",
			"login":          true,
			"username":       session.Get("UserName"),
			"kanji_f_name":   personal[2],
			"kanji_l_name":   personal[3],
			"kana_f_name":    personal[4],
			"kana_l_name":    personal[5],
			"postal_code":    personal[6],
			"address_level1": personal[7],
			"address_level2": personal[8],
			"address_line1":  personal[9],
			"address_line2":  personal[10],
			"organization":   personal[11],
			"csrfToken":      csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

func PersonalInformationInputPost(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName != nil {
		kanji_f_name := c.PostForm("kanji-f-name")
		kanji_l_name := c.PostForm("kanji-l-name")
		kana_f_name := c.PostForm("kana-f-name")
		kana_l_name := c.PostForm("kana-l-name")
		postal_code := c.PostForm("postal-code")
		address_level1 := c.PostForm("address-level1")
		address_level2 := c.PostForm("address-level2")
		address_line1 := c.PostForm("address-line1")
		address_line2 := c.PostForm("address-line2")
		organization := c.PostForm("organization")

		userid := models.GetUserID(UserInfo.UserName)

		if models.PersonalUserIdCheck(userid) == "なし" {
			models.PersonalInsert(userid, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization)
		} else {
			models.PersonalUpdate(userid, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, postal_code, address_level1, address_level2, address_line1, address_line2, organization)
		}
		c.Redirect(302, "/personal-information")
	} else {
		c.Redirect(302, "/")
	}
}

func ReturnPersonalInformationInput(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName != nil {
		userid := models.GetUserID(UserInfo.UserName)
		personal := models.GetReturnPersonal(userid)
		c.HTML(200, "ReturnPersonalInformationInput", gin.H{
			"title":                 "返品先情報",
			"login":                 true,
			"username":              session.Get("UserName"),
			"return_kanji_f_name":   personal[2],
			"return_kanji_l_name":   personal[3],
			"return_kana_f_name":    personal[4],
			"return_kana_l_name":    personal[5],
			"return_phone_number":   personal[6],
			"return_postal_code":    personal[7],
			"return_address_level1": personal[8],
			"return_address_level2": personal[9],
			"return_address_line1":  personal[10],
			"return_address_line2":  personal[11],
			"return_organization":   personal[12],
			"csrfToken":             csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

func ReturnPersonalInformationInputPost(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName != nil {
		kanji_f_name := c.PostForm("kanji-f-name")
		kanji_l_name := c.PostForm("kanji-l-name")
		kana_f_name := c.PostForm("kana-f-name")
		kana_l_name := c.PostForm("kana-l-name")
		phone_number := c.PostForm("phone-number")
		postal_code := c.PostForm("postal-code")
		address_level1 := c.PostForm("address-level1")
		address_level2 := c.PostForm("address-level2")
		address_line1 := c.PostForm("address-line1")
		address_line2 := c.PostForm("address-line2")
		organization := c.PostForm("organization")

		userid := models.GetUserID(UserInfo.UserName)

		if models.ReturnPersonalUserIdCheck(userid) == "なし" {
			models.ReturnPersonalInsert(userid, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, phone_number, postal_code, address_level1, address_level2, address_line1, address_line2, organization)
		} else {
			models.ReturnPersonalUpdate(userid, kanji_f_name, kanji_l_name, kana_f_name, kana_l_name, phone_number, postal_code, address_level1, address_level2, address_line1, address_line2, organization)
		}
		c.Redirect(302, "/personal-information")
	} else {
		c.Redirect(302, "/")
	}
}

func ReturnPersonalInformation(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")
	productid := c.PostForm("productid")
	product_userid := c.PostForm("product_userid")
	product_userid_int, _ := strconv.Atoi(product_userid)
	if UserInfo.UserName != nil && models.CheckArrives(productid) == "1" {
		personal := models.GetReturnPersonal(product_userid_int)
		product := models.GetProduct(productid)
		uname := models.GetUserName(product_userid)
		c.HTML(200, "ReturnPersonalInformation", gin.H{
			"title":                 "返品先情報",
			"login":                 true,
			"username":              session.Get("UserName"),
			"uname":                 uname,
			"ItemName":              product[4],
			"return_kanji_f_name":   personal[2],
			"return_kanji_l_name":   personal[3],
			"return_kana_f_name":    personal[4],
			"return_kana_l_name":    personal[5],
			"return_phone_number":   personal[6],
			"return_postal_code":    personal[7],
			"return_address_level1": personal[8],
			"return_address_level2": personal[9],
			"return_address_line1":  personal[10],
			"return_address_line2":  personal[11],
			"return_organization":   personal[12],
			"csrfToken":             csrf.GetToken(c),
		})
	}
}
