package managed_dns

import (

	// "time"

	"encoding/json"
	"fmt"
	"time"

	"github.com/SebLab98/dynamic_records/records"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(App{})
}

// App provides a range of IP address prefixes (CIDRs) retrieved from cloudflare.
type App struct {
	RecordsRaw []json.RawMessage `json:"records" caddy:"namespace=managed_dns.record inline_key=type"`

	records []records.Record

	ctx    caddy.Context
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "managed_dns",
		New: func() caddy.Module { return new(App) },
	}
}

// Provision sets up the app module.
func (a *App) Provision(ctx caddy.Context) error {
	a.ctx = ctx
	a.logger = ctx.Logger(a)

	// set up the Records to manage
	if a.RecordsRaw != nil {
		vals, err := ctx.LoadModule(a, "RecordsRaw")
		if err != nil {
			return fmt.Errorf("loading Records module: %v", err)
		}
		for _, val := range vals.([]interface{}) {
			a.records = append(a.records, val.(records.Record))
		}
	}

	return nil
}

// Start starts the app module.
func (a App) Start() error {
	for _, r := range a.records {
		go checkerLoop(a.ctx, r)
	}

	return nil
}

// Stop stops the app module.
func (a App) Stop() error {
	return nil
}

// checkerLoop checks and updates records at every check
// interval. It stops when ctx is cancelled.
func checkerLoop(ctx caddy.Context, r records.Record) {
	ticker := time.NewTicker(time.Duration(r.GetConfig().CheckInterval))
	defer ticker.Stop()

	// a.checkIPAndUpdateDNS()

	for {
		select {
		case <-ticker.C:
			// a.checkIPAndUpdateDNS()
		case <-ctx.Done():
			return
		}
	}
}

// Interface guards
var (
	_ caddy.Provisioner = (*App)(nil)
	_ caddy.App         = (*App)(nil)
)
