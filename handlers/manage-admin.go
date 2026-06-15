package handlers

import (
	"fmt"
	"os"
	"strconv"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/crypto/bcrypt"
)

var tableNameA = "admin_master"
var pageNameA = "admin/admin-manager"
var pageTitleA = "Admin"
var listOrderByA = "status ASC,full_name ASC"

// Page for manage admin user and profile

// function for display added admin list
func GetAdminList(c *fiber.Ctx) error {

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

	adminList := []models.AdminList{}

	var total int64
	database.DB.Db.Table(tableNameA).Order(listOrderByA).Limit(limit).Offset(offset).Find(&adminList).Count(&total)

	return c.Render(pageNameA, fiber.Map{
		"Title":     pageTitleA,
		"Subtitle":  pageTitleA,
		"Action":    "List",
		"AlertX":    Alerts,
		"DataList":  adminList,
		"AdminData": adminData,
		"Page":      page,
		"Limit":     limit,
		"Total":     total,
	})
}

// function for Display Add Admin Form
func AddAdminView(c *fiber.Ctx) error {

	AdminSession(c)
	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	return c.Render(pageNameA, fiber.Map{
		"Title":     pageTitleA,
		"Subtitle":  pageTitleA,
		"Action":    "Add",
		"AdminData": adminData,
	})
}

// function for Submit Admin  Form
func AdminPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	username := c.FormValue("username")
	full_name := c.FormValue("full_name")
	role := c.FormValue("role")
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 1041")
	}
	status := int(status1)

	tableID := c.FormValue("tableID")
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 1057")
	}
	getTableID := uint(cid)

	var password = function.PasswordGenerator(8)
	//fmt.Println(password)

	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//////////
	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table(tableNameA).Save(&models.AdminList{Admin_id: getTableID, Username: username, Full_name: full_name, Password: string(hash), Role: role, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := pageTitleA + " Processed successfully with Password - " + password
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitleA + " Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageNameA)
}

// function for Display Edit Admin Form
func EditAdmin(c *fiber.Ctx) error {

	AdminSession(c)
	tableID := c.Params("TID")

	// Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	data := models.AdminList{}
	database.DB.Db.Table(tableNameA).Where("admin_id = ?", tableID).Find(&data)

	return c.Render(pageNameA, fiber.Map{
		"Title":     pageTitleA,
		"Subtitle":  pageTitleA,
		"Action":    "Edit",
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for De Active Admin
func DeleteAdmin(c *fiber.Ctx) error {
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
	result := database.DB.Db.Table(tableNameA).Save(&models.AdminDeleted{Admin_id: getTableID, Status: status})

	Alerts := pageTitleA + " Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitleA + " Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/" + pageNameA)
}

// function for Display profile Form
func ProfileForm(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	if adminData == nil {
		return c.Redirect("/admin/login")
	}
	// Convert the session data to a map
	adminMap := adminData.(map[string]interface{})
	//fmt.Println(adminMap["AdminID"])
	LoginAdminID := adminMap["AdminID"].(uint)
	fmt.Print(LoginAdminID)
	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	data := models.AdminList{}
	database.DB.Db.Table(tableNameA).Where("admin_id = ?", LoginAdminID).Find(&data)

	return c.Render("admin/profile", fiber.Map{
		"Title":     "Profile",
		"Subtitle":  "Profile",
		"AlertX":    Alerts,
		"AdminData": adminData,
		"EditData":  data,
	})
}

// function for Submit profile Form
func ProfileFormPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	full_name := c.FormValue("full_name")
	role := c.FormValue("role")
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 1042")
	}
	status := int(status1)

	tableID := c.FormValue("tableID")
	cid, err := strconv.ParseUint(tableID, 10, 32)
	if err != nil {
		fmt.Println("Error 1055")
	}
	getTableID := uint(cid)

	//////////
	// if GET ID than work update else insert
	// for Full path use- filePath & only file name use file.Filename
	result := database.DB.Db.Table(tableNameA).Save(&models.AdminUpdate{Admin_id: getTableID, Full_name: full_name, Role: role, Status: status})

	//fmt.Println(loginList.Status)
	Alerts := pageTitleA + " update successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = pageTitleA + " Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/profile")
}

// function for Display Change Password Form
func ChangePasswordForm(c *fiber.Ctx) error {

	AdminSession(c)
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	Alerts := sess.Get("AlertX")
	if Alerts != "" {
		sess.Delete("AlertX")
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	return c.Render("admin/change-password", fiber.Map{
		"Title":     "Change Password",
		"Subtitle":  "Change Password",
		"AlertX":    Alerts,
		"AdminData": adminData,
	})
}

// function for post Change Password Form
func ChangePasswordFormPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	new_password := c.FormValue("new_password")
	confirm_password := c.FormValue("confirm_password")
	Alerts := ""
	if new_password == confirm_password {

		tableID := c.FormValue("tableID")
		cid, err := strconv.ParseUint(tableID, 10, 32)
		if err != nil {
			fmt.Println("Error 1056")
		}
		getTableID := uint(cid)

		//fmt.Println(password)

		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)

		//////////
		// if GET ID than work update else insert
		// for Full path use- filePath & only file name use file.Filename
		result := database.DB.Db.Table(tableNameA).Save(&models.AdminPassword{Admin_id: getTableID, Password: string(hash)})

		//fmt.Println(loginList.Status)
		Alerts = "Password update successfully"
		if result.Error != nil {
			//fmt.Println("ERROR in QUERY")
			Alerts = "Password Not Updated"
		}
	} else {
		Alerts = "Password and confirm password not matched"
	}
	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/change-password")
}
