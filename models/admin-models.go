package models

import (
	"time"
)

// Struct for Admin Section
// Here create struct for used in admin functions and methods

// For display login history listing
type AdminLoginHistory struct {
	//gorm.Model
	Token_id    uint   `gorm:"primaryKey"`
	Client_id   uint   `json:"client_id,omitempty"`
	Login_time  string `json:"login_time,omitempty"`
	Logout_time string `json:"logout_time,omitempty"`
	Login_ip    string `json:"login_ip,omitempty"`
	//Count       int64  `json:"count"`
}

// For display Admin User listing
type AdminLoginList struct {
	//gorm.Model
	Admin_id  uint   `gorm:"primaryKey"`
	Full_name string `json:"full_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Status    int    `json:"status,omitempty"`
	Role      string `json:"role,omitempty"`
}

// for manage coin List
type CoinList struct {
	Coin_id         uint   `gorm:"primaryKey"`
	Coin            string `json:"coin,omitempty"`
	Icon            string `json:"icon,omitempty"`
	Status          int    `json:"status,omitempty"`
	Coin_title      string `json:"coin_title,omitempty"`
	Coin_pay_url    string `json:"coin_pay_url,omitempty"`
	Coin_network    string `json:"coin_network,omitempty"`
	Coin_status_url string `json:"coin_status_url,omitempty"`
}

// for manage coin address
type PayCoinAddress struct {
	Coin_title   string `json:"coin_title,omitempty"`
	Coin_id      int    `json:"coin_id,omitempty"`
	Coin_pay_url string `json:"coin_pay_url,omitempty"`
	Address      string `json:"address,omitempty"`
	Lastupdate   string `json:"lastupdate,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Address_id   string `json:"address_id,omitempty"`
	Coin_network string `json:"coin_network,omitempty"`
}

// for update coin address
type CoinAddress struct {
	Coin    string `json:"coin,omitempty"`
	Address string `json:"address,omitempty"`
	Coin_id int    `json:"coin_id,omitempty"`
}

// for fet coin address url
type CoinAddressUrl struct {
	Coin            string `json:"coin,omitempty"`
	Coin_status_url string `json:"coin_status_url,omitempty"`
}

// for delete coin
type CoinDeleted struct {
	Coin_id uint `gorm:"primaryKey"`
	Status  int  `json:"status,omitempty"`
}

// for manage Address Listing
type AddressListing struct {
	Address_id    uint      `gorm:"primaryKey"`
	Coin          string    `json:"coin"`
	Status        int       `json:"status"`
	Address       string    `json:"address"`
	Lastupdate    time.Time `json:"lastupdate"`
	Coin_id       int       `json:"coin_id"`
	Token_details string    `json:"token_details"`
}

// for delete Address
type AddressDeleted struct {
	Address_id uint `gorm:"primaryKey"`
	Status     int  `json:"status,omitempty"`
}

// for update address date
type AddressDateUpdate struct {
	Address_id uint      `gorm:"primaryKey"`
	Lastupdate time.Time `json:"lastupdate"`
}

// for manage Currency
type CurrencyList struct {
	Currency_id             uint   `gorm:"primaryKey"`
	Currency_name           string `json:"currency_name,omitempty"`
	Currency_code           string `json:"currency_code,omitempty"`
	Currency_territory      string `json:"currency_territory,omitempty"`
	Currency_icon_bootstrap string `json:"currency_icon_bootstrap,omitempty"`
	Status                  int    `json:"status,omitempty"`
}

// for delete currency
type CurrencyDeleted struct {
	Currency_id uint `gorm:"primaryKey"`
	Status      int  `json:"status,omitempty"`
}

// for display Crypto Currency List
type CryptoCurrencyList struct {
	Crypto_id            uint   `gorm:"primaryKey"`
	Crypto_code          string `json:"crypto_code,omitempty"`
	Crypto_title         string `json:"crypto_title,omitempty"`
	Crypto_network       string `json:"crypto_network,omitempty"`
	Crypto_network_short string `json:"crypto_network_short,omitempty"`
	Status               int    `json:"status,omitempty"`
}

