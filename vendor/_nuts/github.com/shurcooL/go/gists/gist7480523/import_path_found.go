package gist7480523

import (
	"path/filepath"
)

// An ImportPathFound describes the Import Path found in a GOPATH workspace.
type ImportPathFound struct {
	importPath  string
	gopathEntry string
}

func NewImportPathFound(importPath, gopathEntry string) ImportPathFound {
	return ImportPathFound{
		importPath:  importPath,
		gopathEntry: gopathEntry,
	}
}

func (w *ImportPathFound) ImportPath() string {
	return w.importPath
}

func (w *ImportPathFound) GopathEntry() string {
	return w.gopathEntry
}

func (w *ImportPathFound) FullPath() string {
	return filepath.Join(w.gopathEntry, "src", w.importPath)
}
