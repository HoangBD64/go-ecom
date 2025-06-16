package response

import "github.com/HoangBD64/go-ecom/pkg/domain"

type OrderPayment struct {
	PaymentType  domain.PaymentType `json:"payment_type"`
	PaymentOrder any                `json:"payment_order"`
}
