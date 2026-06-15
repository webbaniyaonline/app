package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Function for Display Payment form in merchant section
func PayView(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)

	merchantData := s.Get("MerchantData")

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	// For display Currency List on List Box
	currencyList := []models.CurrencyList{}
	database.DB.Db.Table("currency").Select("currency_code").Order("currency_id ASC").Where("status = ?", 1).Find(&currencyList)
	//fmt.Println(currencyList)

	return c.Render("merchant-payment", fiber.Map{
		"Title":        "Pay",
		"Subtitle":     "Pay",
		"Action":       "Add",
		"CoinList":     coinList,
		"CurrencyList": currencyList,
		"MerchantData": merchantData,
	})
}

// Function for Display Transaction Listing in merchant section
func TransactionsView(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	// Get query parameters
	// Get query parameters
	searchKey := strings.TrimSpace(c.Query("searchkey", ""))
	searchBy := c.Query("searchby", "transaction_id")
	status := c.Query("status", "")
	date_1st := c.Query("date_1st", "")
	date_2nd := c.Query("date_2nd", "")
	time_period := c.Query("time_period", "")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	searchString := ""
	searchStringFull := ""

	if searchKey != "" && searchBy != "" {
		searchString = " AND " + searchBy + " ILIKE " + "'%" + searchKey + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if status != "" {
		searchString = " AND substatus = " + status
		searchStringFull = searchStringFull + "" + searchString
	}

	if date_1st != "" && date_2nd != "" {
		searchString = " AND createdate >= " + "'" + date_1st + "' AND createdate <= " + "'" + date_2nd + "'"
		searchStringFull = searchStringFull + "" + searchString
	}
	if LoginMerchantID != "" {
		searchString = " AND client_id =  ? "
		searchStringFull = searchStringFull + "" + searchString
	}

	if len(searchStringFull) > 4 {
		searchStringFull = searchStringFull[4:]
	}

	transactionList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where(searchStringFull, LoginMerchantID).Limit(limit).Offset(offset).Find(&transactionList)

	var total int64
	database.DB.Db.Table("transaction").Where(searchStringFull, LoginMerchantID).Count(&total)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		fmt.Println(err)
	}

	return c.Render("merchant-transactions", fiber.Map{
		"Title":           "Transactions List",
		"Subtitle":        "Transactions List",
		"Alert":           Alerts,
		"TransactionList": transactionList,
		"CoinList":        coinList,
		"MerchantData":    merchantData,
		"Page":            page,
		"Limit":           limit,
		"Total":           total,
		"SearchKey":       searchKey,
		"SearchBy":        searchBy,
		"Status":          status,
		"Date_1st":        date_1st,
		"Date_2nd":        date_2nd,
		"Time_period":     time_period,
	})
}

// Function for Display Requested Payment in merchant section
func RequestedPaymentViews(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	transactionList := []models.Invoice_Master{}
	database.DB.Db.Table("invoice").Order("invoice_id desc").Limit(limit).Offset(offset).Where("client_id = ? AND invoice_type = ?", LoginMerchantID, 2).Find(&transactionList)

	var total int64
	database.DB.Db.Table("invoice").Where("client_id = ? AND invoice_type = ?", LoginMerchantID, 2).Count(&total)

	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		fmt.Println("3434343=>>", err)
	}
	var commonURL = os.Getenv("CommonURL")
	//fmt.Println(transactionList)
	return c.Render("requested-payment", fiber.Map{
		"Title":           "Pay Request",
		"Subtitle":        "List",
		"Alert":           Alerts,
		"TransactionList": transactionList,
		"MerchantData":    merchantData,
		"CommonURL":       commonURL,
		"Page":            page,
		"Limit":           limit,
		"Total":           total,
	})
}

// Function for Display Pay Request form in merchant section
func PaymentViews(c *fiber.Ctx) error {
	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	// For display Currency List on List Box
	currencyList := []models.CurrencyList{}
	database.DB.Db.Table("currency").Order("currency_id ASC").Where("status = ?", 1).Find(&currencyList)

	return c.Render("payment-request", fiber.Map{
		"Title":        "Payment Request",
		"Subtitle":     "Send By Email",
		"Action":       "Add",
		"CurrencyList": currencyList,
		"MerchantData": merchantData,
	})
}

