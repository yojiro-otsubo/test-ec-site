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
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	products := models.GetProductTop()
	log.Println(UserInfo.UserName)
	log.Println(loginbool)
	log.Println(UserInfo.logintoken)
	c.HTML(200, "top", gin.H{
		"title":     "top",
		"login":     loginbool,
		"csrfToken": csrf.GetToken(c),
		"products":  products,
		"username":  UserInfo.UserName,
	})

}
