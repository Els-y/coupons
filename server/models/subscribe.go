package models

import (
	"github.com/nats-io/nats.go"
	"strconv"
	"github.com/sirupsen/logrus"
)

type SubscribeAssignCoupon struct {
	SalerName string
	TokenUsername string
	CouponData Coupon
}

func init() {
	nc, nat_err = nats.Connect(NatsUrl)
	if nat_err != nil {
		logrus.Infof("[queue.init] subscribe nats.Connect url error, url: %v, err: %v", NatsUrl, nat_err.Error())
		return
	}

	NatsEncodedConn, nat_err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if nat_err != nil {
		logrus.Infof("[queue.init] subscribe nats.NewEncodedConn error, err: %v", nat_err.Error())
		return
	}
	StartSubScribe()
}

func StartSubScribe() {
	for i := 0; i < AssignCounpon_SubscribeNum; i++ {
		// logrus.Infof("[queue subscribe] start worker: %v subj: %v", "worker" + strconv.Itoa(i), addUser_subj)
		startServiceAddUser(AssignCoupon_Subj, "worker" + strconv.Itoa(i), AssignCoupon_Subj + "_queue")
	}
}

func startServiceAddUser(subj string, name string, queue string) {
	go asyncAddUser(subj, name, queue)
}

func asyncAddUser(subj string, name string, queue string) {
	NatsEncodedConn.QueueSubscribe(subj, queue, func(s *SubscribeAssignCoupon) {
		// logrus.Infof("[queue subscribe] recieve publish worker: %v subj: %v, queue: %v", name, subj, queue)
		AssignCoupon(s.SalerName, s.TokenUsername, &(s.CouponData))
	})
	select {}
}


