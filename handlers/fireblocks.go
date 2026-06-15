package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"template/database"
	"template/fireblocks"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
)

// Function for use fireblock API
// for sandbox
var privateKeyPath = "./fireblocks_secret.key"
var apiKey = "abdb1385-da35-44e9-b66c-e21feae3745f" // itioapi

// for Live
//var privateKeyPath = "./fireblocks_secret_live.key"
//var apiKey = "7cd91e4b-012c-4ea8-97d5-22f6761c0d2a"

// Get user list from fireblock user by API
func UsersView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Printf("Error initializing token provider: %v\n", err)
		//return
	}

	// Example API calls
	// Ensure functions getAccountsPaged and createAccount are correctly defined to use this main function.

	path := "/v1/users"
	respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
	if err != nil {
		return fmt.Errorf("error making GET request to accounts_paged: %w", err)
	}

	//var DataList = string(respBody)
	//fmt.Println("====== !! >>", DataList)
	//////////////////////////

	// Parse the JSON data into the struct
	var fireblocksData []models.FireblocksUsers
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	return c.Render("admin/fireblocks-users", fiber.Map{
		"Title":     "Fire Blocked User List",
		"Subtitle":  "User List",
		"Data":      fireblocksData,
		"AdminData": adminData,
	})
}

// Get Vault Wallet list from fireblock API
func CreateVaultWallet(c *fiber.Ctx) error {

	VID := c.FormValue("VID")
	WID := c.FormValue("WID")

	//fmt.Println(VID, WID)

	// check session
	s, err := store.Get(c)
	if err != nil {
		//panic(err)
		fmt.Println("111", err)
	}
	// Get value

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	path := "/v1/vault/accounts/" + VID + "/" + WID
	//fmt.Println("======", path)
	respBody, err := fireblocks.MakeAPIRequest("POST", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}

	//var DataList = string(respBody)
	//fmt.Println("======", DataList)

	//FireblocksWallet
	var fireblocksData models.FireblocksWallet
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	//fmt.Println("======", fireblocksData.Message)
	if fireblocksData.Message == "" {
		s.Set("Alerts", "Wallet Creates Successfully with ID : "+fireblocksData.Id)

		/// Insert Data in table

		qry := models.AddWallet{Volt_id: VID, Coin: WID, Address: fireblocksData.Address, Legacyaddress: fireblocksData.LegacyAddress, Tag: fireblocksData.Tag}
		result := database.DB.Db.Table("wallet_list").Create(&qry)

		if result.Error != nil {
			fmt.Println("ERROR in QUERY")

		}

		/// End Insert Data
	} else {
		s.Set("Alerts", fireblocksData.Message)
	}

	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("222", err)
	}
	//Alerts := s.Get("Alerts")
	//fmt.Println("==>Message :: ", Alerts)

	//////////////////////////
	return c.Redirect("/admin/vault")
}

// Get Wallet list from fireblock API
func WalletView(c *fiber.Ctx) error {

	VID := c.Params("VID")
	WID := c.Params("WID")

	//fmt.Println(VID, WID)

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	Alerts := sess.Get("Alerts")
	sess.Delete("Alerts")

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	path := "/v1/vault/accounts/" + VID + "/" + WID + "/addresses_paginated"
	respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.FireblocksAddress
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	//return c.Render("vault", fireblocksData)
	return c.Render("admin/wallet", fiber.Map{
		"Title":     "Wallet Address",
		"Alert":     Alerts,
		"VoltID":    VID,
		"AssetID":   WID,
		"AdminData": adminData,
		"Assets":    fireblocksData.Addresses,
	})
}

// Get Vault list from fireblock API
func VoltView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	Alerts := sess.Get("Alerts")
	sess.Delete("Alerts")

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	//#############################//
	var voltID = "1"
	//fmt.Println("LoginVoltID -> ", voltID)
	var fireblocksData models.FireblocksResponse

	if voltID != "" {
		path := "/v1/vault/accounts/" + voltID //+ voltID
		//fmt.Println(path)
		respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(string(respBody))
		// Parse the JSON data into the struct
		//var fireblocksData models.FireblocksResponse
		if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
			fmt.Println(err)
		}

	}

	//fmt.Println(fireblocksData)
	//return c.Render("vault", fireblocksData)
	return c.Render("admin/vault", fiber.Map{
		"Title":       "Wallet List",
		"Subtitle":    "Wallet List",
		"Alert":       Alerts,
		"LoginVoltID": voltID,
		"AdminData":   adminData,
		"ID":          fireblocksData.ID,
		"Name":        fireblocksData.Name,
		"HiddenOnUI":  fireblocksData.HiddenOnUI,
		"AutoFuel":    fireblocksData.AutoFuel,
		"Assets":      fireblocksData.Assets,
	})
}

