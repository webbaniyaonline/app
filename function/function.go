package function

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/models"
	"time"

	"github.com/go-resty/resty/v2"
)

// Function for send email by Email Template
func SendEmail(template_code string, p models.EmailData) error {
	//////////////////////////////////////////

	//Define Variables
	Email := "mailers@itio.in"
	FullName := "Null"
	UserName := "Null"
	Password := "Null"
	Invoice_id := ""
	Invoice_url := ""
	Invoice_url_html := ""
	Amount := ""
	Crypto := ""
	Status := ""
	HashCode := ""
	TransID := ""
	Title := ""
	Details := ""

	// Set Variables
	if p.Email != "" {
		Email = p.Email
	}
	if p.FullName != "" {
		FullName = p.FullName
	}
	if p.Amount != "" {
		Amount = p.Amount
	}
	if p.UserName != "" {
		UserName = p.UserName
	}
	if p.Crypto != "" {
		Crypto = p.Crypto
	}
	if p.Title != "" {
		Title = p.Title
	}
	if p.Details != "" {
		Details = p.Details
	}
	if p.Status != "" {
		Status = p.Status
	}
	if p.HashCode != "" {
		HashCode = p.HashCode
	}
	if p.TransID != "" {
		TransID = p.TransID
	}
	if p.Password != "" {
		Password = p.Password
	}
	if p.Invoice_id != "" {
		Invoice_id = p.Invoice_id
	}
	if p.Invoice_url != "" {
		Invoice_url = p.Invoice_url
		var commonURL = os.Getenv("CommonURL")
		// Split the string by "iid="
		parts := strings.Split(Invoice_url, "iid=")
		iid := parts[1]
		invoice_pdf_url := commonURL + "/invoice-details?iid=" + iid
		Invoice_url_html = "<a href='" + Invoice_url + "'  target='_blank' title='Pay Now'>Pay Now</a><br><a href='" + invoice_pdf_url + "'  target='_blank' title='View Invoice'><img src='{{.CssURLS}}/assets/images/pdficon.png' height='100'></a>"
	}

	if template_code == "2FA-STATUS" {
		//fmt.Println("Invoice_url==>", Invoice_url)
		Invoice_url_html = "<img src=" + p.Details + "  />"
	}

	/////////////////////////////////////////////
	// Get Email Template from  table - email_template
	emailTemplate := models.EmailTemplate{}
	result := database.DB.Db.Table("email_template").Where("template_code = ?", template_code).Find(&emailTemplate)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	// smtp - Details Get from env file
	var SMTPusername = os.Getenv("SMTPusername")
	var SMTPpassword = os.Getenv("SMTPpassword")
	var fromName = os.Getenv("HostName")
	var fromEmail = os.Getenv("FromEmail")
	host := os.Getenv("SMTPhost")
	port := os.Getenv("SMTPport")
	address := host + ":" + port

	// Replace content with email template format
	subject := strings.Replace(emailTemplate.Template_subject, "[sitename]", os.Getenv("SITENAME"), 1)
	subject = strings.Replace(subject, "[fullname]", FullName, 1)
	subject = strings.Replace(subject, "[username]", UserName, 1)
	subject = strings.Replace(subject, "[status]", Status, 1)

	EmailBody := strings.Replace(emailTemplate.Template_desc, "[sitename]", os.Getenv("SITENAME"), 1)
	EmailBody = strings.Replace(EmailBody, "[fullname]", FullName, 1)
	EmailBody = strings.Replace(EmailBody, "[username]", UserName, 1)
	EmailBody = strings.Replace(EmailBody, "[password]", Password, 1)
	EmailBody = strings.Replace(EmailBody, "[invoiceid]", Invoice_id, 1)
	EmailBody = strings.Replace(EmailBody, "[invoiceurl]", Invoice_url_html, 1)
	EmailBody = strings.Replace(EmailBody, "[amount]", Amount, 1)
	EmailBody = strings.Replace(EmailBody, "[crypto]", Crypto, 1)
	EmailBody = strings.Replace(EmailBody, "[sitename1]", os.Getenv("SITENAME"), 1)
	EmailBody = strings.Replace(EmailBody, "[sitename2]", os.Getenv("SITENAME2"), 1)
	EmailBody = strings.Replace(EmailBody, "[status]", Status, 1)
	EmailBody = strings.Replace(EmailBody, "[hashCode]", HashCode, 1)
	EmailBody = strings.Replace(EmailBody, "[transid]", TransID, 1)
	EmailBody = strings.Replace(EmailBody, "[title]", Title, 1)
	EmailBody = strings.Replace(EmailBody, "[details]", Details, 1)

	to := []string{Email}
	// Set up authentication information.
	auth := smtp.PlainAuth("", SMTPusername, SMTPpassword, host)
	// create MessageBody
	msg := []byte(
		"From: " + fromName + ": <" + fromEmail + ">\r\n" +
			"To: " + Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			EmailBody)
	err := smtp.SendMail(address, auth, fromEmail, to, msg)
	if err != nil {
		return err
	}

	return nil
}

