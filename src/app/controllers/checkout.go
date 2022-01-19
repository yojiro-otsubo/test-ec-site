package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/app/models"
	"main/config"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"
	csrf "github.com/utrack/gin-csrf"
)

//-------------------------------------------------- CheckOut --------------------------------------------------
type CheckoutData struct {
	ClientSecret string
}

func CheckOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	UserInfo.UserId = session.Get("UserId")
	UserInfo.StripeAccount = session.Get("StripeAccount")

	productid := c.PostFormArray("item")
	for i := 0; i < len(productid); i++ {
		product := models.GetProduct(productid[i])
		if product[7] == "1" {
			c.Redirect(302, "/")
		}
	}

	amount := c.PostForm("totalAmount")
	amountInt64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Println(err)
	}
	userid := models.GetUserID(UserInfo.UserId)
	if UserInfo.UserId != nil && models.PersonalUserIdCheck(userid) == "あり" {

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

	} else if UserInfo.UserId != nil && models.PersonalUserIdCheck(userid) == "なし" {
		c.Redirect(302, "/personal-information-input")
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
		log.Println(paymentIntent.TransferGroup)
		productid := models.GetProductIdWithTg(paymentIntent.TransferGroup)
		for i := 0; i < len(productid); i++ {
			models.UpdataSoldOutValue(productid[i], "1")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