// Get Wallet list from fireblock API
func WalletListView(c *fiber.Ctx) error {
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")

	Alerts := s.Get("Alerts")

	//#############################//s
	var voltID = s.Get("LoginVoltID").(string)
	WalletList := []models.WalletList{}
	database.DB.Db.Table("wallet_list").Order("wallet_id desc").Where("volt_id = ?", voltID).Find(&WalletList)
	//.Select("login_time")

	s.Delete("Alerts")

	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("333", err)
	}

	//fmt.Println(WalletList)
	//return c.Render("vault", fireblocksData)
	return c.Render("wallet-list", fiber.Map{
		"Title":        "Wallet List",
		"Subtitle":     "Wallet List",
		"ID":           voltID,
		"Alert":        Alerts,
		"WalletList":   WalletList,
		"MerchantData": merchantData,
	})
}

// create Volt Wallet  from fireblock API
func CreateVaultWalletAddress(c *fiber.Ctx) error {

	VID := c.Params("VID")
	WID := c.Params("WID")

	//fmt.Println(VID, WID)

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	//adminData := sess.Get("AdminData")

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	path := "/v1/vault/accounts/" + VID + "/" + WID + "/addresses"
	respBody, err := fireblocks.MakeAPIRequest("POST", path, nil, tokenProvider)
	//respBody = nil
	fmt.Println(respBody)
	if err != nil {
		fmt.Println(err)
	}

	///////////////////////////////////
	Alerts := "Address Generated Successfully "
	sess.Set("Alerts", Alerts)
	if err := sess.Save(); err != nil {
		//panic(err)
		fmt.Println("444", err)
	}
	path = "/admin/wallet/" + VID + "/" + WID
	return c.Redirect(path)
}

// create New Volt Wallet  from fireblock API
func CreateNewVault(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	Alerts := "Account Generated Successfully"
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	LoginMerchantEmail := s.Get("LoginMerchantEmail").(string)
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}

	path := "/v1/vault/accounts"
	MyData := struct {
		Name string `json:"name"`
	}{
		Name: LoginMerchantEmail,
	}

	respBody, err := fireblocks.MakeAPIRequest("POST", path, MyData, tokenProvider)
	if err != nil {
		fmt.Println(err)
		Alerts = "Account Not Generated"
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.CreateVaultAccountResponse
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	s.Set("LoginVoltID", fireblocksData.ID)
	//fmt.Println(fireblocksData.ID)
	///////////////////////////////////

	//fmt.Println(fireblocksData.ID)
	if fireblocksData.ID != "" {

		Voltid := fireblocksData.ID
		LoginID := LoginMerchantID.(uint)
		result := database.DB.Db.Table("client_master").Save(&models.UpdateVolt{Client_id: LoginID, Volt_id: Voltid})

		if result.Error != nil {
			fmt.Println("ERROR in QUERY")
			Alerts = "Account Not Generated - 2"
		}

	}
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}
	return c.Redirect("/wallet-list")
}

// Update Wallet Balance  from fireblock API
func UpdateWalletBalance(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	Alerts := "Balance Updated"
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	VID := c.Params("VID")
	WID := c.Params("WID")

	numberStr := c.Params("AID")

	// Convert string to uint
	number, err := strconv.ParseUint(numberStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid number format")
	}
	// Convert uint64 to uint
	AID := uint(number)

	//fmt.Println(VID, WID, AID)
	//GetWID, err := strconv.ParseUint(AID, 10, 64)
	//GetWID = GetWID.(uint)

	path := "/v1/vault/accounts/" + VID + "/" + WID + "/balance"
	//fmt.Println(path)
	respBody, err := fireblocks.MakeAPIRequest("POST", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
		Alerts = "Account Not Generated"
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.WalletListBalance
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}
	//fmt.Println("<=============>")
	//fmt.Println(fireblocksData)
	//fmt.Println("<=============>")
	///////////////////////////////////

	if fireblocksData.Total != "" {

		result := database.DB.Db.Table("wallet_list").Save(&models.WalletListBalance{Wallet_id: AID, Total: fireblocksData.Total, Available: fireblocksData.Available, Pending: fireblocksData.Pending, Frozen: fireblocksData.Frozen, Lockedamount: fireblocksData.Lockedamount})

		if result.Error != nil {
			fmt.Println("ERROR in QUERY")
			Alerts = "Account Not Generated - 2"
		}

	} else {
		Alerts = "Error : Balance Not Updated"

	}
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}
	return c.Redirect("/wallet-list")
}

// view Vault Accounts from fireblock API
func VaultAccountsView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}

	path := "/v1/vault/accounts_paged"

	respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.FireblocksVoltResponse
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	return c.Render("admin/vault-accounts", fiber.Map{
		"Title":     "Volt List",
		"Subtitle":  "Volt List",
		"AdminData": adminData,
		"Data":      fireblocksData.Accounts,
	})
}

