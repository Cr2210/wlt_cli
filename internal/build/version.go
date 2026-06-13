// Package build holds version info injected via -ldflags at build time.
package build

// Version is overridden at build time via:
//
//	go build -ldflags "-X github.com/weiliantong/cli/internal/build.Version=$(VERSION)"
var Version = "dev"
