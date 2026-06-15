package models

import "time"

// Define struct for payment related

// Structs to represent the JSON data
type OnlineResponse struct {
	ContractMap map[string]bool  `json:"contractMap"`
	TokenInfo   OnlineTokenInfo  `json:"tokenInfo"`
	PageSize    int              `json:"page_size"`
	Code        int              `json:"code"`
	Data        []OnlineDataItem `json:"data"`
}

// Structs for get Online Token info
type OnlineTokenInfo struct {
	TokenId      string `json:"tokenId"`
	TokenAbbr    string `json:"tokenAbbr"`
	TokenName    string `json:"tokenName"`
	TokenDecimal int    `json:"tokenDecimal"`
	TokenCanShow int    `json:"tokenCanShow"`
	TokenType    string `json:"tokenType"`
	TokenLogo    string `json:"tokenLogo"`
	TokenLevel   string `json:"tokenLevel"`
	IssuerAddr   string `json:"issuerAddr"`
	VIP          bool   `json:"vip"`
}

// Structs for get Online Data
type OnlineDataItem struct {
	Amount          string `json:"amount"`
	Status          int    `json:"status"`
	ApprovalAmount  string `json:"approval_amount"`
	BlockTimestamp  int64  `json:"block_timestamp"`
	Block           int    `json:"block"`
	From            string `json:"from"`
	To              string `json:"to"`
	Hash            string `json:"hash"`
	ContractAddress string `json:"contract_address"`
	Confirmed       int    `json:"confirmed"`
	ContractType    string `json:"contract_type"`
	ContractTypeInt int    `json:"contractType"`
	Revert          int    `json:"revert"`
	ContractRet     string `json:"contract_ret"`
	FinalResult     string `json:"final_result"`
	EventType       string `json:"event_type"`
	IssueAddress    string `json:"issue_address"`
	Decimals        int    `json:"decimals"`
	TokenName       string `json:"token_name"`
	ID              string `json:"id"`
	Direction       int    `json:"direction"`
}

// Define a struct that matches the JSON structure
type OnlineTokenData struct {
	TokenId      string `json:"tokenId"`
	TokenName    string `json:"tokenName"`
	TokenAbbr    string `json:"tokenAbbr"`
	TokenCanShow int    `json:"tokenCanShow"`
	TokenType    string `json:"tokenType"`
}

