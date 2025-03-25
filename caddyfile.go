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
	"fmt"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("dynamic_records", parseApp)
}

// parseAppGlobal configures the "dynamic_records" global option from Caddyfile.
// Syntax:
//
//	dynamic_records {
//		<hostname>, <hostname>, ... {
//			provider <name> ...
//			ip_source upnp|simple_http <endpoint>
//			versions ipv4|ipv6
//			check_interval <duration>
//			ttl <duration>
//		}
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
		fmt.Println(d.Val())

		recordType := d.Val()
		modID := "dynamic_records." + recordType
		unm, err := caddyfile.UnmarshalModule(d, modID)
		if err != nil {
			return nil, err
		}
		app.RecordsRaw = caddyconfig.JSONModuleObject(unm, "type", recordType, nil)
	}

	fmt.Println(caddyconfig.JSON(app, nil))

	return httpcaddyfile.App{
		Name:  "dynamic_records",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}
