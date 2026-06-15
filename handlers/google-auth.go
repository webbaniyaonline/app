package handlers

import (
	"encoding/base64"
	"fmt"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// toBase64 converts the QR code bytes to base64 to embed in the HTML img tag

// Function for Google enable 2FA
func EnableTwoFA(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	// Get value
	LoginMerchantEmail := s.Get("LoginMerchantEmail").(string)
	//fmt.Sprintln("LoginMerchantEmail == > ", LoginMerchantEmail)
	Alerts := s.Get("Alert")
	s.Delete("Alert")
	if err := s.Save(); err != nil {
		panic(err)
	}

	// Generate a new TOTP secret for the user
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Crypto Pay",
		AccountName: LoginMerchantEmail,
		Period:      30, // Refreshes every 30 seconds
		SecretSize:  16,
	})
	if err != nil {
		return c.Status(500).SendString("Error generating TOTP key")
	}
	//fmt.Println("key => ", key)
	// Store the TOTP secret in the database (omitted here for brevity)
	// saveTOTPSecret(username, key.Secret())

	// Generate a QR code URL for the user to scan
	qrCodeURL := key.URL()

	// Generate a QR code image
	qrBytes, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).SendString("Error generating QR code")
	}
	// Convert the QR code bytes to base64 for embedding in the HTML
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrBytes)

	// Return the QR code as a PNG image
	return c.Render("enable-2fa", fiber.Map{
		"Title":        "Enable 2FA",
		"QrCodeURL":    qrCodeURL,
		"Secret":       key.Secret(),
		"Qrimage":      qrCodeBase64,
		"Alert":        Alerts,
		"MerchantData": merchantData,
	})
}

// Function for Google enable 2FA
func ActivateTwoFA(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")
	getName := s.Get("LoginMerchantName").(string)
	getEmail := s.Get("LoginMerchantEmail").(string)

	// Parse the incoming JSON data
	var data map[string]interface{}

	// Bind the request body to the 'data' map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Failed to parse body")
	}

	// Access posted data
	totpSecret := data["secret"].(string)
	totpCode := data["code"].(string)
	qrimg := "data:image/png;base64," + data["qrimg"].(string)
	//fmt.Println("Data are -> ", qrimg)

	// Validate the TOTP code (assumes 30-second time window)
	valid := totp.Validate(totpCode, totpSecret)
	if !valid {
		fmt.Println("Invalid 2FA code")
		return c.JSON(fiber.Map{
			"message": "Invalid 2FA code",
			"status":  400,
		})

	} else {

		LoginID := LoginMerchantID.(uint)
		result := database.DB.Db.Table("client_master").Save(&models.Update2FA{Client_id: LoginID, Totp_secret: totpSecret, Totp_status: 1})

		if result.Error != nil {
			fmt.Println("ERROR in QUERY")

		}
		//  Email///
		template_code := "2FA-STATUS"
		//fmt.Println("qrimg==>", qrimg)
		emailData := models.EmailData{FullName: getName, Email: getEmail, UserName: getEmail, HashCode: totpSecret, Details: qrimg}
		//fmt.Println("ERROR 10001 => ", emailData)
		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}
		// Check if "MerchantData" is of type map[string]interface{}
		s.Set("LoginMerchantSecret", totpSecret)
		s.Set("LoginMerchantGoogleStatus", 1)
		if merchantDataMap, ok := merchantData.(map[string]interface{}); ok {
			// Update the specific value in the map
			merchantDataMap["MerchantSecret"] = totpSecret
			merchantDataMap["MerchantGoogleStatus"] = 1

			// Set the updated map back to the session
			s.Set("MerchantData", merchantDataMap)
		}

		Alerts := "2FA Activated"
		// check session
		s.Set("Alert", Alerts) // Set alert
		if err := s.Save(); err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"message": "2FA Activated",
			"status":  200,
		})
	}

}

