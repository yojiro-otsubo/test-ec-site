package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"main/app/models"
	"main/config"
	"math"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/webhook"
	csrf "github.com/utrack/gin-csrf"
	"golang.org/x/crypto/bcrypt"
)

//トップページ
func top(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	products := models.GetProductTop()
	log.Println(products)
	if UserInfo.UserId == nil {
		c.HTML(200, "top", gin.H{
			"title":     "top",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
			"products":  products,
		})
	} else {
		c.HTML(200, "top", gin.H{
			"title":     "top",
			"login":     true,
			"username":  UserInfo.UserId,
			"csrfToken": csrf.GetToken(c),
			"products":  products,
		})
	}
}

//-------------------------------------------------- MyPage --------------------------------------------------

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

//購入履歴
func purchaseHistory(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId == nil {
		c.Redirect(302, "/loginform")
	} else {
		userid := models.GetUserID(UserInfo.UserId)
		products := models.GetProductIdFromPaymentHistory(userid)
		c.HTML(200, "purchaseHistory", gin.H{
			"title":     "purchaseHistory",
			"login":     true,
			"username":  UserInfo.UserId,
			"csrfToken": csrf.GetToken(c),
			"products":  products,
		})
	}
}

//登録済商品一覧
func registeredItems(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId == nil {
		c.Redirect(302, "/")
	} else {
		userid := models.GetUserID(UserInfo.UserId)
		UserProduct := models.GetTheProductOfUserId(userid)
		SoldOutProduct := models.GetSoldOutProductOfUserId(userid)
		c.HTML(200, "registeredItems", gin.H{
			"title":          "registeredItems",
			"login":          true,
			"username":       UserInfo.UserId,
			"csrfToken":      csrf.GetToken(c),
			"products":       UserProduct,
			"SoldOutProduct": SoldOutProduct,
		})
	}
}

//カート
func AddCart(c *gin.Context) {
	productid := c.PostForm("cart")
	log.Println("product = ", productid)
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	if UserInfo.UserId != nil {
		userid := models.GetUserID(UserInfo.UserId)
		models.AddToCart(userid, productid)
		redirecturl := "/product/" + productid
		c.Redirect(302, redirecturl)
	} else {
		c.Redirect(302, "/loginform")
	}
}

func CartPage(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	userid := models.GetUserID(UserInfo.UserId)
	products := models.GetProductFromCartDB(userid, 1.1)
	log.Println(products)

	var totalAmount int
	for _, p := range products {
		i, err := strconv.Atoi(p.Amount)
		if err != nil {
			log.Println(err)
		}
		totalAmount = totalAmount + i
	}
	log.Println(totalAmount)
	if UserInfo.UserId != nil {
		c.HTML(200, "cart", gin.H{
			"title":       "cart",
			"login":       true,
			"products":    products,
			"username":    UserInfo.UserId,
			"csrfToken":   csrf.GetToken(c),
			"totalAmount": totalAmount,
		})
	} else {
		c.Redirect(302, "/loginform")
	}
}

func DeleteItemsInCart(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId != nil {
		userid := models.GetUserID(UserInfo.UserId)
		productid := c.PostForm("delete_item")
		models.DeleteCartItem(userid, productid)
		c.Redirect(302, "/mycart")
	} else {
		c.Redirect(302, "/")
	}
}

//-------------------------------------------------- AUTH --------------------------------------------------
type SessionInfo struct {
	UserId        interface{}
	StripeAccount interface{}
	provisional   interface{}
}

var UserInfo SessionInfo

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
		var check bool
		userid := models.GetUserID(username)
		stripeid, check := models.GetStripeAccountId(userid)
		if check == true {
			session := sessions.Default(c)
			session.Set("UserId", username)
			session.Set("StripeAccount", stripeid)
			session.Save()
			log.Println("username = ", session.Get("UserId"), "///stripeid = ", session.Get("StripeAccount"))
			c.Redirect(302, "/")
		} else {
			session := sessions.Default(c)
			session.Set("UserId", username)
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
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId != nil {
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
		UserInfo.UserId = session.Get(("UserID"))
		if UserInfo.UserId == nil {
			log.Println("ログインしていません")
			c.HTML(http.StatusMovedPermanently, "loginform", gin.H{
				"login_massage": "ログインが必要です",
				"csrfToken":     csrf.GetToken(c),
			})
			c.Abort()
		} else {
			c.Set("UserID", UserInfo.UserId)
			c.Next()
		}
		log.Println("ログインチェック終了")
	}
}

