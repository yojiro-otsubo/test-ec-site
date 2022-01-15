package controllers

import (
	"log"
	"main/app/models"
	"main/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
)

//-------------------------------------------------- StripeAccountLink --------------------------------------------------
//stripe連結アカウント作成、アカウント登録リンクへリダイレクト
func CreateAnExpressAccount(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserId = session.Get("UserId")
	log.Println("Username = ", UserInfo.UserId)
	if UserInfo.UserId != nil && UserInfo.StripeAccount == nil {
		//stripe連結アカウント作成
		stripe.Key = config.Config.StripeKey
		params1 := &stripe.AccountParams{Type: stripe.String("express")}
		result1, _ := account.New(params1)

		session.Set("provisional", result1.ID)
		session.Save()
		log.Println()

		//アカウントリンク作成
		params2 := &stripe.AccountLinkParams{
			Account:    stripe.String(result1.ID),
			RefreshURL: stripe.String("http://localhost:8080/refresh-create-an-express-account"),
			ReturnURL:  stripe.String("http://localhost:8080/ok-create-an-express-account"),
			Type:       stripe.String("account_onboarding"),
		}
		result2, _ := accountlink.New(params2)

		c.Redirect(307, result2.URL)

	} else if UserInfo.UserId == nil && UserInfo.StripeAccount == nil {
		c.Redirect(302, "/loginform")
	} else {
		c.Redirect(302, "/")
	}

}

func OkCreateAnExpressAccount(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.provisional = session.Get("provisional")
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.provisional != nil && UserInfo.UserId != nil {
		//user_id取得
		userid := models.GetUserID(UserInfo.UserId)
		if models.UserIdCheck(userid) == true {
			//stripeアカウント登録
			models.AccountRegist(userid, UserInfo.provisional)
			session.Set("StripeAccount", UserInfo.provisional)
			session.Delete("provisional")
			session.Save()
			c.Redirect(302, "/")

		} else {
			session.Delete("provisional")
			session.Save()
			c.Redirect(302, "/")

		}

	} else {
		c.Redirect(302, "/")
	}
}

func RefreshCreateAnExpressAccount(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.provisional = session.Get("provisional")
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.provisional != nil && UserInfo.UserId != nil {
		session.Delete("provisional")
		session.Save()
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/")
	}
}
