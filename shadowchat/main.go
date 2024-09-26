package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/skip2/go-qrcode"
)

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
	MaxNameChars     int      `json:"max_name_chars"`
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
	MinAmnt float64
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

func main() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("reading config.json")
	//this function execs itself to close the file after the main function is over
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var conf configJson
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		panic(err)
	}

	ScamThresold = conf.MinimumDonation
	MessageMaxChar = conf.MaxMessageChars
	NameMaxChar = conf.MaxNameChars
	rpcURL = conf.RPCWalletURL
	username = conf.WebViewUsername
	password = conf.WebViewPassword
	AlertWidgetRefreshInterval = conf.OBSWidgetRefresh
	enableEmail = conf.EnableEmail
	smtpHost = conf.SMTPServer
	smtpPost = conf.SMTPPort
	smtpUser = conf.SMTPUser
	smtpPass = conf.SMTPPass
	sendTo = conf.SendToEmail
	if conf.Checked == true {
		checked = " checked"
	}

	fmt.Println(fmt.Sprintf("email notifications enabled?: %t", enableEmail))
	fmt.Println(fmt.Sprintf("OBS Alert path: /alert?auth=%s", password))

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/style.css")
	})
	http.HandleFunc("/xmr.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/xmr.svg")
	})

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pay", paymentHandler)
	http.HandleFunc("/check", checkHandler)

}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	payload := strings.NewReader(`{"jsonrpc":"2.0","id":"0","method":"get_address"}`)
	req, _ := http.NewRequest("POST", rpcURL, payload)
	req.Header.Set("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	resp := &getAddress{}

	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		fmt.Println(err.Error())
	}

	var c checkPage
	c.Meta = `<meta http-equiv="Refresh" content="3">`
	c.Addy = resp.Result.Address
	c.PayID = r.FormValue("id")
	c.Name = truncateStrings(r.FormValue("name"), NameMaxChar)
	c.Msg = truncateStrings(r.FormValue("msg"), MessageMaxChar)
	c.Media = r.FormValue("media")
	c.Receipt = "Waiting for payment..."

	payload2 := strings.NewReader(`{"jsonrpc":"2.0", "id":"0","method":"get_transfers","params":{"in":true, "pool":true, "account_index":0}}`)
	req2, _ := http.NewRequest("POST", "http://127.0.0.1:28088/json_rpc", payload2)

	req2.Header.Set("Content-type", "application/json")
	res2, _ := http.DefaultClient.Do(req2)
	resp2 := &GetTransferResponse{}
}

func truncateStrings(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for !utf8.ValidString(s[:n]) {
		n--
	}
	return s[:n]
}
func condenseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
func indexHandler(w http.ResponseWriter, _ *http.Request) {
	var i indexDisplay
	i.MaxChar = MessageMaxChar
	i.MinAmnt = ScamThresold
	i.Checked = checked
	err := indexTemplate.Execute(w, i)
	if err != nil {
		fmt.Println(err)
	}
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	payload := strings.NewReader(`{"jsonrpc":"2.0", "id":"0","method":"make_integrated_address"}`)
	req, err := http.NewRequest("POST", rpcURL, payload)
	if err == nil {
		req.Header.Set("Content-Type", "applications/json")
		res, err := http.DefaultClient.Do(req)
		if err == nil {
			resp := &rpcResponse{}
			if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
				fmt.Println(err.Error())
			}

			var s superChat
			s.Amount = html.EscapeString(r.FormValue("amount"))
			if r.FormValue("amount") == "" {
				s.Amount = fmt.Sprint(ScamThresold)
			}
			if r.FormValue("amount") == "" {
				s.Name = "Anonymous"
			} else {
				s.Name = html.EscapeString(truncateStrings(condenseSpaces(r.FormValue("name")), NameMaxChar))
			}
			s.Message = html.EscapeString(truncateStrings(condenseSpaces(r.FormValue("message")), MessageMaxChar))
			s.Media = html.EscapeString(r.FormValue("media"))
			s.PayID = html.EscapeString(resp.Result.PaymentID)
			s.Address = resp.Result.IntegratedAddress

			params := url.Values{}
			params.Add("id", resp.Result.PaymentID)
			params.Add("name", s.Name)
			params.Add("msg", r.FormValue("message"))
			params.Add("media", condenseSpaces(s.Media))
			params.Add("show", html.EscapeString(r.FormValue("showAmount")))
			s.CheckURL = params.Encode()

			tmp, _ := qrcode.Encode(fmt.Sprintf("monero:%s?tx_amount=%s", resp.Result.IntegratedAddress, s.Amount), qrcode.Low, 320)
			s.QRB64 = base64.StdEncoding.EncodeToString(tmp)
			err := payTemplate.Execute(w, s)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
