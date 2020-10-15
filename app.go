package main

import (
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v71/charge"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stripe/stripe-go/v71"
	billingPortalSession "github.com/stripe/stripe-go/v71/billingportal/session"
	checkoutSession "github.com/stripe/stripe-go/v71/checkout/session"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/create-checkout-session", createCheckoutSession)
	e.POST("/create-portal-session", createBillingPortalSession)
	e.GET("/charge-test", chargeTest)

	e.Logger.Fatal(e.Start("localhost:4242"))
}

type CreateCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func createCheckoutSession(c echo.Context) (err error) {
	productName := os.Getenv("PRODUCT_NAME")
	productPrice, err := strconv.Atoi(os.Getenv("PRODUCT_PRICE"))

	if err != nil {
		panic(err)
	}

	productCurrency := os.Getenv("PRODUCT_CURRENCY")

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(os.Getenv("CUSTOMER_ID")),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(productCurrency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(int64(productPrice)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:3000/success"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
	}

	s, e := checkoutSession.New(params)

	if e != nil {
		return e
	}

	data := CreateCheckoutSessionResponse{
		SessionID: s.ID,
	}

	return c.JSON(http.StatusOK, data)
}

type CreateBillingPortalSessionResponse struct {
	Url string `json:"url"`
}

func createBillingPortalSession(c echo.Context) error {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(os.Getenv("CUSTOMER_ID")),
		ReturnURL: stripe.String("http://localhost:3000"),
	}
	s, err := billingPortalSession.New(params)

	if err != nil {
		return err
	}

	data := CreateBillingPortalSessionResponse{
		Url: s.URL,
	}

	return c.JSON(http.StatusOK, data)
}

func chargeTest(c echo.Context) error {
	ch, err := charge.Get(os.Getenv("CHARGE_ID"), nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ch)
}
