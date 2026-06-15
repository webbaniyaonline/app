package handlers

import (
	"strconv"
	"strings"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
)

// Function for Post Settlement form in merchant section
func SettlementSettingsPost(c *fiber.Ctx) error {
	s, _ := store.Get(c)
	// Get Data from ajax
	req := new(models.SettlementRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	//println("req=>", req)

	alerts := ""
	qrs := ""
	if req.Action == "Update" && req.CoinAddress != "" && req.Coin_id != 0 {

		settlementCoin := models.SettlementCoin{}
		database.DB.Db.Table("coin_list").Where("coin_id = ?", req.Coin_id).Find(&settlementCoin)

		crypto_code := settlementCoin.Coin
		crypto_title := settlementCoin.Coin_title
		crypto_network := settlementCoin.Coin_network

		if req.CoinAddress != "" && crypto_code != "" && crypto_code != "arb" {

			message, valid := function.CryptoAddressValidator(strings.TrimSpace(req.CoinAddress), crypto_code)
			//fmt.Println("message =>", message)
			//fmt.Println("valid =>", valid)
			if !valid {
				//fmt.Println("Not valid =>", message)
				response := models.SettlementResponse{
					Status:  2,
					Message: message,
				}
				return c.JSON(response)
			}

		}
		//fmt.Println(crypto_code, crypto_title, crypto_network)

		getTableID := uint(req.Wallet_id)
		//////////
		client_id := s.Get("LoginMerchantID").(uint)
		database.DB.Db.Table("client_wallet").Save(&models.CryptoWalletList{Wallet_id: getTableID, Client_id: client_id, Assetid: req.Coin_id, Crypto_address: req.CoinAddress, Crypto_code: crypto_code, Crypto_title: crypto_title, Crypto_network: crypto_network, Status: 1})
		alerts = "Settlement Addresses Updated"
		qrs = " Wallet_id: " + strconv.FormatUint(uint64(getTableID), 10) + " Assetid: " + strconv.Itoa(req.Coin_id) + "Crypto_address:" + req.CoinAddress + " Crypto_code:" + crypto_code + " Crypto_title: " + crypto_title + " Crypto_network: " + crypto_network
		updateIp := c.Context().RemoteIP().String()
		function.UpdateMerchantHistory("Settlement", qrs, updateIp, client_id)

	}

	if (req.Action == "deactivated" || req.Action == "activated") && req.Coin_id != 0 {
		coinStatus := "0"
		alerts = "Coin De Activated"
		if req.Action == "activated" {
			coinStatus = "1"
			alerts = "Coin Activated"
		}

		database.DB.Db.Table("client_wallet").Where("assetid = ?", req.Coin_id).UpdateColumns(&models.SettlementStatus{Status: coinStatus})
		qrs = " Change Status of Coin ID : " + strconv.Itoa(req.Coin_id) + " change Status to: " + coinStatus
		updateIp := c.Context().RemoteIP().String()
		client_id := s.Get("LoginMerchantID").(uint)
		function.UpdateMerchantHistory("Settlement", qrs, updateIp, client_id)
	}

	response := models.SettlementResponse{
		Status:  1,
		Message: alerts,
	}

	return c.JSON(response)
}

// function for Display Settlement Form
func SettlementSettingsView(c *fiber.Ctx) error {
	// check session
	MerchantSession(c) // redirect when session not found
	s, _ := store.Get(c)
	merchantData := s.Get("MerchantData")
	LoginMerchantID := s.Get("LoginMerchantID")

	//fmt.Print("LoginMerchantID =>", LoginMerchantID)
	settlementSetting := []models.SettlementSetting{}
	database.DB.Db.Table("coin_list A ").Select("a.coin_id, a.coin, a.coin_title, LOWER(a.coin_network) AS coin_network, b.crypto_address, b.assetid, b.status, b.wallet_id").Joins("LEFT JOIN client_wallet B ON A.coin_id = B.assetid AND B.client_id = ?", LoginMerchantID).Order("a.coin_title asc").Find(&settlementSetting)
	// .Where(" a.status = ?", 1)
	//fmt.Println("settlementSetting => ", settlementSetting)
	return c.Render("settlement-settings", fiber.Map{
		"Title":             "Settlement Settings",
		"Subtitle":          "Settlement Settings",
		"MerchantData":      merchantData,
		"SettlementSetting": settlementSetting,
	})
}
