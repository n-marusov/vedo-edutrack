// Package apiemb embeds the OpenAPI specification so it can be served by the
// HTTP layer (docs endpoint). The spec lives at backend/api/openapi/v1.yaml —
// one level above internal/ — so a dedicated package at backend/api/ is the
// only go:embed-able location.
package apiemb

import _ "embed"

//go:embed openapi/v1.yaml
var Spec []byte
