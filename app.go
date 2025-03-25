package dynamicrecords

import (

	// "time"

	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"

	// "github.com/libdns/libdns"
	"github.com/caddy-dns/cloudflare"
)

func init() {
	caddy.RegisterModule(App{})
	// httpcaddyfile.RegisterGlobalOption("dynamic_records", parseApp)
}

// App provides a range of IP address prefixes (CIDRs) retrieved from cloudflare.
type App struct {
	Records []Record

	RecordsRaw json.RawMessage `json:"record,omitempty" caddy:"namespace=dynamic_records inline_key=type"`

	// The configuration for the DNS provider with which the DNS
	// records will be updated.
	DNSProviderRaw json.RawMessage `json:"dns_provider,omitempty" caddy:"namespace=dns.providers inline_key=name"`

	// How frequently to check the public IP address. Default: 30m
	CheckInterval caddy.Duration `json:"check_interval,omitempty"`

	// The TTL to set on DNS records.
	TTL caddy.Duration `json:"ttl,omitempty"`

	ctx    caddy.Context
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dynamic_records",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the app module.
func (a *App) Provision(ctx caddy.Context) error {
	a.ctx = ctx
	a.logger = ctx.Logger(a)

	return nil
}

// Start starts the app module.
func (a App) Start() error {
	provider := new(cloudflare.Provider)
	fmt.Println(provider.APIToken)

	return nil
}

// Stop stops the app module.
func (a App) Stop() error {
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
