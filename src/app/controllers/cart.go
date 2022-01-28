package controllers

import (
	"log"
	"main/app/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//カート
func AddCart(c *gin.Context) {
	productid := c.PostForm("cart")
	log.Println("product = ", productid)
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	userid := models.GetUserID(UserInfo.UserName)
	struserid := strconv.Itoa(userid)
	product := models.GetProduct(productid)

	if product[7] == "1" {
		c.Redirect(302, "/")
	}

	if UserInfo.UserName != nil && struserid != product[1] {
		userid := models.GetUserID(UserInfo.UserName)
		models.AddToCart(userid, productid)
		redirecturl := "/product/" + productid
		c.Redirect(302, redirecturl)
	} else {
		c.Redirect(302, "/loginform")
	}
}

func CartPage(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
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
	if UserInfo.UserName != nil {
		c.HTML(200, "cart", gin.H{
			"title":       "cart",
			"login":       true,
			"products":    products,
			"username":    UserInfo.UserName,
			"csrfToken":   csrf.GetToken(c),
			"totalAmount": totalAmount,
		})
	} else {
		c.Redirect(302, "/loginform")
	}
}

func DeleteItemsInCart(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName != nil {
		userid := models.GetUserID(UserInfo.UserName)
		productid := c.PostForm("delete_item")
		models.DeleteCartItem(userid, productid)
		c.Redirect(302, "/mycart")
	} else {
		c.Redirect(302, "/")
	}
}
