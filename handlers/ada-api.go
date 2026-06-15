package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"template/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// This page for testing not in use

// Struct to hold the API response
type BalanceResponse struct {
	Balance string `json:"balance"`
}

// Function to fetch balance from CardanoScan API
func GetCardanoBalance(address string) (string, error) {
	// Construct the API URL
	url := "https://api.cardanoscan.io/api/v1/transaction/list?address=" + address + "&pageNo=1&limit=1&order=desc"
	apiKey := os.Getenv("CARDANO_SCAN_API_KEY")

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating HTTP request: %w", err)
	}
	// Send Api in header
	req.Header.Set("apiKey", apiKey)
	// Make request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body: %w", err)
	}

	// Initialize the Response struct
	var responseD models.CardanoResponse

	// Unmarshal the byte data into the struct
	err = json.Unmarshal(body, &responseD)
	if err != nil {
		fmt.Println("Failed to parse JSON")
	}
	//fmt.Println(" Data :", responseD)
	receivedHash := responseD.Data[0].Hash
	receivedTimestamp := responseD.Data[0].Timestamp
	receivedFees := responseD.Data[0].Fees
	receivedFinalResult := responseD.Data[0].Status
	receivedFrom := responseD.Data[0].Outputs[0].Address
	receivedTo := responseD.Data[0].Outputs[1].Address
	receivedAmount := responseD.Data[0].Outputs[1].Value
	receivedAmountNew := cryptoAmountFormat(receivedAmount)
	status := ""
	if receivedFinalResult == true {
		status = "Success"
	} else {
		status = "Declined"
	}

	fmt.Println("receivedHash =>>> ", receivedHash)
	fmt.Println("receivedTimestamp =>>> ", receivedTimestamp)
	fmt.Println("receivedFinalResult =>>> ", status)
	fmt.Println("receivedFees =>>> ", receivedFees)
	fmt.Println("receivedFrom =>>> ", receivedFrom)
	fmt.Println("receivedTo =>>> ", receivedTo)
	fmt.Println("receivedAmount =>>> ", receivedAmountNew)

	return "respBodyv", nil
}

// Handler to get balance by address
func GetBalanceHandler(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Address is required")
	}

	//fmt.Println("address =>", address)

	// Fetch balance from API
	balance, err := GetCardanoBalance(address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to fetch balance: %v", err))
	}

	// Return balance as JSON
	return c.JSON(fiber.Map{
		"address": address,
		"balance": balance,
	})
}
