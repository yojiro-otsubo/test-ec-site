package controllers

import (
	"main/app/models"
	"main/config"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
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
		personal := models.GetPersonal(int_product_userid)
		c.HTML(200, "buyerInfo", gin.H{
			"title":          "BuyerInformation",
			"login":          true,
			"buyerUserId":    product_userid[1],
			"csrfToken":      csrf.GetToken(c),
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
			"username":       UserInfo.UserId,
			"productid":      productid,
		})
	} else {
		c.Redirect(302, "/loginform")
	}
}
