package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func helpTop(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	c.HTML(200, "helpTop", gin.H{
		"title": "ヘルプページ",
		"login": true,
	})
}
