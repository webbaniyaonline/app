package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// Functions for Merchant Login / Registration/ Logout/ Login History/ Support Ticket/ change password and check session

// Configure session middleware for auto session out when page inactivity for 20 minutes
var store = session.New(session.Config{
	Expiration: 20 * time.Minute, // Set session expiration time
	//CookieHTTPOnly: true,            // Secure cookie handling
})

func init() {
	//database.ConnectDb() // for connect Db define in function.go

}

const (
	maxFailedAttempts = 3
	lockDuration      = 1 * time.Minute // Lock duration if needed
)

// Check admin Session exist or not
func MerchantSession(c *fiber.Ctx) error {
	// Session Check

	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	fmt.Println("MerchantData", merchantData)
	if merchantData == nil {

		s.Set("Alert", "Error Session expired") // Set a session key
		if err := s.Save(); err != nil {
			fmt.Println("check data - ", err)
		}
		//fmt.Println("PPPPPPpppppppp")
		// Redirect to another URL
		return c.Redirect("/logout")
	}
	return nil
}

// For Display Login form
func LoginView(c *fiber.Ctx) error {
	//c.Set("Cache-Control", "no-store")
	s, _ := store.Get(c)
	Alerts := s.Get("Alert")
	//fmt.Print("Alerts ", Alerts)
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println("Error Log1001 => ", err)
		}
	}

	//s, _ = store.Get(c)
	merchantData := s.Get("MerchantData")
	if merchantData != nil {
		//fmt.Println("Logged User")
		return c.Redirect("/")
	}

	//return c.Render("login", fiber.Map{})

	return c.Render("login", fiber.Map{
		"Title": "Login Form",
		"Alert": Alerts,
	})
}

// For Submit Login form
func LoginPost(c *fiber.Ctx) error {
	s, _ := store.Get(c)
	// Parses the request body
	getUserName := c.FormValue("username")
	getPassword := c.FormValue("password")

	//fmt.Println(getUserName, getPassword)

	Alerts := ""
	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("status = ? AND username = ?", 1, getUserName).Find(&loginList)
	//fmt.Println("result -> ", loginList)
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "ERROR in QUERY"
	}

	//fmt.Println("Days := ", days)
	if result.RowsAffected == 1 {

		var failedAttempts = loginList.FailedAttempts
		var accountLocked = loginList.AccountLocked
		var lockTime = loginList.LockTime

		//fmt.Println("failedAttempts =>", failedAttempts)
		//fmt.Println("accountLocked =>", accountLocked)

		// Check if account is locked
		if accountLocked {
			fmt.Println("accountLocked =>", accountLocked)
			if lockTime != nil && time.Since(*lockTime) >= lockDuration { //

				fmt.Println("Account unlocked Locked")
				// Unlock the account if the lock duration has passed

				database.DB.Db.Table("client_master").Save(&models.ClientWrongPasswordActivation{Client_id: loginList.Client_id, AccountLocked: false, FailedAttempts: 0, LockTime: nil})

			} else {

				s.Set("Alert", "Account is locked for 15 Minuts") // Set a session key
				if err := s.Save(); err != nil {
					fmt.Println("check data - ", err)
				}
				return c.Redirect("/login")

			}
		}

		// for check password  Expired or not
		days := function.PasswordGeneratedDuration(loginList.Password_created_at)

		// Get Password Period Limit from Env File
		passwordPeriod, _ := strconv.Atoi(os.Getenv("PASSWORDEXPIREDPERIOD")) // Convert string to integer
		if loginList.Status != 1 {
			Alerts = "Account Not Activate / Deleted"
		} else if days > passwordPeriod {
			Alerts = "Password Is Expired. Please change password"
		} else if loginList.Password != "" {
			err := bcrypt.CompareHashAndPassword([]byte(loginList.Password), []byte(getPassword))
			if err == nil {
				s, _ := store.Get(c)
				if loginList.Totp_status == 1 {
					s.Set("LoginMerchantName", loginList.Full_name)
					s.Set("LoginMerchantEmail", getUserName)
					s.Set("LoginMerchantID", loginList.Client_id)
					//Save sessions
					if err := s.Save(); err != nil {
						//panic(err)
						fmt.Println(err)
					}
					return c.Redirect("/verify-2fa")
				}

				// Set key/value
				loginIp := c.Context().RemoteIP().String()
				s.Set("LoginMerchantName", loginList.Full_name)
				s.Set("LoginMerchantID", loginList.Client_id)

				s.Set("LoginMerchantEmail", getUserName)
				s.Set("LoginMerchantStatus", loginList.Status)
				s.Set("LoginMerchantSecret", loginList.Totp_secret)
				s.Set("LoginMerchantGoogleStatus", loginList.Totp_status)
				s.Set("LoginMerchantUserAgent", loginList.User_agent)

				s.Set("MerchantData", map[string]interface{}{
					"MerchantName":         loginList.Full_name,
					"MerchantEmail":        getUserName,
					"MerchantID":           loginList.Client_id,
					"MerchantStatus":       loginList.Status,
					"MerchantSecret":       loginList.Totp_secret,
					"MerchantGoogleStatus": loginList.Totp_status,
					"MerchantUserAgent":    loginList.User_agent,
					"MerchantLoginIP":      loginIp,
				})

				// Format the current time as a string
				login_time := time.Now().Format("2006-01-02 15:04:05")
				qry := models.LoginHistory{Client_id: loginList.Client_id, Login_ip: loginIp, Login_time: login_time}
				result := database.DB.Db.Table("login_history").Select("client_id", "login_ip", "login_time").Create(&qry)
				//fmt.Println("Token_id    ", qry.Token_id)
				if result.Error != nil {
					fmt.Println(result.Error)
				}
				s.Set("LoginToken_id", qry.Token_id)
				//Save sessions
				if err := s.Save(); err != nil {
					//panic(err)
					fmt.Println(err)
				}

				return c.Redirect("/")

			} else {

				Alerts = "Wrong Password"

				failedAttempts++
				if failedAttempts >= maxFailedAttempts {

					database.DB.Db.Table("client_master").Save(&models.ClientWrongPasswordLock{Client_id: loginList.Client_id, AccountLocked: true, FailedAttempts: failedAttempts, LockTime: time.Now()})
					if result.Error != nil {
						fmt.Println("Database error")

					}

					s.Set("Alert", "Account locked due to multiple failed attempts") // Set a session key
					if err := s.Save(); err != nil {
						fmt.Println("check data - ", err)
					}
					return c.Redirect("/login")

				} else {
					//_, err := db.Exec(`UPDATE users SET failed_attempts = $1 WHERE id = $2`, failedAttempts, id)
					result = database.DB.Db.Table("client_master").Omit("AccountLocked", "LockTime").Save(&models.ClientWrongPasswordActivation{Client_id: loginList.Client_id, FailedAttempts: failedAttempts})

					if result.Error != nil {
						fmt.Println("Database error")

					}
				}

			}

		}

	} else {
		Alerts = "Account Not Found / Deactivated"

	}

	return c.Render("login", fiber.Map{
		"Title": "Login Form",
		"Alert": Alerts,
	})
}