// for manage Crypto Currency Data
type CryptoCurrencyData struct {
	Crypto_code    string `json:"crypto_code,omitempty"`
	Crypto_title   string `json:"crypto_title,omitempty"`
	Crypto_network string `json:"crypto_network,omitempty"`
}

// for Delete Crypto Currency
type CryptoCurrencyDeleted struct {
	Crypto_id uint `gorm:"primaryKey"`
	Status    int  `json:"status,omitempty"`
}

// for manage Acquirer
type AcquirerList struct {
	Acquirer_id    uint   `gorm:"primaryKey"`
	Acquirer_title string `json:"acquirer_title,omitempty"`
	Api_key        string `json:"api_key,omitempty"`
	Secret_key     string `json:"secret_key,omitempty"`
	Endpoint_url   string `json:"endpoint_url,omitempty"`
	Callback_url   string `json:"callback_url,omitempty"`
	Success_url    string `json:"success_url,omitempty"`
	Failed_url     string `json:"failed_url,omitempty"`
	Json_data      string `json:"json_data"`
	Status         int    `json:"status,omitempty"`
}

// for delete acquirer
type AcquirerDeleted struct {
	Acquirer_id uint `gorm:"primaryKey"`
	Status      int  `json:"status,omitempty"`
}

// for Fee Details
type FeesDetails struct {
	Client_id             uint   `gorm:"primaryKey"`
	Min_amount_fee        string `json:"min_amount_fee,omitempty"`
	Amount_fee_in_percent string `json:"amount_fee_in_percent,omitempty"`
}

// for Fee Details
type FeesDetail struct {
	Fee_id                uint   `gorm:"primaryKey"`
	Min_amount_fee        string `json:"min_amount_fee,omitempty"`
	Amount_fee_in_percent string `json:"amount_fee_in_percent,omitempty"`
	Client_id             int    `json:"client_id,omitempty"`
}

// for manage Admin
type AdminList struct {
	Admin_id  uint   `gorm:"primaryKey"`
	Username  string `json:"username,omitempty"`
	Full_name string `json:"full_name,omitempty"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"role,omitempty"`
	//Timestamp string `json:"timestamp,omitempty"`
	Status int `json:"status,omitempty"`
}

// for delete Admin
type AdminDeleted struct {
	Admin_id uint `gorm:"primaryKey"`
	Status   int  `json:"status,omitempty"`
}

// for manage Admin password
type AdminPassword struct {
	Admin_id uint   `gorm:"primaryKey"`
	Password string `json:"password,omitempty"`
}

// for update Admin details
type AdminUpdate struct {
	Admin_id  uint   `gorm:"primaryKey"`
	Full_name string `json:"full_name,omitempty"`
	Role      string `json:"role,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// for delete Email Template
type TemplateDeleted struct {
	Template_id uint `gorm:"primaryKey"`
	Status      int  `json:"status,omitempty"`
}

// For COIN ID Listing
type CoinIDDetails struct {
	//gorm.Model
	Id      uint   `gorm:"primaryKey"`
	Coin_id string `json:"coin_id,omitempty"`
	Symbol  string `json:"symbol,omitempty"`
	Name    string `json:"name,omitempty"`
}

// for  manage Developer Guide
type DeveloperGuide struct {
	Id        uint   `gorm:"primaryKey"`
	Title     string `json:"title,omitempty"`
	Heading   string `json:"heading,omitempty"`
	Functions string `json:"functions,omitempty"`
	Used      string `json:"used,omitempty"`
	Status    int    `json:"status,omitempty"`
}

// for delete Developer Guide data
type DeveloperGuideDeleted struct {
	Id     uint `gorm:"primaryKey"`
	Status int  `json:"status,omitempty"`
}
