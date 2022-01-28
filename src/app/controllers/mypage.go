package controllers

import (
	"io"
	"log"
	"main/app/models"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//マイページ
func mypage(c *gin.Context) {
	uname := c.Param("username")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId == uname {
		userid := models.GetUserID(UserInfo.UserId)
		Self_Introduction := models.GetSelfIntroduction(userid)
		email := models.GetUserEmail(userid)
		struserid := strconv.Itoa(userid)
		filepath := "app/static/img/icon/userid" + struserid + "/" + "icon.jpg"
		if f, err := os.Stat(filepath); os.IsNotExist(err) || f.IsDir() {
			log.Println("ファイルは存在しません！")
			c.HTML(200, "mypage", gin.H{
				"title":            "マイページ",
				"login":            true,
				"username":         UserInfo.UserId,
				"userid":           userid,
				"email":            email,
				"csrfToken":        csrf.GetToken(c),
				"SelfIntroduction": Self_Introduction,
				"filepath":         false,
			})
		} else {
			log.Println("存在するファイルです")
			c.HTML(200, "mypage", gin.H{
				"title":            "マイページ",
				"login":            true,
				"username":         UserInfo.UserId,
				"userid":           userid,
				"email":            email,
				"csrfToken":        csrf.GetToken(c),
				"SelfIntroduction": Self_Introduction,
				"filepath":         true,
			})
		}
	} else {
		c.Redirect(302, "/")
	}
}

func mypageDetail(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId != nil {
		userid := models.GetUserID(UserInfo.UserId)
		self_introduction := c.PostForm("self-introduction")
		models.SelfIntroductionRegistration(userid, self_introduction)
		struserid := strconv.Itoa(userid)
		username := models.GetUserName(struserid)
		redirecturl := "/mypage/" + username
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			log.Println(err)
			c.Redirect(302, redirecturl)
			return
		}

		filename := header.Filename
		log.Println("filename =", filename)
		pos := strings.LastIndex(filename, ".")
		log.Println("pos =", filename[pos:])

		create_path := "app/static/img/icon/userid" + struserid + "/" + filename
		mkdir_path := "app/static/img/icon/userid" + struserid

		// mkdir
		err = os.Mkdir(mkdir_path, 0755)
		if err != nil {
			log.Println(err)
		}
		out, err := os.Create(create_path)
		if err != nil {
			log.Println("os.create err = ", err)
		}

		//file copy
		_, err = io.Copy(out, file)
		if err != nil {
			log.Println("io.copy err = ", err)
		}

		defer out.Close()

		newpathjpg := "app/static/img/icon/userid" + struserid + "/" + "icon.jpg"

		if filename[pos:] == ".png" {
			if err := os.Rename(create_path, newpathjpg); err != nil {
				log.Println(err)
			}
		} else {
			if err := os.Rename(create_path, newpathjpg); err != nil {
				log.Println(err)
			}
		}

		c.Redirect(302, redirecturl)

	} else {
		c.Redirect(302, "/")
	}
}
