package models

import (
	"github.com/nats-io/nats.go"
)

const (
	KindCustomerInt = 0
	KindSalerInt    = 1
	KindCustomerStr = "customer"
	KindSalerStr    = "saler"
	NatsUrl = "nats://127.0.0.1:4222"
	AssignCoupon_Subj = "AssignCounpon"
	AssignCounpon_SubscribeNum = 4
)

var KindInt2Str = map[int]string{
	KindCustomerInt: KindCustomerStr,
	KindSalerInt:    KindSalerStr,
}

var nc *nats.Conn
var NatsEncodedConn *nats.EncodedConn
var nat_err error
