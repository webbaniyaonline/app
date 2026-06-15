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

// For Add / Edit / Delete  Acquirer from Admin Section

var tableNameAcq = "acquirer"
var pageNameAcq = "admin/acquirer"
var pageTitleAcq = "Acquirer"
var listOrderByAcq = "status ASC,acquirer_title ASC"

// function for Display Acquirer List
func AcquirerList(c *fiber.Ctx) error {

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

	acquirerList := []models.AcquirerList{}

	var total int64
	database.DB.Db.Table(tableNameAcq).Order(listOrderByAcq).Limit(limit).Offset(offset).Find(&acquirerList).Count(&total)

	//fmt.Println(acquirerList)
	return c.Render(pageNameAcq, fiber.Map{
		"Title":        pageTitleAcq,
		"Subtitle":     pageTitleAcq,
		"Action":       "List",
		"AlertX":       Alerts,
		"AcquirerList": acquirerList,
		"AdminData":    adminData,
		"Page":         page,
		"Limit":        limit,
		"Total":        total,
	})
}

// function for Display Acquirer Form
func AddAcquirerForm(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	return c.Render(pageNameAcq, fiber.Map{
		"Title":     pageTitleAcq,
		"Subtitle":  pageTitleAcq,
		"Action":    "Add",
		"AdminData": adminData,
	})
}

// function for Post Acquirer Form
func AcquirerPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	acquirer_title := c.FormValue("acquirer_title")
	api_key := c.FormValue("api_key")
	secret_key := c.FormValue("secret_key")
	endpoint_url := c.FormValue("endpoint_url")
	callback_url := c.FormValue("callback_url")
	success_url := c.FormValue("success_url")
	failed_url := c.FormValue("failed_url")
	json_data := c.FormValue("json_data")
	if json_data == "" {
		json_data = "{}"
	}
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
	result := database.DB.Db.Table(tableNameAcq).Save(&models.AcquirerList{Acquirer_id: getTableID, Acquirer_title: acquirer_title, Api_key: api_key, Secret_key: secret_key, Endpoint_url: endpoint_url, Callback_url: callback_url, Success_url: success_url, Failed_url: failed_url, Json_data: json_data, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := pageTitleAcq + " Processed successfully"
	if result.Error != nil {
		fmt.Println("ERROR in QUERY", result.Error)
		Alerts = pageTitleAcq + " Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageNameAcq)
}

// function for Post Add / Edit Acquirer Form
func EditAcquirer(c *fiber.Ctx) error {

	AdminSession(c)
	tableID := c.Params("TID")

	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	data := models.AcquirerList{}
	database.DB.Db.Table(tableNameAcq).Where("acquirer_id = ?", tableID).Find(&data)

	return c.Render(pageNameAcq, fiber.Map{
		"Title":     pageTitleAcq,
		"Subtitle":  pageTitleAcq,
		"Action":    "Edit",
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for Delete Acquirer
func DeleteAcquirer(c *fiber.Ctx) error {
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
	result := database.DB.Db.Table(tableNameAcq).Save(&models.AcquirerDeleted{Acquirer_id: getTableID, Status: status})

	Alerts := pageTitleAcq + " Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitleAcq + " Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageNameAcq)

}
