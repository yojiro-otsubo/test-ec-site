package controllers

import (
	"log"
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//トップページ
func top(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	products := models.GetProductTop()
	log.Println(products)
	if UserInfo.UserName == nil {
		c.HTML(200, "top", gin.H{
			"title":     "top",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
			"products":  products,
		})
	} else {
		c.HTML(200, "top", gin.H{
			"title":     "top",
			"login":     true,
			"username":  UserInfo.UserName,
			"csrfToken": csrf.GetToken(c),
			"products":  products,
		})
	}
}