// Function for Merchant Logout
func LogOut(c *fiber.Ctx) error {

	s, err := store.Get(c)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}

	if s.Get("LoginToken_id") != nil {
		LoginToken_id := s.Get("LoginToken_id").(uint)
		//fmt.Println("Logout ID = >", LoginToken_id)
		// Format the current time as a string
		logout_time := time.Now().Format("2006-01-02 15:04:05")
		result := database.DB.Db.Table("login_history").Save(&models.LoginHistoryUpdate{Token_id: LoginToken_id, Logout_time: logout_time})
		if result.Error != nil {
			fmt.Println("ERROR in QUERY")

		}
		s.Delete("LoginToken_id")
	}

	s.Delete("LoginMerchantID")
	s.Delete("LoginMerchantName")
	s.Delete("LoginMerchantEmail")
	s.Delete("LoginMerchantStatus")
	s.Delete("MerchantData")

	// Destroy session
	if err := s.Destroy(); err != nil {
		//panic(err)
		fmt.Println(err)
	}
	// Clear session or cookies
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	c.Cookie(cookie)
	msg, _ := strconv.Atoi(c.Query("v", "1"))
	Alerts := "Session expired. Please log in again."
	if msg == 2 {
		Alerts = ""
	}

	// check session
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}
	return c.Redirect("/login")
}

// For Display registration form
func RegistrationView(c *fiber.Ctx) error {

	return c.Render("registration", fiber.Map{
		"Title": "Registration Form",
		"Alert": "",
	})
}

// For Submit registration form
func RegistrationPost(c *fiber.Ctx) error {
	// Parses the request body
	getName := c.FormValue("name")
	getEmail := c.FormValue("email")

	// Find Duplicate Email in DB

	Alerts := ""
	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("username = ?", getEmail).Find(&loginList)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	receivedId := loginList.Client_id

	if receivedId == 0 {

		// END Find Duplicate Email in DB

		var password = function.PasswordGenerator(8)
		//fmt.Println(password)

		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		qry := models.Client_Master{Username: getEmail, Password: string(hash), Full_name: getName, Status: 1}
		result = database.DB.Db.Table("client_master").Select("username", "full_name", "password", "status").Create(&qry)

		if result.Error != nil {
			fmt.Println(result.Error)
		} else {
			//fmt.Println(result.RowsAffected)
			//fmt.Println(qry.Client_id)
			ClientID := qry.Client_id

			MyData := struct {
				Client_id uint `json:"name"`
			}{
				Client_id: ClientID,
			}
			result = database.DB.Db.Table("client_details").Select("client_id").Create(&MyData)
			if result.Error != nil {
				fmt.Println(result.Error)
			}
			//  Email///
			template_code := "SIGNUP-TO-MEMBER"

			emailData := models.EmailData{FullName: getName, Email: getEmail, UserName: getEmail, Password: password}

			err := function.SendEmail(template_code, emailData)
			if err != nil {
				fmt.Println("issue sending verification email")
			}
			//fmt.Print("Email Error 3")
			s, _ := store.Get(c)

			loginIp := c.Context().RemoteIP().String()
			// Set key/value
			s.Set("LoginMerchantName", getName)
			s.Set("LoginMerchantID", ClientID)
			s.Set("LoginMerchantEmail", getEmail)
			s.Set("LoginVoltID", "")
			s.Set("LoginMerchantStatus", 1)

			s.Set("MerchantData", map[string]interface{}{
				"MerchantName":    getName,
				"MerchantEmail":   getEmail,
				"MerchantID":      ClientID,
				"MerchantStatus":  1,
				"MerchantLoginIP": loginIp,
			})

			// check session
			Alerts = "Your login details have been sent to your registered email. Please check your inbox and complete your profile to enjoy full access to our features. Completing your profile helps us serve you better!"
			s.Set("Alert", Alerts) // Set a session key
			if err := s.Save(); err != nil {
				return err
			}
			//fmt.Println("Password => ", string(hash))
			//fmt.Println("Client ID => ", ClientID)
			function.PasswordHistory(string(hash), ClientID)
			return c.Redirect("/profile")

		}
	} else {
		Alerts = "Duplicate Email ID"

	}

	return c.Render("registration", fiber.Map{
		"Title": "Registration Form",
		"Alert": Alerts,
	})
}

// Function for display transaction list and wallet with balance on merchant dashboard
func IndexView(c *fiber.Ctx) error {

	//MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	if merchantData == nil {
		//fmt.Println("Session Expired102")
		return c.Redirect("/logout?v=2")
	}

	//fmt.Print("merchantData => ", merchantData)
	LoginMerchantID := s.Get("LoginMerchantID")

	transactionList := []models.Transaction_Pay{}
	var total int64
	// fetch query for transaction list
	database.DB.Db.Table("transaction").Order("id desc").Where("status = ? AND client_id = ?", "Success", LoginMerchantID).Limit(10).Find(&transactionList).Count(&total)

	assetList := []models.CoinWithBalance{}
	var totalWallet int64
	// fetch query for wallet with balance
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("client_id = ? AND status = ?", LoginMerchantID, "Success").Group("assetid").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList).Count(&totalWallet)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	// For display Currency List on List Box
	currencyList := []models.CurrencyList{}
	database.DB.Db.Table("currency").Order("currency_id ASC").Where("status = ?", 1).Find(&currencyList)

	// For display Currency List on List Box
	countTrans := models.CountTransactionByStatus{}
	database.DB.Db.Table("transaction").Select("COUNT(*) AS total_transactions, COUNT(CASE WHEN status = 'Success' THEN 1 END) AS total_success, COUNT(CASE WHEN status = 'Declined' THEN 1 END) AS total_failed, COUNT(CASE WHEN status = 'Waiting' THEN 1 END) AS total_process").Where("client_id = ?", LoginMerchantID).Find(&countTrans)

	return c.Render("index", fiber.Map{
		"Title":           "Dashboard",
		"Subtitle":        "Home",
		"TransactionList": transactionList,
		"CoinBalanceList": assetList,
		"MerchantData":    merchantData,
		"CountList":       total,
		"TotalWallet":     totalWallet,
		"CoinList":        coinList,
		"CurrencyList":    currencyList,
		"CountTrans":      countTrans,
	})
}