// Function for generate random password
func PasswordGenerator(passwordLength int) string {
	// Character sets for generating passwords
	lowerCase := "abcdefghijklmnopqrstuvwxyz" // lowercase
	upperCase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // uppercase
	numbers := "0123456789"                   // numbers
	specialChar := "0"                        // special characters
	//specialChar := "!@#$%^&*()_-+={}[/?]"     // special characters

	// Variable for storing password
	password := ""

	// Initialize the random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Generate password character by character
	for n := 0; n < passwordLength; n++ {
		// Generate a random number to choose a character set
		randNum := rng.Intn(4)

		switch randNum {
		case 0:
			randCharNum := rng.Intn(len(lowerCase))
			password += string(lowerCase[randCharNum])
		case 1:
			randCharNum := rng.Intn(len(upperCase))
			password += string(upperCase[randCharNum])
		case 2:
			randCharNum := rng.Intn(len(numbers))
			password += string(numbers[randCharNum])
		case 3:
			randCharNum := rng.Intn(len(specialChar))
			password += string(specialChar[randCharNum])
		}
	}

	return password
}

// GenerateRandomID generates a random ID of the specified length
func GenerateRandomID(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type CurrencyResponse map[string]map[string]float64

// convertCurrencyToCrypto fetches the current crypto exchange rate and converts the given amount of fiat currency to cryptocurrency
func ConvertCurrencyToCrypto(amount, fromCurrency, toCrypto string) (float64, error) {
	client := resty.New()
	fmt.Print("toCrypto", toCrypto)
	if toCrypto == "polygon" {
		toCrypto = "matic-network"
	} else if toCrypto == "binance coin" {
		toCrypto = "binancecoin"
	}
	apiURL := "https://api.coingecko.com/api/v3/simple/price?ids=" + toCrypto + "&vs_currencies=" + fromCurrency

	//fmt.Print(apiURL)

	var apiResponse CurrencyResponse
	_, err := client.R().SetResult(&apiResponse).Get(apiURL)
	if err != nil {
		return 0, err
	}

	usdAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	rate := apiResponse[toCrypto][fromCurrency]
	cryptoAmount := usdAmount / rate

	return cryptoAmount, nil
}

type RequestAddress struct {
	Address string `json:"address"`
	Network string `json:"network"`
}

type ResponseAddress struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// function for Crypto currency address validator
func CryptoAddressValidator(address, currency string) (string, bool) {
	data := RequestAddress{
		Address: address,
		Network: currency,
	}

	// Marshal data to JSON
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling request")
	}
	// Create request
	req, err := http.NewRequest(http.MethodPost, "https://api.checkcryptoaddress.com/wallet-checks", bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error creating request")
	}
	// Set content type header
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")

	}
	defer resp.Body.Close()

	// Read response body (optional)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("responseBody Error : ", err)
	}
	// Print or process response body
	//fmt.Println(string(responseBody))
	// Parse the response
	var responseAddress ResponseAddress
	if err := json.Unmarshal(responseBody, &responseAddress); err != nil {
		fmt.Println("Error parsing API response", err)
	}
	message := responseAddress.Message
	valid := responseAddress.Valid

	//fmt.Println("Result =>", message, valid)

	return message, valid
}

type APIKeyResponse struct {
	Client_id int
}

// Get Merchant ID from API Key
func GetMIDByApikey(apikey string) (uint, string) {

	//fmt.Println("apikey - > ", apikey)

	aPIKeyResponse := APIKeyResponse{}
	database.DB.Db.Table("client_api").Select("client_id").Where("apikey = ? AND status = ?", apikey, 1).Find(&aPIKeyResponse)
	//fmt.Println("aPIKeyResponse - > ", aPIKeyResponse.Client_id)
	uid := uint(aPIKeyResponse.Client_id)
	errorx := ""
	if aPIKeyResponse.Client_id > 0 {
		return uid, errorx
	} else {
		errorx = "API Key Not Found"
		return uid, errorx
	}

}

type NameByMID struct {
	Full_name string
}

// Get Merchant Name  By Merchant ID
func GetNameByMID(MID uint) string {
	nameByMID := NameByMID{}
	database.DB.Db.Table("client_master").Select("full_name").Where("client_id = ? AND status = ?", MID, 1).Find(&nameByMID)
	full_name := string(nameByMID.Full_name)
	return full_name
}

type EmailByMID struct {
	Username string
}

// Get Merchant Email  By Merchant ID
func GetEmailByMID(MID uint) string {
	emailByMID := EmailByMID{}
	database.DB.Db.Table("client_master").Select("username").Where("client_id = ? AND status = ?", MID, 1).Find(&emailByMID)
	username := string(emailByMID.Username)
	return username
}

type MidByHash struct {
	Mid uint
}

