package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Function for Display Failed payment page
func FailedView(c *fiber.Ctx) error {

	transID := c.Params("TransID")
	//fmt.Println(transID)

	//=============Update Transaction status from Transaction id===============
	currentTime := time.Now()
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).UpdateColumns(models.FailedStatus{Status: "Declined", Substatus: 8, Response_timestamp: currentTime})

	//=============Fetch Transaction details from Transaction id===============
	transData := models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).Find(&transData)
	mID := transData.Client_id
	customerrefid := transData.Customerrefid

	//  Start Email Section ///
	template_code := "PAYMENT-STATUS"
	getName := function.GetNameByMID(mID)
	getEmail := function.GetEmailByMID(mID)
	hashCode := transData.Response_hash
	// Convert to string with 6 decimal places
	amount := strconv.FormatFloat(transData.Receivedamount, 'f', 6, 64) // 'f'
	crypto := transData.Receivedcurrency
	status := transData.Status
	details := "Transaction Not Done"

	// Send Email to Merchant
	emailData := models.EmailData{FullName: getName, Email: getEmail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

	err := function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	// Get customer details
	customerdetails := function.GetCustomerEmail(transID)

	// Split the string using "||" as the delimiter
	parts := strings.Split(customerdetails, "||")

	// Variables to store the parts
	var customername, customeremail string
	if len(parts) == 2 {
		customername = parts[0]
		customeremail = parts[1]

		// Send Email to Merchant
		emailData := models.EmailData{FullName: customername, Email: customeremail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
	}

	if os.Getenv("SupportEmail") != "" {
		adminname := os.Getenv("SupportEmail")
		adminnemail := os.Getenv("SupportEmail")

		// Send Email to Merchant
		emailData := models.EmailData{FullName: adminname, Email: adminnemail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
	}
	//  END Email Section ///
	// Fetch Webhook Url from client store
	storeData := models.ClientStore{}
	database.DB.Db.Table("client_store").Where("client_id = ?", mID).Find(&storeData)
	return_url := storeData.Return_url

	redirectURL := ""
	if return_url != "" {
		redirectURL = return_url + "?referanceID=" + transData.Customerrefid + "&transactionID=" + transData.Transaction_id + "&orderID=" + transData.Order_id + "&hash=" + transData.Response_hash + "&amount=" + fmt.Sprintf("%f", transData.Receivedamount) + "&currency=" + transData.Receivedcurrency + "&status=" + transData.Status
	}

	//check wordpress Redirect URL
	wpredirectURL := function.GetRedirecURL(customerrefid)
	if wpredirectURL != "" {
		redirectURL = wpredirectURL + "&transaction_id=" + transData.Transaction_id + "&status=failed" //success/failed
	}

	// For Post response on webhook URL
	webhookurl := storeData.Webhookurl //Get Webhook Url
	if webhookurl != "" {

		// Convert the struct to JSON
		jsonData, err := json.Marshal(transData)
		if err != nil {
			fmt.Println("Error in jsonData")
		}
		// Create the POST request
		resp, err := http.Post(webhookurl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error in resp")
		}
		defer resp.Body.Close()

	}
	return c.Render("failed", fiber.Map{
		"Title":       "Declined Payment",
		"Subtitle":    "Declined Payment",
		"TransID":     transID,
		"TransData":   transData,
		"RedirectURL": redirectURL,
	})
}

// Function for Display Dispute payment page
func DisputeView(c *fiber.Ctx) error {

	transID := c.Params("TransID")
	fmt.Println(transID)

	//=============Update Transaction status from Transaction id===============
	currentTime := time.Now()
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).UpdateColumns(models.FailedStatus{Status: "Declined", Substatus: 9, Response_timestamp: currentTime})

	//=============Fetch Transaction details from Transaction id===============
	transData := models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).Find(&transData)
	mID := transData.Client_id
	customerrefid := transData.Customerrefid

	//  Start Email Section ///
	template_code := "PAYMENT-STATUS"
	getName := function.GetNameByMID(mID)
	getEmail := function.GetEmailByMID(mID)
	hashCode := transData.Response_hash
	// Convert to string with 6 decimal places
	amount := strconv.FormatFloat(transData.Receivedamount, 'f', 6, 64) // 'f'
	crypto := transData.Receivedcurrency
	//fmt.Println("Amount = >> ", amount, amount)
	status := transData.Status
	details := "Transaction Received but Transaction Amount is Underpay / Overpay"

	// Send Email to Merchant
	emailData := models.EmailData{FullName: getName, Email: getEmail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

	err := function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	// Get customer details
	customerdetails := function.GetCustomerEmail(transID)

	// Split the string using "||" as the delimiter
	parts := strings.Split(customerdetails, "||")

	// Variables to store the parts
	var customername, customeremail string
	if len(parts) == 2 {
		customername = parts[0]
		customeremail = parts[1]

		// Send Email to Merchant
		emailData := models.EmailData{FullName: customername, Email: customeremail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
	}

	if os.Getenv("SupportEmail") != "" {
		adminname := os.Getenv("SupportEmail")
		adminnemail := os.Getenv("SupportEmail")

		// Send Email to Merchant
		emailData := models.EmailData{FullName: adminname, Email: adminnemail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
	}
	//  END Email Section ///

	// Fetch Webhook Url from client store
	storeData := models.ClientStore{}
	database.DB.Db.Table("client_store").Where("client_id = ?", mID).Find(&storeData)
	return_url := storeData.Return_url

	redirectURL := ""
	if return_url != "" {
		redirectURL = return_url + "?referanceID=" + transData.Customerrefid + "&transactionID=" + transData.Transaction_id + "&orderID=" + transData.Order_id + "&hash=" + transData.Response_hash + "&amount=" + fmt.Sprintf("%f", transData.Receivedamount) + "&currency=" + transData.Receivedcurrency + "&status=" + transData.Status
	}

	//check wordpress Redirect URL
	wpredirectURL := function.GetRedirecURL(customerrefid)
	if wpredirectURL != "" {
		redirectURL = wpredirectURL + "&transaction_id=" + transData.Transaction_id + "&status=failed" //success/failed
	}

	// For Post response on webhook URL
	webhookurl := storeData.Webhookurl //Get Webhook Url
	if webhookurl != "" {

		// Convert the struct to JSON
		jsonData, err := json.Marshal(transData)
		if err != nil {
			fmt.Println("Error in jsonData")
		}
		// Create the POST request
		resp, err := http.Post(webhookurl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error in resp")
		}
		defer resp.Body.Close()

	}

	return c.Render("dispute", fiber.Map{
		"Title":       "Dispute Payment",
		"Subtitle":    "Dispute Payment",
		"TransID":     transID,
		"TransData":   transData,
		"RedirectURL": redirectURL,
	})
}

// Function for Display Success payment page
func SuccessView(c *fiber.Ctx) error {

	transID := c.Params("TransID")
	fmt.Println(transID)
	//=============Fetch Invoice Details by trackid===============

	transData := models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Where("transaction_id = ?", transID).Find(&transData)
	mID := transData.Client_id
	customerrefid := transData.Customerrefid

	//fmt.Println(" Payment Status Success = ", transData.Status)

	if transData.Status != "Success" {
		return c.Redirect("/failed/" + transID)
	}

	//  Start Email Section ///
	template_code := "PAYMENT-STATUS"
	getName := function.GetNameByMID(mID)
	getEmail := function.GetEmailByMID(mID)
	hashCode := transData.Response_hash
	// Convert to string with 6 decimal places
	amount := strconv.FormatFloat(transData.Receivedamount, 'f', 6, 64) // 'f'
	//fmt.Println("Converted string:", amount)
	crypto := transData.Receivedcurrency

	status := transData.Status
	details := "Payment Received"

	// Send Email to Merchant
	emailData := models.EmailData{FullName: getName, Email: getEmail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

	err := function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	// Get customer details
	customerdetails := function.GetCustomerEmail(transID)

	// Split the string using "||" as the delimiter
	parts := strings.Split(customerdetails, "||")

	// Variables to store the parts
	var customername, customeremail string
	if len(parts) == 2 {
		customername = parts[0]
		customeremail = parts[1]

		// Send Email to Merchant
		emailData := models.EmailData{FullName: customername, Email: customeremail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
	}

	if os.Getenv("SupportEmail") != "" {
		adminname := os.Getenv("SupportEmail")
		adminnemail := os.Getenv("SupportEmail")

		// Send Email to Merchant
		emailData := models.EmailData{FullName: adminname, Email: adminnemail, Status: status, Amount: amount, Crypto: crypto, HashCode: hashCode, TransID: transID, Details: details}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
	}
	//  END Email Section ///

	// Fetch Webhook Url from client store
	storeData := models.ClientStore{}
	database.DB.Db.Table("client_store").Where("client_id = ?", mID).Find(&storeData)
	return_url := storeData.Return_url //Get Return URL

	redirectURL := ""
	if return_url != "" {
		redirectURL = return_url + "?referanceID=" + transData.Customerrefid + "&transactionID=" + transData.Transaction_id + "&orderID=" + transData.Order_id + "&hash=" + transData.Response_hash + "&amount=" + fmt.Sprintf("%f", transData.Receivedamount) + "&currency=" + transData.Receivedcurrency + "&status=" + transData.Status
	}

	//check wordpress Redirect URL
	wpredirectURL := function.GetRedirecURL(customerrefid)
	if wpredirectURL != "" {
		redirectURL = wpredirectURL + "&transaction_id=" + transData.Transaction_id + "&status=success" //success/failed
	}
	// For Post response on webhook URL
	webhookurl := storeData.Webhookurl //Get Webhook Url
	if webhookurl != "" {

		// Convert the struct to JSON
		jsonData, err := json.Marshal(transData)
		if err != nil {
			fmt.Println("Error in jsonData")
		}
		// Create the POST request
		resp, err := http.Post(webhookurl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error in resp")
		}
		defer resp.Body.Close()

	}

	return c.Render("success", fiber.Map{
		"Title":       "Success Payment",
		"Subtitle":    "Success Payment",
		"TransID":     transID,
		"TransData":   transData,
		"RedirectURL": redirectURL,
	})
}

// Function for Display  payment detail on checkout with qr code
func PayQRView(c *fiber.Ctx) error {

	PaymentID := c.Query("iid")

	//=============Fetch Invoice Details by trackid===============
	invoiceData := models.Invoice_Master{}
	database.DB.Db.Table("invoice").Where("trackid = ?", PaymentID).Find(&invoiceData)
	MID := invoiceData.Client_id
	//=============Fetch coin list ===============
	coinList := []models.CoinList{}
	//database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)
	database.DB.Db.Table("coin_list A ").Select("a.coin_id, a.coin, a.icon, a.status, a.coin_title, LOWER(a.coin_network) AS coin_network, a.coin_status_url, a.coin_pay_url ").Joins("LEFT JOIN client_wallet B ON A.coin_id = B.assetid ").Where(" a.status = ? AND B.client_id = ? AND B.status = ?", 1, MID, 1).Order("a.coin_title DESC").Find(&coinList)

	//fmt.Println(invoiceData)
	var commonURL = os.Getenv("CommonURL")
	return c.Render("checkout-pay-views", fiber.Map{
		"Title":       "Checkout",
		"Subtitle":    "Checkout",
		"CoinList":    coinList,
		"InvoiceData": invoiceData,
		"CommonURL":   commonURL,
	})
}
