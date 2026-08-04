// Package api implements the HTTP API layer generated from the OpenAPI spec
// (backend/api/openapi/v1.yaml).
package api

import (
	"net/http"

	rootapi "vedo-edutrack/backend/api"
)

// OpenAPISpecHandler serves the raw OpenAPI spec as YAML so docs tools can
// consume it (Swagger UI loads it from /api/v1/openapi.json).
func OpenAPISpecHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(rootapi.Spec)
}

// DocsHandler serves a minimal Swagger UI page loading the spec from
// /api/v1/openapi.json (development convenience; the generated chi router
// does not serve docs by default).
func DocsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>VEDO EduTrack API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: '/api/v1/openapi.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.presets.AuthorizationPlugin],
      });
    };
  </script>
</body>
</html>`))
}
