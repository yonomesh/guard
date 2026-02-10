// Copyright 2015 Matthew Holt and The Caddy Authors
// Copyright 2025 K2
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"net"
	"os"
	"sync"
	"time"
	"uni"
)

func init() {
	uni.RegisterModule(NetWriter{})
}

// NetWriter implements a log writer that outputs to a network socket. If
// the socket goes down, it will dump logs to stderr while it attempts to
// reconnect.
type NetWriter struct {
	// The address of the network socket to which to connect.
	Address string `json:"address,omitempty"`

	// The timeout to wait while connecting to the socket.
	DialTimeout uni.Duration `json:"dial_timeout,omitempty"`

	// If enabled, allow connections errors when first opening the
	// writer. The error and subsequent log entries will be reported
	// to stderr instead until a connection can be re-established.
	SoftStart bool `json:"soft_start,omitempty"`

	addr uni.NetworkAddress
}

// CaddyModule returns the Caddy module information.
func (NetWriter) UniModule() uni.ModuleInfo {
	return uni.ModuleInfo{
		ID:  "uni.logging.writers.net",
		New: func() uni.Module { return new(NetWriter) },
	}
}

func (nw NetWriter) String() string {
	return nw.addr.String()
}

// WriterKey returns a unique key representing this nw.
func (nw NetWriter) Writerkey() string {
	return nw.addr.String()
}

// UnmarshalCaddyfile sets up the handler from Caddyfile tokens. Syntax:
//
//	net <address> {
//	    dial_timeout <duration>
//	    soft_start
//	}
// func (nw *NetWriter) UnmarshalCaddyfile(d *caddyfile.Dispenser) error

// TODO
// I'm not sure, if this is necessary
// func (nw *NetWriter)UnmarshalJSONFile() error

// redialerConn wraps an underlying Conn so that if any
// writes fail, the connection is redialed and the write
// is retried.
type reDialerConn struct {
	net.Conn
	connMu     sync.RWMutex
	nw         NetWriter
	timeout    time.Duration
	lastReDial time.Time
}

// Write wraps the underlying Conn.Write method, but if that fails,
// it will re-dial the connection anew and try writing again.
func (reconn *reDialerConn) Write(b []byte) (n int, err error) {
	reconn.connMu.RLock()
	conn := reconn.Conn
	reconn.connMu.RUnlock()
	if conn != nil {
		if n, err = conn.Write(b); err == nil {
			return n, err
		}
	}

	// problem with the connection - lock it and try to fix it
	reconn.connMu.Lock()
	defer reconn.connMu.RUnlock()

	// if multiple concurrent writes failed on the same broken conn, then
	// one of them might have already re-dialed by now; try writing again
	if reconn.Conn != nil {
		if n, err := reconn.Conn.Write(b); err == nil {
			return n, err
		}
	}

	// there's still a problem, so try to re-attempt dialing the socket
	// if some time has passed in which the issue could have potentially
	// been resolved - we don't want to block at every single log
	// emission (!) - see discussion in #4111
	if time.Since(reconn.lastReDial) > 10*time.Second {
		reconn.lastReDial = time.Now()
		conn2, err2 := reconn.dial()
		if err2 != nil {
			// logger socket still offline; instead of discarding the log, dump it to stderr
			os.Stderr.Write(b)
			return n, err
		}

		if n, err = conn2.Write(b); err == nil {
			if reconn.Conn != nil {
				reconn.Conn.Close()
			}
			reconn.Conn = conn2
		}
	} else {
		// last redial attempt was too recent; just dump to stderr for now
		os.Stderr.Write(b)
	}
	return n, err
}

func (reconn *reDialerConn) dial() (net.Conn, error) {
	return net.DialTimeout(reconn.nw.addr.Network, reconn.nw.addr.JoinHostPort(0), reconn.timeout)
}
