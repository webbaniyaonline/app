package handlers

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// functions for login/ Logout / check session/ Dashboard listing/ login History/ Member / support ticket

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register(string(""))
	gob.Register(int(0))
	gob.Register(float64(0))
}

// function for display admin login form
func AdminLoginView(c *fiber.Ctx) error {

	sess, err := store.Get(c)
	if err != nil {
		return err
	}
	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			fmt.Println("Error Log1002 => ", err)
		}
	}
	//adminData := sess.Get("AdminData")
	//if adminData != nil {
	//return c.Redirect("/admin/") // if login redirect to dashboard
	//}

	return c.Render("admin/login", fiber.Map{
		"Title":  "Web Admin - Sign in",
		"AlertX": Alerts,
	})
}

// function for post admin login form
func AdminLoginPost(c *fiber.Ctx) error {
	// Parses the request body
	getAdminUserName := c.FormValue("admin-username")
	getAdminPassword := c.FormValue("admin-password")

	//fmt.Println(getUserName, getPassword)
	Alerts := ""
	loginList := models.AdminLoginList{}
	result := database.DB.Db.Table("admin_master").Where("username = ?", getAdminUserName).Find(&loginList)

	//fmt.Println(loginList.Status)

	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "ERROR in QUERY"
	}

	if result.RowsAffected == 1 {

		if loginList.Status != 1 {
			//fmt.Println("Account Not Activate / Deleted")
		} else if loginList.Password != "" {
			//fmt.Println(loginList.Password)
			err := bcrypt.CompareHashAndPassword([]byte(loginList.Password), []byte(getAdminPassword))
			if err == nil {
				loginIp := c.Context().RemoteIP().String()
				// Format the current time as a string
				login_time := time.Now().Format("2006-01-02 15:04:05")
				qry := models.LoginHistory{Client_id: loginList.Admin_id, Login_ip: loginIp, Login_type: 2, Login_time: login_time}
				result := database.DB.Db.Table("login_history").Select("client_id", "login_ip", "login_type", "login_time").Create(&qry)
				//fmt.Println("Token_id    ", qry.Token_id)
				if result.Error != nil {
					fmt.Println(result.Error)
				}

				// Set key/value
				sess, err := store.Get(c)
				if err != nil {
					return err
				}

				sess.Set("AdminData", map[string]interface{}{
					"AdminName":  loginList.Full_name,
					"AdminEmail": getAdminUserName,
					"AdminID":    loginList.Admin_id,
					"AdminRole":  loginList.Role,
				})

				sess.Set("AdminLoginToken_id", qry.Token_id)

				if err := sess.Save(); err != nil {
					return err
				}

				return c.Redirect("/admin/")

			} else {
				//fmt.Println("Wrong Password")
				Alerts = "Wrong Password"
			}

		}

	} else {
		//fmt.Println("Account Not Found")
		Alerts = "Account Not Found"

	}

	return c.Render("admin/login", fiber.Map{
		"Title":  "Web Admin - Sign in",
		"AlertX": Alerts,
		//"Facts":    facts,
	})
}

// Display page with Data on Admin Dashboard
func AdminIndexView(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Fetch Transaction List
	transactionList := []models.Transaction_Pay{}
	var total int64
	database.DB.Db.Table("transaction").Where("status = ?", "Success").Order("id desc").Limit(10).Find(&transactionList).Count(&total)

	//fetch balance
	assetList := []models.CoinWithBalance{}
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("status = ?", "Success").Group("assetid").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	return c.Render("admin/index", fiber.Map{
		"Title":           "Dashboard",
		"Subtitle":        "Home",
		"AdminData":       adminData,
		"TransactionList": transactionList,
		"CoinBalanceList": assetList,
		"CoinList":        coinList,
		"CountList":       total,
	})
}

