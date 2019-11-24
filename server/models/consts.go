package models

const (
	KindCustomerInt = 0
	KindSalerInt    = 1
	KindCustomerStr = "customer"
	KindSalerStr    = "saler"
	PageSize        = 20
)

var KindInt2Str = map[int]string{
	KindCustomerInt: KindCustomerStr,
	KindSalerInt:    KindSalerStr,
}
