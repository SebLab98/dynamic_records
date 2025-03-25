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

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
)

func init() {
}

type Record interface {
	// GetIPs(context.Context, IPVersions) ([]net.IP, error)
	GetConfig() BaseRecord
	// GetValue() string
	// GetCheckInterval() caddy.Duration
	// CreateRecord()
	// UpdateRecord()
	// DeleteRecord()
}

type BaseRecord struct {
	// The configuration for the DNS provider with which the DNS records will be updated.
	DNSProviderRaw json.RawMessage `json:"dns_provider" caddy:"namespace=dns.providers inline_key=name"`

	// How frequently to check the public IP address. Default: 30m
	CheckInterval caddy.Duration `json:"check_interval"`

	// The TTL to set on DNS records.
	TTL caddy.Duration `json:"ttl"`

	// Whether to leave record after app stops or not
	DeleteAfterUse bool `json:"delete_after_use"`

	dnsProvider any
	record      libdns.Record
}

func (r *BaseRecord) UnmarshalDirective(d *caddyfile.Dispenser) error {
	switch d.Val() {
	case "check_interval":
		if !d.NextArg() {
			return d.ArgErr()
		}
		dur, err := caddy.ParseDuration(d.Val())
		if err != nil {
			return err
		}
		r.CheckInterval = caddy.Duration(dur)

	case "provider":
		if !d.NextArg() {
			return d.ArgErr()
		}
		provName := d.Val()
		modID := "dns.providers." + provName
		unm, err := caddyfile.UnmarshalModule(d, modID)
		if err != nil {
			return err
		}
		r.DNSProviderRaw = caddyconfig.JSONModuleObject(unm, "name", provName, nil)

	case "ttl":
		if !d.NextArg() {
			return d.ArgErr()
		}
		dur, err := caddy.ParseDuration(d.Val())
		if err != nil {
			return err
		}
		r.TTL = caddy.Duration(dur)

	default:
		return d.ArgErr()
	}

	return nil
}