// Function for display Withdraw View
func WithdrawView(c *fiber.Ctx) error {

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

	feeList := models.FeesDetails{}
	database.DB.Db.Table("client_fees").Where("client_id = ?", LoginMerchantID).Order("client_id ASC").Find(&feeList)

	currencyList := []models.CryptoWalletList{}
	database.DB.Db.Table("client_wallet").Where("client_id = ?", LoginMerchantID).Find(&currencyList)

	//fmt.Println(feeList)
	//fmt.Println(currencyList)

	return c.Render("withdraw", fiber.Map{
		"Title":           "Withdraw",
		"Subtitle":        "Withdraw Crypto Balance",
		"CoinBalanceList": assetList,
		"MerchantData":    merchantData,
		"TotalWallet":     totalWallet,
		"CoinList":        coinList,
		"FeeList":         feeList,
		"CurrencyList":    currencyList,
	})
}

// Function for Submit Withdraw Form
func WithdrawFormPost(c *fiber.Ctx) error {

	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	asset_ID, err := strconv.ParseInt(c.FormValue("assetID"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	assetID := int(asset_ID)
	withdraw_amount := c.FormValue("withdraw_amount")
	withdrawAmount, err := strconv.ParseFloat(withdraw_amount, 64)
	if err != nil {
		fmt.Printf("Error converting form value to float: %v", err)
	}
	//withdrawAmountMinus :=-withdraw_amount
	withdrawAmountMinus := -withdrawAmount // covert with negative value

	assetList := models.CoinWithBalance{}
	// fetch query for wallet with balance
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("client_id = ? AND assetid = ? AND status = ?", LoginMerchantID, assetID, "Success").Group("assetid").Having("COUNT(assetid) > ?", 0).Find(&assetList)

	//balance := assetList.Balance
	balance, err := strconv.ParseFloat(assetList.Balance, 64)
	if err != nil {
		fmt.Printf("Error converting form value to float: %v", err)
	}

	feeList := models.FeesDetails{}
	database.DB.Db.Table("client_fees").Where("client_id = ?", LoginMerchantID).Order("client_id ASC").Find(&feeList)
	amount_fee_in_percent := feeList.Amount_fee_in_percent
	amountFeeInPercent, err := strconv.ParseFloat(amount_fee_in_percent, 64)
	if err != nil {
		fmt.Printf("Error converting form value to float: %v", err)
	}

	min_amount_fee := feeList.Min_amount_fee
	minAmountFee, err := strconv.ParseFloat(min_amount_fee, 64)
	if err != nil {
		fmt.Printf("Error converting form value to float: %v", err)
	}

	currencyList := models.CryptoWalletList{}
	database.DB.Db.Table("client_wallet").Where("client_id = ? AND assetid = ?", LoginMerchantID, assetID).Find(&currencyList)

	fees := (withdrawAmount * (amountFeeInPercent / 100))

	if minAmountFee > fees {
		fees = minAmountFee
	}

	// Convert the float value to its negative equivalent
	feesMinus := -fees // covert with negative value

	amountWithFees := withdrawAmount + fees

	Crypto_code := currencyList.Crypto_code
	destinationAddress := currencyList.Crypto_address
	destinationAddressFee := ""
	noteWithdraw := "Sending " + Crypto_code + " to Addresses - " + currencyList.Crypto_address
	noteFee := "Withdraw Fees"
	status := "Waiting"
	transTypeWithdraw := "Withdraw"
	transTypeFees := "Withdraw Fee"
	Ip := c.Context().RemoteIP().String()
	currentTime := time.Now()
	// Format the current time as a string
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// Generate randomly Transaction ID
	transIDWithdraw, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	// Generate randomly Transaction ID
	tranIDFees, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	if amountWithFees > balance {
		return c.Redirect("/withdraw")
	}

	// fmt.Println("=================DATAS=======================")
	// fmt.Println("transIDWithdraw => ", transIDWithdraw)
	// fmt.Println("tranIDFees => ", tranIDFees)
	// fmt.Println("assetID => ", assetID)
	// fmt.Println("clientID => ", LoginMerchantID)

	// fmt.Println("Withdraw Amount => ", withdrawAmount)
	// fmt.Println("fees = > ", fees)
	// fmt.Println("Crypto Code = > ", Crypto_code)
	// fmt.Println("Destination Address = > ", destinationAddress)
	// fmt.Println("Note Withdraw = > ", noteWithdraw)
	// fmt.Println("Note Fees = > ", noteFee)
	// fmt.Println("status = > ", status)
	// fmt.Println("transTypeWithdraw = > ", transTypeWithdraw)
	// fmt.Println("transTypeFees = > ", transTypeFees)
	// fmt.Println("IP = > ", Ip)
	// fmt.Println("formattedTime = > ", formattedTime)
	// fmt.Println("Balance = > ", balance)
	// fmt.Println("AmountWithFees = > ", AmountWithFees)

	qry := models.Transaction_Pay{Client_id: LoginMerchantID.(uint), Transaction_id: transIDWithdraw, Assetid: assetID, Requestedamount: withdrawAmount, Requestedcurrency: Crypto_code, Convertedcurrency: Crypto_code, Convertedamount: withdrawAmount, Receivedcurrency: Crypto_code, Receivedamount: withdrawAmountMinus, Transaction_type: transTypeWithdraw, Ip: Ip, Note: noteWithdraw, Status: status, Destinationaddress: destinationAddress, Createdate: formattedTime}
	result := database.DB.Db.Table("transaction").Select("client_id", "transaction_id", "assetid", "requestedamount", "requestedcurrency", "convertedamount", "convertedcurrency", "receivedcurrency", "receivedamount", "transaction_type", "ip", "note", "status", "destinationaddress", "createdate").Create(&qry)
	if result.Error != nil {
		fmt.Println("ERROR in QUERY", result.Error)

	}
	qryfees := models.Transaction_Pay{Client_id: LoginMerchantID.(uint), Transaction_id: tranIDFees, Assetid: assetID, Requestedamount: fees, Requestedcurrency: Crypto_code, Convertedamount: fees, Convertedcurrency: Crypto_code, Receivedcurrency: Crypto_code, Receivedamount: feesMinus, Transaction_type: transTypeFees, Ip: Ip, Note: noteFee, Status: status, Destinationaddress: destinationAddressFee, Createdate: formattedTime, Customerrefid: transIDWithdraw}
	resultfees := database.DB.Db.Table("transaction").Select("client_id", "transaction_id", "assetid", "requestedamount", "requestedcurrency", "convertedamount", "convertedcurrency", "receivedcurrency", "receivedamount", "transaction_type", "ip", "note", "status", "destinationaddress", "createdate", "customerrefid").Create(&qryfees)
	if resultfees.Error != nil {
		fmt.Println("ERROR in QUERY", result.Error)

	}

	Alerts := "Withdraw Request Generated"
	// check session
	s.Set("Alerts", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/merchant-transactions")
}

// For Display Merchant Login History
func LoginHistoryView(c *fiber.Ctx) error {

	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	if merchantData == nil {
		return c.Redirect("/logout")
	}
	LoginMerchantID := s.Get("LoginMerchantID")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	loginHistory := []models.LoginHistory{}
	database.DB.Db.Table("login_history").Order("token_id desc").Limit(limit).Offset(offset).Where(&models.LoginHistory{Login_type: 1, Client_id: LoginMerchantID.(uint)}).Find(&loginHistory)

	var total int64
	database.DB.Db.Table("login_history").Where(&models.LoginHistory{Login_type: 1, Client_id: LoginMerchantID.(uint)}).Count(&total)

	return c.Render("logged-history", fiber.Map{
		"Title":        "Login History",
		"Subtitle":     "Login History",
		"LoginHistory": loginHistory,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
		"MerchantData": merchantData,
	})
}

// For Display Customer Listing
func CustomerView(c *fiber.Ctx) error {
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	if merchantData == nil {
		return c.Redirect("/logout")
	}
	LoginMerchantID := s.Get("LoginMerchantID")

	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	customerData := []models.CustomerList{}
	result := database.DB.Db.Table("customer").Select("customer_name", "customer_email", "COUNT(*) AS total_customer").Limit(limit).Offset(offset).Where(&models.LoginHistory{Client_id: LoginMerchantID.(uint)}).Group("customer_email, customer_name").Find(&customerData)
	if result.Error != nil {
		fmt.Println(result.Error)
		//return c.Redirect("/logoutxx")
	}
	var total int64
	database.DB.Db.Table("customer").Select("customer_name", "customer_email", "COUNT(*) AS total_customer").Where(&models.CustomerList{Client_id: LoginMerchantID.(uint)}).Group("customer_email, customer_name").Count(&total)

	return c.Render("customer", fiber.Map{
		"Title":        "Customer",
		"Subtitle":     "Customer",
		"CustomerData": customerData,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
		"MerchantData": merchantData,
	})
}

// For Display Merchant API Key
func ApiKeyView(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	if merchantData == nil {
		return c.Redirect("/logout")
	}
	LoginMerchantID := s.Get("LoginMerchantID")

	// Fetch Data from client API
	var total int64
	clientAPI := []models.ClientAPI{}
	database.DB.Db.Table("client_api").Order("id desc").Where("status = ? AND Client_id=?", 1, LoginMerchantID.(uint)).Find(&clientAPI).Count(&total)

	// Fetch Data from client store
	clientStore := models.ClientStore{}
	database.DB.Db.Table("client_store").Where("Client_id=?", LoginMerchantID.(uint)).Find(&clientStore)
	// check session
	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	return c.Render("api-key", fiber.Map{
		"Title":        "API Key",
		"Subtitle":     "API Key",
		"ClientAPI":    clientAPI,
		"ClientStore":  clientStore,
		"Count":        total,
		"MerchantData": merchantData,
		"Alert":        Alerts,
	})
}

// For fetch Merchant API Key
func GetApiKey(c *fiber.Ctx) error {
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")
	currentTime := time.Now()
	// Format the current time as a string
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	// Generate randomly Transaction ID
	apikey, err := function.GenerateRandomID(8) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	database.DB.Db.Table("client_api").Save(&models.ClientAPI{Client_id: LoginMerchantID.(uint), Apikey: apikey, Status: 1, Timestamp: formattedTime})

	Alerts := "API Generated Successfully"
	// check session
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}
	return c.Redirect("/api-key")
}

// function for display merchant Reset Password Form
func ResetPasswordView(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}
	hashcode := ""
	emailx := ""

	if c.Params("HID") != "" {
		hashcode = c.Params("HID")
		//emailx = function.GetEmailByHashID(hashcode)

		emailByHash := models.EmailByHash{}
		database.DB.Db.Table("password_change_request").Select("password_email").Where("password_hash = ? AND timestamp >= NOW() - INTERVAL '30 minutes'", hashcode).Find(&emailByHash)
		emailx = string(emailByHash.Password_email)
		if emailx == "" {
			hashcode = ""
			Alerts = "Reset password link expired, Please generate new"
		}

	}
	//fmt.Println(" Email !!=> ", emailx)
	return c.Render("reset-password", fiber.Map{
		"Title":    "Reset Password",
		"Alert":    Alerts,
		"Hashcode": hashcode,
		"Emailx":   emailx,
	})
}

// function for Submit merchant Reset Password Form
func ResetPasswordPost(c *fiber.Ctx) error {
	// Parses the request body
	getEmail := c.FormValue("email")

	Alerts := "Reset password link has been sent to your registered Email ID"
	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("username = ?", getEmail).Find(&loginList)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	receivedId := loginList.Client_id
	//fmt.Println("XXXc ", receivedId)

	if receivedId > 0 {
		var token_id = function.PasswordGenerator(16)

		result := database.DB.Db.Table("password_change_request").Save(&models.PasswordChangeRequest{Password_hash: token_id, Password_email: getEmail, Mid: receivedId})
		if result.Error != nil {
			fmt.Println("ERROR in QUERY")

		}

		// var hash []byte
		// // func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		// hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// result := database.DB.Db.Table("client_master").Save(&models.ClientPassword{Client_id: loginList.Client_id, Password: string(hash)})
		// if result.Error != nil {
		// 	//fmt.Println("ERROR in QUERY")
		// 	Alerts = "Password Not Updated"
		// } else {
		// 	qrs := " Change Password request : New Password :" + password
		// 	updateIp := c.Context().RemoteIP().String()
		// 	function.UpdateMerchantHistory("RequestPass", qrs, updateIp, loginList.Client_id)

		// }
		//  Email///
		template_code := "RESET-PASSWORD-LINK"
		var commonURL = os.Getenv("CommonURL")
		PassURL := commonURL + "/reset-password/" + token_id
		hashCode := "<a href='" + PassURL + "' title='Reset Password' target='_blank'>" + PassURL + "</a>"
		emailData := models.EmailData{FullName: loginList.Full_name, Email: getEmail, HashCode: hashCode}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}

		//

	} else {
		Alerts = "Email Not Matched with our Record. Please Check and Try with Correct"
		// check session
		s, _ := store.Get(c)
		// check session
		s.Set("Alert", Alerts) // Set a session key
		if err := s.Save(); err != nil {
			return err
		}
		return c.Redirect("/reset-password")

	}

	return c.Render("reset-password", fiber.Map{
		"Title": "Reset-Password",
		"Alert": Alerts,
	})
}

