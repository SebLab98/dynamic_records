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

package managed_dns

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("managed_dns", parseApp)
}

// parseApp configures the "managed_dns" global option from Caddyfile.
// Syntax:
//
//	managed_dns {
//		<record_type> <zone> <name> {
//			provider <name> ...
//			check_interval <duration>
//			ttl <duration>
//			...
//		}
//		...
//	}
//
// If <names...> are omitted after <zone>, then "@" will be assumed.
func parseApp(d *caddyfile.Dispenser, _ interface{}) (interface{}, error) {
	app := new(App)

	// consume the option name
	if !d.Next() {
		return nil, d.ArgErr()
	}

	// handle the block
	for d.NextBlock(0) {
		recordType := d.Val()
		modID := "managed_dns.record." + recordType
		unm, err := caddyfile.UnmarshalModule(d, modID)
		if err != nil {
			return nil, err
		}
		app.RecordsRaw = append(app.RecordsRaw, caddyconfig.JSONModuleObject(unm, "type", recordType, nil))
	}

	return httpcaddyfile.App{
		Name:  "managed_dns",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}
