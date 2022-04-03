package resource

import (
	"embed"
)

//go:embed public/*
var Public embed.FS

//go:embed font/*
var Font embed.FS
