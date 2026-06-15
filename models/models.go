package models

import "time"

type Datas struct {
	MID   uint
	MName string
}

// Here create struct for all functions and methods calls in projects

// Struct For display Merchant login history listing
type LoginHistory struct {
	//gorm.Model
	Token_id    uint   `gorm:"primaryKey"`
	Client_id   uint   `json:"client_id,omitempty"`
	Login_time  string `json:"login_time,omitempty"`
	Logout_time string `json:"logout_time,omitempty"`
	Login_ip    string `json:"login_ip,omitempty"`
	Login_type  int    `json:"login_type,omitempty"`
	Resource    string `json:"resource,omitempty"`
	//Count       int64  `json:"count"`
}

// Struct For Update login history
type LoginHistoryUpdate struct {
	//gorm.Model
	Token_id    uint   `gorm:"primaryKey"`
	Logout_time string `json:"logout_time,omitempty"`
}

// Struct For Change Password request
type PasswordChangeRequest struct {
	//gorm.Model
	Id             uint   `gorm:"primaryKey"`
	Mid            uint   `json:"mid,omitempty"`
	Password_hash  string `json:"password_hash,omitempty"`
	Password_email string `json:"password_email,omitempty"`
}

// Struct For Email Template Listing
type EmailTemplate struct {
	//gorm.Model
	Template_id      uint   `gorm:"primaryKey"`
	Template_code    string `json:"template_code,omitempty"`
	Template_subject string `json:"template_subject,omitempty"`
	Template_desc    string `json:"template_desc,omitempty"`
	Status           int    `json:"status,omitempty"`
}

// Struct For Update login history
type UpdateHistory struct {
	//gorm.Model
	Update_id   uint   `gorm:"primaryKey"`
	Client_id   uint   `json:"client_id,omitempty"`
	Update_type string `json:"update_type,omitempty"`
	Update_desc string `json:"update_desc,omitempty"`
	Update_ip   string `json:"c_ip,omitempty"`
	Updated_on  string `json:"updated_on,omitempty"`
}

// Struct For Display Transaction List
type TransactionList struct {
	//gorm.Model
	Transactionid    uint   `gorm:"primaryKey"`
	Client_id        uint   `json:"client_id,omitempty"`
	Walletid         string `json:"walletid,omitempty"`
	Transactiontype  string `json:"transactiontype,omitempty"`
	Amount           string `json:"amount,omitempty"`
	Transaction_hash string `json:"transaction_hash,omitempty"`
	Status           string `json:"status,omitempty"`
	Timestamp        string `json:"timestamp,omitempty"`
}

// Struct For Display Wallet  List
type WalletList struct {
	//gorm.Model
	Wallet_id     uint   `gorm:"primaryKey"`
	Volt_id       string `json:"volt_id,omitempty"`
	Coin          string `json:"coin,omitempty"`
	Address       string `json:"address,omitempty"`
	Legacyaddress string `json:"legacyaddress,omitempty"`
	Tag           string `json:"tag,omitempty"`
	Total         string `json:"total,omitempty"`
	Status        string `json:"status,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
	Available     string `json:"available,omitempty"`
	Pending       string `json:"pending,omitempty"`
	Frozen        string `json:"frozen,omitempty"`
	Lockedamount  string `json:"lockedamount,omitempty"`
}

// Struct For Fetch Wallet Balance
type WalletListBalance struct {
	//gorm.Model
	Wallet_id    uint   `gorm:"primaryKey"`
	Total        string `json:"total,omitempty"`
	Pending      string `json:"pending,omitempty"`
	Frozen       string `json:"frozen,omitempty"`
	Lockedamount string `json:"lockedamount,omitempty"`
	Available    string `json:"available,omitempty"`
}

// Struct For Login Details
type LoginList struct {
	//gorm.Model
	Client_id           uint `gorm:"primaryKey"`
	Full_name           string
	Password            string
	Status              int
	Totp_secret         string
	Totp_status         int
	User_agent          string
	Password_created_at string
	FailedAttempts      int
	AccountLocked       bool
	LockTime            *time.Time
}

// Struct For Merchant Details
type MemberDetails struct {
	//gorm.Model
	Client_id     uint   `gorm:"primaryKey"`
	Username      string `json:"username,omitempty"`
	Full_name     string `json:"full_name,omitempty"`
	Status        int    `json:"status,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
	Title         string `json:"title,omitempty"`
	Gender        string `json:"gender,omitempty"`
	Birth_date    string `json:"birth_date,omitempty"`
	Country_code  string `json:"country_code,omitempty"`
	Mobile        string `json:"mobile,omitempty"`
	Address_line1 string `json:"address_line1,omitempty"`
	Address_line2 string `json:"address_line2,omitempty"`
	//City          string `json:"city,omitempty"`
	//State         string `json:"state,omitempty"`
	//Country       string `json:"country,omitempty"`
	//Pincode       string `json:"pincode,omitempty"`
}

