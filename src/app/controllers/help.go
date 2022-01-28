package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func helpTop(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")

	c.HTML(200, "helpTop", gin.H{
		"title":    "ヘルプページ",
		"username": UserInfo.UserName,
		"login":    true,
	})
}
func helpSellItemGuide(c *gin.Context) {
	c.HTML(200, "helpSellItemGuide", gin.H{
		"title":    "出品ガイド",
		"username": UserInfo.UserName,
		"login":    true,
	})
}
func helpBuyItemGuide(c *gin.Context) {
	c.HTML(200, "helpBuyItemGuide", gin.H{
		"title":    "購入ガイド",
		"username": UserInfo.UserName,
		"login":    true,
	})
}
func helpRulesAndManners(c *gin.Context) {
	c.HTML(200, "helpRulesAndManners", gin.H{
		"title":    "マナー・ルール",
		"username": UserInfo.UserName,
		"login":    true,
	})
}
func helpReturnGuide(c *gin.Context) {
	c.HTML(200, "helpReturnGuide", gin.H{
		"title":    "返品ガイド",
		"username": UserInfo.UserName,
		"login":    true,
	})
}
func helpInquiry(c *gin.Context) {
	c.HTML(200, "helpInquiry", gin.H{
		"title":    "お問い合わせ",
		"username": UserInfo.UserName,
		"login":    true,
	})
}