//ログインフォーム
func LoginForm(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId == nil {
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
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId == nil {
		c.HTML(200, "signupform", gin.H{
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

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

//-------------------------------------------------- Regist Product --------------------------------------------------

//商品登録フォーム
func SellItemsForm(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	UserInfo.StripeAccount = session.Get("StripeAccount")

	if UserInfo.UserId != nil && UserInfo.StripeAccount != nil {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     true,
			"username":  UserInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	} else if UserInfo.UserId != nil && UserInfo.StripeAccount == nil {
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/loginform")
	}
}

func ItemRegist(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	UserInfo.StripeAccount = session.Get("StripeAccount")

	if UserInfo.UserId != nil && models.CheckStripeAccountId(UserInfo.StripeAccount) == true {
		var err error
		//post data
		item := c.PostForm("itemname")
		description := c.PostForm("item-description")
		amount := c.PostForm("price")
		amountInt, _ := strconv.Atoi(amount)
		amountInt64, _ := strconv.ParseInt(amount, 10, 64)

		file, header, err := c.Request.FormFile("image")
		if err != nil {
			log.Println(err)
			c.Redirect(302, "/test")
			return
		}

		filename := header.Filename
		log.Println("filename =", filename)
		pos := strings.LastIndex(filename, ".")
		log.Println("pos =", filename[pos:])

		//get user id
		userid := models.GetUserID(UserInfo.UserId)
		//regist userid and get productid(pk)
		productid := models.RegistUserIdAndGetProductId(userid, amountInt, item, description)
		//change int to str
		strproductid := strconv.Itoa(productid)

		create_path := "app/static/img/item/productid" + strproductid + "/" + filename

		mkdir_path := "app/static/img/item/productid" + strproductid

		// mkdir
		err = os.Mkdir(mkdir_path, 0755)
		if err != nil {
			log.Println(err)
		}

		//create img file
		filepath := []string{create_path}
		out, err := os.Create(filepath[0])
		if err != nil {
			log.Println("os.create err = ", err)
		}

		//file copy
		_, err = io.Copy(out, file)
		if err != nil {
			log.Println("io.copy err = ", err)
		}

		defer out.Close()

		//newpathpng := "app/static/img/item/productid" + strproductid + "/" + "1.png"
		newpathjpg := "app/static/img/item/productid" + strproductid + "/" + "1.jpg"

		if filename[pos:] == ".png" {
			if err := os.Rename(create_path, newpathjpg); err != nil {
				log.Println(err)
			}
		} else {
			if err := os.Rename(create_path, newpathjpg); err != nil {
				log.Println(err)
			}
		}

		//create stripe product
		stripe.Key = config.Config.StripeKey
		params := &stripe.ProductParams{
			Name:        stripe.String(item),
			Description: stripe.String(description),
			Images:      stripe.StringSlice(filepath),
		}
		result, _ := product.New(params)

		//create stripe price
		params2 := &stripe.PriceParams{
			Product:    stripe.String(result.ID),
			UnitAmount: stripe.Int64(amountInt64),
			Currency:   stripe.String("jpy"),
		}
		p, _ := price.New(params2)

		//regist for productsdb
		models.RegistProduct(productid, result.ID, p.ID)

		c.Redirect(302, "/")

	} else if UserInfo.UserId != nil && UserInfo.StripeAccount == nil {
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/loginform")
	}

}

//-------------------------------------------------- Test --------------------------------------------------
func test(c *gin.Context) {

	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	array := c.PostFormArray("item")

	for i := 0; i < len(array); i++ {
		log.Println(array[i])
		log.Printf("%T", array[i])
	}
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
	if UserInfo.UserId == nil {
		c.HTML(200, "test", gin.H{
			"title":     "test",
			"login":     false,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.HTML(200, "test", gin.H{
			"title":     "test",
			"login":     true,
			"username":  UserInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	}
}

//-------------------------------------------------- ProductPage --------------------------------------------------
func ProductPage(c *gin.Context) {
	productNumber := c.Param("number")
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	product := models.GetProduct(productNumber)
	log.Println(product)
	username := models.GetUserName(product[1])
	if product[7] == "1" {
		c.Redirect(302, "/")
	}
	f, err := strconv.ParseFloat(product[6], 64)
	if err != nil {
		log.Println(err)
	}
	f = f * 1.1
	taxamount := int(math.Round(f))
	log.Println(taxamount)

	if UserInfo.UserId == nil {
		c.HTML(200, "product", gin.H{
			"title":           "product",
			"login":           false,
			"csrfToken":       csrf.GetToken(c),
			"ProductId":       product[0],
			"ProductUsername": username,
			"StripeProductId": product[2],
			"StripePriceId":   product[3],
			"ItemName":        product[4],
			"Description":     product[5],
			"Amount":          taxamount,
		})
	} else {
		c.HTML(200, "product", gin.H{
			"title":           "product",
			"login":           true,
			"username":        UserInfo.UserId,
			"csrfToken":       csrf.GetToken(c),
			"ProductId":       product[0],
			"ProductUsername": username,
			"StripeProductId": product[2],
			"StripePriceId":   product[3],
			"ItemName":        product[4],
			"Description":     product[5],
			"Amount":          taxamount,
		})
	}
}

//-------------------------------------------------- CheckOut --------------------------------------------------
type CheckoutData struct {
	ClientSecret string
}

func CheckOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")

	productid := c.PostFormArray("item")
	amount := c.PostForm("totalAmount")
	amountInt64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Println(err)
	}

	if UserInfo.UserId != nil && models.CheckStripeAccountId(UserInfo.StripeAccount) == true {
		userid := models.GetUserID(UserInfo.UserId)
		var transferGroup string

		for {
			transferGroup = ""
			transferGroup = "tg_" + RandString(25)
			if models.CheckTransferGroup(transferGroup) == true {
				break
			}
		}
		stripe.Key = config.Config.StripeKey
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(amountInt64),
			Currency: stripe.String(string(stripe.CurrencyJPY)),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			TransferGroup: stripe.String(transferGroup),
		}
		result, _ := paymentintent.New(params)

		for i := 0; i < len(productid); i++ {
			log.Println(productid[i])
			models.AddTransferGroup(userid, productid[i], transferGroup)
		}

		c.HTML(200, "checkout", gin.H{
			"ClientSecret": result.ClientSecret,
			"pk":           config.Config.PK,
		})
	} else {
		c.Redirect(302, "/loginform")
	}
}

//-------------------------------------------------- Payment Completion --------------------------------------------------
func PaymentCompletion(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	if UserInfo.UserId != nil {
		c.HTML(200, "paymentCompletion", gin.H{
			"title":     "paymentCompletion",
			"login":     true,
			"username":  UserInfo.UserId,
			"csrfToken": csrf.GetToken(c),
		})
	} else {
		c.Redirect(302, "/")
	}
}

func handleWebhook(c *gin.Context) {
	stripe.Key = config.Config.StripeKey
	const MaxBodyBytes = int64(65536)
	var w http.ResponseWriter = c.Writer
	c.Request.Body = http.MaxBytesReader(w, c.Request.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		log.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		log.Printf("⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endpointSecret := config.Config.EPS
	event, err = webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		log.Printf("Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		//var paymentIntent stripe.PaymentIntent
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			log.Printf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		productid := models.GetProductIdWithTg(paymentIntent.TransferGroup)
		for i := 0; i < len(productid); i++ {
			models.UpdataSoldOutValue(productid[i], "1")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

//-------------------------------------------------- BuyerInfo --------------------------------------------------
func BuyerInfo(c *gin.Context) {

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
	render.AddFromFiles("test", "app/views/test.html")
	render.AddFromFiles("product", "app/views/base.html", "app/views/product.html")
	render.AddFromFiles("cart", "app/views/base.html", "app/views/mypage/cart.html")
	render.AddFromFiles("checkout", "app/views/checkout.html")
	render.AddFromFiles("paymentCompletion", "app/views/base.html", "app/views/paymentCompletion.html")

	return render
}

//スタートウェブサーバー（main.goから呼び出し)

func StartWebServer() {

	r := gin.Default()
	r.HTMLRender = createMultitemplate()
	r.Static("/static", "app/static")

	store := cookie.NewStore([]byte("secret"))
	CSRFGroup := r.Group("/")
	CSRFGroup.Use(sessions.Sessions("mysession", store))
	CSRFGroup.Use(csrf.Middleware(csrf.Options{
		Secret: RandString(10),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	CSRFGroup.POST("/test", test)

	//topページ
	CSRFGroup.GET("/", top)
	//マイページ
	CSRFGroup.GET("/mypage/:username", mypage)
	//購入履歴一覧
	CSRFGroup.GET("/purchase-history", purchaseHistory)
	//登録した自分のアイテム一覧
	CSRFGroup.GET("/registered-items", registeredItems)
	//商品登録フォーム
	CSRFGroup.GET("/sell-items-form", SellItemsForm)
	CSRFGroup.POST("/itemregist", ItemRegist)
	//アカウントリンク
	CSRFGroup.GET("/create-an-express-account", CreateAnExpressAccount)
	CSRFGroup.GET("/ok-create-an-express-account", OkCreateAnExpressAccount)
	CSRFGroup.GET("/refresh-create-an-express-account", RefreshCreateAnExpressAccount)
	//商品ページ
	CSRFGroup.GET("/product/:number", ProductPage)
	//購入処理
	CSRFGroup.POST("/checkout", CheckOutHandler)
	//支払い完了
	CSRFGroup.GET("/payment-completion", PaymentCompletion)
	r.POST("/webhook", handleWebhook)
	//カート
	CSRFGroup.POST("/addcart", AddCart)
	CSRFGroup.GET("/mycart", CartPage)
	CSRFGroup.POST("/delete-cart", DeleteItemsInCart)
	//購入者情報
	CSRFGroup.POST("/buyer/:number", BuyerInfo)

	//ログインフォーム
	CSRFGroup.GET("/loginform", LoginForm)

	//ログイン処理
	CSRFGroup.POST("/login", Login)

	//ユーザー登録フォーム
	CSRFGroup.GET("/signupform", SignupForm)

	//ユーザー登録処理
	CSRFGroup.POST("/registration", registration)

	//ログアウト処理
	CSRFGroup.GET("/logout", Logout)

	//RUNサーバー
	r.Run(fmt.Sprintf(":%d", config.Config.Port))

}