// Admin Logout Function for session destroy
func AdminLogOut(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}

	if sess.Get("AdminLoginToken_id") != nil {
		AdminLoginToken_id := sess.Get("AdminLoginToken_id").(uint)
		//fmt.Println("Logout ID = >", AdminLoginToken_id)
		// Format the current time as a string
		logout_time := time.Now().Format("2006-01-02 15:04:05")
		result := database.DB.Db.Table("login_history").Save(&models.LoginHistoryUpdate{Token_id: AdminLoginToken_id, Logout_time: logout_time})
		if result.Error != nil {
			fmt.Println("ERROR in QUERY")

		}
		sess.Delete("LoginToken_id")
	}

	sess.Delete("AdminData")

	// Destroy session
	if err := sess.Destroy(); err != nil {
		//panic(err)
		fmt.Println(err)
	}
	// Clear session or cookies
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	c.Cookie(cookie)

	return c.Redirect("/admin/login")
}

// Check admin Session exist or not
func AdminSession(c *fiber.Ctx) error {
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	if adminData == nil {
		sess.Set("AlertX", "Error - Session expired") // Set a session key
		if err := sess.Save(); err != nil {
			fmt.Println("check data - ", err)
		}
		return c.Redirect("/admin/login")
	}

	return nil
}

// function for admin login history Listing
func AdminLoginHistory(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	//fmt.Println("check data - ", adminData)
	if adminData == nil {
		return c.Redirect("/admin/login")
	}
	// Convert the session data to a map
	adminMap := adminData.(map[string]interface{})
	//fmt.Println(adminMap["AdminID"])
	LoginAdminID := adminMap["AdminID"].(uint)
	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	loginHistory := []models.LoginHistory{}
	database.DB.Db.Table("login_history").Order("token_id desc").Limit(limit).Offset(offset).Where(&models.LoginHistory{Login_type: 2, Client_id: LoginAdminID}).Find(&loginHistory)

	var total int64
	database.DB.Db.Table("login_history").Where(&models.LoginHistory{Login_type: 2, Client_id: LoginAdminID}).Count(&total)

	//fmt.Println(loginHistory)
	return c.Render("admin/login-history", fiber.Map{
		"Title":        "Login History",
		"Subtitle":     "Login History",
		"LoginHistory": loginHistory,
		"AdminData":    adminData,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
	})
}

// Function for Member Listing in admin section
func AdminMembersView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Get query parameters
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "transaction_id")
	status := c.Query("status", "")
	order := c.Query("order", "desc")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	searchString := ""
	searchStringFull := ""

	if search != "" && sortBy != "" {
		searchString = " AND " + sortBy + " ILIKE " + "'%" + strings.TrimSpace(search) + "%'"
		searchStringFull = searchStringFull + "" + searchString
	}

	if status != "" {
		searchString = " AND status = " + status
		searchStringFull = searchStringFull + "" + searchString
	}

	if len(searchStringFull) > 4 {
		searchStringFull = searchStringFull[4:]
	}

	memberList := []models.MemberDetails{}
	database.DB.Db.Table("client_master as a ").Where(searchStringFull).Select("a.client_id, a.Username, a.Full_name,a.status,a.timestamp, b.title, b.gender, b.country_code, b.mobile, b.address_line1, b.address_line2").Joins("left join client_details as b on b.client_id = a.client_id").Order("a.client_id DESC").Limit(limit).Offset(offset).Find(&memberList)

	var total int64
	database.DB.Db.Table("client_master").Where(searchStringFull).Count(&total)
	//fmt.Println(memberList)

	feeList := []models.FeesDetails{}
	database.DB.Db.Table("client_fees").Order("client_id DESC").Find(&feeList)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/members", fiber.Map{
		"Title":      "Registered Members",
		"Subtitle":   "Registered Members",
		"AlertX":     Alerts,
		"AdminData":  adminData,
		"MemberList": memberList,
		"FeeList":    feeList,
		"Page":       page,
		"Limit":      limit,
		"Total":      total,
		"Search":     search,
		"SortBy":     sortBy,
		"Order":      order,
		"status":     status,
	})
}

