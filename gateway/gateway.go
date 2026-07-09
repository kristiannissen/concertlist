// Package gateway is a thin, non-internal seam between Vercel's Go
// serverless-function entry point (api/index.go) and the app's internal/
// packages.
//
// Vercel's Go builder compiles api/*.go by generating a wrapper file and
// building it as an ad-hoc "command-line-arguments" package with a
// synthetic import path (it shows up in build errors as "handler/api"
// rather than the real module path). Because that synthetic path isn't
// recognized as being rooted at this module, Go's internal-package
// visibility rule rejects any import of internal/... made directly from
// api/index.go, even though the code lives in the same repo.
//
// Routing api/index.go's dependency through this ordinary, non-internal
// package avoids that: gateway is a real package resolved normally under
// github.com/kristiannissen/concertlist, so its own import of internal/...
// is unaffected and type-checks fine. api/index.go then only needs to
// import gateway, never internal/... directly.
package gateway

import (
	"net/http"

	"github.com/kristiannissen/concertlist/internal/adapters"
)

// NewRouter returns the application's HTTP router.
func NewRouter() *http.ServeMux {
	return adapters.NewRouter()
}
