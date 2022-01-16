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
	UserInfo.UserId = session.Get("UserId")
	userid := models.GetUserID(UserInfo.UserId)
	struserid := strconv.Itoa(userid)
	product := models.GetProduct(productid)

	if UserInfo.UserId != nil && struserid != product[1] {
		userid := models.GetUserID(UserInfo.UserId)
		models.AddToCart(userid, productid)
		redirecturl := "/product/" + productid
		c.Redirect(302, redirecturl)
	} else {
		c.Redirect(302, "/")
	}
}

func CartPage(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	userid := models.GetUserID(UserInfo.UserId)
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
	if UserInfo.UserId != nil {
		c.HTML(200, "cart", gin.H{
			"title":       "cart",
			"login":       true,
			"products":    products,
			"username":    UserInfo.UserId,
			"csrfToken":   csrf.GetToken(c),
			"totalAmount": totalAmount,
		})
	} else {
		c.Redirect(302, "/loginform")
	}
}

func DeleteItemsInCart(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId != nil {
		userid := models.GetUserID(UserInfo.UserId)
		productid := c.PostForm("delete_item")
		models.DeleteCartItem(userid, productid)
		c.Redirect(302, "/mycart")
	} else {
		c.Redirect(302, "/")
	}
}
