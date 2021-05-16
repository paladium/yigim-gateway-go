package yigim_test

import (
	"net/http"
	"testing"

	"github.com/paladium/yigim-gateway-go"
	"github.com/stretchr/testify/assert"
)

func TestEncodePaymentCreate(t *testing.T) {
	req, err := http.NewRequest("GET", "http://test.com/api", nil)
	assert.Nil(t, err)
	paymentCreate := yigim.PaymentCreate{
		Reference:   "test",
		Type:        "SMS",
		Token:       "",
		Save:        "y",
		Amount:      100,
		Currency:    "944",
		Biller:      "TEST",
		Description: "Card linking",
		Template:    "TMPL01",
		Language:    "en",
		Callback:    "http://test.com",
		Extra: map[string]string{
			"id": "test",
		},
	}
	paymentCreate.Encode(req, &yigim.Configuration{
		Secret: "test",
	})
	assert.Equal(t, "amount=100&biller=TEST&callback=http%3A%2F%2Ftest.com&currency=944&description=Card+linking&extra=id%3Dtest%3B&language=en&reference=test&save=y&template=TMPL01&token=&type=SMS", req.URL.RawQuery)
	assert.Equal(t, "XKxwnKA23rNeg73ShEcdkw==", req.Header["X-Signature"][0])
}

func TestEncodePaymentStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "http://test.com/api", nil)
	assert.Nil(t, err)
	paymentStatus := yigim.PaymentStatus{
		Reference: "test",
	}
	paymentStatus.Encode(req, &yigim.Configuration{
		Secret: "test",
	})
	assert.Equal(t, "reference=test", req.URL.RawQuery)
	assert.Equal(t, "BaZxxmrv6hJMwIt26m0wuw==", req.Header["X-Signature"][0])
}

func TestEncodePaymentRefund(t *testing.T) {
	req, err := http.NewRequest("GET", "http://test.com/api", nil)
	assert.Nil(t, err)
	paymentRefund := yigim.PaymentRefund{
		Reference: "test",
	}
	paymentRefund.Encode(req, &yigim.Configuration{
		Secret: "test",
	})
	assert.Equal(t, "reference=test", req.URL.RawQuery)
	assert.Equal(t, "BaZxxmrv6hJMwIt26m0wuw==", req.Header["X-Signature"][0])
}

func TestEncodePaymentExecute(t *testing.T) {
	req, err := http.NewRequest("GET", "http://test.com/api", nil)
	assert.Nil(t, err)
	paymentExecute := yigim.PaymentExecute{
		Reference: "test",
		Type:      "SMS",
		Token:     "test",
		Amount:    100,
		Currency:  "944",
		Biller:    "TEST",
	}
	paymentExecute.Encode(req, &yigim.Configuration{
		Secret: "test",
	})
	assert.Equal(t, "amount=100&biller=TEST&currency=944&reference=test&token=test&type=SMS", req.URL.RawQuery)
	assert.Equal(t, "Cyt2oVnqen3I//2pmCiCDw==", req.Header["X-Signature"][0])
}
