package models

type ManifestObj struct {
	Manifest Manifest `json:"manifest"`
	Config string `json:"config"`
}

//omitempty
type Manifest struct {
	Config Layer `json:"config"`
	Layers []Layer `json:"layers"`
	MediaType string `json:"mediaType"`
	SchemaVersion int `json:"schemaVersion"`
}

type Layer struct {
	Digest string `json:"digest"`
	MediaType string `json:"mediaType"`
	Size int `json:"size"`
}

type ClairLayer struct {
	Name string
	Digest string
	ParentName string
}