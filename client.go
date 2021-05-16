package yigim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type yigimClient struct {
	config *Configuration
}

func NewClient(config *Configuration) YigimClient {
	return &yigimClient{
		config: config,
	}
}

// Execute the request
func (client *yigimClient) executeRequest(req *http.Request) ([]byte, error) {
	httpClient := http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot execute the request")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot read the body")
	}
	return body, nil
}

// Build the request by attaching all the nessesary parameters and headers
func (client *yigimClient) buildRequest(path string, data IEncodable) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", client.config.Address, path), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot build the request")
	}
	data.Encode(req, client.config)
	req.Header.Add("X-Merchant", client.config.Merchant)
	req.Header.Add("X-Type", "JSON")
	return req, nil
}

// Create executed create action
//This command is used for initial payment registration and optionally may start card linking
//procedure. Returns URL to redirect cardholder to for entering of card information.
func (client *yigimClient) Create(paymentCreate *PaymentCreate) (*PaymentCreateResult, error) {
	req, err := client.buildRequest("payment/create", paymentCreate)
	if err != nil {
		return nil, err
	}
	response, err := client.executeRequest(req)
	if err != nil {
		return nil, err
	}
	var result PaymentCreateResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshal the response")
	}
	return &result, nil
}

// GetRef returns unique reference
func (client *yigimClient) GetRef() string {
	return fmt.Sprintf("ref-%s", uuid.New().String()[0:10])
}

// Status returns the status of a transaction by a reference
func (client *yigimClient) Status(paymentStatus *PaymentStatus) (*PaymentStatusResult, error) {
	req, err := client.buildRequest("payment/status", paymentStatus)
	if err != nil {
		return nil, err
	}
	response, err := client.executeRequest(req)
	if err != nil {
		return nil, err
	}
	var result PaymentStatusResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshal the response")
	}
	return &result, nil
}

// Refund the previous payment
func (client *yigimClient) Refund(paymentRefund *PaymentRefund) (*PaymentRefundResult, error) {
	req, err := client.buildRequest("payment/refund", paymentRefund)
	if err != nil {
		return nil, err
	}
	response, err := client.executeRequest(req)
	if err != nil {
		return nil, err
	}
	var result PaymentRefundResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshal the response")
	}
	return &result, nil
}

// Execute a payment using a saved token
func (client *yigimClient) Execute(paymentExecute *PaymentExecute) (*PaymentExecuteResult, error) {
	req, err := client.buildRequest("payment/execute", paymentExecute)
	if err != nil {
		return nil, err
	}
	response, err := client.executeRequest(req)
	if err != nil {
		return nil, err
	}
	var result PaymentExecuteResult
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshal the response")
	}
	return &result, nil
}
