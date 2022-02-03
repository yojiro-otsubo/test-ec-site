package controllers

import (
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func Search(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	category := c.PostForm("category")
	f_price := c.PostForm("f")
	l_price := c.PostForm("l")
	var products []models.Product
	if category == "全て" {
		products = models.GetProductSearchAll(f_price, l_price)
	} else {
		products = models.GetProductSearch(category, f_price, l_price)

	}

	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "Search", gin.H{
		"title":     "Search",
		"login":     loginbool,
		"csrfToken": csrf.GetToken(c),
		"products":  products,
		"username":  UserInfo.UserName,
		"c":         category,
		"f":         f_price,
		"l":         l_price,
	})

}
