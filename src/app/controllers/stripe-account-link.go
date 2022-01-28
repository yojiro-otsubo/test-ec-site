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
	"github.com/stripe/stripe-go/v72/customer"
)

//-------------------------------------------------- StripeAccountLink --------------------------------------------------
//stripe連結アカウント作成、アカウント登録リンクへリダイレクト
func CreateAnExpressAccount(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.UserName = session.Get("UserName")
	log.Println("Username = ", UserInfo.UserName)
	if UserInfo.UserName != nil && UserInfo.StripeAccount == nil {
		//stripe連結アカウント作成
		stripe.Key = config.Config.StripeKey
		params1 := &stripe.AccountParams{
			Country:      stripe.String("JP"),
			Type:         stripe.String("express"),
			BusinessType: stripe.String("individual"),
			BusinessProfile: &stripe.AccountBusinessProfileParams{
				URL: stripe.String("https://www.google.com/"),
			},
		}
		result1, _ := account.New(params1)

		session.Set("provisional", result1.ID)
		session.Save()
		log.Println(result1.ID)

		//アカウントリンク作成
		params2 := &stripe.AccountLinkParams{
			Account:    stripe.String(result1.ID),
			RefreshURL: stripe.String("http://localhost:8080/refresh-create-an-express-account"),
			ReturnURL:  stripe.String("http://localhost:8080/ok-create-an-express-account"),
			Type:       stripe.String("account_onboarding"),
		}
		result2, _ := accountlink.New(params2)

		c.Redirect(307, result2.URL)

	} else if UserInfo.UserName == nil {
		c.Redirect(302, "/loginform")
	} else {
		c.Redirect(302, "/")
	}

}

func OkCreateAnExpressAccount(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.provisional = session.Get("provisional")
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.provisional != nil && UserInfo.UserName != nil {
		//user_id取得
		userid := models.GetUserID(UserInfo.UserName)
		if models.UserIdCheck(userid) == true {
			//stripeアカウント登録
			models.AccountRegist(userid, UserInfo.provisional)
			stripeid, _ := models.GetStripeAccountId(userid)
			cp := &stripe.CustomerParams{}
			cp.SetStripeAccount(stripeid)
			customer.New(cp)

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
	UserInfo.UserName = session.Get("UserName")
	if UserInfo.provisional != nil && UserInfo.UserName != nil {
		session.Delete("provisional")
		session.Save()
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/")
	}
}
