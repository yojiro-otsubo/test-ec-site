package controllers

import (
	"log"
	"main/app/models"
	"main/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//-------------------------------------------------- Test --------------------------------------------------
func test(c *gin.Context) {

	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "test", gin.H{
		"title":     "test",
		"username":  UserInfo.UserName,
		"csrfToken": csrf.GetToken(c),
		"login":     loginbool,
	})

}

func testToken(c *gin.Context) {
	session := sessions.Default(c)
	pass := c.PostForm("pass")
	if config.Config.AdminPass == pass {
		session.Set("test", config.Config.AdminToken)
		session.Save()
		c.Redirect(302, "/test/top")
	} else {
		c.Redirect(302, "/")
	}
}

func testTop(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.Test = session.Get("test")
	if UserInfo.Test == config.Config.AdminToken {
		log.Println("pass ok")
		c.HTML(200, "testTop", gin.H{
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")

	}
}

func testInquiry(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.Test = session.Get("test")
	if UserInfo.Test == config.Config.AdminToken {
		inquiry := models.GetInquiryAll()
		log.Println("pass ok")
		c.HTML(200, "testInquiry", gin.H{
			"csrfToken": csrf.GetToken(c),
			"inquiry":   inquiry,
		})
	} else {
		c.Redirect(302, "/")

	}
}
