package controllers

import (
	"fmt"
	"main/config"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

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
	CSRFGroup.GET("/test", test)

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

	//発送
	//CSRFGroup.POST("/sipping-success", SippingSuccess)
	//顧客情報
	//CSRFGroup.POST("/buyer-information", BuyerInformation)
	//カート
	CSRFGroup.POST("/addcart", AddCart)
	CSRFGroup.GET("/mycart", CartPage)
	CSRFGroup.POST("/delete-cart", DeleteItemsInCart)

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
