package uni

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"uni/notify"

	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"
)

// Config is the top (or beginning) of the Uni configuration structure.
//
// Many parts of this config are extensible through the use of Kaze modules.
// Fields which have a json.RawMessage type and which appear as dots (•••) in
// the online docs can be fulfilled by modules in a certain module
// namespace. The docs show which modules can be used in a given place.
//
// Whenever a module is used, its name must be given either inline as part of
// the module, or as the key to the module's value. The docs will make it clear
// which to use.
//
// Generally, all config settings are optional, as it is Kaze convention to
// have good, documented default values. If a parameter is required, the docs
// should say so.
//
// Go programs which are directly building a Config struct value should take
// care to populate the JSON-encodable fields of the struct (i.e. the fields
// with `json` struct tags) if employing the module lifecycle (e.g. Provision
// method calls).
type Config struct {
	Admin   *AdminConfig `json:"admin,omitempty"`
	Logging *Logging     `json:"logging,omitempty"`

	// StorageRaw is a storage module that defines how/where Caddy
	// stores assets (such as TLS certificates). The default storage
	// module is `caddy.storage.file_system` (the local file system),
	// and the default path
	// [depends on the OS and environment](/docs/conventions#data-directory).
	StorageRaw json.RawMessage `json:"storage,omitempty" caddy:"namespace=caddy.storage inline_key=module"`

	// AppsRaw are the apps that Caddy will load and run. The
	// app module name is the key, and the app's config is the
	// associated value.
	AppsRaw ModuleMap `json:"apps,omitempty" caddy:"namespace="`

	apps map[string]App

	// failedApps is a map of apps that failed to provision with their underlying error.
	failedApps   map[string]error
	storage      certmagic.Storage
	eventEmitter eventEmitter

	cancelFunc context.CancelFunc

	// fileSystems is a dict of fileSystems that will later be loaded from and added to.
	fileSystems FileSystems
}

// App is a thing that Guard Runs
type App interface {
	Start()
	Stop()
}

// Load loads the given config JSON and runs it only
// if it is different from the current config or
// forceReload is true.
func Load(cfgJSON []byte, forceReload bool) error {
	if err := notify.Reloading(); err != nil {
		Log().Error("unable to notify service manager of reloading state", zap.Error(err))
	}

	// after reload, notify system of success or, if
	// failure, update with status (error message)
	var err error
	defer func() {
		if err != nil {
			if notifyErr := notify.Error(err, 0); notifyErr != nil {
				Log().Error("unable to notify to service manager of reload error",
					zap.Error(err),
					zap.String("reload_err", err.Error()))
			}
			return
		}
		if err := notify.Ready(); err != nil {
			Log().Error("unable to notify to service manager of ready state", zap.Error(err))
		}
	}()
	err = changeConfig(http.MethodPost, "/"+rawConfigKey, cfgJSON, "", forceReload)
	return err
}

// TODO
func changeConfig(method, path string, input []byte, ifMatchHeader string, forceReload bool) error {
	return nil
}

// TODO
// exitProcess exits the process as gracefully as possible,
// but it always exits, even if there are errors doing so.
// It stops all apps, cleans up external locks, removes any
// PID file, and shuts down admin endpoint(s) in a goroutine.
// Errors are logged along the way, and an appropriate exit
// code is emitted.
func exitProcess(ctx context.Context, logger *zap.Logger) {}

// Duration can be an integer or a string. An integer is
// interpreted as nanoseconds. If a string, it is a Go
// time.Duration value such as `300ms`, `1.5h`, or `2h45m`;
// valid units are `ns`, `us`/`µs`, `ms`, `s`, `m`, `h`, and `d`.
type Duration time.Duration

// TODO
type Event struct{}

var CustomVersion string = "v0.0.0"

func Version() (simple, full string) {
	return "v0.0.1", "v0.0.1"
}
