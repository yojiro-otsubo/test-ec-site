package controllers

import (
	"log"
	"main/app/models"
	"main/app/others"
	"main/config"
	"net/http"
	"net/mail"
	"net/smtp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

//-------------------------------------------------- AUTH --------------------------------------------------

type email_detail struct {
	from     string
	username string
	password string
	to       string
	sub      string
	msg      string
}

func gmailSend(m email_detail) error {
	smtpSvr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", m.username, m.password, "smtp.gmail.com")
	if err := smtp.SendMail(smtpSvr, auth, m.from, []string{m.to}, []byte(m.msg)); err != nil {
		return err
	}
	return nil
}

func SendMail(c *gin.Context) {

	emailaddr := c.PostForm("email")
	username := c.PostForm("username")
	pk := c.PostForm("pk")
	url := "http://localhost:8080/registration/" + pk
	msg := username + " 様 \n" + "仮登録が完了しました。\n" + "下記URLから本登録を完了してください。\n" + url
	sub := "【サービス名】本会員登録"

	log.Println("from:", config.Config.GmailAddr, "username:", config.Config.GmailUser, "password", config.Config.GmailPass)
	m := email_detail{
		from:     config.Config.GmailAddr,
		username: config.Config.GmailUser,
		password: config.Config.GmailPass,
		to:       emailaddr,
		sub:      sub,
		msg:      msg,
	}

	if err := gmailSend(m); err != nil {
		log.Println(err)
	}

	c.HTML(200, "SignupInputCheck", gin.H{
		"email":     emailaddr,
		"username":  username,
		"csrfToken": csrf.GetToken(c),
		"login":     false,
		"sendok":    true,
	})

}

func SignupSuccess(c *gin.Context) {
	c.HTML(200, "SignupSuccess", gin.H{
		"csrfToken": csrf.GetToken(c),
		"login":     false,
		"sendok":    true,
	})
}

func Registration(c *gin.Context) {
	pk := c.Param("pk")
	username, email, pass := models.GetKariUserALL(pk)
	models.UserRegistration(username, email, pass)
	models.DeleteKariUser(pk)
	c.Redirect(302, "/signup-success")
}

func SignupInputCheck(c *gin.Context) {
	karinumber := c.Param("karinumber")
	username, email := models.KariUserCheck(karinumber)
	c.HTML(200, "SignupInputCheck", gin.H{
		"email":      email,
		"username":   username,
		"karinumber": karinumber,
		"csrfToken":  csrf.GetToken(c),
		"login":      false,
	})

}

//ユーザー登録
func kari_registration(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	//emailチェック
	e, err := mail.ParseAddress(email)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signupform", gin.H{
			"bad_email": "正しくないメールアドレス",
			"email":     email,
			"csrfToken": csrf.GetToken(c),
			"login":     false,
		})
		return
	}

	//passwordハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	//usernameとemail存在チェック
	if models.UserCheck(username) == true && models.EmailCheck(e.Address) == true {
		//登録
		pk := models.KariUserRegistration(username, e.Address, string(hash))
		redirectURL := "/signupinputcheck/" + pk
		c.Redirect(302, redirectURL)
	} else if models.UserCheck(username) == false && models.EmailCheck(e.Address) == true {
		c.HTML(http.StatusBadRequest, "signupform", gin.H{
			"username_status": "は既に使われたユーザーネーム",
			"username":        username,
			"login":           false,
			"csrfToken":       csrf.GetToken(c),
		})
		return
	} else if models.UserCheck(username) == true && models.EmailCheck(e.Address) == false {
		c.HTML(http.StatusBadRequest, "signupform", gin.H{
			"email_status": "は既に使われたメールアドレス",
			"email":        email,
			"login":        false,
			"csrfToken":    csrf.GetToken(c),
		})
		return
	} else {
		c.HTML(http.StatusBadRequest, "signupform", gin.H{
			"username_status": "は既に使われたユーザーネーム",
			"username":        username,
			"email_status":    "は既に使われたメールアドレス",
			"email":           email,
			"csrfToken":       csrf.GetToken(c),
			"login":           false,
		})
		return
	}

}

//ログイン処理
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if models.LoginCheck(username, password) == true {
		var check bool
		userid := models.GetUserID(username)
		stripeid, check := models.GetStripeAccountId(userid)
		var login_token string
		for {
			login_token = ""
			login_token = others.RandString(30)
			bool_token := models.TokenCheck(login_token)
			if bool_token == true {
				break
			}
		}

		if check == true {
			session := sessions.Default(c)
			session.Set("UserName", username)
			session.Set("StripeAccount", stripeid)
			session.Set("logintoken", login_token)
			session.Save()
			log.Println("username = ", session.Get("UserName"), "///stripeid = ", session.Get("StripeAccount"))
			c.Redirect(302, "/")
		} else {
			session := sessions.Default(c)
			session.Set("UserName", username)
			session.Set("logintoken", login_token)
			session.Save()
			log.Println("username = ", session.Get("UserId"), "///stripeid = 未登録")

			c.Redirect(302, "/")
		}

	} else {
		c.HTML(http.StatusBadRequest, "loginform", gin.H{
			"login_status": "ユーザーネームまたはパスワードが違います",
			"login":        false,
			"csrfToken":    csrf.GetToken(c),
		})
	}

}

//ログアウト処理
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")

	if UserInfo.UserName != nil {
		session.Clear()
		session.Save()
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/")
	}

}

//ログインフォーム
func LoginForm(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName == nil {
		c.HTML(200, "loginform", gin.H{
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

//登録フォーム
func SignupForm(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.UserName == nil {
		c.HTML(200, "signupform", gin.H{
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}
