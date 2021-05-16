package yigim_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/paladium/yigim-gateway-go"
	"github.com/stretchr/testify/assert"
)

func getAppConfig() *yigim.Configuration {
	return &yigim.Configuration{
		Secret:   os.Getenv("YIGIM_SECRET"),
		Merchant: os.Getenv("YIGIM_MERCHANT"),
		Address:  "https://sandbox.api.pay.yigim.az",
	}
}

func TestYigimClientCreate(t *testing.T) {
	client := yigim.NewClient(getAppConfig())
	reference := client.GetRef()
	result, err := client.Create(&yigim.PaymentCreate{
		Reference:   reference,
		Type:        "SMS",
		Save:        "y",
		Amount:      100,
		Currency:    "944",
		Biller:      os.Getenv("YIGIM_BILLER"),
		Description: "Test",
		Template:    "TPL0001",
		Language:    "en",
		Callback:    "http://test.com",
	})
	assert.Nil(t, err)
	assert.Equal(t, yigim.ResponseCode0, result.Code)
	assert.Equal(t, "OK (No error)", result.Message)
	assert.Regexp(t, regexp.MustCompile(`https:\/\/sandbox\.pay\.yigim\.az\/payment\/[0-9A-Z]+`), result.URL)
}

func TestYigimClientStatus(t *testing.T) {
	client := yigim.NewClient(getAppConfig())
	reference := client.GetRef()
	//Create the payment first, then extract the payment status
	_, err := client.Create(&yigim.PaymentCreate{
		Reference:   reference,
		Type:        "SMS",
		Save:        "y",
		Amount:      100,
		Currency:    "944",
		Biller:      os.Getenv("YIGIM_BILLER"),
		Description: "Test",
		Template:    "TPL0001",
		Language:    "en",
		Callback:    "http://test.com",
	})
	assert.Nil(t, err)
	result, err := client.Status(&yigim.PaymentStatus{
		Reference: reference,
	})
	assert.Nil(t, err)
	assert.Equal(t, yigim.ResponseCode0, result.Code)
	assert.Equal(t, "SMS", result.Type)
	assert.Equal(t, float64(1), result.Amount)
	assert.Equal(t, "AZN", result.Currency)
}

func TestYigimClientExecute(t *testing.T) {
	client := yigim.NewClient(getAppConfig())
	reference := client.GetRef()
	//Create the payment first, then extract the payment status
	result, err := client.Execute(&yigim.PaymentExecute{
		Reference: reference,
		Type:      "SMS",
		Token:     os.Getenv("YIGIM_TOKEN"),
		Amount:    100,
		Currency:  "944",
		Biller:    os.Getenv("YIGIM_BILLER"),
	})
	assert.Nil(t, err)
	assert.Equal(t, yigim.ResponseCode0, result.Code)
	assert.Equal(t, "OK (No error)", result.Message)
}

func TestYigimClientRefund(t *testing.T) {
	client := yigim.NewClient(getAppConfig())
	reference := client.GetRef()
	//Create the payment first, then extract the payment status
	_, err := client.Execute(&yigim.PaymentExecute{
		Reference: reference,
		Type:      "SMS",
		Token:     os.Getenv("YIGIM_TOKEN"),
		Amount:    100,
		Currency:  "944",
		Biller:    os.Getenv("YIGIM_BILLER"),
	})
	assert.Nil(t, err)
	result, err := client.Refund(&yigim.PaymentRefund{
		Reference: reference,
	})
	assert.Nil(t, err)
	assert.Equal(t, "OK (No error)", result.Message)
}