// Struct For Update Member Status
type MemberStatusUpdate struct {
	//gorm.Model
	Client_id uint `gorm:"primaryKey"`
	Status    int  `json:"status,omitempty"`
}

// Struct For Update Fireblock Volt
type UpdateVolt struct {
	//gorm.Model
	Client_id uint `gorm:"primaryKey"`
	Volt_id   string
}

// Struct For Update Google 2FA
type Update2FA struct {
	//gorm.Model
	Client_id   uint `gorm:"primaryKey"`
	Totp_secret string
	Totp_status int
}

// Struct For Manage Merchant Details
type Client_Master struct {
	//gorm.Model
	Client_id uint `gorm:"primaryKey"`
	Username  string
	Full_name string
	Password  string
	Status    int
	Resource  string
}

// for change merchant password
type ClientPassword struct {
	Client_id           uint   `gorm:"primaryKey"`
	Password            string `json:"password,omitempty"`
	Password_created_at string `json:"password_created_at,omitempty"`
}

// for change merchant Account active - deactive after wrong password
type ClientWrongPasswordActivation struct {
	Client_id      uint       `gorm:"primaryKey"`
	FailedAttempts int        `json:"failed_attempts,omitempty"`
	AccountLocked  bool       `json:"account_locked,omitempty"`
	LockTime       *time.Time `json:"lock_time,omitempty"`
}

// for change merchant Account active - deactive after wrong password
type ClientWrongPasswordLock struct {
	Client_id      uint      `gorm:"primaryKey"`
	FailedAttempts int       `json:"failed_attempts,omitempty"`
	AccountLocked  bool      `json:"account_locked,omitempty"`
	LockTime       time.Time `json:"lock_time,omitempty"`
}

type ApiBody struct {
	BodyData string
}

// for Display Vault List Fireblock
type VaultList struct {
	Id         int           `json:"id,omitempty"`
	Name       string        `json:"name,omitempty"`
	HiddenOnUI string        `json:"hiddenOnUI,omitempty"`
	AutoFuel   string        `json:"autoFuel,omitempty"`
	Assets     []AddressList `json:"assets"`
}

// For Manage Address List Fireblock
type AddressList struct {
	Id           string `json:"id,omitempty"`
	Total        string `json:"total,omitempty"`
	Balance      string `json:"balance,omitempty"`
	LockedAmount string `json:"lockedAmount,omitempty"`
	Available    string `json:"available,omitempty"`
	Pending      string `json:"pending,omitempty"`
	Frozen       string `json:"frozen,omitempty"`
	Staked       string `json:"staked,omitempty"`
	BlockHeight  string `json:"blockHeight,omitempty"`
}

// For Manage User API Fireblock
type UserApi struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Role      string `json:"role,omitempty"`
	Email     string `json:"email,omitempty"`
	Enabled   string `json:"enabled,omitempty"`
}

// Define the struct for the nested reward info
type RewardInfo struct {
	PendingRewards string `json:"pendingRewards"`
}

// Define the struct for each asset
type Asset struct {
	ID           string      `json:"id"`
	Total        string      `json:"total"`
	Balance      string      `json:"balance"`
	LockedAmount string      `json:"lockedAmount"`
	Available    string      `json:"available"`
	Pending      string      `json:"pending"`
	Frozen       string      `json:"frozen"`
	Staked       string      `json:"staked"`
	BlockHeight  string      `json:"blockHeight"`
	RewardInfo   *RewardInfo `json:"rewardInfo,omitempty"`
}

// Define the main struct for the JSON response
type FireblocksResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	HiddenOnUI bool    `json:"hiddenOnUI"`
	AutoFuel   bool    `json:"autoFuel"`
	Assets     []Asset `json:"assets"`
}

