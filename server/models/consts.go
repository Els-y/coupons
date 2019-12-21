package models

import (
	"github.com/nats-io/nats.go"
)

const (
	KindCustomerInt            = 0
	KindSalerInt               = 1
	KindCustomerStr            = "customer"
	KindSalerStr               = "saler"
	NatsUrl                    = "nats:4222"
	AssignCoupon_Subj          = "AssignCounpon"
	AssignCounpon_SubscribeNum = 4
)

var KindInt2Str = map[int]string{
	KindCustomerInt: KindCustomerStr,
	KindSalerInt:    KindSalerStr,
}

var Nc *nats.Conn
var NatsEncodedConn *nats.EncodedConn
var Nat_err error
