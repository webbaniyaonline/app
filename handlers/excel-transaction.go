package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// function for generate  Excel with Transaction List
func TransExcel(c *fiber.Ctx) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		//return
	}

	// Set value of a cell
	f.SetCellValue("Sheet1", "A1", "S.No.")
	f.SetCellValue("Sheet1", "B1", "TransactionID")
	f.SetCellValue("Sheet1", "C1", "Asset")
	f.SetCellValue("Sheet1", "D1", "Requested Amount")
	f.SetCellValue("Sheet1", "E1", "Converted Amount")
	f.SetCellValue("Sheet1", "F1", "Received Amount")
	f.SetCellValue("Sheet1", "G1", "Status")
	f.SetCellValue("Sheet1", "H1", "Trans Type")
	f.SetCellValue("Sheet1", "I1", "Txhash")
	f.SetCellValue("Sheet1", "J1", "Timestamp")

	// check session
	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")

	mid := c.Query("mid") // Get from Admin
	if mid != "" {
		LoginMerchantID = mid
	}

	transactionList := []models.Transaction_Pay{}
	if LoginMerchantID == nil {
		database.DB.Db.Table("transaction").Order("id desc").Find(&transactionList)
	} else {
		database.DB.Db.Table("transaction").Order("id desc").Where("client_id = ?", LoginMerchantID).Find(&transactionList)
	}
	//fmt.Println(transactionList)

	// Fill the Excel sheet with some data
	i := 1
	for _, p := range transactionList {
		i = i + 1

		RequestedCurrencyFull := strings.ToUpper(p.Requestedcurrency) + " " + fmt.Sprintf("%f", p.Requestedamount)
		ConvertedAmountFull := strings.ToUpper(p.Convertedcurrency) + " " + fmt.Sprintf("%f", p.Convertedamount)
		ReceivedAmountFull := strings.ToUpper(p.Convertedcurrency) + " " + fmt.Sprintf("%f", p.Receivedamount)
		SubStatus := function.GetSubStatusByStatusID(p.Substatus)

		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), i-1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), p.Transaction_id)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), strings.ToUpper(p.Receivedcurrency))
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), RequestedCurrencyFull)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), ConvertedAmountFull)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), ReceivedAmountFull)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), SubStatus)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i), p.Transaction_type)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i), p.Response_hash)
		f.SetCellValue("Sheet1", "J"+strconv.Itoa(i), p.Createdate)
	}

	// Set the active sheet
	f.SetActiveSheet(index)

	// Save the file to a buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return err
	}

	// Set the appropriate headers for downloading an Excel file
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=Transaction-List.xlsx")
	c.Set("Content-Length", string(len(buf.Bytes())))

	// Send the file
	return c.SendStream(buf)
}
