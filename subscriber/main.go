// Copyright 2012-2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"time"
	"github.com/Els-y/coupons/subscriber/models"
	"github.com/Els-y/coupons/subscriber/pkgs/setting"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type SubscribeAssignCoupon struct {
	SalerName string
	TokenUsername string
	CouponName string
	CouponStock int
}

func init() {
	setting.Setup()
	models.Setup() 
}

func main() {
	// var urls = "nats:4222"
	// var subj = "temp"
	// var queue = "my-queue"
	
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, nat_err := nats.Connect(models.NatsUrl, opts...)
	if nat_err != nil {
		logrus.Infof("[subscribe] subscribe nats.Connect url error, url: %v, err: %v", models.NatsUrl, nat_err.Error())
		return 
	}

	NatsEncodedConn, nat_err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if nat_err != nil {
		logrus.Infof("[subscribe] subscribe nats.NewEncodedConn error, err: %v", nat_err.Error())
		return
	}

	// i := 0

	logrus.Infof("[subscribe] subscribe worker subj: %v, queue: %v", models.AssignCoupon_Subj, models.AssignCoupon_Subj + "_queue")
	NatsEncodedConn.QueueSubscribe(models.AssignCoupon_Subj, models.AssignCoupon_Subj + "_queue", func(s *SubscribeAssignCoupon) {
		logrus.Infof("[queue subscribe] recieve publish subj: %v, queue: %v, info: %+v", models.AssignCoupon_Subj, models.AssignCoupon_Subj + "_queue", s)
		models.AssignCoupon(s.SalerName, s.TokenUsername, s.CouponName, s.CouponStock)
	})

	select{}

	// nc.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
	// 	var cur_msg string
	// 	cur_msg = string(msg.Data)
	// 	infors := strings.Split(cur_msg, ".")
	// 	log.Printf("%s", cur_msg)
	// 	switch infors[2] {
	// 		case "decre":
	// 			isSuccess, err := models.DecreaseAmount(infors[0], infors[1])
	// 			_, err2 := models.GetCoupon(infors[0], infors[1])
	// 			if err == nil{
	// 				if err2 == nil {
	// 					// log.Printf("%d", query_coupon.Left)
	// 					log.Printf("good!")
	// 				} else {
	// 					log.Printf("query coupon wrong! %d", isSuccess)
	// 				}
					
	// 			}
	// 		case "check":
	// 			log.Printf("Check is unfinished")
	// 		default:
	// 			log.Printf("wrong")
	// 	}
	// 	// models.db.save(&coupon)
		
	// 	log.Printf("%s", infors[0])

	// 	i += 1
	// 	printMsg(msg, i)
	// })
	// nc.Flush()

	// if err := nc.LastError(); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Listening on [%s]", subj)

	// // Setup the interrupt handler to drain so we don't miss
	// // requests when scaling down.
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
	// log.Println()
	// log.Printf("Draining...")
	// nc.Drain()
	// log.Fatalf("Exiting")
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
