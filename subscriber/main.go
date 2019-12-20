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
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
	"strings"
	"github.com/Els-y/coupons/subscriber/models"
	"github.com/Els-y/coupons/subscriber/pkgs/setting"
	
	"github.com/nats-io/nats.go"
)

// NOTE: Can test with demo servers.
// nats-qsub -s demo.nats.io <subject> <queue>
// nats-qsub -s demo.nats.io:4443 <subject> <queue> (TLS version)

func usage() {
	log.Printf("Usage: nats-qsub [-s server] [-creds file] [-t] <subject> <queue>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}

func init() {
	setting.Setup()
	models.Setup()
}

func main() {
	var urls = "nats:4222"
	var subj = "temp"
	var queue = "my-queue"
	
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)


	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	i := 0

	nc.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		var cur_msg string
		cur_msg = string(msg.Data)
		infors := strings.Split(cur_msg, ".")
		log.Printf("%s", cur_msg)
		switch infors[2] {
			case "decre":
				isSuccess, err := models.DecreaseAmount(infors[0], infors[1])
				_, err2 := models.GetCoupon(infors[0], infors[1])
				if err == nil{
					if err2 == nil {
						// log.Printf("%d", query_coupon.Left)
						log.Printf("good!")
					} else {
						log.Printf("query coupon wrong! %d", isSuccess)
					}
					
				}
			case "check":
				log.Printf("Check is unfinished")
			default:
				log.Printf("wrong")
		}
		// models.db.save(&coupon)
		
		log.Printf("%s", infors[0])

		i += 1
		printMsg(msg, i)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")
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
