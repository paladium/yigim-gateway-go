package yigim

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// IEncodable type of model that can be encoded within query params
type IEncodable interface {
	Encode(*http.Request, *Configuration)
}

// PaymentCreate includes the parameters to initialise the payment - link the card
type PaymentCreate struct {
	//Unique reference number
	Reference string `json:"reference"`
	//Type of payment - SMS or DMS
	Type string `json:"type"`
	//Unique token to use for card linking
	Token string `json:"token"`
	//Save the card or not - y or n
	Save string `json:"save"`
	//Amount multiplied by 100, so 50.55 = 5055
	Amount int `json:"amount"`
	//Currency to charge in, default: 944
	Currency string `json:"currency"`
	//Current biller
	Biller string `json:"biller"`
	//Description for the payment
	Description string `json:"description"`
	//Template to use
	Template string `json:"template"`
	//Language code: az,en or ru
	Language string `json:"language"`
	//Callback url to call after successful card linking
	Callback string `json:"callback"`
	//Extra information
	Extra map[string]string `json:"extra"`
}

// Make a base64-md5 signature
func generateSignature(value string) string {
	md5Bytes := md5.Sum([]byte(value))
	return base64.StdEncoding.EncodeToString(md5Bytes[:])
}

// Encode the payment create into the request
func (model *PaymentCreate) Encode(req *http.Request, appConfig *Configuration) {
	query := req.URL.Query()
	query.Add("reference", model.Reference)
	query.Add("type", model.Type)
	query.Add("token", model.Token)
	query.Add("save", model.Save)
	query.Add("amount", strconv.Itoa(model.Amount))
	query.Add("currency", model.Currency)
	query.Add("biller", model.Biller)
	query.Add("description", model.Description)
	query.Add("template", model.Template)
	query.Add("language", model.Language)
	query.Add("callback", model.Callback)
	extra := ""
	for key, value := range model.Extra {
		extra += fmt.Sprintf("%s=%s;", key, value)
	}
	query.Add("extra", extra)

	req.URL.RawQuery = query.Encode()
	//Now encode
	signature := generateSignature(strings.Join([]string{
		model.Reference,
		model.Type,
		model.Token,
		model.Save,
		strconv.Itoa(model.Amount),
		model.Currency,
		model.Biller,
		model.Description,
		model.Template,
		model.Language,
		model.Callback,
		extra,
		appConfig.Secret,
	}, ""))
	req.Header.Add("X-Signature", signature)
}

// PaymentStatusCode represents a status code response from the server
type PaymentStatusCode string

//Different codes
const (
	// //Newly created transaction, waiting for card data input
	CodeS0 = PaymentStatusCode("S0")

	//Pending DMS transaction (pre-authorized, call 'charge' or 'cancel')
	CodeS1 = PaymentStatusCode("S1")

	//Transaction is in progress
	CodeS2 = PaymentStatusCode("S2")

	//Unknown error
	CodeS3 = PaymentStatusCode("S3")

	//Reversed transaction (cancelled)
	CodeS4 = PaymentStatusCode("S4")

	//Refunded transaction
	CodeS5 = PaymentStatusCode("S5")

	//System malfunction
	CodeS7 = PaymentStatusCode("S7")

	//Approved
	Code00 = PaymentStatusCode("00")

	//Decline, refer to issuer
	Code01 = PaymentStatusCode("01")

	//Decline, expired card
	Code02 = PaymentStatusCode("02")

	//Decline, invalid amount
	Code03 = PaymentStatusCode("03")

	//Decline, inactive card
	Code04 = PaymentStatusCode("04")

	//Decline, insufficient funds
	Code05 = PaymentStatusCode("05")

	//Decline, suspected fraud
	Code06 = PaymentStatusCode("06")

	//Decline, exceeds withdrawal limit
	Code07 = PaymentStatusCode("07")

	//Format error
	Code08 = PaymentStatusCode("08")
)

// ResponseCode response regarding the api code
type ResponseCode int

// Different payment codes
const (
	ResponseCode0 = ResponseCode(0)
	ResponseCode1 = ResponseCode(1)
	ResponseCode2 = ResponseCode(2)
	ResponseCode3 = ResponseCode(3)
	ResponseCode4 = ResponseCode(4)
	ResponseCode5 = ResponseCode(5)
	ResponseCode6 = ResponseCode(6)
)

// SuccessCodes contains a list of results that considered successfull
var SuccessCodes = []ResponseCode{
	ResponseCode0,
}

// PaymentCreateResult contains the result of payment create operation
type PaymentCreateResult struct {
	URL     string       `json:"url"`
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
}