// Function for Post Pay Request form in merchant section
func PaymentRequestPost(c *fiber.Ctx) error {

	// Get Data from ajax
	////////////////////////////////////
	s, _ := store.Get(c)
	Alerts := ""
	CID := s.Get("LoginMerchantID").(uint)
	Ip := c.Context().RemoteIP().String()

	price_currency := c.FormValue("price_currency")

	price_amount := c.FormValue("price_amount")
	requestedamount, err := strconv.ParseFloat(price_amount, 64)
	if err != nil {
		fmt.Println(err)
	}

	sender_name := c.FormValue("sender_name")
	sender_email := c.FormValue("sender_email")
	description := c.FormValue("order_description")

	// Generate randomly Transaction ID
	trackID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	qry := models.Invoice_Master{Client_id: CID, Requestedamount: requestedamount, Requestedcurrency: price_currency, Name: sender_name, Email: sender_email, Description: description, Ip: Ip, Trackid: trackID, Invoice_type: 2}
	result := database.DB.Db.Table("invoice").Select("client_id", "requestedamount", "requestedcurrency", "name", "email", "description", "ip", "trackid", "invoice_type").Create(&qry)
	invoice_id := strconv.FormatUint(uint64(qry.Invoice_id), 10)
	//fmt.Println(invoice_id)
	if result.Error != nil {
		Alerts = "Error : Invoice Not Created "

	} else {

		//  Email///
		template_code := "INVOICE-PAYMENT"
		var commonURL = os.Getenv("CommonURL")
		invoice_url := commonURL + "/pay?iid=" + trackID
		amount := price_amount + " " + price_currency
		emailData := models.EmailData{FullName: sender_name, Email: sender_email, Invoice_id: trackID, Invoice_url: invoice_url, Amount: amount}
		err = function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		} else {
			fmt.Println("Mail Going")
		}

		Alerts = "Payment Request Generated Successfully with Pay ID :" + invoice_id
	}
	///////////////////////==================
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}

	//fmt.Print(response)

	return c.Redirect("/requested-payment")
}

// Function for Display Requested Payment in merchant section
func PayLinksListViews(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	transactionList := []models.Invoice_Master{}
	database.DB.Db.Table("invoice").Order("invoice_id desc").Limit(limit).Offset(offset).Where("client_id = ? AND invoice_type = ?", LoginMerchantID, 1).Find(&transactionList)

	var total int64
	database.DB.Db.Table("invoice").Where("client_id = ? AND invoice_type = ?", LoginMerchantID, 1).Count(&total)

	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		fmt.Println("3434343=>>", err)
	}
	var commonURL = os.Getenv("CommonURL")
	//fmt.Println(transactionList)
	return c.Render("pay-links", fiber.Map{
		"Title":           "Pay Links",
		"Subtitle":        "List",
		"Alert":           Alerts,
		"TransactionList": transactionList,
		"MerchantData":    merchantData,
		"CommonURL":       commonURL,
		"Page":            page,
		"Limit":           limit,
		"Total":           total,
	})
}

// Function for Display Pay Request form in merchant section
func PayLinksViews(c *fiber.Ctx) error {
	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	// For display Currency List on List Box
	currencyList := []models.CurrencyList{}
	database.DB.Db.Table("currency").Order("currency_id ASC").Where("status = ?", 1).Find(&currencyList)

	return c.Render("pay-link-form", fiber.Map{
		"Title":        "Pay Link",
		"Subtitle":     "Generate",
		"Action":       "Add",
		"CurrencyList": currencyList,
		"MerchantData": merchantData,
	})
}

