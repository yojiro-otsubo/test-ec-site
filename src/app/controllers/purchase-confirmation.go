package controllers

import (
	"log"
	"main/app/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func PurchaseConfirmationCart(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)

	userid := models.GetUserID(UserInfo.UserName)
	if loginbool == true && models.PersonalUserIdCheck(userid) == "あり" {
		userid := models.GetUserID(UserInfo.UserName)
		products := models.GetProductFromCartDB(userid, 1.1)
		log.Println(products)

		var totalAmount int
		for _, p := range products {
			i, err := strconv.Atoi(p.Amount)
			if err != nil {
				log.Println(err)
			}
			totalAmount = totalAmount + i
		}
		log.Println(totalAmount)
		if loginbool == true {
			c.HTML(200, "PurchaseConfirmation", gin.H{
				"title":       "PurchaseConfirmation",
				"login":       true,
				"products":    products,
				"username":    UserInfo.UserName,
				"csrfToken":   csrf.GetToken(c),
				"totalAmount": totalAmount,
			})
		}
	} else if loginbool == true && models.PersonalUserIdCheck(userid) == "なし" {
		c.Redirect(302, "/personal-information-input")
	} else {
		c.Redirect(302, "/loginform")
	}
}
