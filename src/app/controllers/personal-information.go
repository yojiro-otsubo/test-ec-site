package controllers

import (
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//マイページ
func PersonalInformation(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId != nil && models.CheckStripeAccountId(UserInfo.StripeAccount) == true {
		c.HTML(200, "personalInformation", gin.H{
			"title":    "個人情報",
			"login":    true,
			"username": session.Get("UserId"),
		})
	} else {
		c.Redirect(302, "/")
	}
}
