/*
Copyright 2016 SPIFFE Authors

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

package api

import (
	"github.com/gravitational/teleport/lib/auth/api/protogen"

	"github.com/gravitational/trace"
	"google.golang.org/grpc"
)

func NewClient(conn *grpc.ClientConn) (*Client, error) {
	if conn == nil {
		return nil, trace.BadParameter("missing parameter conn")
	}
	return &Client{AuditClient: protogen.NewAuditClient(conn), conn: conn}, nil
}

// Client is GRPC based Workload service client
type Client struct {
	protogen.AuditClient
	conn *grpc.ClientConn
}

// Close closes underlying connection
func (c *Client) Close() error {
	return c.conn.Close()
}
