# Yigim gateway SDK
The package provides ability to access the Yigim payment's api via the Go language.

# Usage
To install run:
```bash
go get github.com/paladium/yigim-gateway-go
```

## Importing
```go
import yigim "github.com/paladium/yigim-gateway-go"
```

## Start card linking
```go
import yigim "github.com/paladium/yigim-gateway-go"

client := yigim.NewClient(&yigim.Configuration{
    Secret:   "YOUR_SECRET",
    Merchant: "YOUR_MERCHANT",
    Address:  "https://sandbox.api.pay.yigim.az",
})
result, err := client.Create(&yigim.PaymentCreate{
    Reference:   reference,
    Type:        "SMS",
    Save:        "y",
    Amount:      100,
    Currency:    "944",
    Biller:      "[YOUR_BILLER]",
    Description: "Test",
    Template:    "TPL0001",
    Language:    "en",
    Callback:    "[YOUR_WEBHOOK_URL]",
})
if err != nil{
    panic(err)
}
```

## Testing
To run tests, open the file ```run_tests.sh``` and set your variables, after that run:
```bash
sh ./run_tests.sh
```
