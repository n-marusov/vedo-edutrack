package mockhub

import (
	_ "embed"
)

//go:embed schema.graphql
var embeddedSchemaSDL string
