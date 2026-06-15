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

// Define for struct for payment / pay link
type APIResponsePaylinkSuccess struct {
	Status      string
	ReferanceID string
	PayURL      string
}
type APIResponseFailed struct {
	Status string
	Error  string
}

type APIResponseBalanceSuccess struct {
	Receivedcurrency string
	Balance          string
	Timestamp        time.Time
}

type CreateLink struct {
	ProductName         string  `json:"ProductName"`
	Description         string  `json:"Description"`
	Currency            string  `json:"Currency"`
	Amount              float64 `json:"Amount"`
	CustomerName        string  `json:"CustomerName"`
	CustomerEmail       string  `json:"CustomerEmail"`
	OrderID             string  `json:"OrderID"`
	Is_fee_paid_by_user bool    `json:"is_fee_paid_by_user"`
	Return_url          string  `json:"return_url"`
}

// Function for Generate Pay Link By API
func ApiPaymentLink(c *fiber.Ctx) error {
	apiError := ""

	//fmt.Println("Calling")
	// Retrieve a specific header
	apikey := c.Get("Apikey")
	//fmt.Println("apikey", apikey)
	link := new(CreateLink)

	// Parse the request body into the User struct
	if err := c.BodyParser(link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON !!",
		})
	}
	//fmt.Println("BODY", link)
	MID, errorx := function.GetMIDByApikey(apikey)
	// Retrieve data from header
	price_currency := strings.TrimSpace(link.Currency)

	requestedamount := link.Amount

	//fmt.Println("Amount - > ".requestedamount)

	//  get values from array and set in variable
	productName := strings.TrimSpace(link.ProductName)
	description := strings.TrimSpace(link.Description)
	customerName := strings.TrimSpace(link.CustomerName)
	customerEmail := strings.TrimSpace(link.CustomerEmail)
	orderID := strings.TrimSpace(link.OrderID)
	return_url := strings.TrimSpace(link.Return_url)
	is_fee_paid_by_user := link.Is_fee_paid_by_user

	//fmt.Println(customerName, customerEmail, orderID, is_fee_paid_by_user)

	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx
	} else if price_currency == "" {
		apiError = "Currency Not Found"
	} else if link.Amount == 0.0 {
		apiError = "Amount Not Found"
	} else if productName == "" {
		apiError = "Product Name Name Not Found"
	} else if description == "" {
		apiError = "Description Not Found"
	} else if orderID == "" {
		apiError = "orderID Not Found"
	} else {

		// Generate randomly Transaction ID
		trackID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate random ID",
			})
		}

		Ip := c.Context().RemoteIP().String() // Get Current IP
		qry := models.Invoice_Master{Client_id: MID, Requestedamount: requestedamount, Requestedcurrency: price_currency, Product_name: productName, Description: description, Ip: Ip, Trackid: trackID, Name: customerName, Email: customerEmail, Order_id: orderID, Is_fee_paid_by_user: is_fee_paid_by_user, Return_url: return_url}
		result := database.DB.Db.Table("invoice").Select("client_id", "requestedamount", "requestedcurrency", "product_name", "description", "ip", "trackid", "name", "email", "order_id", "is_fee_paid_by_user", "return_url").Create(&qry)
		invoice_id := strconv.FormatUint(uint64(qry.Invoice_id), 10)
		fmt.Println(invoice_id)
		if result.Error != nil {
			fmt.Println("Data Not Inserted")

		}
		var commonURL = os.Getenv("CommonURL")
		PayURL := commonURL + "/pay?iid=" + trackID // create pay url
		status := "Ok"

		response := APIResponsePaylinkSuccess{
			Status:      status,
			ReferanceID: trackID,
			PayURL:      PayURL,
		}

		return c.JSON(response)
	}
	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// Function for Get Balance By API
func ApiBalanceByCrypto(c *fiber.Ctx) error {

	apiError := ""
	apikey := ""
	currency := ""
	// Retrieve a specific header
	apikey = strings.TrimSpace(c.Get("Apikey"))
	currency = strings.TrimSpace(strings.ToLower(c.Query("Currency")))
	//currency := strings.ToLower("LTC")

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx

	} else {

		assetList := []APIResponseBalanceSuccess{}
		var totalWallet int64
		// fetch query for wallet with balance
		if currency != "" {
			database.DB.Db.Table("transaction").Select("assetid, receivedcurrency, SUM(receivedamount)  as balance , now() as timestamp").Where("client_id = ? AND status = ? AND receivedcurrency = ?", MID, "Success", currency).Group("assetid,receivedcurrency").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList).Count(&totalWallet)

		} else {
			database.DB.Db.Table("transaction").Select("assetid, receivedcurrency, SUM(receivedamount)  as balance , now() as timestamp").Where("client_id = ? AND status = ?", MID, "Success").Group("assetid,receivedcurrency").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList).Count(&totalWallet)

		}
		//fmt.Println(c.JSON(assetList))
		return c.JSON(assetList)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// Function for Get Customer By API
