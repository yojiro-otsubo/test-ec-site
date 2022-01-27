package controllers

import (
	"log"
	"main/app/models"
	"math"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//-------------------------------------------------- ProductPage --------------------------------------------------
func ProductPage(c *gin.Context) {
	productNumber := c.Param("number")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	product := models.GetProduct(productNumber)
	log.Println(product)
	username := models.GetUserName(product[1])
	if product[7] == "1" {
		c.Redirect(302, "/")
	}
	f, err := strconv.ParseFloat(product[6], 64)
	if err != nil {
		log.Println(err)
	}
	f = f * 1.1
	taxamount := int(math.Round(f))
	log.Println(taxamount)

	if UserInfo.UserId == nil {
		c.HTML(200, "product", gin.H{
			"title":           "product",
			"login":           false,
			"csrfToken":       csrf.GetToken(c),
			"ProductId":       product[0],
			"ProductUsername": username,
			"StripeProductId": product[2],
			"StripePriceId":   product[3],
			"ItemName":        product[4],
			"Description":     product[5],
			"username":        UserInfo.UserId,
			"Amount":          taxamount,
		})
	} else {
		userid := models.GetUserID(UserInfo.UserId)
		if models.CheckCart(userid, product[0]) == true {
			c.HTML(200, "product", gin.H{
				"title":           "product",
				"login":           true,
				"username":        UserInfo.UserId,
				"csrfToken":       csrf.GetToken(c),
				"ProductId":       product[0],
				"ProductUsername": username,
				"UserId":          product[1],
				"StripeProductId": product[2],
				"StripePriceId":   product[3],
				"ItemName":        product[4],
				"Description":     product[5],
				"Amount":          taxamount,
				"cart":            true,
			})
		} else {
			c.HTML(200, "product", gin.H{
				"title":           "product",
				"login":           true,
				"username":        UserInfo.UserId,
				"csrfToken":       csrf.GetToken(c),
				"ProductId":       product[0],
				"ProductUsername": username,
				"UserId":          product[1],
				"StripeProductId": product[2],
				"StripePriceId":   product[3],
				"ItemName":        product[4],
				"Description":     product[5],
				"Amount":          taxamount,
				"cart":            false,
			})
		}

	}
}

func ProductImage(c *gin.Context) {
	productNumber := c.Param("number")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	c.HTML(200, "image", gin.H{
		"productid": productNumber,
	})
}

func UserProductPage(c *gin.Context) {
	user_id := c.Param("number")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	products := models.GetAllProductOfUserId(user_id)
	product_username := models.GetUserName(user_id)
	title := product_username + "さんのアイテム"
	if UserInfo.UserId == nil {
		c.HTML(200, "UserProductPage", gin.H{
			"title":           title,
			"login":           false,
			"csrfToken":       csrf.GetToken(c),
			"products":        products,
			"productUsername": title,
		})
	} else {
		c.HTML(200, "UserProductPage", gin.H{
			"title":           title,
			"login":           true,
			"username":        UserInfo.UserId,
			"csrfToken":       csrf.GetToken(c),
			"products":        products,
			"productUsername": title,
		})
	}

}
