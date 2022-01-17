package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//マイページ
func PersonalInformation(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserId = session.Get("UserId")
	//userid := models.GetUserID(UserInfo.UserId)
	if UserInfo.UserId != nil {

		c.HTML(200, "personalInformation", gin.H{
			"title":    "個人情報",
			"login":    true,
			"username": session.Get("UserId"),
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

		c.HTML(200, "PersonalInformationInput", gin.H{
			"title":    "お届け先入力",
			"login":    true,
			"username": session.Get("UserId"),
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

	} else {
		c.Redirect(302, "/")
	}
}
