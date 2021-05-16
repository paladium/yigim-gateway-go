package yigim

// YigimClient client for working with yigim API
type YigimClient interface {
	Create(paymentCreate *PaymentCreate) (*PaymentCreateResult, error)
	GetRef() string
	Status(paymentStatus *PaymentStatus) (*PaymentStatusResult, error)
	Refund(paymentRefund *PaymentRefund) (*PaymentRefundResult, error)
	Execute(paymentExecute *PaymentExecute) (*PaymentExecuteResult, error)
}
