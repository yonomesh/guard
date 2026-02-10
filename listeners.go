// Copyright 2015 Matthew Holt and The Caddy Authors
// Copyright 2025 k2 <skrik2@outlook.com>
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

package uni

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// listenFdsStart is the first file descriptor number for systemd socket activation.
// File descriptors 0, 1, 2 are reserved for stdin, stdout, stderr.
const listenFdsStart = 3

// NetworkAddress represents one or more network addresses.
// It contains the individual components for a parsed network
// address of the form accepted by ParseNetworkAddress().
type NetworkAddress struct {
	// Should be a network value accepted by Go's net package or
	// by a plugin providing a listener for that network type.
	Network string

	// The "main" part of the network address is the host, which
	// often takes the form of a hostname, DNS name, IP address,
	// or socket path.
	Host string

	// For addresses that contain a port, ranges are given by
	// [StartPort, EndPort]; i.e. for a single port, StartPort
	// and EndPort are the same. For no port, they are 0.
	StartPort uint
	EndPort   uint
}

func (na NetworkAddress) port() string {
	if na.StartPort == na.EndPort {
		return strconv.FormatUint(uint64(na.StartPort), 10)
	}
	return fmt.Sprintf("%d-%d", na.StartPort, na.EndPort)
}

// String reconstructs the address string for human display.
// The output can be parsed by ParseNetworkAddress(). If the
// address is a unix socket, any non-zero port will be dropped.
func (na NetworkAddress) String() string {
	if na.Network == "tcp" && (na.Host != "" || na.port() != "") {
		na.Network = "" // omit default network value for brevity
	}
	return JoinNetworkAddress(na.Network, na.Host, na.port())
}

// IsUnixNetwork returns true if na.Network is
// unix, unixgram, or unixpacket.
func (na NetworkAddress) IsUnixNetwork() bool {
	return IsUnixNetwork(na.Network)
}

// IsFdNetwork returns true if na.Network is
// fd or fdgram.
func (na NetworkAddress) IsFdNetwork() bool {
	return IsFdNetwork(na.Network)
}

// JoinHostPort is like net.JoinHostPort, but where the port
// is StartPort + offset.
func (na NetworkAddress) JoinHostPort(offset uint) string {
	if na.IsUnixNetwork() || na.IsFdNetwork() {
		return na.Host
	}
	return net.JoinHostPort(na.Host, strconv.FormatUint(uint64(na.StartPort+offset), 10))
}

// JoinNetworkAddress combines network, host, and port into a single
// address string of the form accepted by ParseNetworkAddress(). For
// unix sockets, the network should be "unix" (or "unixgram" or
// "unixpacket") and the path to the socket should be given as the
// host parameter.
func JoinNetworkAddress(network, host, port string) string {
	var a string
	if network != "" {
		a = network + "/"
	}
	if (host != "" && port == "") || IsUnixNetwork(network) || IsFdNetwork(network) {
		a += host
	} else if port != "" {
		a += net.JoinHostPort(host, port)
	}
	return a
}

// IsUnixNetwork returns true if the netw is a unix network.
func IsUnixNetwork(netw string) bool {
	return strings.HasPrefix(netw, "unix")
}

// IsFdNetwork returns true if the netw is a fd network.
func IsFdNetwork(netw string) bool {
	return strings.HasPrefix(netw, "fd")
}