// Function for Member Details in admin section
func AdminMembersDetailsView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	MID := c.Params("MID")

	memberList := models.MemberDetails{}
	database.DB.Db.Table("client_master as a ").Select("a.client_id, a.Username, a.Full_name,a.status,a.timestamp, b.title, b.gender, b.country_code, b.mobile, b.birth_date, b.address_line1, b.address_line2").Joins("join client_details as b on b.client_id = a.client_id AND b.client_id=?", MID).Find(&memberList)
	//fmt.Println("=======>", memberList)
	feeList := models.FeesDetails{}
	database.DB.Db.Table("client_fees").Where("client_id = ? ", MID).Find(&feeList)

	cryptoWalletList := []models.CryptoWalletList{}
	database.DB.Db.Table("client_wallet").Order("status ASC,crypto_code ASC").Where("client_id = ? ", MID).Find(&cryptoWalletList)

	clientAPI := []models.ClientAPI{}
	database.DB.Db.Table("client_api").Order("id desc").Where("status = ? AND Client_id=?", 1, MID).Find(&clientAPI)

	assetList := []models.CoinWithBalance{}
	database.DB.Db.Table("transaction").Select("assetid, SUM(receivedamount)  as balance").Where("client_id = ? AND status = ?", MID, "Success").Group("assetid").Having("COUNT(assetid) > ?", 0).Order("assetid ASC").Find(&assetList)

	// Display coin list in List box
	coinList := []models.CoinList{}
	database.DB.Db.Table("coin_list").Order("coin ASC").Where("status = ?", 1).Find(&coinList)

	// fetch query for transaction list
	transactionList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where("status = ? AND client_id = ?", "Success", MID).Limit(50).Find(&transactionList)
	// fetch query for transaction list
	withdrawList := []models.Transaction_Pay{}
	database.DB.Db.Table("transaction").Order("id desc").Where("transaction_type = ? AND client_id = ?", "Withdraw Fee", MID).Limit(50).Find(&withdrawList)
	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	//fmt.Println(memberList)
	//fmt.Println(feeList)

	return c.Render("admin/members-details", fiber.Map{
		"Title":            "Members Details",
		"Subtitle":         "Members Details",
		"AlertX":           Alerts,
		"AdminData":        adminData,
		"MemberList":       memberList,
		"FeeList":          feeList,
		"CryptoWalletList": cryptoWalletList,
		"ClientAPI":        clientAPI,
		"CoinBalanceList":  assetList,
		"CoinList":         coinList,
		"TransactionList":  transactionList,
		"WithdrawList":     withdrawList,
	})
}

// Function for Member Active / De Active
func AdminMemberStatus(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)

	// Parses the from Query

	status1, err := strconv.ParseInt(c.Query("key"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	status := int(status1)

	cid, err := strconv.ParseUint(c.Query("mid"), 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	mid := uint(cid)

	result := database.DB.Db.Table("client_master").Save(&models.MemberStatusUpdate{Client_id: mid, Status: status})

	Alerts := "Member Status Updated"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Member Status Not Updated"
	}

	sess.Set("AlertX", Alerts)
	if err := sess.Save(); err != nil {
		//panic(err)
		fmt.Println(err)
	}

	return c.Redirect("/admin/members")
}

// Function for Support Ticket Listing - Admin
func SupportTicketListing(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	ticketList := []models.Support_Ticket{}
	database.DB.Db.Table("support-ticket").Order("ticket_id desc").Limit(limit).Offset(offset).Find(&ticketList)

	var total int64
	database.DB.Db.Table("support-ticket").Count(&total)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/support-ticket", fiber.Map{
		"Title":      "Support Ticket",
		"Subtitle":   "Support Ticket",
		"AlertX":     Alerts,
		"AdminData":  adminData,
		"TicketList": ticketList,
		"Page":       page,
		"Limit":      limit,
		"Total":      total,
	})
}

