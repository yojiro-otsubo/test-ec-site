package controllers

import (
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func helpTop(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpTop", gin.H{
		"title":    "ヘルプページ",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}
func helpSellItemGuide(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpSellItemGuide", gin.H{
		"title":    "出品ガイド",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}
func helpBuyItemGuide(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpBuyItemGuide", gin.H{
		"title":    "購入ガイド",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}
func helpRulesAndManners(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpRulesAndManners", gin.H{
		"title":    "マナー・ルール",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}
func helpReturnGuide(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpReturnGuide", gin.H{
		"title":    "返品ガイド",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}

func helpShippingGuide(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpShippingGuide", gin.H{
		"title":    "発送ガイド",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}
func helpInquiry(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	c.HTML(200, "helpInquiry", gin.H{
		"title":    "お問い合わせ",
		"username": UserInfo.UserName,
		"login":    loginbool,
	})
}
