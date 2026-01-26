package templates

import (
	"embed"
	"html/template"
	"sync"
)

//go:embed files/*
var templateFS embed.FS

var (
	templateCache = make(map[string]*template.Template)
	cacheMutex    sync.RWMutex
)
