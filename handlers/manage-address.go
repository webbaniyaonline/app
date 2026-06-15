package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"template/database"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// For Add / Edit / Delete / List  Coin Address from Admin Section

// function for Display Address List
func GetAddressList(c *fiber.Ctx) error {

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

	var total int64
	// For Address
	coinAddress := []models.AddressListing{}
	database.DB.Db.Table("coin_address").Order("status ASC,coin ASC").Limit(limit).Offset(offset).Find(&coinAddress)
	database.DB.Db.Table("coin_address").Count(&total)
	//fmt.Println("total => ", total)
	// For Address

	addressUrl := []models.CoinList{}
	database.DB.Db.Table("coin_list").Select("coin", "coin_id", "coin_network", "coin_status_url").Find(&addressUrl)

	return c.Render("admin/manage-address", fiber.Map{
		"Title":       "Manage Address",
		"Subtitle":    "Manage Address",
		"Action":      "List",
		"AlertX":      Alerts,
		"CoinAddress": coinAddress,
		"AdminData":   adminData,
		"Page":        page,
		"Limit":       limit,
		"Total":       total,
		"AddressUrl":  addressUrl,
	})
}

// function for Display Address Form
func AddAddressView(c *fiber.Ctx) error {

	AdminSession(c)
	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")
	cid := c.Query("CID")
	cname := c.Query("CNAME")

	return c.Render("admin/manage-address", fiber.Map{
		"Title":     "Coin List",
		"Subtitle":  "Coin List",
		"Action":    "Add",
		"CID":       cid,
		"CNAME":     cname,
		"AdminData": adminData,
	})
}

// function for Display Address Form for Edit
func EditAddress(c *fiber.Ctx) error {

	AdminSession(c)
	TID := c.Params("TID")

	// // Session Check
	sess, _ := store.Get(c)
	adminData := sess.Get("AdminData")

	coinData := models.AddressListing{}
	database.DB.Db.Table("coin_address").Where("address_id = ?", TID).Find(&coinData)

	return c.Render("admin/manage-address", fiber.Map{
		"Title":     "Manage Address",
		"Subtitle":  "Manage Address",
		"Action":    "Edit",
		"CoinData":  coinData,
		"CID":       coinData.Coin_id,
		"CNAME":     coinData.Coin,
		"AdminData": adminData,
	})
}

// function for Post Add / Edit Address Form
func AddressPost(c *fiber.Ctx) error {

	AdminSession(c)
	// Parses the request body
	address := c.FormValue("address")
	token_details := c.FormValue("token_details")
	status1, err := strconv.ParseInt(c.FormValue("status"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")

	}
	status := int(status1)
	coid := c.FormValue("Tableid")
	cid, err := strconv.ParseUint(coid, 10, 32)
	if err != nil {
		fmt.Println("Error 105")
	}
	address_id := uint(cid)

	coinid1, err := strconv.ParseInt(c.FormValue("coinid"), 10, 32)
	if err != nil {
		fmt.Println("Error 104")

	}
	coinid := int(coinid1)

	coinname := strings.ToLower(c.FormValue("coinname"))
	currentTime := time.Now()
	//formattedTime := currentTime.Format("2006-01-02 15:04:05")

	result := database.DB.Db.Table("coin_address").Save(&models.AddressListing{Address_id: address_id, Address: address, Status: status, Coin_id: coinid, Coin: coinname, Lastupdate: currentTime, Token_details: token_details})

	Alerts := "Processed successfully"
	if result.Error != nil {
		fmt.Println("ERROR in QUERY", result.Error)
		Alerts = "Address Not Updated"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/manage-address")
}

// function for Delete Address
func DeleteAddress(c *fiber.Ctx) error {
	AdminSession(c)
	id := c.Params("TID")

	// Convert string to uint
	number, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid number format")
	}
	// Convert uint64 to uint
	address_id := uint(number)

	status := 2
	result := database.DB.Db.Table("coin_address").Save(&models.AddressDeleted{Address_id: address_id, Status: status})

	Alerts := "Address Deleted successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Address Not Deleted"
	}

	// check session
	sess, _ := store.Get(c)
	sess.Set("AlertX", Alerts) // Set a session key
	if err := sess.Save(); err != nil {
		return err
	}

	return c.Redirect("/admin/manage-address")
}
