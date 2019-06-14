package bhandler

type BillingHandler interface {
	GetType() byte
	GetResponse(request *BillingData) *BillingData
}
