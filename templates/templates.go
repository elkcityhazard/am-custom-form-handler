package templates

import "embed"

//go:embed templates
var templatesFS embed.FS

// GetTemplates returns the embedded templates
func GetTemplates() *embed.FS {
	return &templatesFS
}
