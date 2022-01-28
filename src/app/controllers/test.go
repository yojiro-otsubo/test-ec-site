package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//-------------------------------------------------- Test --------------------------------------------------
func test(c *gin.Context) {

	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	/*
		out, err := os.Create("app/static/img/item/test.txt")
		if err != nil {
			log.Println(err)
		}
		defer out.Close()
	*/
	/*
		err := os.Mkdir("app/static/img/item/test", 0755)
		if err != nil {
			log.Println(err)
		}
	*/
	if UserInfo.UserName == nil {
		c.HTML(200, "test", gin.H{
			"title":     "test",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "test", gin.H{
			"title":     "test",
			"login":     true,
			"username":  UserInfo.UserName,
			"csrfToken": csrf.GetToken(c),
		})
	}
}
