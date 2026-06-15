package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"template/database"
	"template/function"
	"template/models"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

// Function for For convert crypto to USDT in Merchant section http://localhost:3000/crypto-exchange
func CryptoExchangeView(c *fiber.Ctx) error {
	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	LoginMerchantID := s.Get("LoginMerchantID")

	assetList := []models.CoinWithBalance{}
	var totalWallet int64
	// fetch query for wallet with balance
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("client_id = ? AND status = ?", LoginMerchantID, "Success").Group("assetid").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList).Count(&totalWallet)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	return c.Render("crypto-exchange", fiber.Map{
		"Title":           "Crypto Exchange",
		"Subtitle":        "Crypto Exchange",
		"MerchantData":    merchantData,
		"CoinBalanceList": assetList,
		"CoinList":        coinList,
	})
}

// function for Post Crypto Exchange Data
func CryptoExchangePost(c *fiber.Ctx) error {
	// check session
	MerchantSession(c) // redirect when session not found
	//s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")

	//LoginMerchantID := s.Get("LoginMerchantID")

	price_amount := c.FormValue("price_amount")
	price_currency := c.FormValue("price_currency")

	fmt.Print("price_amount =>> ", price_amount)
	fmt.Print("price_currency =>> ", price_currency)

	return c.Redirect("/crypto-exchange")
}

// function for get Crypto Exchange Rate
func CryptoExchangeRate(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	if s.Get("LoginMerchantID") == nil {
		fmt.Println("Session Expired101")
		//return c.Redirect("/login", 301)
	}
	//loginMerchantEmail := s.Get("LoginMerchantEmail")

	// Parse the incoming JSON data
	var data map[string]interface{}

	// Bind the request body to the 'data' map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Failed to parse body")
	}

	// Access posted data
	price_amount := data["amount"].(string)
	price_currency := strings.ToLower(data["fromCurrency"].(string))
	toCurrency := data["toCurrency"].(string)

	if toCurrency == "polygon" {
		toCurrency = "matic-network"
	} else if toCurrency == "binance coin" {
		toCurrency = "binancecoin"
	}
	toCurrencyx := "tather," + price_currency // tocurrency , from currency

	mainCurrency := "usd" // for defaul currency

	client := resty.New()
	apiURL := "https://api.coingecko.com/api/v3/simple/price?ids=" + toCurrencyx + "&vs_currencies=" + mainCurrency

	//fmt.Print("apiURL =>> ", apiURL)

	var apiResponse function.CurrencyResponse
	_, err := client.R().SetResult(&apiResponse).Get(apiURL)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error",
			"status":  400,
		})

	}
	usdAmount, err := strconv.ParseFloat(price_amount, 64)
	if err != nil {
		fmt.Println("ENO-1010 ", err)
	}

	rate := apiResponse[price_currency][mainCurrency]
	cryptoAmount := usdAmount * rate

	return c.JSON(fiber.Map{
		"message": cryptoAmount,
		"status":  200,
	})

	///////////////////////////

}
