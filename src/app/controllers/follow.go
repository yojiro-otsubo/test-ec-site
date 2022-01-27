package controllers

import (
	"log"
	"main/app/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Follow(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId != nil {
		follow_user_id := c.PostForm("user_id")
		int_follow_user_id, _ := strconv.Atoi(follow_user_id)
		log.Println("follow_user_id = ", follow_user_id)
		my_user_id := models.GetUserID(UserInfo.UserId)
		models.AddToFollow(my_user_id, int_follow_user_id)
		redirect_url := "/userproduct/" + follow_user_id
		c.Redirect(302, redirect_url)

	} else {
		c.Redirect(302, "loginform")
	}
}

func DeleteFollow(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId != nil {
		follow_user_id := c.PostForm("user_id")
		int_follow_user_id, _ := strconv.Atoi(follow_user_id)
		log.Println("follow_user_id = ", follow_user_id)
		my_user_id := models.GetUserID(UserInfo.UserId)
		models.DeleteFollow(my_user_id, int_follow_user_id)
		redirect_url := c.PostForm("redirect_url")
		log.Println(redirect_url)
		c.Redirect(302, redirect_url)

	} else {
		c.Redirect(302, "loginform")
	}
}
