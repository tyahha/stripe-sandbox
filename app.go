package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
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

	s, e := session.New(params)

	if e != nil {
		return e
	}

	data := CreateCheckoutSessionResponse{
		SessionID: s.ID,
	}

	return c.JSON(http.StatusOK, data)
}