// For Get Fireblock Response
type FireblocksVoltResponse struct {
	Accounts []FireblocksResponse `json:"accounts"`
	Paging   Paging               `json:"paging"`
}
type Paging struct {
	// Add fields if there are any in the actual JSON
}

// Define the main struct for the JSON response
type FireblocksUsers struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Email     string `json:"email,omitempty"`
	Enabled   bool   `json:"enabled"`
}

// Define the struct for each asset
type Address struct {
	AssetId           string `json:"assetId"`
	Address           string `json:"address"`
	Description       string `json:"description,omitempty"`
	Tag               string `json:"tag,omitempty"`
	Type              string `json:"type"`
	AddressFormat     string `json:"addressFormat"`
	LegacyAddress     string `json:"legacyAddress,omitempty"`
	EnterpriseAddress string `json:"enterpriseAddress,omitempty"`
	Bip44AddressIndex int    `json:"bip44AddressIndex"`
	UserDefined       bool   `json:"userDefined"`
}

// Define the main struct for the JSON response
type FireblocksAddress struct {
	Addresses []Address `json:"addresses"`
}

// Define the main struct for the JSON response
type FireblocksWallet struct {
	Id            string `json:"id,omitempty"`
	Address       string `json:"address,omitempty"`
	LegacyAddress string `json:"legacyAddress,omitempty"`
	Tag           string `json:"tag,omitempty"`
	Message       string `json:"message,omitempty"`
	Code          int    `json:"code,omitempty"`
}
type CreateVaultAccountResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// for manage merchant Profile
type ProfileData struct {
	Client_id     uint   `gorm:"primaryKey"`
	Gender        string `json:"gender"`
	BirthDate     string `json:"birth_date"`
	CountryCode   string `json:"country_code"`
	Mobile        string `json:"mobile"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	Pincode       string `json:"pincode"`
	Profile_image string `json:"profile_image"`
}

// Struct for alerts
type Alert struct {
	Alert string
}

// Struct For Add Wallet
type AddWallet struct {
	//Wallet_id     uint   `json:"wallet_id"`
	Volt_id       string `json:"volt_id"`
	Coin          string `json:"coin"`
	Address       string `json:"address,omitempty"`
	Legacyaddress string `json:"legacyAddress,omitempty"`
	Tag           string `json:"tag,omitempty"`
}

// Struct for display wallet List
type GetWallet struct {
	Wallet_id     uint   `json:"wallet_id"`
	Volt_id       string `json:"volt_id"`
	Coin          string `json:"coin"`
	Address       string `json:"address,omitempty"`
	Legacyaddress string `json:"legacyAddress,omitempty"`
	Tag           string `json:"tag,omitempty"`
}

// Struct For send Envoice Data
type EmailData struct {
	Email       string
	FullName    string
	UserName    string
	Password    string
	Invoice_id  string
	Invoice_url string
	Amount      string
	Status      string
	HashCode    string
	TransID     string
	Crypto      string
	Title       string
	Details     string
}

/// Struct for Now Payment //////////

// For Send Api Request
type TransferRequestNP struct {
	Price_amount      float64 `json:"price_amount"`
	Price_currency    string  `json:"price_currency"`
	Pay_currency      string  `json:"pay_currency"`
	Ipn_callback_url  string  `json:"ipn_callback_url"`
	Order_id          string  `json:"order_id"`
	Order_description string  `json:"order_description"`
	//Success_url       string  `json:"success_url"`
	//Cancel_url        string  `json:"cancel_url"`
}

// For Send Api Request
type TransferRequestCrypto struct {
	Price_amount      string `json:"price_amount"`
	Price_currency    string `json:"price_currency"`
	Order_id          string `json:"order_id"`
	Order_description string `json:"order_description"`
	Ipn_callback_url  string `json:"ipn_callback_url"`
	Success_url       string `json:"success_url"`
	Cancel_url        string `json:"cancel_url"`
}

// For Get Api Response
type TransResponceNP struct {
	Payment_id               string  `json:"payment_id"`
	Payment_status           string  `json:"payment_status"`
	Pay_address              string  `json:"pay_address"`
	Price_amount             float64 `json:"price_amount"`
	Price_currency           string  `json:"price_currency"`
	Pay_amount               float64 `json:"pay_amount"`
	Amount_received          float64 `json:"amount_received"`
	Pay_currency             string  `json:"pay_currency"`
	Order_id                 string  `json:"order_id"`
	Order_description        string  `json:"order_description"`
	Payin_extra_id           string  `json:"payin_extra_id"`
	Ipn_callback_url         string  `json:"ipn_callback_url"`
	Created_at               string  `json:"created_at"`
	Updated_at               string  `json:"updated_at"`
	Purchase_id              string  `json:"purchase_id"`
	Smart_contract           string  `json:"smart_contract"`
	Network                  string  `json:"network"`
	Network_precision        string  `json:"network_precision"`
	Time_limit               string  `json:"time_limit"`
	Burning_percent          string  `json:"burning_percent"`
	Expiration_estimate_date string  `json:"expiration_estimate_date"`
	Is_fixed_rate            bool    `json:"is_fixed_rate"`
	Is_fee_paid_by_user      bool    `json:"is_fee_paid_by_user"`
	Valid_until              string  `json:"valid_until"`
	Product                  string  `json:"product"`
	Success                  string  `json:"success"`
	Code                     string  `json:"code,omitempty"`
	Message                  string  `json:"message,omitempty"`
}

type TransResponcecryptoNP struct {
	Payment_id               string  `json:"payment_id"`
	Payment_status           string  `json:"payment_status"`
	Pay_address              string  `json:"pay_address"`
	Price_amount             string  `json:"price_amount"`
	Price_currency           string  `json:"price_currency"`
	Pay_amount               float64 `json:"pay_amount"`
	Amount_received          float64 `json:"amount_received"`
	Pay_currency             string  `json:"pay_currency"`
	Order_id                 string  `json:"order_id"`
	Order_description        string  `json:"order_description"`
	Payin_extra_id           string  `json:"payin_extra_id"`
	Ipn_callback_url         string  `json:"ipn_callback_url"`
	Created_at               string  `json:"created_at"`
	Updated_at               string  `json:"updated_at"`
	Purchase_id              string  `json:"purchase_id"`
	Smart_contract           string  `json:"smart_contract"`
	Network                  string  `json:"network"`
	Network_precision        string  `json:"network_precision"`
	Time_limit               string  `json:"time_limit"`
	Burning_percent          string  `json:"burning_percent"`
	Expiration_estimate_date string  `json:"expiration_estimate_date"`
	Is_fixed_rate            bool    `json:"is_fixed_rate"`
	Is_fee_paid_by_user      bool    `json:"is_fee_paid_by_user"`
	Valid_until              string  `json:"valid_until"`
	Product                  string  `json:"product"`
	Success                  string  `json:"success"`
	Code                     string  `json:"code,omitempty"`
	Message                  string  `json:"message,omitempty"`
	Invoice_id               string  `json:"invoice_id,omitempty"`
	Token_id                 string  `json:"token_id,omitempty"`
	Invoice_url              string  `json:"invoice_url,omitempty"`
	Request_json             string  `json:"request_json,omitempty"`
	Ip                       string  `json:"ip,omitempty"`
	Id                       string  `json:"id,omitempty"`
}

// For Fetch Table
type Transaction_MasterNP struct {
	//gorm.Model
	Tid               uint    `gorm:"primaryKey"`
	Client_id         uint    `json:"client_id,omitempty"`
	Payment_id        string  `json:"payment_id,omitempty"`
	Payment_status    string  `json:"payment_status,omitempty"`
	Pay_address       string  `json:"pay_address,omitempty"`
	Price_amount      float64 `json:"price_amount,omitempty"`
	Pay_amount        float64 `json:"pay_amount,omitempty"`
	Amount_received   float64 `json:"amount_received,omitempty"`
	Price_currency    string  `json:"price_currency,omitempty"`
	Pay_currency      string  `json:"pay_currency,omitempty"`
	Order_id          string  `json:"order_id,omitempty"`
	Order_description string  `json:"order_description,omitempty"`
	Created_at        string  `json:"created_at,omitempty"`
	Updated_at        string  `json:"updated_at,omitempty"`
	Invoice_id        string  `json:"invoice_id,omitempty"`
	Token_id          string  `json:"token_id,omitempty"`
	Invoice_url       string  `json:"invoice_url,omitempty"`
	Request_json      string  `json:"request_json,omitempty"`
	Ip                string  `json:"ip,omitempty"`
}

/// End Struct for Now Payment //////////

/// Start Struct for Fireblocks //////////

type Transaction_Master struct {
	//gorm.Model
	Id                 uint    `gorm:"primaryKey"`
	Transaction_id     string  `json:"transaction_id,omitempty"`
	Client_id          uint    `json:"client_id,omitempty"`
	Volt_id            string  `json:"volt_id,omitempty"`
	Assetid            string  `json:"assetid,omitempty"`
	Amount             float64 `json:"amount,omitempty"`
	Status             string  `json:"status,omitempty"`
	Operation          string  `json:"operation,omitempty"`
	Customerrefid      string  `json:"customerrefid,omitempty"`
	Createdate         string  `json:"createdate,omitempty"`
	Transaction_type   string  `json:"transaction_type,omitempty"`
	Ip                 string  `json:"ip,omitempty"`
	Note               string  `json:"note,omitempty"`
	Source             string  `json:"source,omitempty"`
	Requestedamount    float64 `json:"requestedamount,omitempty"`
	Netamount          float64 `json:"netamount,omitempty"`
	Amountusd          float64 `json:"amountusd,omitempty"`
	Fee                float64 `json:"fee,omitempty"`
	Networkfee         float64 `json:"networkfee,omitempty"`
	Substatus          string  `json:"substatus,omitempty"`
	Txhash             string  `json:"txhash,omitempty"`
	Lastupdated        string  `json:"lastupdated,omitempty"`
	Destinationaddress string  `json:"destinationaddress,omitempty"`
	Createdby          string  `json:"createdby,omitempty"`
}

// Structs to match the JSON schema
type OneTimeAddress struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
}

type Destination struct {
	OneTimeAddress OneTimeAddress `json:"oneTimeAddress"`
	Type           string         `json:"type"`
}

type Source struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Struct for transfer Request
type TransferRequest struct {
	AssetId       string      `json:"assetId"`
	Source        Source      `json:"source"`
	Destination   Destination `json:"destination"`
	Amount        float64     `json:"amount"`
	FeeLevel      string      `json:"feeLevel"`
	Note          string      `json:"note"`
	Operation     string      `json:"operation"`
	CustomerRefId string      `json:"customerRefId"`
}

// Struct for get transactions response
type TransResponce struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Struct for get transactions response
type TransGetResponce struct {
	Id              uint    `gorm:"primaryKey"`
	Requestedamount float64 `json:"requestedAmount"`
	Netamount       float64 `json:"netAmount"`
	Amountusd       float64 `json:"amountUSD"`
	Fee             float64 `json:"fee"`
	Networkfee      float64 `json:"networkFee"`
	//Lastupdated     int     `json:"lastUpdated"`
	Txhash             string `json:"txHash"`
	Status             string `json:"status"`
	Substatus          string `json:"subStatus"`
	Destinationaddress string `json:"destinationAddress"`
	Createdby          string `json:"createdBy"`
}

// Struct for transactions
type Transaction_Pay struct {
	//gorm.Model
	Id                  uint      `gorm:"primaryKey"`
	Transaction_id      string    `json:"transaction_id,omitempty"`
	Client_id           uint      `json:"client_id,omitempty"`
	Assetid             int       `json:"assetid,omitempty"`
	Transaction_type    string    `json:"transaction_type,omitempty"`
	Requestedamount     float64   `json:"requestedamount,omitempty"`
	Requestedcurrency   string    `json:"requestedcurrency,omitempty"`
	Convertedamount     float64   `json:"convertedamount,omitempty"`
	Convertedcurrency   string    `json:"convertedcurrency,omitempty"`
	Receivedamount      float64   `json:"receivedamount,omitempty"`
	Receivedcurrency    string    `json:"receivedcurrency,omitempty"`
	Status              string    `json:"status,omitempty"`
	Substatus           int       `json:"substatus,omitempty"`
	Customerrefid       string    `json:"customerrefid,omitempty"`
	Note                string    `json:"note,omitempty"`
	Createdate          string    `json:"createdate,omitempty"`
	Destinationaddress  string    `json:"destinationaddress,omitempty"`
	Ip                  string    `json:"ip,omitempty"`
	Approved_by         string    `json:"approved_by,omitempty"`
	Approver_id         uint      `json:"approver_id,omitempty"`
	Approver_comment    string    `json:"approver_comment,omitempty"`
	Approved_date       string    `json:"approved_date,omitempty"`
	Response_hash       string    `json:"response_hash,omitempty"`
	Response_from       string    `json:"response_from,omitempty"`
	Response_to         string    `json:"response_to,omitempty"`
	Response_timestamp  time.Time `json:"response_timestamp,omitempty"`
	Response_json       string    `json:"response_json,omitempty"`
	Order_id            string    `json:"order_id,omitempty"`
	Is_fee_paid_by_user bool      `json:"is_fee_paid_by_user,omitempty"`
}

// Struct for transactions
type Transaction_Update struct {
	//gorm.Model
	Id               uint    `gorm:"primaryKey"`
	Receivedamount   float64 `json:"receivedamount,omitempty"`
	Status           string  `json:"status,omitempty"`
	Substatus        int     `json:"substatus,omitempty"`
	Response_hash    string  `json:"response_hash,omitempty"`
	Approved_by      string  `json:"approved_by,omitempty"`
	Approver_id      uint    `json:"approver_id,omitempty"`
	Approver_comment string  `json:"approver_comment,omitempty"`
	Approved_date    string  `json:"approved_date,omitempty"`
}

// Struct for update transactions withdraw
type Transaction_Withdraw_Update struct {
	//gorm.Model
	Id               uint   `gorm:"primaryKey"`
	Status           string `json:"status,omitempty"`
	Substatus        int    `json:"substatus,omitempty"`
	Response_hash    string `json:"response_hash,omitempty"`
	Approved_by      string `json:"approved_by,omitempty"`
	Approver_id      uint   `json:"approver_id,omitempty"`
	Approver_comment string `json:"approver_comment,omitempty"`
	Approved_date    string `json:"approved_date,omitempty"`
}

// Update Failed Status
type FailedStatus struct {
	Status             string    `json:"status,omitempty"`
	Substatus          int       `json:"substatus,omitempty"`
	Response_timestamp time.Time `json:"response_timestamp,omitempty"`
}

// Get Coin Balance
type CoinWithBalance struct {
	Assetid int    `json:"assetid,omitempty"`
	Balance string `json:"balance,omitempty"`
}

// for manage Pay Reques
type Invoice_Master struct {
	//gorm.Model
	Invoice_id          uint    `gorm:"primaryKey"`
	Client_id           uint    `json:"client_id,omitempty"`
	Name                string  `json:"name,omitempty"`
	Email               string  `json:"email,omitempty"`
	Description         string  `json:"description,omitempty"`
	Requestedamount     float64 `json:"requestedamount,omitempty"`
	Requestedcurrency   string  `json:"requestedcurrency,omitempty"`
	Status              string  `json:"status,omitempty"`
	Createdate          string  `json:"createdate,omitempty"`
	Ip                  string  `json:"ip,omitempty"`
	Trackid             string  `json:"trackid,omitempty"`
	Product_name        string  `json:"product_name,omitempty"`
	Invoice_type        int     `json:"invoice_type,omitempty"`
	Order_id            string  `json:"order_id,omitempty"`
	Is_fee_paid_by_user bool    `json:"is_fee_paid_by_user,omitempty"`
	Return_url          string  `json:"return_url"`
}

// for manage Pay Reques
type Invoice_Data struct {
	//gorm.Model
	Invoice_id          uint    `gorm:"primaryKey"`
	Client_id           uint    `json:"client_id,omitempty"`
	Requestedamount     float64 `json:"requestedamount,omitempty"`
	Requestedcurrency   string  `json:"requestedcurrency,omitempty"`
	Order_id            string  `json:"order_id,omitempty"`
	Is_fee_paid_by_user bool    `json:"is_fee_paid_by_user,omitempty"`
}

// Struct for manage support ticket
type Support_Ticket struct {
	//gorm.Model
	Ticket_id      uint   `gorm:"primaryKey"`
	Client_id      uint   `json:"client_id,omitempty"`
	Ticket_subject string `json:"ticket_subject,omitempty"`
	Ticket_desc    string `json:"ticket_desc,omitempty"`
	Status         int    `json:"status,omitempty"`
	Timestamp      string `json:"timestamp,omitempty"`
}

// Struct for manage support ticket reply
type Support_Ticket_Reply struct {
	//gorm.Model
	Reply_id   uint   `gorm:"primaryKey"`
	Ticket_id  int64  `json:"ticket_id,omitempty"`
	Reply_desc string `json:"reply_desc,omitempty"`
	Reply_by   string `json:"reply_by,omitempty"`
	Type       string `json:"type,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}

// Struct for manage crypto wallet
type CryptoWalletList struct {
	Wallet_id      uint   `gorm:"primaryKey"`
	Client_id      uint   `json:"client_id,omitempty"`
	Crypto_code    string `json:"crypto_code,omitempty"`
	Crypto_title   string `json:"crypto_title,omitempty"`
	Crypto_network string `json:"crypto_network,omitempty"`
	Crypto_address string `json:"crypto_address,omitempty"`
	Assetid        int    `json:"assetid,omitempty"`
	Status         int    `json:"status,omitempty"`
}

// for manage Settlement Setting
type SettlementSetting struct {
	Coin_id        int    `json:"coin_id,omitempty"`
	Coin           string `json:"coin,omitempty"`
	Coin_title     string `json:"coin_title,omitempty"`
	Coin_network   string `json:"coin_network,omitempty"`
	Crypto_address string `json:"crypto_address,omitempty"`
	Assetid        int    `json:"assetid,omitempty"`
	Status         int    `json:"status,omitempty"`
	Wallet_id      int    `json:"wallet_id,omitempty"`
}

// for Delete crypto wallet
type CryptoWalletDeleted struct {
	Wallet_id uint `gorm:"primaryKey"`
	Status    int  `json:"status,omitempty"`
}

// For Network List
type Network struct {
	Crypto_code          string `db:"crypto_code"`
	Crypto_network_short string `db:"Crypto_network_short"`
	Crypto_id            string `db:"crypto_id"`
}

// For Coin By Network Details
type NetworkCoin struct {
	Coin         string `db:"coin"`
	Coin_network string `db:"coin_network"`
	Coin_id      string `db:"coin_id"`
}

// For Client API
type ClientAPI struct {
	//gorm.Model
	ID        uint   `gorm:"primaryKey"`
	Client_id uint   `json:"client_id,omitempty"`
	Apikey    string `json:"apikey,omitempty"`
	Status    int    `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// For Manage Client Store
type ClientStore struct {
	//gorm.Model
	ID             uint    `gorm:"primaryKey"`
	Client_id      uint    `json:"client_id,omitempty"`
	Webhookurl     string  `json:"webhookurl,omitempty"`
	Return_url     string  `json:"return_url,omitempty"`
	Success_margin float64 `json:"success_margin,omitempty"`
}

// for Count Transaction By Status
type CountTransactionByStatus struct {
	Total_transactions int `db:"total_transactions"`
	Total_success      int `db:"total_success"`
	Total_failed       int `db:"total_failed"`
	Total_process      int `db:"total_process"`
}

// for Customer Data
type CustomerData struct {
	Customer_id    uint   `gorm:"primaryKey"`
	Client_id      uint   `json:"client_id,omitempty"`
	Customer_name  string `json:"customer_name,omitempty"`
	Customer_email string `json:"customer_email,omitempty"`
	Customer_tid   string `json:"customer_tid,omitempty"`
}

// for Customer Listing
type CustomerList struct {
	Customer_name  string `json:"customer_name,omitempty"`
	Customer_email string `json:"customer_email,omitempty"`
	Total_customer string `json:"total_customer,omitempty"`
	Client_id      uint   `json:"client_id,omitempty"`
}

// for Settlement Request
type SettlementRequest struct {
	Coin_id     int    `json:"coin_id,omitempty"`
	Wallet_id   int    `json:"wallet_id,omitempty"`
	CoinAddress string `json:"coinAddress,omitempty"`
	Action      string `json:"action,omitempty"`
}

// Get get Settlement Response
type SettlementResponse struct {
	Status  int    `json:"status" form:"sender_name"`
	Message string `json:"message" form:"sender_email"`
}

// Get Settlement Coin Details
type SettlementCoin struct {
	Coin         string `json:"coin,omitempty"`
	Coin_title   string `json:"coin_title,omitempty"`
	Coin_network string `json:"coin_network,omitempty"`
}

// Get Settlement Status
type SettlementStatus struct {
	Status string `json:"status" form:"sender_name"`
}

// Get Merchant Email  By Reset Password Hash ID
type EmailByHash struct {
	Password_email string
}
