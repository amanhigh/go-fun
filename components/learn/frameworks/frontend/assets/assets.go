package assets

import "embed"

// FS contains the frontend demo static assets.
//
//go:embed css/* js/* images/*
var FS embed.FS