// Function for Post Pay Request form in merchant section
func PayLinkPost(c *fiber.Ctx) error {

	// Get Data from ajax
	////////////////////////////////////
	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	Alerts := ""
	CID := s.Get("LoginMerchantID").(uint)
	Ip := c.Context().RemoteIP().String()

	price_currency := c.FormValue("price_currency")

	price_amount := c.FormValue("price_amount")
	requestedamount, err := strconv.ParseFloat(price_amount, 64)
	if err != nil {
		fmt.Println(err)
	}

	product_name := c.FormValue("product_name")
	description := c.FormValue("order_description")

	// Generate randomly Transaction ID
	trackID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}
	payLink := ""

	qry := models.Invoice_Master{Client_id: CID, Requestedamount: requestedamount, Requestedcurrency: price_currency, Product_name: product_name, Description: description, Ip: Ip, Trackid: trackID, Invoice_type: 1}
	result := database.DB.Db.Table("invoice").Select("client_id", "requestedamount", "requestedcurrency", "product_name", "description", "ip", "trackid", "invoice_type").Create(&qry)
	invoice_id := strconv.FormatUint(uint64(qry.Invoice_id), 10)
	//fmt.Println(invoice_id)
	if result.Error != nil {
		Alerts = "Error : Pay Link Not Created "

	} else {

		var commonURL = os.Getenv("CommonURL")
		payLink = commonURL + "/pay?iid=" + trackID
		Alerts = "Pay Link Generated Successfully with Pay ID :" + invoice_id
	}
	///////////////////////==================
	fmt.Println(payLink)
	// For display Currency List on List Box
	currencyList := []models.CurrencyList{}
	database.DB.Db.Table("currency").Order("currency_id ASC").Where("status = ?", 1).Find(&currencyList)

	return c.Render("pay-link-form", fiber.Map{
		"Title":        "Pay Link",
		"Subtitle":     "Generate",
		"Action":       "Add",
		"Alert":        Alerts,
		"CurrencyList": currencyList,
		"PayLink":      payLink,
		"MerchantData": merchantData,
	})
}

