package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//マイページ
func mypage(c *gin.Context) {
	uname := c.Param("username")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId == uname {
		c.HTML(200, "mypage", gin.H{
			"title":    uname,
			"login":    true,
			"username": session.Get("UserId"),
		})
	} else {
		c.Redirect(302, "/")
	}
}
