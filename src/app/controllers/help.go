package controllers

import (
	"log"
	"main/app/models"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
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
	UserInfo.InquiryId = session.Get("InquiryId")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	if UserInfo.InquiryId != nil {
		c.HTML(200, "helpInquiry", gin.H{
			"title":     "お問い合わせ",
			"username":  UserInfo.UserName,
			"login":     loginbool,
			"csrfToken": csrf.GetToken(c),
			"post":      true,
		})
	} else {
		c.HTML(200, "helpInquiry", gin.H{
			"title":     "お問い合わせ",
			"username":  UserInfo.UserName,
			"login":     loginbool,
			"csrfToken": csrf.GetToken(c),
			"post":      false,
		})
	}

}

func PostInquiry(c *gin.Context) {
	var err error
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	UserInfo.InquiryId = session.Get("InquiryId")
	log.Println("UserInfo.InquiryId", UserInfo.InquiryId)
	name := c.PostForm("name")
	email := c.PostForm("email")
	inquiry := c.PostForm("inquiry")

	t := time.Now()
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println(err)
	}
	timeTokyo := t.In(tokyo) //.Format(time.RFC3339)
	log.Println("timeTokyo = ", timeTokyo)
	yesterday := t.In(tokyo).Add(-24 * time.Hour)
	log.Println("yesterday = ", yesterday)

	if UserInfo.InquiryId != nil {
		log.Println("UserInfo.InquiryId = あり")
		d, databool := models.GetInquiry(UserInfo.InquiryId)
		date := stringToTime(d)
		log.Println("date =", date, "d =", d)
		if databool == true {
			if date.Before(yesterday) == true {
				log.Println("24時間経過")
				inquiry_id := models.InsertInquiry(name, email, inquiry, timeToString(timeTokyo))
				session.Delete("InquiryId")
				session.Set("InquiryId", inquiry_id)
				session.Save()
				c.Redirect(302, "/help/inquiry")
			} else {
				log.Println("24時間経過していません。")
				c.Redirect(302, "/help/inquiry")
			}
			c.Redirect(302, "/help/inquiry")
		}

	} else {
		log.Println("UserInfo.InquiryId = なし")

		inquiry_id := models.InsertInquiry(name, email, inquiry, timeToString(timeTokyo))
		session.Set("InquiryId", inquiry_id)
		session.Save()
		c.Redirect(302, "/help/inquiry")
	}

}

func stringToTime(str string) time.Time {
	var layout = "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, str)
	return t
}
func timeToString(t time.Time) string {
	var layout = "2006-01-02 15:04:05"
	str := t.Format(layout)
	return str
}