// Secure3DStatus contains status of 3ds authorization
type Secure3DStatus string

// Different status codes
const (
	Secure3D00 = Secure3DStatus("00")
	Secure3D10 = Secure3DStatus("10")
	Secure3D20 = Secure3DStatus("20")
	Secure3D21 = Secure3DStatus("21")
	Secure3D22 = Secure3DStatus("22")
	Secure3D23 = Secure3DStatus("23")
	Secure3D24 = Secure3DStatus("24")
	Secure3D25 = Secure3DStatus("25")
	Secure3D30 = Secure3DStatus("30")
)

// PaymentStatus is used to retrieve status for a transaction
type PaymentStatus struct {
	//Unique reference number
	Reference string `json:"reference"`
}

// Encode the payment status into the request
func (model *PaymentStatus) Encode(req *http.Request, appConfig *Configuration) {
	query := req.URL.Query()
	query.Add("reference", model.Reference)
	req.URL.RawQuery = query.Encode()
	//Now encode
	signature := generateSignature(strings.Join([]string{
		model.Reference,
		appConfig.Secret,
	}, ""))
	req.Header.Add("X-Signature", signature)
}

// PaymentStatusResult contains the result of the payment
type PaymentStatusResult struct {
	Reference      string            `json:"reference"`
	Datetime       string            `json:"datetime"`
	Type           string            `json:"type"`
	Token          string            `json:"token"`
	Pan            *string           `json:"pan"`
	Expiry         *string           `json:"expiry"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	Biller         string            `json:"biller"`
	System         string            `json:"system"`
	Issuer         *string           `json:"issuer"`
	Rrn            string            `json:"rrn"`
	Secure3DStatus *Secure3DStatus   `json:"3ds" bson:"3ds"`
	Approval       string            `json:"approval"`
	Status         PaymentStatusCode `json:"status"`
	Code           ResponseCode      `json:"code"`
	Message        string            `json:"message"`
	Extra          []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"extra"`
}

// PaymentRefund refunfs the previous locked amount for card verification
type PaymentRefund struct {
	Reference string `json:"reference"`
}

// Encode the payment refund into the request
func (model *PaymentRefund) Encode(req *http.Request, appConfig *Configuration) {
	query := req.URL.Query()
	query.Add("reference", model.Reference)
	req.URL.RawQuery = query.Encode()
	//Now encode
	signature := generateSignature(strings.Join([]string{
		model.Reference,
		appConfig.Secret,
	}, ""))
	req.Header.Add("X-Signature", signature)
}

// PaymentRefundResult contains the result of payment refund
type PaymentRefundResult struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
}

// PaymentExecute stores data to execute a payment using a token
type PaymentExecute struct {
	//Unique payment (order) ID for future usage
	Reference string `json:"reference"`
	//SMS - Single (default), DMS - Dual Message System
	Type string `json:"type"`
	//Card unique token (optional)
	Token string `json:"token"`
	//Amount multiplied by 100, so 50.55 = 5055
	Amount int `json:"amount"`
	//Numeric ISO4712 currency code
	Currency string `json:"currency"`
	//Current biller
	Biller string `json:"biller"`
}

// Encode the payment status into the request
func (model *PaymentExecute) Encode(req *http.Request, appConfig *Configuration) {
	query := req.URL.Query()
	query.Add("reference", model.Reference)
	query.Add("type", model.Type)
	query.Add("token", model.Token)
	query.Add("amount", strconv.Itoa(model.Amount))
	query.Add("currency", model.Currency)
	query.Add("biller", model.Biller)
	req.URL.RawQuery = query.Encode()
	//Now encode
	signature := generateSignature(strings.Join([]string{
		model.Reference,
		model.Type,
		model.Token,
		strconv.Itoa(model.Amount),
		model.Currency,
		model.Biller,
		appConfig.Secret,
	}, ""))
	req.Header.Add("X-Signature", signature)
}

// PaymentExecuteResult contains the result of payment execute action
type PaymentExecuteResult struct {
	Reference      string            `json:"reference"`
	Datetime       string            `json:"datetime"`
	Type           string            `json:"type"`
	Pan            *string           `json:"pan"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	Biller         string            `json:"biller"`
	System         string            `json:"system"`
	Issuer         *string           `json:"issuer"`
	Rrn            string            `json:"rrn"`
	Secure3DStatus *Secure3DStatus   `json:"3ds" bson:"3ds"`
	Approval       string            `json:"approval"`
	Status         PaymentStatusCode `json:"status"`
	Code           ResponseCode      `json:"code"`
	Message        string            `json:"message"`
}
