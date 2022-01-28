package controllers

import (
	"io"
	"log"
	"main/app/models"
	"main/config"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	csrf "github.com/utrack/gin-csrf"
)

//-------------------------------------------------- Regist Product --------------------------------------------------

//商品登録フォーム
func SellItemsForm(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	userid := models.GetUserID(UserInfo.UserName)
	stripeid, _ := models.GetStripeAccountId(userid)

	if loginbool == true && UserInfo.StripeAccount == stripeid && models.ReturnPersonalUserIdCheck(userid) == "あり" {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     true,
			"username":  UserInfo.UserName,
			"csrfToken": csrf.GetToken(c),
			"stripeid":  true,
			"personal":  true,
		})
	} else if loginbool == true && UserInfo.StripeAccount != stripeid && models.ReturnPersonalUserIdCheck(userid) == "あり" {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     true,
			"username":  UserInfo.UserName,
			"csrfToken": csrf.GetToken(c),
			"stripeid":  false,
			"personal":  true,
		})
	} else if loginbool == true && UserInfo.StripeAccount == stripeid && models.ReturnPersonalUserIdCheck(userid) == "なし" {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     true,
			"username":  UserInfo.UserName,
			"csrfToken": csrf.GetToken(c),
			"stripeid":  true,
			"personal":  false,
		})
	} else if loginbool == true && UserInfo.StripeAccount != stripeid && models.ReturnPersonalUserIdCheck(userid) == "なし" {
		c.HTML(200, "SellItems", gin.H{
			"title":     "SellItems",
			"login":     true,
			"username":  UserInfo.UserName,
			"csrfToken": csrf.GetToken(c),
			"stripeid":  false,
			"personal":  false,
		})
	} else {
		c.Redirect(302, "loginform")
	}
}

func ItemRegist(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	if loginbool == true && models.CheckStripeAccountId(UserInfo.StripeAccount) == true {
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
		userid := models.GetUserID(UserInfo.UserName)
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

	} else if loginbool == true && UserInfo.StripeAccount == nil {
		c.Redirect(302, "/")
	} else {
		c.Redirect(302, "/loginform")
	}

}

//登録済商品一覧
func registeredItems(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)

	if loginbool == false {
		c.Redirect(302, "/")
	} else {
		userid := models.GetUserID(UserInfo.UserName)
		UserProduct := models.GetTheProductOfUserId(userid)
		SoldOutProduct := models.GetSoldOutProductOfUserId(userid)
		SippingOkProduct := models.GetSippingOkProductOfUserId(userid)
		c.HTML(200, "registeredItems", gin.H{
			"title":            "registeredItems",
			"login":            true,
			"username":         UserInfo.UserName,
			"csrfToken":        csrf.GetToken(c),
			"products":         UserProduct,
			"SoldOutProduct":   SoldOutProduct,
			"SippingOkProduct": SippingOkProduct,
		})
	}
}
func ItemDelete(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserName = session.Get("UserName")
	UserInfo.StripeAccount = session.Get("StripeAccount")
	UserInfo.logintoken = session.Get("logintoken")
	loginbool := models.LoginTokenCheck(UserInfo.UserName, UserInfo.logintoken)
	productid := c.PostForm("productid")
	userid := models.GetUserID(UserInfo.UserName)
	if loginbool == true && models.CheckDeliveryStatusProductId(productid) == "なし" {
		models.DeleteProduct(userid, productid)
		c.Redirect(302, "/registered-items")
	}
	c.Redirect(302, "/")
}
