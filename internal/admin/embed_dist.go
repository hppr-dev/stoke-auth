package admin

import (
	"embed"
	"io/fs"
	"path"
	"strings"
)

//go:embed dist/*
var dist embed.FS

var Pages distFS = distFS{ dist }

type distFS struct {
	dist embed.FS
}

func (fs distFS) Open(name string) (fs.File, error) {
	if
		strings.HasSuffix(name, "user") || 
		strings.HasSuffix(name, "group") || 
		strings.HasSuffix(name, "claim") || 
		strings.HasSuffix(name, "key") || 
		strings.HasSuffix(name, "monitor") {
		return fs.dist.Open(path.Join("dist", "index.html"))
	}
	return fs.dist.Open(path.Join("dist", name))
}
