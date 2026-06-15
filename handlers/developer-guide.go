package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"template/database"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// For Add / Edit / Delete Developer Guide from Admin Sections

// function for Display Developer Guide List
func DeveloperGuide(c *fiber.Ctx) error {

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
	limit, _ := strconv.Atoi(c.Query("limit", "1000")) //os.Getenv("PagingSize")
	offset := (page - 1) * limit

	guideList := []models.DeveloperGuide{}

	var total int64
	database.DB.Db.Table("developer_guide").Order("title ASC,heading ASC").Limit(limit).Offset(offset).Find(&guideList).Count(&total)

	//fmt.Println(currencyList)
	return c.Render("admin/developer-guide", fiber.Map{
		"Title":     "Developer Guide",
		"Subtitle":  "Developer Guide",
		"Action":    "List",
		"AlertX":    Alerts,
		"GuideList": guideList,
		"AdminData": adminData,
		"Page":      page,
		"Limit":     limit,
		"Total":     total,
	})
}

// function for Display Developer Guide Form
func AddGuideView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	return c.Render("admin/developer-guide", fiber.Map{
		"Title":     "Developer Guide",
		"Subtitle":  "Developer Guide",
		"Action":    "Add",
		"AdminData": adminData,
	})
}

// function for Post Add / Edit Developer Guide Form
func GuidePost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	title := strings.TrimSpace(c.FormValue("title"))
	heading := strings.TrimSpace(c.FormValue("heading"))
	functions := strings.TrimSpace(c.FormValue("functions"))
	used := strings.TrimSpace(c.FormValue("used"))
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
	result := database.DB.Db.Table("developer_guide").Save(&models.DeveloperGuide{Id: getTableID, Title: title, Heading: heading, Functions: functions, Used: used, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := "Developer Guide Processed successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Developer Guide Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/add-guide")
}

// function for Post Add / Edit Developer Guide Form
func EditGuide(c *fiber.Ctx) error {

	AdminSession(c)
	tableID := c.Params("TID")

	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	data := models.DeveloperGuide{}
	database.DB.Db.Table("developer_guide").Where("id = ?", tableID).Find(&data)

	return c.Render("admin/developer-guide", fiber.Map{
		"Title":     "Developer Guide",
		"Subtitle":  "Developer Guide",
		"Action":    "Edit",
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for Developer Guide
func DeleteGuide(c *fiber.Ctx) error {
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
	result := database.DB.Db.Table("developer_guide").Save(&models.DeveloperGuideDeleted{Id: getTableID, Status: status})

	Alerts := "Developer Guide Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Developer Guide Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/developer-guide")

}
