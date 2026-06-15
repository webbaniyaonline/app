package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
)

// Function for Now Payment API

var apiKeyNP = "Y3WR3PA-TG04W8G-HDTMVZG-Z3PCWYD"
var apiPath = "https://api-sandbox.nowpayments.io"
var callbackURL = "https://itio.in/nowpayments/callback.php"
var successURL = "https://itio.in/nowpayments/responce.php"
var failedURL = "https://itio.in/nowpayments/responce.php"

// Function for display Coin List
func AddCryptoView(c *fiber.Ctx) error {
	//VID := c.Params("VID")

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	//voltID := s.Get("LoginVoltID")
	//fmt.Println("LoginVoltID -> ", voltID, LoginMerchantID)
	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("3434343=>>", err)
	}
	if LoginMerchantID == nil {
		fmt.Println("Session Expired106")
		return c.Redirect("/login")
	}

	return c.Render("add-crypto-np", fiber.Map{
		"Title":        "Add Crypto",
		"Subtitle":     "Add Crypto",
		"Alert":        Alerts,
		"MerchantData": merchantData,
	})
}

// Function for Generate Coin
func AddCryptoPost(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	Alerts := "Transfer Process"
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	//voltID := s.Get("LoginVoltID").(string)
	//LoginMerchantEmail := s.Get("LoginMerchantEmail").(string)
	if LoginMerchantID == nil {
		fmt.Println("Session Expired105")
		return c.Redirect("/login")
	}

	randomID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	// Convert string to uint
	price_currency := c.FormValue("price_currency")
	price_amount := c.FormValue("price_amount")
	pay_currency := c.FormValue("pay_currency")
	order_description := c.FormValue("order_description")

	getAmount, err := strconv.ParseFloat(price_amount, 64)
	if err != nil {
		fmt.Println(err)
	}
	// Convert uint64 to uint
	//AID := uint(number)

	MyData := models.TransferRequestNP{
		Price_amount:      getAmount,
		Price_currency:    price_currency,
		Pay_currency:      pay_currency,
		Ipn_callback_url:  callbackURL,
		Order_id:          randomID,
		Order_description: order_description,
	}

	path := "/v1/payment"
	//tokenProvider := "NULL"
	respBody, err := MakeAPIRequestNP("POST", path, MyData)
	if err != nil {
		fmt.Println(err)
		Alerts = "Transaction Failed"
	}
	fmt.Println("RRR - >", string(respBody))
	// Parse the JSON data into the struct
	var npData models.TransResponceNP
	if err := json.Unmarshal([]byte(respBody), &npData); err != nil {
		fmt.Println(err)
	}
	fmt.Println(npData)

	if npData.Payment_id != "" {

		TransID := npData.Payment_id
		TransStatus := npData.Payment_status
		Message := "Transaction Processed with Status - " + TransStatus + " And Trans Id - " + TransID
		Alerts = Message
		//CID := int(LoginMerchantID)
		CID := s.Get("LoginMerchantID").(uint)
		//fmt.Println("=====>>>", CID)
		Ip := c.Context().RemoteIP().String()

		qry := models.Transaction_MasterNP{Client_id: CID, Payment_id: TransID, Payment_status: TransStatus, Pay_address: npData.Pay_address, Price_amount: npData.Price_amount, Pay_amount: npData.Pay_amount, Amount_received: npData.Amount_received, Price_currency: npData.Price_currency, Pay_currency: npData.Pay_currency, Order_id: npData.Order_id, Ip: Ip, Order_description: npData.Order_description, Created_at: npData.Created_at, Updated_at: npData.Updated_at}
		result := database.DB.Db.Table("transaction_master_nowpayments").Select("client_id", "payment_id", "payment_status", "pay_address", "price_amount", "pay_amount", "amount_received", "price_currency", "pay_currency", "order_id", "order_description", "created_at", "updated_at", "ip").Create(&qry)
		//fmt.Println(result)

		if result.Error != nil {
			fmt.Println(result.Error)
		}

	} else {

		Message := "Transaction Failed with Status - " + npData.Code + " And Message - " + npData.Message
		Alerts = Message
	}
	///////////////////////==================
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}
	return c.Redirect("/transactions-np")
}

// Function for view coin Request Form
func RequestCryptoView(c *fiber.Ctx) error {
	//VID := c.Params("VID")

	// check session
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	//voltID := s.Get("LoginVoltID")
	//fmt.Println("LoginVoltID -> ", voltID, LoginMerchantID)
	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("3434343=>>", err)
	}
	if LoginMerchantID == nil {
		fmt.Println("Session Expired107")
		return c.Redirect("/login")
	}

	return c.Render("request-crypto-np", fiber.Map{
		"Title":        "Request Crypto",
		"Subtitle":     "Request Crypto",
		"Alert":        Alerts,
		"MerchantData": merchantData,
	})
}

