package assets

import "embed"

// FS contains the Kohan static assets.
//
//go:embed css/* js/*
var FS embed.FS