// For Post Merchant Change Password form
func ResetPasswordSubmit(c *fiber.Ctx) error {

	// check session

	hashcode := c.FormValue("hashcode")
	if hashcode == "" {
		//fmt.Println("Session Expired120")
		return c.Redirect("/logout")
	}

	// Parses the request body
	new_password := c.FormValue("new_password")
	confirm_password := c.FormValue("confirm_password")
	Alerts := ""

	//fmt.Print("+++++=>", LoginMerchantID, new_password, confirm_password)
	if new_password == confirm_password {

		//getTableID := s.Get("LoginMerchantID").(uint) //c.FormValue("tableID")
		getTableID := function.GetMidByHashID(hashcode) //c.FormValue("tableID")

		//fmt.Println(password)

		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)

		//////////
		// if GET ID than work update else insert
		// for Full path use- filePath & only file name use file.Filename
		curr_date := time.Now().Format("2006-01-02")
		result := database.DB.Db.Table("client_master").Save(&models.ClientPassword{Client_id: getTableID, Password: string(hash), Password_created_at: curr_date})
		qrs := "Change Password by Merchant with hashcode - " + hashcode
		//fmt.Println(loginList.Status)
		Alerts = "Password Change successfully"
		if result.Error != nil {
			//fmt.Println("ERROR in QUERY")
			Alerts = "Password Not Updated"
			return c.Redirect("/login")
		} else {

			updateIp := c.Context().RemoteIP().String()
			function.UpdateMerchantHistory("Password", qrs, updateIp, getTableID)
			//fmt.Println("Password => ", string(hash))
			//fmt.Println("Client ID => ", getTableID)
			function.PasswordHistory(string(hash), getTableID)

		}
	} else {
		Alerts = "Password and confirm password not matched"
	}
	// check session
	s, _ := store.Get(c)

	//fmt.Print("Message=>", Alerts)
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/login")
}

