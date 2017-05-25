/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package state

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/teleport/lib/auth/api"
	"github.com/gravitational/teleport/lib/auth/api/protogen"
	"github.com/gravitational/teleport/lib/events"
	"github.com/gravitational/teleport/lib/session"
	"github.com/gravitational/trace"

	"github.com/codahale/hdrhistogram"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	errNotSupported = trace.BadParameter("method not supported")
)

const (
	// MaxQueueSize determines how many logging events to queue in-memory
	// before start dropping them (probably because logging server is down)
	MaxQueueSize = 1000
)

// CachingAuditLog implements events.IAuditLog on the recording machine (SSH server)
// It captures the local recording and forwards it to the AuditLog network server
type CachingAuditLog struct {
	server    events.IAuditLog
	queue     chan msg
	closeC    chan int
	closeOnce sync.Once
	client    protogen.AuditClient
}

// msg structure is used to transfer logging calls from the calling thread into
// asynchronous queue
type msg struct {
	eventType string
	fields    events.EventFields
	sid       session.ID
	namespace string
	reader    io.Reader
}

// MakeCachingAuditLog creaets a new & fully initialized instance of the alog
func MakeCachingAuditLog(logServer events.IAuditLog) *CachingAuditLog {

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})
	conn, err := grpc.Dial("127.0.0.1:3089", grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.WithTimeout(100*time.Second))
	if err != nil {
		panic(err)
	}

	client, err := api.NewClient(conn)
	if err != nil {
		panic(err)
	}

	ll := &CachingAuditLog{
		server: logServer,
		closeC: make(chan int),
		client: client,
	}
	// start the queue:
	if logServer != nil {
		ll.queue = make(chan msg, MaxQueueSize+1)
		go ll.run()
	}
	return ll
}

// run thread is picking up logging events and tries to forward them
// to the logging server
func (ll *CachingAuditLog) run() {
	clt, err := ll.client.SessionChunks(context.TODO())
	if err != nil {
		panic(err)
	}
	hist := hdrhistogram.New(1, 60000, 3)
	requests := 0
	lastReport := time.Now()
	for ll.server != nil {
		select {
		case <-ll.closeC:
			return
		case msg := <-ll.queue:
			if msg.fields != nil {
				err = ll.server.EmitAuditEvent(msg.eventType, msg.fields)
			} else if msg.reader != nil {
				bytes, err := ioutil.ReadAll(msg.reader)
				if err != nil {
					log.Warnf("%v", err)
				}
				start := time.Now()
				err = clt.Send(&protogen.SessionChunk{Namespace: msg.namespace, SessionID: string(msg.sid), Chunk: bytes})
				if err != nil {
					log.Warnf("%v", err)
				}
				requests += 1
				hist.RecordValue(int64(time.Now().Sub(start) / time.Microsecond))
				if time.Now().Sub(lastReport) > 10*time.Second {
					diff := time.Now().Sub(lastReport) / time.Second
					lastReport = time.Now()
					fmt.Printf("client histogram\n")
					for _, quantile := range []float64{25, 50, 75, 90, 95, 99, 100} {
						fmt.Printf("%v\t%v microseconds\n", quantile, hist.ValueAtQuantile(quantile))
					}
					fmt.Printf("%v requests/sec\n", requests/int(diff))

				}
				//err = ll.server.PostSessionChunk(msg.namespace, msg.sid, msg.reader)

			}
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (ll *CachingAuditLog) post(m msg) error {
	select {
	case ll.queue <- m:
	default:
		//log.Warnf("Audit log cannot keep up. Dropping event '%v' queue length %v", m.eventType, len(ll.queue))
	}
	return nil

}

func (ll *CachingAuditLog) Close() error {
	ll.closeOnce.Do(func() {
		close(ll.closeC)
	})
	return nil
}

func (ll *CachingAuditLog) EmitAuditEvent(eventType string, fields events.EventFields) error {
	return ll.post(msg{eventType: eventType, fields: fields})
}

func (ll *CachingAuditLog) PostSessionChunk(namespace string, sid session.ID, reader io.Reader) error {
	return ll.post(msg{sid: sid, reader: reader, namespace: namespace})
}

func (ll *CachingAuditLog) GetSessionChunk(string, session.ID, int, int) ([]byte, error) {
	return nil, errNotSupported
}
func (ll *CachingAuditLog) GetSessionEvents(string, session.ID, int) ([]events.EventFields, error) {
	return nil, errNotSupported
}
func (ll *CachingAuditLog) SearchEvents(time.Time, time.Time, string) ([]events.EventFields, error) {
	return nil, errNotSupported
}
