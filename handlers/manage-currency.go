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

var tableName = "currency"
var pageName = "admin/currency"
var pageTitle = "Currency"
var listOrderBy = "currency_id ASC"

// For Add / Edit / Delete / List  Currency from Admin Section

// function for Display Currency List
func GetCurrencyList(c *fiber.Ctx) error {

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

	currencyList := []models.CurrencyList{}

	var total int64
	database.DB.Db.Table(tableName).Order(listOrderBy).Limit(limit).Offset(offset).Find(&currencyList).Count(&total)

	//fmt.Println(currencyList)
	return c.Render(pageName, fiber.Map{
		"Title":        pageTitle,
		"Subtitle":     pageTitle,
		"Action":       "List",
		"AlertX":       Alerts,
		"CurrencyList": currencyList,
		"AdminData":    adminData,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
	})
}

// function for Display Currency Form
func AddCurrencyView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	return c.Render(pageName, fiber.Map{
		"Title":     pageTitle,
		"Subtitle":  pageTitle,
		"Action":    "Add",
		"AdminData": adminData,
	})
}

// function for Post Add / Edit Currency Form
func CurrencyPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	currency_name := c.FormValue("currency_name")
	currency_code := c.FormValue("currency_code")
	currency_territory := c.FormValue("currency_territory")
	currency_icon_bootstrap := c.FormValue("currency_icon_bootstrap")
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
	result := database.DB.Db.Table(tableName).Save(&models.CurrencyList{Currency_id: getTableID, Currency_name: currency_name, Currency_code: currency_code, Currency_territory: currency_territory, Currency_icon_bootstrap: currency_icon_bootstrap, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := pageTitle + " Processed successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitle + " Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageName)
}

// function for Post Add / Edit Currency Form
func EditCurrency(c *fiber.Ctx) error {

	AdminSession(c)
	tableID := c.Params("TID")

	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	data := models.CurrencyList{}
	database.DB.Db.Table(tableName).Where("currency_id = ?", tableID).Find(&data)

	return c.Render(pageName, fiber.Map{
		"Title":     pageTitle,
		"Subtitle":  pageTitle,
		"Action":    "Edit",
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for Delete Currency
func DeleteCurrency(c *fiber.Ctx) error {
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
	result := database.DB.Db.Table(tableName).Save(&models.CurrencyDeleted{Currency_id: getTableID, Status: status})

	Alerts := pageTitle + " Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitle + " Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageName)

}
