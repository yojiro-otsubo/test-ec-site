package controllers

import (
	"log"
	"main/app/models"
	"main/app/others"
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
	UserInfo.UserName = session.Get("UserName")
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

	if UserInfo.UserName == nil {
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
			"username":        UserInfo.UserName,
			"Amount":          taxamount,
		})
	} else {
		userid := models.GetUserID(UserInfo.UserName)
		if models.CheckCart(userid, product[0]) == true {
			c.HTML(200, "product", gin.H{
				"title":           "product",
				"login":           true,
				"username":        UserInfo.UserName,
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
				"username":        UserInfo.UserName,
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
	UserInfo.UserName = session.Get("UserName")

	c.HTML(200, "image", gin.H{
		"productid": productNumber,
	})
}

func UserProductPage(c *gin.Context) {
	user_id := c.Param("number")
	int_user_id, _ := strconv.Atoi(user_id)
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	products := models.GetAllProductOfUserId(user_id)
	product_username := models.GetUserName(user_id)
	title := product_username + "さんのアイテム"
	count_follower := models.CountFollower(int_user_id)
	my_user_id := models.GetUserID(UserInfo.UserName)
	count_product := models.CountProduct(int_user_id)
	filepathbool := others.IconFilePathCheck(user_id)
	self_introduction := models.GetSelfIntroduction(int_user_id)
	log.Println(filepathbool)
	if UserInfo.UserName == nil {
		c.HTML(200, "UserProductPage", gin.H{
			"title":            title,
			"login":            false,
			"csrfToken":        csrf.GetToken(c),
			"products":         products,
			"productUserId":    user_id,
			"productUsername":  product_username,
			"countFollower":    count_follower,
			"countProduct":     count_product,
			"filepath":         filepathbool,
			"SelfIntroduction": self_introduction,
		})
	} else {
		if models.CheckFollow(my_user_id, int_user_id) == "なし" {
			c.HTML(200, "UserProductPage", gin.H{
				"title":            title,
				"login":            true,
				"username":         UserInfo.UserName,
				"csrfToken":        csrf.GetToken(c),
				"products":         products,
				"productUsername":  product_username,
				"productUserId":    user_id,
				"countFollower":    count_follower,
				"follow":           false,
				"countProduct":     count_product,
				"filepath":         filepathbool,
				"SelfIntroduction": self_introduction,
			})
		} else {
			c.HTML(200, "UserProductPage", gin.H{
				"title":            title,
				"login":            true,
				"username":         UserInfo.UserName,
				"csrfToken":        csrf.GetToken(c),
				"products":         products,
				"productUsername":  product_username,
				"productUserId":    user_id,
				"countFollower":    count_follower,
				"countProduct":     count_product,
				"follow":           true,
				"filepath":         filepathbool,
				"SelfIntroduction": self_introduction,
			})
		}

	}

}