// For  Merchant Support ticket Listing
func SupportTicket(c *fiber.Ctx) error {
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	LoginMerchantID := s.Get("LoginMerchantID")

	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	ticketList := []models.Support_Ticket{}
	var total int64
	database.DB.Db.Table("support-ticket").Where("client_id = ?", LoginMerchantID).Order("ticket_id desc").Limit(limit).Offset(offset).Find(&ticketList).Count(&total)

	return c.Render("support-tickets", fiber.Map{
		"Title":        "Support Ticket",
		"Subtitle":     "Support Ticket",
		"Action":       "List",
		"Alert":        Alerts,
		"TicketList":   ticketList,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
		"MerchantData": merchantData,
	})
}

// For Display Merchant Support ticket form
func AddSupportTicket(c *fiber.Ctx) error {
	MerchantSession(c) // redirect when session not found
	// check session
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	return c.Render("support-tickets", fiber.Map{
		"Title":        "Support Ticket",
		"Subtitle":     "Support Ticket",
		"Action":       "Add",
		"Alert":        Alerts,
		"MerchantData": merchantData,
	})
}

// For  Merchant Support ticket Details
func SupportTicketDetails(c *fiber.Ctx) error {

	ticketID := c.Query("tid")
	//fmt.Println("ticketID => ", ticketID)
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}
	// fetch data from support ticket
	supportList := models.Support_Ticket{}
	database.DB.Db.Table("support-ticket").Where("client_id = ? AND ticket_id = ? ", LoginMerchantID, ticketID).Find(&supportList)

	// fetch data from support ticket
	replyList := []models.Support_Ticket_Reply{}
	database.DB.Db.Table("support-ticket-reply").Where("ticket_id = ? ", ticketID).Order("reply_id desc").Find(&replyList)

	return c.Render("support-details", fiber.Map{
		"Title":        "Support Details",
		"Subtitle":     "Support Details",
		"Alert":        Alerts,
		"Modal":        1,
		"SupportList":  supportList,
		"ReplyList":    replyList,
		"MerchantData": merchantData,
	})
}

