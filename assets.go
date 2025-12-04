package assets

import (
    "embed"
)

// FS contains embedded templates and static assets.
//go:embed templates/*.html web/static/*.css static/*
var FS embed.FS
