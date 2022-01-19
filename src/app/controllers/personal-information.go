package controllers

import (
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//マイページ
func PersonalInformation(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId != nil {
		userid := models.GetUserID(UserInfo.UserId)
		personal := models.GetPersonal(userid)
		c.HTML(200, "personalInformation", gin.H{
			"title":          "個人情報",
			"login":          true,
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
			"username":       session.Get("UserId"),
			"csrfToken":      csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

func PersonalInformationInput(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId != nil {
		userid := models.GetUserID(UserInfo.UserId)
		personal := models.GetPersonal(userid)
		c.HTML(200, "PersonalInformationInput", gin.H{
			"title":          "お届け先入力",
			"login":          true,
			"username":       session.Get("UserId"),
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
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId != nil {
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

		userid := models.GetUserID(UserInfo.UserId)

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