// For Post Merchant Support ticket
func SubmitSupportTicket(c *fiber.Ctx) error {

	ticket_subject := c.FormValue("ticket_subject")
	ticket_desc := c.FormValue("ticket_desc")
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")

	client_id := s.Get("LoginMerchantID").(uint)

	// for Full path use- filePath & only file name use file.Filename
	supportTicket := models.Support_Ticket{Client_id: client_id, Ticket_subject: ticket_subject, Ticket_desc: ticket_desc}
	result := database.DB.Db.Table("support-ticket").Omit("Status", "timestamp").Save(&supportTicket)

	//fmt.Println(loginList.Status)
	//////////////Email/////////////////
	ticket_id := strconv.FormatUint(uint64(supportTicket.Ticket_id), 10)

	template_code := "SUPPORT-TICKET"
	receiverEmail := os.Getenv("SupportEmail")
	ticket_subject = "Ticket # " + ticket_id + " - " + ticket_subject

	emailData := models.EmailData{Email: receiverEmail, Title: ticket_subject, Details: ticket_desc}
	err := function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	//////////////////End Email/////////

	Alerts := "Support Ticket Submitted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Support Ticket Not Submitted"
	}

	// check session
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}
	return c.Redirect("/support-tickets")
}

// For Post Merchant Support ticket Reply
func SubmitSupportTicketReply(c *fiber.Ctx) error {

	tid := c.FormValue("ticket_id")
	ticket_id, err := strconv.ParseInt(tid, 10, 64) // base 10 and int64 as result type
	if err != nil {
		fmt.Println("Error:", err)
	}
	reply_desc := c.FormValue("reply_desc")

	//fmt.Println(ticket_id, reply_desc)
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")

	reply_by := s.Get("LoginMerchantName").(string)
	usertype := "Merchant"

	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table("support-ticket-reply").Omit("timestamp").Save(&models.Support_Ticket_Reply{Ticket_id: ticket_id, Reply_by: reply_by, Type: usertype, Reply_desc: reply_desc})

	//////////////Email/////////////////

	template_code := "SUPPORT-REPLY"
	ticket_subject := "Re Ticket # " + tid
	SenderEmail := os.Getenv("SupportEmail")
	emailData := models.EmailData{Email: SenderEmail, Title: ticket_subject, Details: reply_desc}
	err = function.SendEmail(template_code, emailData)
	if err != nil {
		fmt.Println("issue sending verification email")
	}

	//////////////////End Email/////////

	//fmt.Println(loginList.Status)
	Alerts := "Ticket replied successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Ticket Not replied"
	}

	// check session

	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/support-details?tid=" + tid)
}

// For Display Merchant Change Password form
func ChangePasswordView(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")

	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	return c.Render("change-password", fiber.Map{
		"Title":           "Change Password",
		"Subtitle":        "Change Password",
		"Alert":           Alerts,
		"LoginMerchantID": LoginMerchantID,
		"MerchantData":    merchantData,
	})
}

// For Post Merchant Change Password form
func ChangePasswordPost(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")
	// Parses the request body
	new_password := c.FormValue("new_password")
	confirm_password := c.FormValue("confirm_password")
	Alerts := ""

	//fmt.Print("+++++=>", LoginMerchantID, new_password, confirm_password)
	if new_password == confirm_password {

		getTableID := s.Get("LoginMerchantID").(uint) //c.FormValue("tableID")

		//fmt.Println(password)

		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)

		//////////
		// if GET ID than work update else insert
		// for Full path use- filePath & only file name use file.Filename
		curr_date := time.Now().Format("2006-01-02")
		result := database.DB.Db.Table("client_master").Save(&models.ClientPassword{Client_id: getTableID, Password: string(hash), Password_created_at: curr_date})
		qrs := "Change Password by Merchant from profile"
		//fmt.Println(loginList.Status)
		Alerts = "Password update successfully"
		if result.Error != nil {
			//fmt.Println("ERROR in QUERY")
			Alerts = "Password Not Updated"
		} else {

			updateIp := c.Context().RemoteIP().String()
			function.UpdateMerchantHistory("Password", qrs, updateIp, getTableID)
			//fmt.Println("Password => ", string(hash))
			//fmt.Println("Client ID => ", getTableID)
			function.PasswordHistory(string(hash), getTableID)
		}
	} else {
		Alerts = "Password and confirm password not matched"
	}
	// check session
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/profile")
}

// For Post Merchant Store Details
func StorePost(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")

	// Parses the request body

	tableID := c.FormValue("tableID")
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	getTableID := uint(cid)
	webhookurl := c.FormValue("webhookurl")
	return_url := c.FormValue("return_url")
	successmargin := c.FormValue("success_margin")

	//convert value from string to float
	success_margin, err := strconv.ParseFloat(successmargin, 64)
	if err != nil {
		fmt.Printf("Error converting form value to float: %v", err)
	}

	Alerts := ""

	Client_id := s.Get("LoginMerchantID").(uint) //c.FormValue("tableID")

	//////////
	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table("client_store").Save(&models.ClientStore{ID: getTableID, Client_id: Client_id, Webhookurl: webhookurl, Return_url: return_url, Success_margin: success_margin})
	qrs := " Webhookurl: " + webhookurl + " Return_url: " + return_url + " Success_margin: " + strconv.FormatFloat(success_margin, 'f', 2, 64)
	//fmt.Println(loginList.Status)
	Alerts = "Store data update successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Store data Not Updated"
	} else {
		updateIp := c.Context().RemoteIP().String()
		function.UpdateMerchantHistory("Store", qrs, updateIp, Client_id)
	}

	// check session
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/api-key")
}

// function for display profile update form
func ProfileView(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	Alerts := s.Get("Alert")
	s.Delete("Alert")
	if err := s.Save(); err != nil {
		fmt.Println("Error - 90001")
	}

	profile := models.ProfileData{}
	database.DB.Db.Table("client_details").Where("client_id = ?", LoginMerchantID).Find(&profile)

	return c.Render("profile", fiber.Map{
		"Title":        "Profile",
		"Subtitle":     "Profile",
		"Alert":        Alerts,
		"Profile":      profile,
		"MerchantData": merchantData,
	})
}