// Function for Post Payment form and display qr code in merchant section
func PayDataPost(c *fiber.Ctx) error {

	// Get Data from ajax
	req := new(models.PayRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	//fmt.Println("req=>", req)

	// Generate randomly Transaction ID
	transID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}
	///////////////////////////
	// Fetch Coin Data With Address
	coinList := models.PayCoinAddress{}
	database.DB.Db.Table("coin_list as a ").Select("b.address, b.lastupdate, a.coin_title, a.coin_id, a.icon, a.coin_network,a.coin_pay_url, b.address_id").Joins("left join coin_address as b on b.coin_id = a.coin_id").Where(" b.lastupdate < NOW() - INTERVAL '5 minutes' AND b.status = ? AND b.coin_id = ? ", 1, req.Crypto_id).Limit(1).Find(&coinList) //INTERVAL '1 hour' change for 1 hour
	//fmt.Println("coinList => ", coinList)
	if coinList.Address == "" {
		//fmt.Println("Try AGAIN")
		database.DB.Db.Table("coin_list as a ").Select("b.address, b.lastupdate, a.coin_title, a.coin_id, a.icon, a.coin_network,a.coin_pay_url, b.address_id").Joins("left join coin_address as b on b.coin_id = a.coin_id").Where(" b.status = ? AND b.coin_id = ? ", 1, req.Crypto_id).Limit(1).Find(&coinList) //INTERVAL '1 hour' change for 1 hour

	}
	// Fetch Data for Order List
	invoiceList := models.Invoice_Data{}
	database.DB.Db.Table("invoice").Where("trackid = ?", req.Customerrefid).Find(&invoiceList)
	//fmt.Println("invoiceList = ", invoiceList)
	price_amount := fmt.Sprintf("%f", invoiceList.Requestedamount)
	price_currency := strings.ToLower(invoiceList.Requestedcurrency)

	///////////////////////////

	if price_amount == "" || price_currency == "" || req.Cid == "" {
		return c.Redirect("/logout")

	}

	cryptoAmount, err := function.ConvertCurrencyToCrypto(price_amount, price_currency, strings.ToLower(coinList.Coin_title))
	//fmt.Println("Get Convert Currency => ", cryptoAmount)
	if err != nil {
		fmt.Println("Static Crypto Value")
		cryptoAmount = 0.00012
	}

	// Convert float64 to string
	//convertAMT := strconv.FormatFloat(cryptoAmount, 'f', -1, 64)

	//qr_code := strings.ToLower(coinList.Coin_pay_url) + ":" + coinList.Address + "?amount=" + convertAMT + "&label=Order%20ID%" + transID + "&message=ITIO" + transID
	qr_code := coinList.Address

	////////////////////////////////////
	s, _ := store.Get(c)
	var CID uint = 0
	if req.Pay_type == 1 {
		CID = req.Client_id
		//fmt.Print("inv -> ", CID)

	} else {
		CID = s.Get("LoginMerchantID").(uint)
		//fmt.Print("Pay -> ", CID)
	}

	// convert to float
	requestedAmount, err := strconv.ParseFloat(req.Price_amount, 64)
	if err != nil {
		fmt.Println(err)
	}

	requestedCurrency := price_currency
	convertedAmountV := strconv.FormatFloat(cryptoAmount, 'f', 6, 64)
	// convert string to float value
	convertedAmount, err := strconv.ParseFloat(convertedAmountV, 64)
	if err != nil {
		fmt.Println(" Error convert string to float value :")
	}

	//fmt.Println("convertedAmount ->", convertedAmount)
	convertedcurrency := req.Cid

	assetId := int(req.Crypto_id)
	receivedcurrency := req.Cid

	status := "Waiting"
	Transaction_type := "Collection"
	Note := "Sending " + strings.ToUpper(receivedcurrency) + " to Addresses - " + coinList.Address
	Ip := c.Context().RemoteIP().String()
	currentTime := time.Now()
	// Format the current time as a string
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	Customerrefid := req.Customerrefid
	if Customerrefid == "" {
		Customerrefid = transID
	}
	qry := models.Transaction_Pay{Client_id: CID, Transaction_id: transID, Assetid: assetId, Requestedamount: requestedAmount, Requestedcurrency: requestedCurrency, Convertedamount: convertedAmount, Convertedcurrency: convertedcurrency, Receivedcurrency: receivedcurrency, Customerrefid: Customerrefid, Transaction_type: Transaction_type, Ip: Ip, Note: Note, Status: status, Destinationaddress: coinList.Address, Createdate: formattedTime, Order_id: invoiceList.Order_id, Is_fee_paid_by_user: invoiceList.Is_fee_paid_by_user}
	result := database.DB.Db.Table("transaction").Select("client_id", "transaction_id", "assetid", "requestedamount", "requestedcurrency", "convertedamount", "convertedcurrency", "receivedcurrency", "customerrefid", "transaction_type", "ip", "note", "status", "destinationaddress", "createdate", "order_id", "is_fee_paid_by_user").Create(&qry)

	if result.Error != nil {
		fmt.Println("ERROR in QUERY qry ", result.Error)

	}

	customerData := models.CustomerData{Customer_name: req.Sender_name, Customer_email: req.Sender_email, Customer_tid: transID, Client_id: CID}
	result = database.DB.Db.Table("customer").Select("customer_name", "customer_email", "customer_tid", "client_id").Create(&customerData)

	if result.Error != nil {
		fmt.Println("ERROR in QUERY customerData ", result.Error)

	}

	///////////////////////////////////

	address := coinList.Address
	coinicon := coinList.Icon
	coinnetwork := coinList.Coin_network
	coin_id := coinList.Coin_id
	coin_pay_url := coinList.Coin_pay_url

	response := models.PayResponse{
		Qr_code:      qr_code,
		Address:      address,
		Amount:       convertedAmount,
		Transid:      transID,
		Coinicon:     coinicon,
		Coinnetwork:  coinnetwork,
		Coin_id:      coin_id,
		Coin_pay_url: coin_pay_url,
	}

	//fmt.Println("Response => ", response)

	aid, err := strconv.ParseUint(coinList.Address_id, 10, 32)
	if err != nil {
		fmt.Println("Error 105 XXX")
	}
	address_id := uint(aid)

	currentTime = time.Now()
	database.DB.Db.Table("coin_address").Save(&models.AddressDateUpdate{Address_id: address_id, Lastupdate: currentTime}).Where("address = ?", address)
	return c.JSON(response)
}