// Function for Display Merchant Log File
func AdminMemberLogs(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	mid, _ := strconv.Atoi(c.Query("mid"))
	//fmt.Println("Mid => ", mid)
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	logsList := []models.UpdateHistory{}
	var total int64 // for count table data
	if mid == 0 {
		database.DB.Db.Table("update_history").Order("update_id DESC").Limit(limit).Offset(offset).Find(&logsList)
		database.DB.Db.Table("update_history").Count(&total)
	} else {
		database.DB.Db.Table("update_history").Where("client_id = ?", mid).Order("update_id DESC").Limit(limit).Offset(offset).Find(&logsList)
		database.DB.Db.Table("update_history").Where("client_id = ?", mid).Count(&total)
	}

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/merchant-logs", fiber.Map{
		"Title":     "Merchant Logs",
		"Subtitle":  "Merchant Logs",
		"AlertX":    Alerts,
		"AdminData": adminData,
		"LogsList":  logsList,
		"Page":      page,
		"Limit":     limit,
		"Total":     total,
		"Mid":       mid,
	})
}

// function for Post Add / Edit Fees
func FeesPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	min_amount_fee := c.FormValue("min_amount_fee")
	amount_fee_in_percent := c.FormValue("amount_fee_in_percent")

	clientID, err := strconv.ParseUint(c.FormValue("client_id"), 10, 32)
	if err != nil {
		fmt.Println("Error 1041")
	}
	client_id := uint(clientID)
	client_idx := int(clientID)

	checkList := models.FeesDetails{}
	result := database.DB.Db.Table("client_fees").Where("client_id = ?", client_id).Find(&checkList)

	//fmt.Println(checkList.Client_id)

	if result.RowsAffected == 1 {
		fmt.Print("Update")
		result = database.DB.Db.Table("client_fees").Save(&models.FeesDetails{Client_id: client_id, Min_amount_fee: min_amount_fee, Amount_fee_in_percent: amount_fee_in_percent}).Where("client_id = ?", client_id)

	} else {
		fmt.Print("Add")

		result = database.DB.Db.Table("client_fees").Save(&models.FeesDetail{Min_amount_fee: min_amount_fee, Amount_fee_in_percent: amount_fee_in_percent, Client_id: client_idx})
	}
	//////////

	Alerts := " Fee imposed Successfully"
	if result.Error != nil {
		fmt.Println("ERROR in QUERY", result.Error)
		Alerts = "Fees Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/members")
}

