package mq

import (
	"github.com/Els-y/coupons/server/pkgs/setting"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var Nc *nats.Conn
var NatsEncodedConn *nats.EncodedConn

func Setup() {
	Nc, err := nats.Connect(setting.NatsSetting.Host)
	if err != nil {
		logrus.WithError(err).Fatal("mq.Setup err")
	}

	NatsEncodedConn, err = nats.NewEncodedConn(Nc, nats.JSON_ENCODER)
	if err != nil {
		logrus.WithError(err).Fatal("subscribe nats.NewEncodedConn error")
	}
}
