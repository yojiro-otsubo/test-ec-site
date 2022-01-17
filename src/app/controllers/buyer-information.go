package controllers

import (
	"log"
	"main/app/models"
	"main/config"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	csrf "github.com/utrack/gin-csrf"
)

func BuyerInformation(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	UserInfo.StripeAccount = session.Get("StripeAccount")

	//userid := models.GetUserID(UserInfo.UserId)
	if UserInfo.UserId != nil && models.CheckStripeAccountId(UserInfo.StripeAccount) == true {
		stripe.Key = config.Config.StripeKey
		productid := c.PostForm("productid")
		product_userid := models.GetProduct(productid)
		int_product_userid, _ := strconv.Atoi(product_userid[1])
		stripe_id, _ := models.GetStripeAccountId(int_product_userid)

		a, err := account.GetByID(stripe_id, nil)
		if err != nil {
			log.Println(err)
		}
		log.Println("body = ", string(a.LastResponse.RawJSON))

		c.HTML(200, "buyerInfo", gin.H{
			"title":       "BuyerInformation",
			"login":       true,
			"buyerUserId": product_userid[1],
			"csrfToken":   csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/loginform")
	}
}
