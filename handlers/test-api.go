package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler to get balance by address
func TestAPI(c *fiber.Ctx) error {
	receivedFrom := ""
	receivedTo := ""
	receivedHash := ""
	receivedFinalResult := ""
	fetchTimestamp := ""
	responseTimestamp := ""
	receivedAmountNew := 0.00
	var body []byte
	status_address := c.Params("address")

	// URL of the external site to fetch JSON from
	apiKey := os.Getenv("POLYGON_API_KEY")
	url := "https://api.polygonscan.com/api?module=account&action=txlist&address=" + status_address + "&startblock=0&endblock=99999999&sort=desc&page=1&offset=1&apikey=" + apiKey
	//fmt.Println("url => ", url)
	//////////////////////////////////////
	resp, err := http.Get(url)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to fetch data")
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to read response body")
	}
	//fmt.Println("body => ", string(body))
	// Initialize the Response struct
	var responseD models.CoinXAddressData

	// Unmarshal the byte data into the struct
	err = json.Unmarshal(body, &responseD)
	if err != nil {
		//return c.Status(http.StatusInternalServerError).SendString("Failed to parse JSON")
		fmt.Println("Failed to parse JSON")
	}

	//fmt.Println("responseD => ", responseD)

	// Check if data not found
	if responseD.Status == "0" {
		response := StatusResponse{
			Hashcode:       "",
			Payment_status: "",
			Payment_id:     "",
		}
		return c.JSON(response)
	}

	receivedAmount := responseD.Data[0].Value

	// convert string to float value
	receivedAmt, err := strconv.ParseFloat(receivedAmount, 64)
	if err != nil {
		fmt.Println(" Error convert string to float value :")
	}
	// Convert the integer to a float64 with 18 decimal places
	AmountInFloat := float64(receivedAmt) / 1000000000000000000

	// Format the float to 6 decimal places as a string
	formattedResult := strconv.FormatFloat(AmountInFloat, 'f', 6, 64)

	// convert string to float value
	receivedAmountNew, err = strconv.ParseFloat(formattedResult, 64)
	if err != nil {
		fmt.Println(" Error convert string to float value :")
	}

	//fmt.Println("receivedAmountNew", receivedAmountNew)
	receivedFrom = responseD.Data[0].From
	receivedTo = responseD.Data[0].To
	receivedHash = responseD.Data[0].Hash

	fmt.Println("receivedHash => ", receivedHash)
	fmt.Println("receivedAmountNew => ", receivedAmountNew)
	fmt.Println("receivedFrom => ", receivedFrom)
	fmt.Println("receivedTo => ", receivedTo)
	fmt.Println("TimeStamp => ", responseD.Data[0].TimeStamp)

	// Convert the string to an int64
	unixTime, err := strconv.ParseInt(responseD.Data[0].TimeStamp, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int64:", err)
	}
	// Convert Unix timestamp (seconds) to time.Time
	timestamp := time.Unix(unixTime, 0)

	receivedFinalResult = responseD.Status
	if receivedFinalResult == "1" {
		receivedFinalResult = "Success"
	} else {
		receivedFinalResult = "Declined"
	}

	fetchTimestamp = "2024-10-10 07:00:09"
	//fetchTimestamp = coinAddress.Lastupdate.Format("2006-01-02 15:04:05")
	// Format the time to "2006-01-02 15:04:05"
	responseTimestamp = timestamp.Format("2006-01-02 15:04:05")

	fmt.Println("fetchTimestamp => ", fetchTimestamp)
	fmt.Println("responseTimestamp => ", responseTimestamp)

	transactionDetails := "{}"
	return c.JSON(transactionDetails)

}
