package controllers

import (
	"fmt"
	"log"
	"main/app/models"
	"main/config"
	"net/http"
	"net/mail"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

//トップページ
func top(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")
	if LoginInfo.UserId == nil {
		c.HTML(200, "top", gin.H{
			"title":     "top",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "top", gin.H{
			"title":     "top",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

//-------------------------------------------------- MyPage --------------------------------------------------

//マイページ
func mypage(c *gin.Context) {
	uname := c.Param("username")
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")

	if LoginInfo.UserId == uname {
		c.HTML(200, "mypage", gin.H{
			"title":    uname,
			"login":    true,
			"username": session.Get("UserId"),
		})
	} else {
		c.Redirect(302, "/")
	}
}

//購入履歴
func purchaseHistory(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")

	if LoginInfo.UserId == nil {
		c.HTML(200, "purchaseHistory", gin.H{
			"title":     "purchaseHistory",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "purchaseHistory", gin.H{
			"title":     "purchaseHistory",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

//登録済商品一覧
func registeredItems(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")

	if LoginInfo.UserId == nil {
		c.HTML(200, "registeredItems", gin.H{
			"title":     "registeredItems",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "registeredItems", gin.H{
			"title":     "registeredItems",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

//商品登録フォーム
func SellItemsForm(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")

	if LoginInfo.UserId == nil {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

func CreateUpFile(username interface{}) {
}

//商品登録処理
func ItemRegistrationProcess(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")
	item_name := c.PostForm("itemname")
	item_description := c.PostForm("item-description")
	price := c.PostForm("price")

	models.ItemRegistrationDB(LoginInfo.UserId, item_name, item_description, price)
	//CreateUpFile(LoginInfo.UserId, item_name)

	if LoginInfo.UserId != nil {

	} else {
		c.Redirect(302, "/")
	}
}

//お客様情報入力フォーム
func UserDetailedInformationForm(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")

	if LoginInfo.UserId == nil {
		c.HTML(200, "UserDetailedInformation", gin.H{
			"title":     "UserDetailedInformation",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "UserDetailedInformation", gin.H{
			"title":     "UserDetailedInformation",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

//支払い方法入力フォーム
func PaymentInfoFrom(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")

	if LoginInfo.UserId == nil {
		c.HTML(200, "PaymentInfo", gin.H{
			"title":     "PaymentInfo",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "PaymentInfo", gin.H{
			"title":     "PaymentInfo",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

//-------------------------------------------------- AUTH --------------------------------------------------
type SessionInfo struct {
	UserId interface{}
}

var LoginInfo SessionInfo

//ユーザー登録
func registration(c *gin.Context) {
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
		models.UserRegistration(username, e.Address, string(hash))
		c.Redirect(302, "/")
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
		session := sessions.Default(c)
		session.Set("UserId", username)
		session.Save()
		c.Redirect(302, "/")
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
	LoginInfo.UserId = session.Get("UserId")
	if LoginInfo.UserId != nil {
		session.Clear()
		session.Save()
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/")
	}

}

//セッションチェック関数（仮）
func sessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		LoginInfo.UserId = session.Get(("UserID"))
		if LoginInfo.UserId == nil {
			log.Println("ログインしていません")
			c.HTML(http.StatusMovedPermanently, "loginform", gin.H{
				"login_massage": "ログインが必要です",
				"csrfToken":     csrf.GetToken(c),
			})
			c.Abort()
		} else {
			c.Set("UserID", LoginInfo.UserId)
			c.Next()
		}
		log.Println("ログインチェック終了")
	}
}

func LoginForm(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")
	if LoginInfo.UserId == nil {
		c.HTML(200, "loginform", gin.H{
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

func SignupForm(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")
	if LoginInfo.UserId == nil {
		c.HTML(200, "signupform", gin.H{
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

//-------------------------------------------------- WebServer --------------------------------------------------

//マルチテンプレート作成
func createMultitemplate() multitemplate.Renderer {
	render := multitemplate.NewRenderer()
	render.AddFromFiles("top", "app/views/base.html", "app/views/top.html")
	render.AddFromFiles("loginform", "app/views/base.html", "app/views/loginForm.html")
	render.AddFromFiles("signupform", "app/views/base.html", "app/views/signupForm.html")
	render.AddFromFiles("mypage", "app/views/base.html", "app/views/mypage/mypage.html")
	render.AddFromFiles("purchaseHistory", "app/views/base.html", "app/views/mypage/purchaseHistory.html")
	render.AddFromFiles("registeredItems", "app/views/base.html", "app/views/mypage/RegisteredItems.html")
	render.AddFromFiles("SellItems", "app/views/base.html", "app/views/mypage/Sellitem.html")
	render.AddFromFiles("UserDetailedInformation", "app/views/base.html", "app/views/mypage/UserDetailedInfoForm.html")
	render.AddFromFiles("PaymentInfo", "app/views/base.html", "app/views/mypage/PaymentInfo.html")
	render.AddFromFiles("test", "app/views/base.html", "app/views/mypage/test.html")

	return render
}

//スタートウェブサーバー（main.goから呼び出し)
func StartWebServer() {

	r := gin.Default()
	r.HTMLRender = createMultitemplate()
	r.Static("/static", "app/static")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(csrf.Middleware(csrf.Options{
		Secret: RandString(10),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	r.GET("/test", test)

	//topページ
	r.GET("/", top)
	//マイページ
	r.GET("/mypage/:username", mypage)
	r.GET("/purchase-history", purchaseHistory)
	r.GET("/registered-items", registeredItems)
	//商品登録フォーム
	r.GET("/sell-items-form", SellItemsForm)
	//商品登録処理
	r.POST("item-registration-process", ItemRegistrationProcess)
	r.GET("/user-detailed-information", UserDetailedInformationForm)
	r.GET("/payment-info", PaymentInfoFrom)

	//stripe処理
	r.GET("/create-an-express-account", CreateAnExpressAccount)

	//ログインフォーム
	r.GET("/loginform", LoginForm)

	//ログイン処理
	r.POST("/login", Login)

	//ユーザー登録フォーム
	r.GET("/signupform", SignupForm)

	//ユーザー登録処理
	r.POST("/registration", registration)

	//ログアウト処理
	r.GET("/logout", Logout)

	//RUNサーバー
	r.Run(fmt.Sprintf(":%d", config.Config.Port))

}

func CreateAnExpressAccount(c *gin.Context) {
	stripe.Key = config.Config.StripeKey
	params1 := &stripe.AccountParams{Type: stripe.String("express")}
	result1, _ := account.New(params1)
	log.Fatalln(result1.ID)

	params2 := &stripe.AccountLinkParams{
		Account:    stripe.String(result1.ID),
		RefreshURL: stripe.String("http://localhost:8080/"),
		ReturnURL:  stripe.String("http://localhost:8080/test"),
		Type:       stripe.String("account_onboarding"),
	}
	result2, _ := accountlink.New(params2)

	log.Fatalln(result2.URL)

	c.Redirect(302, result2.URL)
}

func test(c *gin.Context) {
	session := sessions.Default(c)
	LoginInfo.UserId = session.Get("UserId")
	if LoginInfo.UserId == nil {
		c.HTML(200, "test", gin.H{
			"title":     "test",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "test", gin.H{
			"title":     "test",
			"login":     true,
			"username":  LoginInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}
