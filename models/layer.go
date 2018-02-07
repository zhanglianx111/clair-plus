package models

import (
	"github.com/coreos/clair/api/v1"
	"github.com/coreos/clair/utils/types"
)

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

type Vulner struct {
	ImageName string `json:"ImageName"`
	Vulners []V `json:"Vulners"`
}

type V struct {
	Name 	       string  `json:"Name"`
	Description    string  `json:"Description"`
	Package        Package `json:"Package"`
	FixedByVersion string  `json:"FixedByVersion"`
	Link           string  `json:"Link"`
	Layer          string  `json:"Layer"`
	Severity       string  `json:"Severity"`
}

type Package struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

type VulnerabilityInfo struct {
	Vulnerability v1.Vulnerability
	Feature       v1.Feature
	Severity      types.Priority
}