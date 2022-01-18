package controllers

import (
	"log"
	"main/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SippingSuccess(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	UserInfo.StripeAccount = session.Get("StripeAccount")
	productid := c.PostForm("productid")
	if UserInfo.UserId != nil && models.CheckDeliveryStatusProductId(productid) == "なし" {
		models.InsertSipping(productid)
		log.Println("なし")
		c.Redirect(302, "/registered-items")
	} else if UserInfo.UserId != nil && models.CheckDeliveryStatusProductId(productid) == "あり" {
		log.Println("あり")
		c.Redirect(302, "/registered-items")
	} else {
		c.Redirect(302, "/loginform")
	}
}

func ArrivalSuccess(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	UserInfo.StripeAccount = session.Get("StripeAccount")
	productid := c.PostForm("productid")
	if UserInfo.UserId != nil && models.CheckDeliveryStatusProductId(productid) == "あり" {
		log.Println("あり")
		models.UpdateArrives(productid)
		c.Redirect(302, "/purchase-history")
	} else if UserInfo.UserId != nil && models.CheckDeliveryStatusProductId(productid) == "なし" {
		c.Redirect(302, "/purchase-history")

	} else {
		c.Redirect(302, "/loginform")

	}
}