// function for post profile update form
func ProfilePost(c *fiber.Ctx) error {
	// Parses the request body
	getGender := c.FormValue("gender")
	getBirthDate := c.FormValue("birth_date")
	getCountryCode := c.FormValue("country_code")
	getMobile := c.FormValue("mobile")
	getAddressLine1 := c.FormValue("address_line1")
	getAddressLine2 := c.FormValue("address_line2")
	getCity := c.FormValue("city")
	getState := c.FormValue("state")
	getCountry := c.FormValue("country")
	getPincode := c.FormValue("pincode")

	filePath := ""
	filename := c.FormValue("profileimage")

	// Retrieve the image file from the form
	file, err := c.FormFile("imageUpload")
	if err != nil {
		//fmt.Println("Image Not Found")
	} else {
		//fmt.Println("Image Found")
		// Read the file content
		fileContent, err := file.Open()
		if err != nil {
			fmt.Println("Error reading file")
		}
		defer fileContent.Close()

		// Save the file to the server
		filePath = fmt.Sprintf("./views/img/%s", file.Filename)
		filename = file.Filename
		bytes, err := ioutil.ReadAll(fileContent)
		if err != nil {
			fmt.Println("Error reading file content")
		}
		err = ioutil.WriteFile(filePath, bytes, 0644)
		if err != nil {
			fmt.Println("Error saving file")
		}

		// Respond with success
		//fmt.Println("Image uploaded successfully")
	}

	//fmt.Println("LoginMerchantID $$$$$$$$$$$")
	s, _ := store.Get(c)
	LoginMerchantID := s.Get("LoginMerchantID").(uint)

	//fmt.Println("LoginMerchantID =>>>>>>>>", LoginMerchantID)
	//fmt.Println("File Name", filename)
	//fmt.Println("File Path", filePath)
	//"Profile_image"

	result := database.DB.Db.Table("client_details").Save(&models.ProfileData{Client_id: LoginMerchantID, Gender: getGender, BirthDate: getBirthDate, CountryCode: getCountryCode, Mobile: getMobile, AddressLine1: getAddressLine1, AddressLine2: getAddressLine2, City: getCity, State: getState, Country: getCountry, Pincode: getPincode, Profile_image: filename})
	qrs := " Gender: " + getGender + " BirthDate: " + getBirthDate + " CountryCode: " + getCountryCode + " Mobile: " + getMobile + " AddressLine1: " + getAddressLine1 + " AddressLine2: " + getAddressLine2 + " City: " + getCity + " State: " + getState + " Country: " + getCountry + " Pincode: " + getPincode + " Profile_image:" + filename
	Alerts := "Profile Updated successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Profile Not Updated"
	} else {
		updateIp := c.Context().RemoteIP().String()
		function.UpdateMerchantHistory("Profile", qrs, updateIp, LoginMerchantID)
	}

	s.Set("Alert", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println(err)
	}

	return c.Redirect("/profile")
}

// For Add / Edit / Delete Crypto Wallet from Admin Section

// function for Display Currency List
func GetCryptoWalletList(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	LoginMerchantID := s.Get("LoginMerchantID")
	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	currencyList := []models.CryptoWalletList{}

	var total int64
	database.DB.Db.Table("client_wallet").Where("client_id = ?", LoginMerchantID).Order("status ASC,crypto_code ASC").Limit(limit).Offset(offset).Find(&currencyList).Count(&total)

	//fmt.Println(currencyList)
	return c.Render("crypto-wallet", fiber.Map{
		"Title":        "Settlement Wallet",
		"Subtitle":     "Settlement Wallet",
		"Action":       "List",
		"Alert":        Alerts,
		"CurrencyList": currencyList,
		"MerchantData": merchantData,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
	})
}

// function for Display Currency Form
func AddCryptoWalletView(c *fiber.Ctx) error {
	//fmt.Println("ADD WALLET")
	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	currencyList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Select("coin", "coin_title").Group("coin, coin_title").Where("status = ?", 1).Find(&currencyList)

	return c.Render("crypto-wallet", fiber.Map{
		"Title":        "Settlement Addresses",
		"Subtitle":     "Settlement Addresses",
		"Action":       "Add",
		"MerchantData": merchantData,
		"CurrencyList": currencyList,
	})
}

// function for Post Add / Edit Currency Form
func CryptoWalletPost(c *fiber.Ctx) error {
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	//merchantData := s.Get("MerchantData")

	// Parses the request body
	crypto_code := c.FormValue("crypto_code")
	crypto_title := c.FormValue("crypto_title")
	crypto_network := c.FormValue("crypto_network")
	crypto_address := c.FormValue("crypto_address")
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	status := int(status1)

	cryptoid, err := strconv.ParseInt(c.FormValue("crypto_id"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	crypto_id := int(cryptoid)

	tableID := c.FormValue("tableID")
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	getTableID := uint(cid)
	//////////
	client_id := s.Get("LoginMerchantID").(uint)
	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table("client_wallet").Save(&models.CryptoWalletList{Wallet_id: getTableID, Client_id: client_id, Crypto_code: crypto_code, Crypto_title: crypto_title, Crypto_network: crypto_network, Crypto_address: crypto_address, Status: status, Assetid: crypto_id})

	//fmt.Println(loginList.Status)
	Alerts := "Settlement Addresses Processed successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Settlement Addresses Not Updated"
	}

	// check session

	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/crypto-wallet")
}

// function for Post Add / Edit Currency Form
func EditCryptoWallet(c *fiber.Ctx) error {

	tableID := c.Params("TID")

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	data := models.CryptoWalletList{}
	database.DB.Db.Table("client_wallet").Where("wallet_id = ?", tableID).Find(&data)

	// For display Currency List on List Box
	currencyList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Select("coin", "coin_title").Group("coin, coin_title").Where("status = ?", 1).Find(&currencyList)
	fmt.Println(data)
	return c.Render("crypto-wallet", fiber.Map{
		"Title":        "Settlement Addresses",
		"Subtitle":     "Settlement Addresses",
		"Action":       "Edit",
		"MerchantData": merchantData,
		"EditData":     data,
		"CurrencyList": currencyList,
	})
}

// function for Delete Currency
func DeleteCryptoWallet(c *fiber.Ctx) error {
	AdminSession(c)
	id := c.Params("TID")

	// Convert string to uint
	tableID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid number format")
	}
	// Convert uint64 to uint
	getTableID := uint(tableID)

	status := 2
	result := database.DB.Db.Table("client_wallet").Save(&models.CryptoWalletDeleted{Wallet_id: getTableID, Status: status})

	Alerts := "Settlement Addresses Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Settlement Addresses Not Deleted"
	}

	// check session
	s, _ := store.Get(c)
	s.Set("Alert", Alerts) // Set a session key
	if err := s.Save(); err != nil {
		return err
	}

	return c.Redirect("/crypto-wallet")

}