// Structs for get ETH Response
type ETHResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		BlockNumber       string `json:"blockNumber"`
		BlockHash         string `json:"blockHash"`
		BlockTimestamp    string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		TransactionIndex  string `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		Value             string `json:"value"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		Input             string `json:"input"`
		MethodID          string `json:"methodId"`
		FunctionName      string `json:"functionName"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		TxReceiptStatus   string `json:"txreceipt_status"`
		GasUsed           string `json:"gasUsed"`
		Confirmations     string `json:"confirmations"`
		IsError           string `json:"isError"`
	} `json:"result"`
}

// Struct for update Transaction Online
type TransactionUpdateOnline struct {
	//gorm.Model
	Id                 uint      `gorm:"primaryKey"`
	Receivedamount     float64   `json:"receivedamount,omitempty"`
	Receivedcurrency   string    `json:"receivedcurrency,omitempty"`
	Status             string    `json:"status,omitempty"`
	Substatus          int       `json:"substatus,omitempty"`
	Response_hash      string    `json:"response_hash,omitempty"`
	Response_from      string    `json:"response_from,omitempty"`
	Response_to        string    `json:"response_to,omitempty"`
	Response_timestamp time.Time `json:"response_timestamp,omitempty"`
	Response_json      string    `json:"response_json,omitempty"`
}

// Struct for generate pay request
type PayRequest struct {
	Cid            string `json:"cid" form:"cid"`
	Price_currency string `json:"price_currency" form:"price_currency"`
	Price_amount   string `json:"price_amount" form:"price_amount"`
	Sender_name    string `json:"sender_name" form:"sender_name"`
	Sender_email   string `json:"sender_email" form:"sender_email"`
	Client_id      uint   `json:"client_id" form:"client_id"`
	Pay_type       int    `json:"pay_type" form:"pay_type"`
	Crypto_id      int    `json:"crypto_id" form:"crypto_id"`
	Customerrefid  string `json:"customerrefid" form:"customerrefid"`
}

// Struct for gate pay response
type PayResponse struct {
	Qr_code      string  `json:"qr_code"`
	Address      string  `json:"address"`
	Amount       float64 `json:"amount"`
	Transid      string  `json:"transid"`
	Coinicon     string  `json:"coinicon"`
	Coinnetwork  string  `json:"coinnetwork"`
	Coin_id      int     `json:"coin_id"`
	Coin_pay_url string  `json:"coin_pay_url"`
}

// CardanoTransaction represents the structure of the response from CardanoScan API.

// Root structure
type CardanoResponse struct {
	Data []CardanoTransaction `json:"transactions"`
}

// Transaction structure
type CardanoTransaction struct {
	Hash      string          `json:"hash"`
	Fees      string          `json:"fees"`
	Timestamp time.Time       `json:"timestamp"`
	Outputs   []CardanoOutput `json:"outputs"`
	Status    bool            `json:"status"`
}

// Input structure
type CardanoInput struct {
	Address string `json:"address"`
	Index   int    `json:"index"`
	TxID    string `json:"txId"`
	Value   string `json:"value"`
}

// Output structure
type CardanoOutput struct {
	Address string `json:"address"`
	Value   string `json:"value"`
}

// TTL structure
type CardanoTTL struct {
	Slot      int       `json:"slot"`
	Timestamp time.Time `json:"timestamp"`
}

//Struct fo BTC Response
type BTCAddressInfo struct {
	Txs []BTCTransaction `json:"txs"`
}
type BTCTransaction struct {
	Hash        string      `json:"hash"`
	Time        int64       `json:"time"`
	BlockHeight int         `json:"block_height"`
	DoubleSpend bool        `json:"double_spend"`
	Inputs      []BTCInput  `json:"inputs"`
	Out         []BTCOutput `json:"out"`
}

type BTCInput struct {
	Addr    string  `json:"addr"`
	PrevOut PrevOut `json:"prev_out"`
}
type PrevOut struct {
	Addr string `json:"addr"`
}

type BTCOutput struct {
	Value int    `json:"value"`
	Addr  string `json:"addr"`
}

// Main struct for the doge address data
type DogeAddressData struct {
	Address string      `json:"address"`
	TxRefs  []DogeTxRef `json:"txrefs"`
}

// Struct for each transaction reference (txrefs)
type DogeTxRef struct {
	TxHash      string    `json:"tx_hash"`
	Value       int64     `json:"value"`
	Spent       bool      `json:"spent,omitempty"` // Omit if not present
	Spent_by    string    `json:"spent_by"`
	Confirmed   time.Time `json:"confirmed"` // Parsed as time
	DoubleSpend bool      `json:"double_spend"`
}

// Main struct for the Lite address data
type LiteAddressData struct {
	Address string      `json:"address"`
	TxRefs  []LiteTxRef `json:"txrefs"`
}

// Struct for each transaction reference (txrefs)
type LiteTxRef struct {
	TxHash      string    `json:"tx_hash"`
	Value       int64     `json:"value"`
	Spent       bool      `json:"spent,omitempty"` // Omit if not present
	Spent_by    string    `json:"spent_by"`
	Confirmed   time.Time `json:"confirmed"` // Parsed as time
	DoubleSpend bool      `json:"double_spend"`
}

// Main struct for the DASH address data
type DashAddressData struct {
	Address string      `json:"address"`
	TxRefs  []DashTxRef `json:"txrefs"`
}

// Struct for each transaction reference (txrefs)
type DashTxRef struct {
	TxHash      string    `json:"tx_hash"`
	Value       int64     `json:"value"`
	Spent       bool      `json:"spent,omitempty"` // Omit if not present
	Spent_by    string    `json:"spent_by"`
	Confirmed   time.Time `json:"confirmed"` // Parsed as time
	DoubleSpend bool      `json:"double_spend"`
}

// Response is the MATIC - Polygon
type CoinXAddressData struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []PolygonResult `json:"result"`
}

// Result represents each transaction detail
type PolygonResult struct {
	TimeStamp string `json:"timeStamp"`
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
}