// view Vault Wallet from fireblock API
func CreateVaultWalletView(c *fiber.Ctx) error {
	VID := c.Params("VID")

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	Alerts := ""

	return c.Render("admin/create-wallet", fiber.Map{
		"Title":     "Create Wallet",
		"Subtitle":  "Create Wallet",
		"Alert":     Alerts,
		"VoltID":    VID,
		"AdminData": adminData,
	})
}

// Display Transfer from fireblock API
func TransferView(c *fiber.Ctx) error {
	//VID := c.Params("VID")

	// check session
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	voltID := s.Get("LoginVoltID")
	fmt.Println("LoginVoltID -> ", voltID, LoginMerchantID)
	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("3434343=>>", err)
	}
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	return c.Render("transfer", fiber.Map{
		"Title":        "Transfer",
		"Subtitle":     "Transfer",
		"Alert":        Alerts,
		"MerchantData": merchantData,
	})
}

// Submit Transfer from fireblock API
func TransferPost(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	Alerts := "Transfer Process"
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	voltID := s.Get("LoginVoltID").(string)
	LoginMerchantEmail := s.Get("LoginMerchantEmail").(string)
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}

	randomID, err := function.GenerateRandomID(16) // 16 bytes will give us a 32 character hex string
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate random ID",
		})
	}

	// Convert string to uint
	amt := c.FormValue("amount")

	getAmount, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		fmt.Println(err)
	}
	// Convert uint64 to uint
	//AID := uint(number)

	fmt.Println("random No -", randomID)
	AssetId := c.FormValue("assetId")
	Amount := getAmount
	Address := c.FormValue("address")
	FeeLevel := "HIGH"
	Note := "Sending " + AssetId + " to Addresses - " + Address
	Operation := "TRANSFER"
	CustomerRefId := randomID
	Tag := randomID
	Transaction_type := c.FormValue("transaction_type")

	MyData := models.TransferRequest{
		AssetId: AssetId,
		Source: models.Source{
			Type: "VAULT_ACCOUNT",
			ID:   voltID,
		},
		Destination: models.Destination{
			OneTimeAddress: models.OneTimeAddress{
				Address: Address,
				Tag:     Tag,
			},
			Type: "ONE_TIME_ADDRESS",
		},
		Amount:        Amount,
		FeeLevel:      FeeLevel,
		Note:          Note,
		Operation:     Operation,
		CustomerRefId: CustomerRefId,
	}

	path := "/v1/transactions"

	respBody, err := fireblocks.MakeAPIRequest("POST", path, MyData, tokenProvider)
	if err != nil {
		fmt.Println(err)
		Alerts = "Transaction Failed"
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.TransResponce
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData.ID)
	if fireblocksData.ID != "" {

		TransID := fireblocksData.ID
		TransStatus := fireblocksData.Status
		Message := "Transaction Processed with Status - " + TransStatus + " And Trans Id - " + TransID
		Alerts = Message
		//CID := int(LoginMerchantID)
		CID := s.Get("LoginMerchantID").(uint)
		//fmt.Println("=====>>>", CID)
		Ip := c.Context().RemoteIP().String()
		qry := models.Transaction_Master{Client_id: CID, Transaction_id: TransID, Volt_id: voltID, Assetid: AssetId, Amount: Amount, Operation: Operation, Customerrefid: CustomerRefId, Status: TransStatus, Transaction_type: Transaction_type, Ip: Ip, Note: Note, Source: LoginMerchantEmail}
		result := database.DB.Db.Table("transaction_master").Select("client_id", "transaction_id", "volt_id", "assetid", "amount", "operation", "customerrefid", "status", "transaction_type", "ip", "note", "source").Create(&qry)
		//fmt.Println(result)

		if result.Error != nil {
			fmt.Println(result.Error)
		}

		receivedId := qry.Id
		//fmt.Println("XXX ", receivedId)

		path := "/v1/transactions/" + TransID
		respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
		if err != nil {
			return fmt.Errorf("error making GET request to accounts_paged: %w", err)
		}

		// Parse the JSON data into the struct
		var transGetResponce models.TransGetResponce
		if err := json.Unmarshal([]byte(respBody), &transGetResponce); err != nil {
			fmt.Println(err)
		}

		result = database.DB.Db.Table("transaction_master").Save(&models.TransGetResponce{Id: receivedId, Substatus: transGetResponce.Substatus, Status: transGetResponce.Status, Txhash: transGetResponce.Txhash, Requestedamount: transGetResponce.Requestedamount, Netamount: transGetResponce.Netamount, Amountusd: transGetResponce.Amountusd, Fee: transGetResponce.Fee, Networkfee: transGetResponce.Networkfee, Destinationaddress: transGetResponce.Destinationaddress, Createdby: transGetResponce.Createdby})

		if result.Error != nil {
			fmt.Println("ERROR in QUERY", result.Error)
		}

	} else {
		TransMessage := fireblocksData.Message
		Message := "Transaction failed with message - " + TransMessage
		Alerts = Message

	}
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}
	return c.Redirect("/transactions")
}