// Function for Google disable 2FA
func DeactivateTwoFA(c *fiber.Ctx) error {

	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginID := s.Get("LoginMerchantID").(uint)
	totpSecret := ""
	result := database.DB.Db.Table("client_master").Save(&models.Update2FA{Client_id: LoginID, Totp_secret: totpSecret, Totp_status: 0})

	if result.Error != nil {
		fmt.Println("ERROR in QUERY")

	}

	// Check if "MerchantData" is of type map[string]interface{}
	s.Set("LoginMerchantSecret", totpSecret)
	s.Set("LoginMerchantGoogleStatus", 0)
	if merchantDataMap, ok := merchantData.(map[string]interface{}); ok {
		// Update the specific value in the map
		merchantDataMap["MerchantSecret"] = totpSecret
		merchantDataMap["MerchantGoogleStatus"] = 0

		// Set the updated map back to the session
		s.Set("MerchantData", merchantDataMap)
	}

	Alerts := "2FA Dectivated"
	// check session
	s.Set("Alert", Alerts) // Set alert
	if err := s.Save(); err != nil {
		return err
	}
	return c.Redirect("/profile")
}

// Function for Google verify 2FA
func VerifyTwoFA(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	if s.Get("LoginMerchantID") == nil {
		fmt.Println("Session Expired101")
		//return c.Redirect("/login", 301)
	}
	LoginMerchantName := s.Get("LoginMerchantName")
	LoginMerchantID := s.Get("LoginMerchantID")
	//fmt.Println("LLLLLLLLLLLl", LoginMerchantID, LoginMerchantName)
	Alerts := s.Get("Alert")
	if Alerts != "" {
		s.Delete("Alert")
		if err := s.Save(); err != nil {
			panic(err)
		}
	}

	return c.Render("verify-2fa", fiber.Map{
		"Title":             "Verify 2FA",
		"Alert":             Alerts,
		"LoginMerchantName": LoginMerchantName,
		"LoginMerchantID":   LoginMerchantID,
	})
}

// Function for Submit 2FA data
func VerifyTwoFAPost(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	if s.Get("LoginMerchantID") == nil {
		fmt.Println("Session Expired101")
		//return c.Redirect("/login", 301)
	}
	loginMerchantEmail := s.Get("LoginMerchantEmail")

	// Parse the incoming JSON data
	var data map[string]interface{}

	// Bind the request body to the 'data' map
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).SendString("Failed to parse body")
	}

	// Access posted data
	totpCode := data["code"].(string)
	cid := data["cid"].(string)

	fmt.Print("Code => ", totpCode)
	fmt.Print("cid => ", cid)

	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("status = ? AND client_id = ?", 1, cid).Find(&loginList)
	if result.Error != nil {
		fmt.Println("ERROR in QUERY")
	}
	totpSecret := loginList.Totp_secret
	// Validate the TOTP code (assumes 30-second time window)
	valid := totp.Validate(totpCode, totpSecret)
	if !valid {
		fmt.Println("Invalid 2FA code")
		return c.JSON(fiber.Map{
			"message": "Invalid 2FA code",
			"status":  400,
		})

	} else {
		fmt.Println("valid 2FA code !!!")
		/////////////////////////
		// Set key/value
		loginIp := c.Context().RemoteIP().String()
		s.Set("LoginMerchantName", loginList.Full_name)
		s.Set("LoginMerchantID", loginList.Client_id)
		s.Set("LoginMerchantEmail", loginMerchantEmail)
		s.Set("LoginMerchantStatus", loginList.Status)
		s.Set("LoginMerchantSecret", loginList.Totp_secret)
		s.Set("LoginMerchantGoogleStatus", loginList.Totp_status)
		s.Set("LoginMerchantUserAgent", loginList.User_agent)

		s.Set("MerchantData", map[string]interface{}{
			"MerchantName":         loginList.Full_name,
			"MerchantEmail":        loginMerchantEmail,
			"MerchantID":           loginList.Client_id,
			"MerchantStatus":       loginList.Status,
			"MerchantSecret":       loginList.Totp_secret,
			"MerchantGoogleStatus": loginList.Totp_status,
			"MerchantUserAgent":    loginList.User_agent,
			"MerchantLoginIP":      loginIp,
		})

		//Save sessions
		if err := s.Save(); err != nil {
			panic(err)
		}

		qry := models.LoginHistory{Client_id: loginList.Client_id, Login_ip: loginIp}
		result := database.DB.Db.Table("login_history").Select("client_id", "login_ip").Create(&qry)
		if result.Error != nil {
			fmt.Println(result.Error)
		}

		return c.JSON(fiber.Map{
			"message": "valid 2FA code",
			"status":  200,
		})

		///////////////////////////

	}

}