// Function for Submit coin Request Form
func RequestCryptoPost(c *fiber.Ctx) error {
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")

	Alerts := "Transfer Process"

	randomID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	// Convert string to uint
	price_currency := c.FormValue("price_currency")
	price_amount := c.FormValue("price_amount")
	order_description := c.FormValue("order_description")
	sender_name := c.FormValue("sender_name")
	sender_email := c.FormValue("sender_email")

	MyData := models.TransferRequestCrypto{
		Price_amount:      price_amount,
		Price_currency:    price_currency,
		Order_id:          randomID,
		Order_description: order_description,
		Ipn_callback_url:  callbackURL,
		Success_url:       successURL,
		Cancel_url:        failedURL,
	}

	path := "/v1/invoice"

	//tokenProvider := "NULL"
	respBody, err := MakeAPIRequestNP("POST", path, MyData)
	if err != nil {
		fmt.Println(err)
		Alerts = "Transaction Failed"
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var npData models.TransResponcecryptoNP
	if err := json.Unmarshal([]byte(respBody), &npData); err != nil {
		fmt.Println(err)
	}
	//fmt.Println(npData)
	//return c.Redirect("/transactions-np")

	if npData.Order_id != "" {

		TransID := npData.Order_id
		TransStatus := npData.Payment_status
		Message := "Transaction Processed with Status - " + TransStatus + " And Order ID - " + TransID
		Alerts = Message
		//CID := int(LoginMerchantID)
		CID := s.Get("LoginMerchantID").(uint)
		//fmt.Println("=====>>>", CID)
		Ip := c.Context().RemoteIP().String()
		payment_status := "waiting"

		getAmount, err := strconv.ParseFloat(npData.Price_amount, 64)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(" AMT2-> ", getAmount)
		qry := models.Transaction_MasterNP{Client_id: CID, Price_amount: getAmount, Price_currency: npData.Price_currency, Order_id: npData.Order_id, Ip: Ip, Order_description: npData.Order_description, Created_at: npData.Created_at, Updated_at: npData.Updated_at, Payment_status: payment_status, Invoice_id: npData.Id, Token_id: npData.Token_id, Invoice_url: npData.Invoice_url, Request_json: string(respBody)}
		result := database.DB.Db.Table("transaction_master_nowpayments").Select("client_id", "price_amount", "price_currency", "order_id", "order_description", "created_at", "updated_at", "ip", "payment_status", "invoice_id", "token_id", "invoice_url", "request_json").Create(&qry)
		//fmt.Println(result)

		if result.Error != nil {
			fmt.Println(result.Error)
		}
		//fmt.Println(" Email ", sender_email, sender_name, npData.Invoice_url, npData.Id)

		//  Email///
		template_code := "REQUEST-MONEY"

		emailData := models.EmailData{FullName: sender_name, Email: sender_email, Invoice_id: TransID, Invoice_url: npData.Invoice_url}

		err = function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		} else {
			fmt.Println("Mail Going")
		}

	} else {

		Message := "Transaction Failed with Status - " + npData.Code + " And Message - " + npData.Message
		Alerts = Message
	}
	///////////////////////==================
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}
	return c.Redirect("/transactions-np")
}

// Function for view Transaction List
func TransactionsNPView(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	// Get query parameters for page and limit
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
	}
	pageLimit := "10"
	limit, err := strconv.Atoi(c.Query("limit", pageLimit))
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid limit number")
	}

	// Calculate offset
	offset := (page - 1) * limit

	transactionList := []models.Transaction_MasterNP{}
	database.DB.Db.Table("transaction_master_nowpayments").Order("tid desc").Limit(limit).Offset(offset).Where("client_id = ?", LoginMerchantID).Find(&transactionList)

	var total int64
	database.DB.Db.Table("transaction_master_nowpayments").Count(&total)

	//fmt.Println(total)

	// Prepare pagination data
	totalPage := total / 10
	//fmt.Println(totalPage)
	nextPage := page + 1
	prevPage := page - 1
	if page == 1 {
		prevPage = 0
	}

	if page >= int(totalPage+1) {
		nextPage = 0
	}

	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		fmt.Println("3434343=>>", err)
	}

	//fmt.Println(transactionList)
	return c.Render("transactions-np", fiber.Map{
		"Title":           "Transactions",
		"Subtitle":        "Transactions",
		"Alert":           Alerts,
		"TransactionList": transactionList,
		"MerchantData":    merchantData,
		"NextPage":        nextPage,
		"PrevPage":        prevPage,
		"Limit":           limit,
		"Count":           total,
	})
}

var httpClient = &http.Client{} // Reuse HTTP client

// Function for execute Api Request
func MakeAPIRequestNP(method, path string, body interface{}) ([]byte, error) {
	var url = apiPath + path

	var reqBodyBytes []byte
	if body != nil {
		var err error
		reqBodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	//req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-API-KEY", apiKeyNP)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return respBody, nil
}