// function for display admin Reset Password Form
func AdminResetPasswordView(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	Alerts := s.Get("Alertx")
	if Alerts != "" {
		s.Delete("Alertx")
		if err := s.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	return c.Render("admin/reset-password", fiber.Map{
		"Title":  "Reset Password",
		"Alertx": Alerts,
	})
}

// function for Submit admin Reset Password Form
func AdminResetPasswordPost(c *fiber.Ctx) error {
	// Parses the request body
	getEmail := c.FormValue("email")

	Alerts := "New Generated password has been sent to your registered Email ID"
	loginList := models.AdminList{}
	result := database.DB.Db.Table("admin_master").Where("username = ?", getEmail).Find(&loginList)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	receivedId := loginList.Admin_id
	//fmt.Println("VVV ", receivedId)

	if receivedId > 0 {
		var password = function.PasswordGenerator(8)
		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		result := database.DB.Db.Table("admin_master").Save(&models.AdminPassword{Admin_id: loginList.Admin_id, Password: string(hash)})
		if result.Error != nil {
			//fmt.Println("ERROR in QUERY")
			Alerts = "Password Not Updated"
		}
		//  Email///
		template_code := "RESTORE-PASSWORD"

		emailData := models.EmailData{FullName: loginList.Full_name, Email: getEmail, UserName: getEmail, Password: password}

		err := function.SendEmail(template_code, emailData)
		if err != nil {
			fmt.Println("issue sending verification email")
		}

	} else {
		Alerts = "Email Not Matched with our Record. Please Check and Try with Correct"
		// check session
		sess, _ := store.Get(c)
		// check session
		sess.Set("Alertx", Alerts) // Set a session key
		if err := sess.Save(); err != nil {
			return err
		}
		return c.Redirect("/admin/reset-password")

	}

	return c.Render("admin/reset-password", fiber.Map{
		"Title":  "Reset-Password",
		"Alertx": Alerts,
	})
}

// Function for Coin ID Listing in admin section
func AdminCoinIDView(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// Get query parameters
	search := c.Query("search", "")

	//fmt.Println("search => ", search)

	searchStringFull := ""

	if search != "" {
		searchStringFull = " coin_id ILIKE " + "'%" + search + "%' OR  symbol ILIKE " + "'%" + search + "%'"
	}

	coinList := []models.CoinIDDetails{}
	database.DB.Db.Table("exchange_coinid").Where(searchStringFull).Limit(200).Find(&coinList)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	//fmt.Println(transactionList)
	return c.Render("admin/coinid", fiber.Map{
		"Title":     "Search COIN ID",
		"Subtitle":  "Search COIN ID",
		"AlertX":    Alerts,
		"AdminData": adminData,
		"CoinList":  coinList,
	})
}

// For  Admin Support ticket Details
func AdminSupportTicketDetails(c *fiber.Ctx) error {

	ticketID := c.Query("tid")
	//fmt.Println("ticketID => ", ticketID)

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// fetch data from support ticket
	supportList := models.Support_Ticket{}
	database.DB.Db.Table("support-ticket").Where("ticket_id = ? ", ticketID).Find(&supportList)

	// fetch data from support ticket
	replyList := []models.Support_Ticket_Reply{}
	database.DB.Db.Table("support-ticket-reply").Where("ticket_id = ? ", ticketID).Order("reply_id desc").Find(&replyList)

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			//panic(err)
			fmt.Println(err)
		}
	}

	return c.Render("admin/support-details", fiber.Map{
		"Title":       "Support Details",
		"Subtitle":    "Support Details",
		"Modal":       1,
		"SupportList": supportList,
		"ReplyList":   replyList,
		"AlertX":      Alerts,
		"AdminData":   adminData,
	})
}

// For Post Admin Support ticket reply form
func AdminSubmitSupportTicketReply(c *fiber.Ctx) error {
	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	//fmt.Println("check data - ", adminData)
	if adminData == nil {
		return c.Redirect("/admin/login")
	}
	// Convert the session data to a map
	adminMap := adminData.(map[string]interface{})
	//fmt.Println(adminMap["AdminID"])
	LoginAdminName := adminMap["AdminName"].(string)

	tid := c.FormValue("ticket_id")

	// Convert string to uint64 first
	mid, err := strconv.ParseUint(c.FormValue("client_id"), 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Cast uint64 to uint
	merchantID := uint(mid)

	SenderEmail := function.GetEmailByMID(merchantID)

	ticket_id, err := strconv.ParseInt(tid, 10, 64) // base 10 and int64 as result type
	if err != nil {
		fmt.Println("Error:", err)
	}

	reply_desc := c.FormValue("reply_desc")

	reply_by := LoginAdminName
	usertype := "Support"

	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table("support-ticket-reply").Omit("timestamp").Save(&models.Support_Ticket_Reply{Ticket_id: ticket_id, Reply_by: reply_by, Type: usertype, Reply_desc: reply_desc})

	//////////////Email/////////////////

	template_code := "SUPPORT-REPLY"
	ticket_subject := "Re Ticket # " + tid

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

	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/support-details?tid=" + tid)
}
