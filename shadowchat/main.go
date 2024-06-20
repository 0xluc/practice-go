package main

import "text/template"

var ScamThresold float64 = 0.005 // minimum donation
var MediaMin float64 = 0.025     // unused
var MessageMaxChar int = 250
var NameMaxChar int = 25
var rpcURL string = "http://127.0.0.1:28088/json_rpc"
var username string = "admin"                // chat log /view page
var AlertWidgetRefreshInterval string = "10" // seconds

// tbis is the password for both the /view page and OBS /alert page
// exmaple OBS url: https://example.com/alert?auth=adminadmin
var password string = "adminadmin"
var checked string = ""

// Email settings
var enableEmail bool = false
var smtpHost string = "smtp.mail.com"
var smtpPost string = "587"
var smtpUser string = "example@mail.com"
var smtpPass string = "password"
var sendTo = []string{"example@email.com"} // Comma separeted recipient list

var indexTemplate *template.Template
var payTemplate *template.Template
var checkTemplate *template.Template
var alertTemplate *template.Template
var viewTemplate *template.Template
var topWidgetTemplate *template.Template

type configJson struct {
	MinimumDonation  float64  `json:"minimum_donation"`
	MaxMessageChars  int      `json:"max_message_chars"`
	maxNameChars     int      `json:"max_name_chars"`
	RPCWalletURL     string   `json:"rpc_wallet_url"`
	WebViewUsername  string   `json:"web_view_username"`
	WebViewPassword  string   `json:"web_view_password"`
	OBSWidgetRefresh string   `json:"obs_widget_refresh"`
	Checked          bool     `json:"checked"`
	EnableEmail      bool     `json:"enable_email"`
	SMTPServer       string   `json:"smtp_server"`
	SMTPPort         string   `json:"smtp_port"`
	SMTPUser         string   `json:"smtp_user"`
	SMTPPass         string   `json:"smtp_pass"`
	SendToEmail      []string `json:"send_to_email"`
}

type checkPage struct {
	Addy     string
	PayID    string
	Received float64
	Meta     string
	Name     string
	Msg      string
	Receipt  string
	Media    string
}

type superChat struct {
	Name     string
	Message  string
	Media    string
	Amount   string
	Address  string
	QRB64    string
	PayID    string
	CheckURL string
}

type csvLog struct {
	ID            string
	Name          string
	Message       string
	Amount        string
	DisplayToddle string
	Refresh       string
}

type indexDisplay struct {
	MaxChar int
	minAmnt float64
	Checked string
}

type viewPageData struct {
	ID      []string
	Name    []string
	Message []string
	Amount  []string
	Display []string
}

type rpcResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		IntegratedAddress string `json:"integrated_address"`
		PaymentID         string `json:"payment_id"`
	} `json:"result"`
}

type getAddress struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Address   string `json:"address"`
		Addresses []struct {
			Address      string `json:"address"`
			AddressIndex int    `json:"address_index"`
			Label        string `json:"label"`
			Used         bool   `json:"used"`
		} `json:"addresses"`
	} `json:"result"`
}

type MoneroPrice struct {
	Monero struct {
		Usd float64 `json:"usd"`
	} `json:"monero"`
}

type GetTransferResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		In []struct {
			Address         string  `json:"address"`
			Amount          int64   `json:"amount"`
			Amounts         []int64 `json:"amounts"`
			Confirmations   int     `json:"confirmations"`
			DoubleSpendSeen bool    `json:"double_spend_seen"`
			Fee             int     `json:"fee"`
			Height          int     `json:"height"`
			Locked          bool    `json:"locked"`
			Note            string  `json:"note"`
			PaymentID       string  `json:"payment_id"`
			SubaddrIndex    struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_index"`
			SubaddrIndices []struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_indices"`
			SuggestedConfirmationsThreshold int    `json:"suggested_confirmations_threshold"`
			Timestamp                       int    `json:"timestamp"`
			Txid                            string `json:"txid"`
			Type                            string `json:"type"`
			UnlockTime                      int    `json:"unlock_time"`
		} `json:"in"`
		Pool []struct {
			Address         string  `json:"address"`
			Amount          int64   `json:"amount"`
			Amounts         []int64 `json:"amounts"`
			DoubleSpendSeen bool    `json:"double_spend_seen"`
			Fee             int     `json:"fee"`
			Height          int     `json:"height"`
			Locked          bool    `json:"locked"`
			Note            string  `json:"note"`
			PaymentID       string  `json:"payment_id"`
			SubaddrIndex    struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_index"`
			SubaddrIndices []struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
			} `json:"subaddr_indices"`
			SuggestedConfirmationsThreshold int    `json:"suggested_confirmations_threshold"`
			Timestamp                       int    `json:"timestamp"`
			Txid                            string `json:"txid"`
			Type                            string `json:"type"`
			UnlockTime                      int    `json:"unlock_time"`
		} `json:"pool"`
	} `json:"result"`
}

func main() {}
