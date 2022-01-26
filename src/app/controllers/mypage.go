package controllers

import (
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//マイページ
func mypage(c *gin.Context) {
	uname := c.Param("username")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId == uname {
		userid := models.GetUserID(UserInfo.UserId)
		email := models.GetUserEmail(userid)
		c.HTML(200, "mypage", gin.H{
			"title":    "マイページ",
			"login":    true,
			"username": UserInfo.UserId,
			"email":    email,
		})
	} else {
		c.Redirect(302, "/")
	}
}
