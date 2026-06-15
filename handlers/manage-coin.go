package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// For Add / Edit / Delete / List  Coin from Admin Section

// function for Display Coin List
func GetCoinList(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	// Get query parameters for page and limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", os.Getenv("PagingSize")))
	offset := (page - 1) * limit

	coinList := []models.CoinList{}
	var total int64
	database.DB.Db.Table("coin_list").Order("status ASC,coin ASC").Limit(limit).Offset(offset).Find(&coinList).Count(&total)

	// For Address
	coinAddress := []models.CoinAddress{}
	database.DB.Db.Table("coin_address").Where("status = ?", 1).Order("coin ASC").Find(&coinAddress)
	//fmt.Println(coinList)
	return c.Render("admin/coin-list", fiber.Map{
		"Title":       "Coin List",
		"Subtitle":    "Coin List",
		"Action":      "List",
		"AlertX":      Alerts,
		"CoinList":    coinList,
		"CoinAddress": coinAddress,
		"AdminData":   adminData,
		"Page":        page,
		"Limit":       limit,
		"Total":       total,
	})
}

// function for Display Coin Form
func AddCoinView(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	// For display Currency List on List Box
	currencyList := []models.CryptoCurrencyData{}
	database.DB.Db.Table("crypto_currency").Select("crypto_code", "crypto_title").Group("crypto_code, crypto_title").Order("crypto_code ASC").Where("status = ?", 1).Find(&currencyList)

	return c.Render("admin/coin-list", fiber.Map{
		"Title":        "Coin List",
		"Subtitle":     "Coin List",
		"Action":       "Add",
		"AdminData":    adminData,
		"CurrencyList": currencyList,
	})
}

// function for Display Coin Form for Edit
func EditCoin(c *fiber.Ctx) error {

	AdminSession(c)
	TID := c.Params("TID")
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	coinData := models.CoinList{}
	database.DB.Db.Table("coin_list").Where("coin_id = ?", TID).Find(&coinData)

	// For display Currency List on List Box
	currencyList := []models.CryptoCurrencyData{}
	database.DB.Db.Table("crypto_currency").Select("crypto_code", "crypto_title").Group("crypto_code, crypto_title").Order("crypto_code ASC").Where("status = ?", 1).Find(&currencyList)

	coinCode := strings.ToUpper(coinData.Coin)
	return c.Render("admin/coin-list", fiber.Map{
		"Title":         "Coin List",
		"Subtitle":      "Coin List",
		"Action":        "Edit",
		"CoinID":        coinData.Coin_id,
		"CoinCode":      coinCode,
		"CoinTitle":     coinData.Coin_title,
		"CoinNetwork":   coinData.Coin_network,
		"CoinPayUrl":    coinData.Coin_pay_url,
		"CoinIcon":      coinData.Icon,
		"CoinStatus":    coinData.Status,
		"CoinStatusUrl": coinData.Coin_status_url,
		"AdminData":     adminData,
		"CurrencyList":  currencyList,
	})
}

// function for Post Add / Edit Coin Form
func CoinPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	coin := strings.ToLower(c.FormValue("crypto_code"))
	coin_title := c.FormValue("crypto_title")
	coin_network := c.FormValue("crypto_network")
	coin_pay_url := strings.ToLower(c.FormValue("coin_pay_url"))
	coin_status_url := strings.ToLower(c.FormValue("coin_status_url"))
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
		//return c.Status(fiber.StatusBadRequest).SendString("Invalid number format 11")
	}
	status := int(status1)
	coid := c.FormValue("coinId")
	cid, err := strconv.ParseUint(coid, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
		//return c.Status(fiber.StatusBadRequest).SendString("Invalid number format 22")
	}
	coin_id := uint(cid)

	///////////  UPLOAD IMAGE/////////

	coinimg := c.FormValue("coinImg")
	uploadfile := ""
	file, err := c.FormFile("icon")
	if err != nil {
		//return c.Status(fiber.StatusBadRequest).SendString("Failed to read file")
		uploadfile = coinimg
	} else {

		uploadDir := "./views/images"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		filePath := fmt.Sprintf("%s/%s", uploadDir, file.Filename)
		err = c.SaveFile(file, filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
		}
		uploadfile = file.Filename
	}

	//////////

	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table("coin_list").Save(&models.CoinList{Coin_id: coin_id, Coin: coin, Icon: uploadfile, Status: status, Coin_title: coin_title, Coin_network: coin_network, Coin_pay_url: coin_pay_url, Coin_status_url: coin_status_url})

	//fmt.Println(loginList.Status)
	Alerts := "Processed successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Coin Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/coin-list")

}

// function for Delete Coin
func DeleteCoin(c *fiber.Ctx) error {
	AdminSession(c)
	id := c.Params("TID")

	// Convert string to uint
	number, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid number format")
	}
	// Convert uint64 to uint
	coin_id := uint(number)

	// convert to float
	//var item models.CoinList
	//database.DB.Db.Table("coin_list").First(&item, id)
	//database.DB.Db.Table("coin_list").Delete(&item)
	status := 2
	result := database.DB.Db.Table("coin_list").Save(&models.CoinDeleted{Coin_id: coin_id, Status: status})

	Alerts := "Coin Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Coin Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/coin-list")

}
