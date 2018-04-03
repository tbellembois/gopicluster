package handlers

import (
	"html/template"
)

// Env is a structure used to pass objects throughout the application.
type Env struct {
	Templates map[string]*template.Template
	Address   string
	Port      string
	DryRun    bool
}
