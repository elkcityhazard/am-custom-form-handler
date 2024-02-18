package models

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	DataMap   map[string]interface{}
	CSRFToken string
	Flash     string
	Error     string
}