func ApiCustomer(c *fiber.Ctx) error {
	apiError := ""
	dateFrom := ""
	dateTo := ""
	searchQuery := ""

	// Retrieve a specific header
	apikey := strings.TrimSpace(c.Get("Apikey"))
	//apikey := "76419b7b23017e61"
	dateFrom = strings.TrimSpace(c.Query("DateFrom")) // Get data from url
	dateTo = strings.TrimSpace(c.Query("DateTo"))     // Get data from url

	//fmt.Println("Get DATAs Are : ", dateFrom, dateTo)

	// convert limit value from string to integer
	default_limit, err := strconv.Atoi(c.Query("Limit", "100"))
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Converted integer:", default_limit)
	}

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		//fmt.Println(errorx)
		apiError = errorx

	} else {

		if default_limit > 500 {
			default_limit = 500
		}

		if dateFrom != "" && dateTo != "" {
			searchQuery = " timestamp BETWEEN '" + dateFrom + "' AND '" + dateTo + "' AND "
		}
		assetList := []models.CustomerList{}
		database.DB.Db.Table("customer").Select("customer_name", "customer_email", "COUNT(*) AS total_customer").Limit(default_limit).Where(searchQuery+"client_id = ?", MID).Group("customer_email, customer_name").Find(&assetList)
		return c.JSON(assetList)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// Function for Get Pay Links By API
func ApiCheckouts(c *fiber.Ctx) error {
	apiError := ""
	// Retrieve a specific header
	apikey := strings.TrimSpace(c.Get("Apikey"))

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		//fmt.Println(errorx)
		apiError = errorx

	} else {

		dataList := []models.Invoice_Master{}
		database.DB.Db.Table("invoice").Where("status=? AND client_id = ?", 1, MID).Order("invoice_id desc").Limit(100).Find(&dataList)

		return c.JSON(dataList)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// API Function for Get Transaction details by Transaction ID
func ApiTransactionByTransID(c *fiber.Ctx) error {

	apiError := ""
	// Retrieve a specific header
	apikey := strings.TrimSpace(c.Get("Apikey"))
	TransID := strings.TrimSpace(c.Params("TransID"))
	//fmt.Println("TransID==>", TransID)

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		//fmt.Println(errorx)
		apiError = errorx
	} else if TransID == "" {
		apiError = "TransID Not Found"

	} else {

		transData := models.Transaction_Pay{}
		database.DB.Db.Table("transaction").Where("transaction_id = ? AND client_id = ? ", TransID, MID).Omit("client_id", "assetid", "response_json").Find(&transData)

		return c.JSON(transData)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// API Function for Get Transaction details by Reference ID
func ApiTransactionByReferenceID(c *fiber.Ctx) error {

	apiError := ""
	// Retrieve a specific header
	apikey := strings.TrimSpace(c.Get("Apikey"))
	ReferenceID := strings.TrimSpace(c.Params("ReferenceID"))
	//fmt.Println("TransID==>", ReferenceID)

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		//fmt.Println(errorx)
		apiError = errorx
	} else if ReferenceID == "" {
		apiError = "ReferenceID Not Found"

	} else {

		transData := models.Transaction_Pay{}
		database.DB.Db.Table("transaction").Where("customerrefid = ? AND client_id = ? ", ReferenceID, MID).Omit("client_id", "assetid", "response_json").Find(&transData)

		return c.JSON(transData)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// API Function for Get Transaction details by Order ID
func ApiTransactionByOrderID(c *fiber.Ctx) error {

	apiError := ""
	// Retrieve a specific header
	apikey := strings.TrimSpace(c.Get("Apikey"))
	OrderID := strings.TrimSpace(c.Params("OrderID"))
	//fmt.Println("TransID==>", OrderID)

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx
	} else if OrderID == "" {
		apiError = "OrderID Not Found"

	} else {

		transData := models.Transaction_Pay{}
		database.DB.Db.Table("transaction").Where("order_id = ? AND client_id = ? ", OrderID, MID).Omit("client_id", "assetid", "response_json").Find(&transData)

		return c.JSON(transData)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}

// API Function for Get Transaction List last 100
func ApiTransactionList(c *fiber.Ctx) error {

	apiError := ""
	dateFrom := ""
	dateTo := ""
	searchQuery := ""

	// Retrieve a specific header
	apikey := strings.TrimSpace(c.Get("Apikey"))
	//apikey := "76419b7b23017e61"
	dateFrom = strings.TrimSpace(c.Query("DateFrom")) // Get data from url
	dateTo = strings.TrimSpace(c.Query("DateTo"))     // Get data from url

	//fmt.Println("Get DATAs Are : ", dateFrom, dateTo)

	// convert limit value from string to integer
	default_limit, err := strconv.Atoi(c.Query("Limit", "100"))
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Converted integer:", default_limit)
	}

	MID, errorx := function.GetMIDByApikey(apikey)
	//fmt.Println(MID)
	if errorx != "" {
		fmt.Println(errorx)
		apiError = errorx

	} else {

		if default_limit > 1000 {
			default_limit = 1000
		}

		if dateFrom != "" && dateTo != "" {
			searchQuery = " createdate BETWEEN '" + dateFrom + "' AND '" + dateTo + "' AND "
		}

		transData := []models.Transaction_Pay{}
		database.DB.Db.Table("transaction").Where(searchQuery+"client_id = ? ", MID).Omit("client_id", "assetid", "response_json").Order("id DESC").Limit(default_limit).Find(&transData)

		return c.JSON(transData)
	}

	status := "Error"
	response := APIResponseFailed{
		Status: status,
		Error:  apiError,
	}

	return c.JSON(response)
}
