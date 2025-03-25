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

package records

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
	"go.uber.org/zap"

	dynamicdns "github.com/mholt/caddy-dynamicdns"
)

func init() {
	caddy.RegisterModule(Address{})
}

type Address struct {
	Config BaseRecord `json:"config"`

	Hostname string `json:"hostname"`

	// The sources from which to get the server's public IP address.
	// Multiple sources can be specified for redundancy.
	// Default: simple_http
	IPSourcesRaw []json.RawMessage `json:"ip_sources" caddy:"namespace=dynamic_dns.ip_sources inline_key=source"`

	// The IP versions to enable. By default, both "ipv4" and "ipv6" will be enabled.
	// To disable IPv6, specify {"ipv6": false}.
	Versions dynamicdns.IPVersions `json:"versions"`

	ipSources []dynamicdns.IPSource

	ctx    caddy.Context
	logger *zap.Logger
}

func (Address) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "managed_dns.record.address",
		New: func() caddy.Module { return new(Address) },
	}
}

// func (s *Address) Provision(ctx caddy.Context) error {
// 	s.logger = ctx.Logger(s)
// 	s.logger.Warn("No static IPs configured")

// 	return nil
// }

func (s *Address) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // skip directive name

	if !d.NextArg() {
		return d.ArgErr()
	}
	s.Hostname = d.Val()

	// handle the block
	for d.NextBlock(0) {
		switch d.Val() {
		case "ip_source":
			if !d.NextArg() {
				return d.ArgErr()
			}
			sourceType := d.Val()
			modID := "dynamic_dns.ip_sources." + sourceType
			unm, err := caddyfile.UnmarshalModule(d, modID)
			if err != nil {
				return err
			}
			s.IPSourcesRaw = append(s.IPSourcesRaw, caddyconfig.JSONModuleObject(unm, "source", sourceType, nil))

		case "versions":
			args := d.RemainingArgs()
			if len(args) == 0 {
				return d.Errf("Must specify at least one version")
			}

			// Set up defaults; if versions are specified,
			// both versions start as false, then flipped
			// to true otherwise.
			falseBool := false
			s.Versions = dynamicdns.IPVersions{
				IPv4: &falseBool,
				IPv6: &falseBool,
			}

			trueBool := true
			for _, arg := range args {
				switch arg {
				case "ipv4":
					s.Versions.IPv4 = &trueBool
				case "ipv6":
					s.Versions.IPv6 = &trueBool
				default:
					return d.Errf("Unsupported version: '%s'", arg)
				}
			}

		default:
			err := s.Config.UnmarshalDirective(d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Address) Provision(ctx caddy.Context) error {
	// set up the DNS provider module
	fmt.Println("Test")

	r := a.Config

	if len(r.DNSProviderRaw) == 0 {
		return fmt.Errorf("a DNS provider is required")
	}
	val, err := ctx.LoadModule(r, "DNSProviderRaw")
	fmt.Println("test")
	if err != nil {
		return fmt.Errorf("loading DNS provider module: %v", err)
	}
	r.dnsProvider = val

	fmt.Println("test")

	recordGetter := r.dnsProvider.(libdns.RecordSetter)
	fmt.Println("Test 2", recordGetter)
	// fmt.Println(recordGetter.GetRecords(ctx, "seblab.net"))

	return nil
}

func (a Address) GetConfig() BaseRecord {
	return a.Config
}

// Interface guards
var (
	_ Record = (*Address)(nil)
	// _ caddy.Provisioner     = (*Address)(nil)
	_ caddyfile.Unmarshaler = (*Address)(nil)
)
