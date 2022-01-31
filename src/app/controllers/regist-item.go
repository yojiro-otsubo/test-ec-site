package controllers

import (
	"io"
	"log"
	"main/app/models"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func CreateImage(file multipart.File, header *multipart.FileHeader, image_number, strproductid string) {
	var err error
	filename := header.Filename
	log.Println("filename =", filename)
	pos := strings.LastIndex(filename, ".")
	log.Println("pos =", filename[pos:])

	create_path := "app/static/img/item/productid" + strproductid + "/" + filename

	mkdir_path := "app/static/img/item/productid" + strproductid

	// mkdir
	err = os.Mkdir(mkdir_path, 0755)
	if err != nil {
		log.Println(err)
	}

	//create img file
	//filepath := []string{create_path}
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

	//newpathpng := "app/static/img/item/productid" + strproductid + "/" + "1.png"
	newpathjpg := "app/static/img/item/productid" + strproductid + "/" + image_number + ".jpg"

	if filename[pos:] == ".png" {
		if err := os.Rename(create_path, newpathjpg); err != nil {
			log.Println(err)
		}
	} else {
		if err := os.Rename(create_path, newpathjpg); err != nil {
			log.Println(err)
		}
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
		//amountInt64, _ := strconv.ParseInt(amount, 10, 64)
		//get user id
		userid := models.GetUserID(UserInfo.UserName)
		//regist userid and get productid(pk)
		productid := models.RegistUserIdAndGetProductId(userid, amountInt, item, description)
		//change int to str
		strproductid := strconv.Itoa(productid)

		file, header, err := c.Request.FormFile("image")
		if err != nil {
			log.Println(err)
			c.Redirect(302, "/sell-items-form")
			return
		}
		CreateImage(file, header, "1", strproductid)
		file2, header2, err := c.Request.FormFile("image2")
		if err != nil {
			log.Println(err)
			c.Redirect(302, "/sell-items-form")
			return
		}
		CreateImage(file2, header2, "2", strproductid)
		file3, header3, err := c.Request.FormFile("image3")
		if err != nil {
			log.Println(err)
			c.Redirect(302, "/sell-items-form")
			return
		}
		CreateImage(file3, header3, "3", strproductid)
		/*
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

			//regist for productsdb*/
		models.RegistProduct(productid, "null", "null")

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
