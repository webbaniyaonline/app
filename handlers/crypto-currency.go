package handlers

import (
	"fmt"
	"os"
	"strconv"
	"template/database"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// For Add / Edit / Delete Crypto Currency from Admin Section

// function for Display Currency List
func GetCryptoCurrencyList(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
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

	currencyList := []models.CryptoCurrencyList{}

	var total int64
	database.DB.Db.Table("crypto_currency").Order("crypto_network_short ASC,crypto_code ASC").Limit(limit).Offset(offset).Find(&currencyList).Count(&total)

	//fmt.Println(currencyList)
	return c.Render("admin/crypto-currency", fiber.Map{
		"Title":        "Crypto Currency",
		"Subtitle":     "Crypto Currency",
		"Action":       "List",
		"AlertX":       Alerts,
		"CurrencyList": currencyList,
		"AdminData":    adminData,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
	})
}

// function for Display Currency Form view
func AddCryptoCurrencyView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	return c.Render("admin/crypto-currency", fiber.Map{
		"Title":     "Crypto Currency",
		"Subtitle":  "Crypto Currency",
		"Action":    "Add",
		"AdminData": adminData,
	})
}

// function for Post Add / Edit Currency Form
func CryptoCurrencyPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	crypto_code := c.FormValue("crypto_code")
	crypto_title := c.FormValue("crypto_title")
	crypto_network := c.FormValue("crypto_network")
	crypto_network_short := c.FormValue("crypto_network_short")
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")
	}
	status := int(status1)

	tableID := c.FormValue("tableID")
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	getTableID := uint(cid)
	//////////

	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table("crypto_currency").Save(&models.CryptoCurrencyList{Crypto_id: getTableID, Crypto_code: crypto_code, Crypto_title: crypto_title, Crypto_network: crypto_network, Crypto_network_short: crypto_network_short, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := "Crypto Currency Processed successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Crypto Currency Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/crypto-currency")
}

// function for Post Add / Edit Currency Form
func EditCryptoCurrency(c *fiber.Ctx) error {

	AdminSession(c)
	tableID := c.Params("TID")

	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	data := models.CryptoCurrencyList{}
	database.DB.Db.Table("crypto_currency").Where("crypto_id = ?", tableID).Find(&data)

	return c.Render("admin/crypto-currency", fiber.Map{
		"Title":     "Crypto Currency",
		"Subtitle":  "Crypto Currency",
		"Action":    "Edit",
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for Delete Currency
func DeleteCryptoCurrency(c *fiber.Ctx) error {
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
	result := database.DB.Db.Table("crypto_currency").Save(&models.CryptoCurrencyDeleted{Crypto_id: getTableID, Status: status})

	Alerts := "Crypto Currency Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Crypto Currency Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/crypto-currency")

}