// Get Merchant ID  By hashcode
func GetMidByHashID(Hash string) uint {
	midByHash := MidByHash{}
	database.DB.Db.Table("password_change_request").Select("mid").Where("password_hash = ?", Hash).Find(&midByHash)
	mid := midByHash.Mid
	return mid
}

type ReturnURLInvoice struct {
	Return_url string
}

// Get WP Get Redirec URL
func GetRedirecURL(Customerrefid string) string {
	data := ReturnURLInvoice{}
	database.DB.Db.Table("invoice").Select("return_url").Where("trackid = ?", Customerrefid).Find(&data)
	return_url := string(data.Return_url)
	return return_url
}

type CustomerEmailData struct {
	Customer_name  string
	Customer_email string
}

// Get Customer EMAIL from transaction id
func GetCustomerEmail(TransID string) string {
	data := CustomerEmailData{}
	database.DB.Db.Table("customer").Select("customer_name", "customer_email").Where("customer_tid = ?", TransID).Find(&data)
	return_data := string(data.Customer_name + "||" + data.Customer_email)
	return return_data
}

// Struct for Get Get Converted Amount By TransID
type ConvertedAmount struct {
	Convertedamount float64
}

// Function for Get Converted Amount By Transaction ID
func GetConvertedAmountByTransID(TransID string) float64 {
	convertedAmount := ConvertedAmount{}
	database.DB.Db.Table("transaction").Select("convertedamount").Where("transaction_id = ? ", TransID).Find(&convertedAmount)
	getamt := convertedAmount.Convertedamount
	return getamt
}

// Function for Get status by id

func GetStatusByStatusID(StatusID int) string {

	status := "Hi"
	if StatusID == 1 || StatusID == 2 || StatusID == 3 {
		status = "Success"
	} else if StatusID == 8 || StatusID == 9 {
		status = "Declined"
	} else if StatusID == 0 {
		status = "Waiting"
	} else {
		status = "Waiting"
	}
	return status
}

// Function for Get sub status by id
func GetSubStatusByStatusID(StatusID int) string {

	status := ""
	if StatusID == 1 {
		status = "FullPay"
	} else if StatusID == 2 {
		status = "OverPay"
	} else if StatusID == 3 {
		status = "UnderPay"
	} else if StatusID == 8 {
		status = "Declined"
	} else if StatusID == 9 {
		status = "Dispute"
	} else {
		status = "Waiting"
	}
	return status
}

// Function to check if the received amount is within 2% of the invoice amount
func IsPaymentSuccess(invoice, received float64, mid int64) bool {
	// Calculate the absolute difference
	diff := math.Abs(invoice - received)

	// Calculate the percentage difference
	percentageDiff := (diff / invoice) * 100

	// Get Merchant Success Margin Ratio
	// Fetch Webhook Url from client store
	storeData := models.ClientStore{}
	database.DB.Db.Table("client_store").Where("client_id = ?", mid).Find(&storeData)
	success_margin := storeData.Success_margin

	//fmt.Println("success_margin ==>", success_margin)

	if success_margin == 0 {
		success_margin = 2.0 // Set Default Value
	}
	//fmt.Println("success_margin set  ==>", success_margin)
	// Check if the percentage difference is less than or equal to 2%
	return percentageDiff <= success_margin
}

// Function for Update Merchant History
func UpdateMerchantHistory(updatetype, uptadedesc, ip string, Client_id uint) bool {

	// Format the current time as a string
	update_time := time.Now().Format("2006-01-02 15:04:05")
	// Create Query and Insert into DB
	qry := models.UpdateHistory{Client_id: Client_id, Update_ip: ip, Update_type: updatetype, Update_desc: uptadedesc, Updated_on: update_time}
	database.DB.Db.Table("update_history").Select("client_id", "Update_ip", "Update_type", "Update_desc", "Updated_on").Create(&qry)

	return true
}

type PasswordHistoryData struct {
	Client_id       uint
	Hashed_password string
}

// Function for Update Merchant History
func PasswordHistory(Password string, Client_id uint) bool {

	// Create Query and Insert into DB
	qry := PasswordHistoryData{Client_id: Client_id, Hashed_password: Password}
	database.DB.Db.Table("password_history").Select("client_id", "hashed_password").Create(&qry)

	return true
}

// Function for Update Merchant History
func PasswordGeneratedDuration(PasswordDate string) int {
	//days := 1
	// Parse the dates
	layout := "2006-01-02"
	// Parse the date string
	parsedDate, _ := time.Parse(time.RFC3339, PasswordDate)
	// Format the date to "YYYY-MM-DD"
	lastPasswordDate := parsedDate.Format(layout)        // convert date format like 2006-01-02
	todayDate := time.Now().Format(layout)               // convert date format like 2006-01-02
	startDate, _ := time.Parse(layout, lastPasswordDate) // convert date type string to time
	endDate, _ := time.Parse(layout, todayDate)          // convert date type string to time

	// Calculate the difference
	duration := endDate.Sub(startDate)
	days := int(duration.Hours() / 24)
	//fmt.Println("days : -  ", days)

	return days
}
