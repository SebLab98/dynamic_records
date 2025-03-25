// Copyright 2020 Matthew Holt
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dynamicrecords

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
	"go.uber.org/zap"

	dynamicdns "github.com/mholt/caddy-dynamicdns"
)

func init() {
	caddy.RegisterModule(Address{})
}

// IPSource is a type that can get IP addresses.
type Record interface {
	// GetIPs(context.Context, IPVersions) ([]net.IP, error)
}

type Address struct {
	// The sources from which to get the server's public IP address.
	// Multiple sources can be specified for redundancy.
	// Default: simple_http
	IPSourcesRaw []json.RawMessage `json:"ip_sources,omitempty" caddy:"namespace=dynamic_dns.ip_sources inline_key=source"`

	// The IP versions to enable. By default, both "ipv4" and "ipv6" will be enabled.
	// To disable IPv6, specify {"ipv6": false}.
	Versions dynamicdns.IPVersions `json:"versions,omitempty"`

	Hostname string `json:"hostname,omitempty"`

	ipSources   []dynamicdns.IPSource
	dnsProvider libdns.RecordSetter

	ctx    caddy.Context
	logger *zap.Logger
}

func (Address) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dynamic_records.address",
		New: func() caddy.Module { return new(Address) },
	}
}

func (s *Address) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s)

	s.logger.Warn("No static IPs configured")

	return nil
}

func (s *Address) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // skip directive name

	for d.NextArg() {
		fmt.Println("Address", d.Val())

		s.Hostname = d.Val()
	}

	return nil
}

// Interface guards
var (
	_ Record                = (*Address)(nil)
	_ caddy.Provisioner     = (*Address)(nil)
	_ caddyfile.Unmarshaler = (*Address)(nil)
)