// function for Get Network List when select crypto on dropdown
func GetNetwork(c *fiber.Ctx) error {

	crypto_code := c.Query("crypto_code")
	tab := c.Query("tab")
	fmt.Println("crypto_code > ", tab)
	if tab == "coin_list" {
		var networks []models.NetworkCoin
		database.DB.Db.Table("coin_list").Select("coin", "coin_network", "coin_id").Where("coin = ?", crypto_code).Find(&networks)
		fmt.Println("coin_list", networks)
		return c.JSON(networks)
	} else {
		var networks []models.Network
		database.DB.Db.Table("crypto_currency").Select("crypto_code", "crypto_network_short", "crypto_id").Where("crypto_code = ?", crypto_code).Find(&networks)
		fmt.Println("crypto_currency", networks)
		return c.JSON(networks)
	}

}

// function for process Social Media Login
func MerchantSocialLogin(c *fiber.Ctx) error {

	s, _ := store.Get(c)

	Alerts := ""
	getEmail := s.Get("SocialMerchantEmail").(string)
	getName := s.Get("SocialMerchantName").(string)
	resource := s.Get("SocialType").(string)

	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("username = ?", getEmail).Find(&loginList)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	receivedId := loginList.Client_id

	if receivedId == 0 {
		// END Find Duplicate Email in DB

		var password = function.PasswordGenerator(8)
		//fmt.Println(password)

		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		qry := models.Client_Master{Username: getEmail, Password: string(hash), Full_name: getName, Status: 1, Resource: resource}
		result = database.DB.Db.Table("client_master").Select("username", "full_name", "password", "status", "resource").Create(&qry)

		if result.Error != nil {
			fmt.Println(result.Error)
		} else {
			//fmt.Println(result.RowsAffected)
			//fmt.Println(qry.Client_id)
			ClientID := qry.Client_id

			MyData := struct {
				Client_id uint `json:"name"`
			}{
				Client_id: ClientID,
			}
			result = database.DB.Db.Table("client_details").Select("client_id").Create(&MyData)
			if result.Error != nil {
				fmt.Println(result.Error)
			}
			//  Email///
			template_code := "SIGNUP-TO-MEMBER"

			emailData := models.EmailData{FullName: getName, Email: getEmail, UserName: getEmail, Password: password}

			err := function.SendEmail(template_code, emailData)
			if err != nil {
				fmt.Println("issue sending verification email")
			}
			fmt.Print("Email Error 3")
			s, _ := store.Get(c)

			loginIp := c.Context().RemoteIP().String()
			// Set key/value
			s.Set("LoginMerchantName", getName)
			s.Set("LoginMerchantID", ClientID)
			s.Set("LoginMerchantEmail", getEmail)
			s.Set("LoginVoltID", "")
			s.Set("LoginMerchantStatus", 1)

			s.Set("MerchantData", map[string]interface{}{
				"MerchantName":    getName,
				"MerchantEmail":   getEmail,
				"MerchantID":      ClientID,
				"MerchantStatus":  1,
				"MerchantLoginIP": loginIp,
			})

			// check session
			Alerts = "Your login details have been sent to your registered email. Please check your inbox and complete your profile to enjoy full access to our features. Completing your profile helps us serve you better!"
			s.Set("Alert", Alerts) // Set a session key
			if err := s.Save(); err != nil {
				return err
			}

			return c.Redirect("/profile")

		}

	} else {
		//Login Process
		loginList := models.LoginList{}
		result := database.DB.Db.Table("client_master").Where("status = ? AND username = ?", 1, getEmail).Find(&loginList)

		//fmt.Println("Login Details", loginList)

		if result.Error != nil {
			//fmt.Println("ERROR in QUERY")
			Alerts = "ERROR in QUERY"
		}

		if result.RowsAffected == 1 {

			if loginList.Status != 1 {
				//fmt.Println("Account Not Activate / Deleted")
			} else {

				s, _ := store.Get(c)
				if loginList.Totp_status == 1 {
					s.Set("LoginMerchantName", loginList.Full_name)
					s.Set("LoginMerchantEmail", getEmail)
					s.Set("LoginMerchantID", loginList.Client_id)
					//Save sessions
					if err := s.Save(); err != nil {
						//panic(err)
						fmt.Println(err)
					}
					return c.Redirect("/verify-2fa")
				}

				// Set key/value
				loginIp := c.Context().RemoteIP().String()
				s.Set("LoginMerchantName", loginList.Full_name)
				s.Set("LoginMerchantID", loginList.Client_id)

				s.Set("LoginMerchantEmail", getEmail)
				s.Set("LoginMerchantStatus", loginList.Status)
				s.Set("LoginMerchantSecret", loginList.Totp_secret)
				s.Set("LoginMerchantGoogleStatus", loginList.Totp_status)
				s.Set("LoginMerchantUserAgent", loginList.User_agent)

				s.Set("MerchantData", map[string]interface{}{
					"MerchantName":         loginList.Full_name,
					"MerchantEmail":        getEmail,
					"MerchantID":           loginList.Client_id,
					"MerchantStatus":       loginList.Status,
					"MerchantSecret":       loginList.Totp_secret,
					"MerchantGoogleStatus": loginList.Totp_status,
					"MerchantUserAgent":    loginList.User_agent,
					"MerchantLoginIP":      loginIp,
				})

				//Save sessions
				if err := s.Save(); err != nil {
					//panic(err)
					fmt.Println(err)
				}

				qry := models.LoginHistory{Client_id: loginList.Client_id, Login_ip: loginIp, Resource: resource}
				result := database.DB.Db.Table("login_history").Select("client_id", "login_ip", "resource").Create(&qry)
				if result.Error != nil {
					fmt.Println(result.Error)
				}
				fmt.Println("Password => ", loginList.Password)
				fmt.Println("Client ID => ", loginList.Client_id)
				function.PasswordHistory(loginList.Password, loginList.Client_id)

				return c.Redirect("/")

			}

		} else {
			Alerts = "Account Not Found / Deactivated"

		}
	}

	return nil
}
