package handlers

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"template/database"
	"template/function"
	"template/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf/v2"
)

// Page for Generate Transaction list data into PDF format

func TransPDF(c *fiber.Ctx) error {

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
	pdfBytes := generatePDFG(transactionList)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=transaction-details.pdf")
	return c.Send(pdfBytes)

}

// Function to generate PDF
func generatePDFG(transactionList []models.Transaction_Pay) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Transaction Details")

	pdf.Ln(12) // Line break

	pdf.SetFont("Arial", "", 6)
	//header := []string{"TransactionID", "Asset", "Trans Type", "Amount", "Source", "Destination", "Status", "Timestamp"}

	// Print header
	//for _, str := range header {
	pdf.CellFormat(10, 7, "S.No.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, 7, "TransactionID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(10, 7, "Asset", "1", 0, "C", false, 0, "")
	pdf.CellFormat(15, 7, "Trans Type", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Requested Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Converted Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Received Amount", "1", 0, "C", false, 0, "")
	pdf.CellFormat(15, 7, "Status", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Timestamp", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	//}
	//pdf.Ln(-1)

	// Print data
	i := 0
	for _, p := range transactionList {
		i = i + 1
		ConvertedamountFull := strings.ToUpper(p.Convertedcurrency) + " " + fmt.Sprintf("%f", p.Convertedamount)
		ReceivedamountFull := strings.ToUpper(p.Convertedcurrency) + " " + fmt.Sprintf("%f", p.Receivedamount)
		RequestedcurrencyFull := strings.ToUpper(p.Requestedcurrency) + " " + fmt.Sprintf("%f", p.Requestedamount)
		SubStatus := function.GetSubStatusByStatusID(p.Substatus)
		pdf.CellFormat(10, 7, fmt.Sprintf("%d", i), "1", 0, "", false, 0, "")
		pdf.CellFormat(45, 7, p.Transaction_id, "1", 0, "", false, 0, "")
		pdf.CellFormat(10, 7, strings.ToUpper(p.Receivedcurrency), "1", 0, "", false, 0, "")
		pdf.CellFormat(15, 7, p.Transaction_type, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 7, RequestedcurrencyFull, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 7, ConvertedamountFull, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 7, ReceivedamountFull, "1", 0, "", false, 0, "")
		pdf.CellFormat(15, 7, SubStatus, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, p.Createdate, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	buf := new(bytes.Buffer)
	err := pdf.Output(buf)
	if err != nil {
		log.Fatalf("Failed to create PDF: %s", err)
	}

	return buf.Bytes()
}
func TransNPPdf(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")

	transactionList := []models.Transaction_MasterNP{}
	database.DB.Db.Table("transaction_master_nowpayments").Order("tid desc").Where("client_id = ?", LoginMerchantID).Find(&transactionList)
	//fmt.Println(transactionList)
	pdfBytes := generateNPPDFG(transactionList)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=transaction-details-np.pdf")
	return c.Send(pdfBytes)

}

// Function to generate PDF
func generateNPPDFG(transactionList []models.Transaction_MasterNP) []byte {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Transaction Details")

	pdf.Ln(12) // Line break

	pdf.SetFont("Arial", "", 6)
	//header := []string{"TransactionID", "Asset", "Trans Type", "Amount", "Source", "Destination", "Status", "Timestamp"}

	// Print header
	//for _, str := range header {
	pdf.CellFormat(8, 7, "S.No.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(15, 7, "Payment ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Order ID", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Original price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Amount sent / received", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "Status", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Created", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Last update", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	//}
	//pdf.Ln(-1)

	// Print data
	i := 0
	for _, p := range transactionList {
		i = i + 1
		var vamt = fmt.Sprintf("%f", p.Price_amount) + "" + p.Price_currency
		pdf.CellFormat(8, 7, fmt.Sprintf("%d", i), "1", 0, "", false, 0, "")
		pdf.CellFormat(15, 7, p.Payment_id, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 7, p.Order_id, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 7, vamt, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%f", p.Pay_amount), "1", 0, "", false, 0, "")
		pdf.CellFormat(12, 7, p.Payment_status, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, p.Created_at, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, p.Updated_at, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	buf := new(bytes.Buffer)
	err := pdf.Output(buf)
	if err != nil {
		log.Fatalf("Failed to create PDF: %s", err)
	}

	return buf.Bytes()
}
