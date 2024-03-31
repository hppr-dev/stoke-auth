package admin

import (
	"embed"
	"io/fs"
	"path"
)

//go:embed dist/*
var dist embed.FS

var Pages distFS = distFS{ dist }

type distFS struct {
	dist embed.FS
}

func (fs distFS) Open(name string) (fs.File, error) {
	return fs.dist.Open(path.Join("dist", name))
}
