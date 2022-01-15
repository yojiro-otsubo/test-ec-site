package controllers

import (
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//購入履歴
func purchaseHistory(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId == nil {
		c.Redirect(302, "/loginform")
	} else {
		userid := models.GetUserID(UserInfo.UserId)
		products := models.GetProductIdFromPaymentHistory(userid)
		c.HTML(200, "purchaseHistory", gin.H{
			"title":     "purchaseHistory",
			"login":     true,
			"username":  UserInfo.UserId,
			"csrfToken": csrf.GetToken(c),
			"products":  products,
		})
	}
}
