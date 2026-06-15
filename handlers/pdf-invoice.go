package handlers

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf/v2"
)

func InvoicePDF(c *fiber.Ctx) error {

	trackid := c.Query("iid")

	row := models.Invoice_Master{}
	database.DB.Db.Table("invoice").Where("trackid = ? AND invoice_type = ? AND status = ?", trackid, 2, 1).Find(&row)
	//fmt.Println("Data => ", row)
	RequesteAmount := strings.ToUpper(row.Requestedcurrency) + " " + fmt.Sprintf("%f", row.Requestedamount)
	ReceiverName := function.GetNameByMID(row.Client_id)
	ReceiverEmail := function.GetEmailByMID(row.Client_id)
	inputDate := row.Createdate

	// Parse the input date string
	parsedDate, err := time.Parse(time.RFC3339, inputDate)
	if err != nil {
		return c.Status(400).SendString("Invalid date format")
	}

	// Format the parsed date
	formattedDate := parsedDate.Format("02-01-2006") // Example: "24-09-2024"

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font for the document
	pdf.SetFont("Arial", "B", 16)

	// Create the invoice title
	pdf.Cell(40, 20, "Receipt")

	// Move to the next line
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(128, 128, 128) // text color
	// Create the invoice ID
	pdf.Cell(40, 10, "#"+trackid)
	// Move to the next line
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.SetTextColor(0, 0, 0) // Red text color for total
	// Add invoice details
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(50, 10, "Payment Date")
	pdf.Cell(50, 10, "Payment Method")
	pdf.Cell(50, 10, "Network")
	pdf.Cell(50, 10, "Amount Paid")
	pdf.Ln(12)
	// Add invoice Data
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(50, 10, formattedDate)
	pdf.Cell(50, 10, "Crypto Transfer")
	pdf.Cell(50, 10, "Base")
	pdf.Cell(50, 10, RequesteAmount)

	pdf.Ln(12)
	pdf.SetTextColor(128, 128, 128) // text color
	// Draw another line after invoice details
	pdf.Line(10, 70, 200, 70) // Line from x=10, y=70 to x=200, y=70
	pdf.Ln(12)

	pdf.SetTextColor(0, 0, 0) // Red text color for total
	// Add invoice details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 10, "Billed to:")
	pdf.Cell(60, 10, "From:")

	pdf.Ln(12)
	// Add invoice Data
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 7, row.Name)
	pdf.Cell(60, 7, ReceiverName)
	pdf.Ln(8)

	pdf.Cell(60, 7, row.Email)
	pdf.Cell(60, 7, ReceiverEmail)

	pdf.Ln(12)
	// Set font for the document
	pdf.SetFont("Arial", "B", 16)

	// Create the invoice title
	pdf.Cell(40, 20, "Payment Information")
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.SetTextColor(0, 0, 0) // Red text color for total

	// Add invoice details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(80, 10, "Description")
	pdf.Cell(30, 10, "Quantity")
	pdf.Cell(40, 10, "Unit Price")
	pdf.Cell(40, 10, "Amount")
	pdf.Ln(12)
	// Set font for the document
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(128, 128, 128) // text color
	// Draw another line after invoice details
	pdf.Line(10, 150, 200, 150) // Line from x=10, y=70 to x=200, y=70
	pdf.SetTextColor(0, 0, 0)   // Red text color for total
	pdf.Cell(80, 10, "Pay Request")
	pdf.Cell(30, 10, "1")
	pdf.Cell(40, 10, RequesteAmount)
	pdf.Cell(40, 10, RequesteAmount)
	pdf.Ln(10)
	pdf.SetTextColor(128, 128, 128) // text color
	// Draw another line after invoice details
	pdf.Line(110, 160, 200, 160) // Line from x=10, y=70 to x=200, y=70
	pdf.SetTextColor(0, 0, 0)    // Red text color for total
	pdf.Cell(80, 10, "")
	pdf.Cell(30, 10, "")
	pdf.Cell(40, 10, "Subtotal")
	pdf.Cell(40, 10, RequesteAmount)
	pdf.Ln(12)
	pdf.SetTextColor(128, 128, 128) // text color
	// Draw another line after invoice details
	pdf.Line(110, 170, 200, 170) // Line from x=10, y=70 to x=200, y=70
	pdf.SetTextColor(0, 0, 0)    // Red text color for total
	pdf.Cell(80, 10, "")
	pdf.Cell(30, 10, "")
	pdf.Cell(40, 10, "Total")
	pdf.Cell(40, 10, RequesteAmount)
	pdf.Ln(12)

	// Draw a horizontal line at the bottom of the page
	//pageHeight := 297.0  // A4 page height in mm
	//marginBottom := 20.0 // Margin from the bottom (you can adjust this)

	// Set the draw color (e.g., black for the bottom line)
	pdf.SetDrawColor(0, 0, 0)

	// Draw the line near the bottom of the page
	//pdf.Line(10, pageHeight-marginBottom, 200, pageHeight-marginBottom)
	// Set font for the document

	pdf.SetFont("Arial", "B", 10)
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.Ln(12)
	pdf.Ln(10)
	var SITENAME = os.Getenv("SITENAME")

	// Create the invoice title
	pdf.Cell(40, 10, "Payment was securely processed by "+SITENAME)

	// Create a buffer to hold the PDF data
	var buf bytes.Buffer

	// Write the PDF data to the buffer
	err = pdf.Output(&buf)
	if err != nil {
		return c.Status(500).SendString("Failed to generate PDF")
	}

	// Set the response header to serve a PDF inline (view in browser)
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "inline; filename=invoice.pdf")

	// Send the PDF file to the browser for viewing
	return c.SendStream(&buf)

}
