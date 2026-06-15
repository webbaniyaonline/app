package handlers

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

// Define a struct to hold the JSON-RPC response
// Struct to hold the result of the signature fetch
type SignatureResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  []struct {
		Signature string `json:"signature"`
	} `json:"result"`
	Id int `json:"id"`
}

// Struct to hold the detailed transaction info
type TransactionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Slot int `json:"slot"`
		Meta struct {
			Err          interface{} `json:"err"`
			Fee          int         `json:"fee"`
			PreBalances  []int64     `json:"preBalances"`
			PostBalances []int64     `json:"postBalances"`
		} `json:"meta"`
		Transaction struct {
			Message struct {
				AccountKeys  []string `json:"accountKeys"`
				Instructions []struct {
					ProgramId string `json:"programId"`
					Data      string `json:"data"`
				} `json:"instructions"`
			} `json:"message"`
		} `json:"transaction"`
	} `json:"result"`
	Id int `json:"id"`
}

// Solana RPC endpoint
const solanaRpcUrl = "https://api.mainnet-beta.solana.com"

// Set up an HTTP client using Resty (or you can use net/http)

// Handler to get balance by address
func TestCode(c *fiber.Ctx) error {
	client := resty.New()
	address := c.Params("address")

	// Step 1: Get transaction signatures
	signaturePayload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getConfirmedSignaturesForAddress2",
		"params":  []interface{}{address, map[string]int{"limit": 5}}, // Fetch 5 transactions
	}

	signatureResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(signaturePayload).
		SetResult(&SignatureResponse{}).
		Post(solanaRpcUrl)

	if err != nil {
		log.Printf("Error fetching transaction signatures: %v", err)
		return c.Status(500).SendString("Failed to fetch transaction signatures")
	}

	signatures := signatureResp.Result().(*SignatureResponse).Result

	if len(signatures) == 0 {
		return c.Status(404).SendString("No transactions found for the given address")
	}

	// Step 2: Fetch transaction details for each signature
	var transactionDetails []TransactionResponse
	for _, signature := range signatures {
		txPayload := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "getTransaction",
			"params":  []interface{}{signature.Signature, map[string]string{"encoding": "json"}},
		}

		txResp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(txPayload).
			SetResult(&TransactionResponse{}).
			Post(solanaRpcUrl)

		if err != nil {
			log.Printf("Error fetching transaction details: %v", err)
			continue
		}

		transactionDetails = append(transactionDetails, *txResp.Result().(*TransactionResponse))
	}

	// Return the transaction details as JSON
	return c.JSON(transactionDetails)

}
