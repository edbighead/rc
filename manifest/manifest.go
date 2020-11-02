package manifest

import "time"

// Manifest is used to map the response from Docker Registry API v2
type Manifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
}

// TagData represents only necessary fields from maniest
type TagData struct {
	Name      string
	Version   string
	CreatedAt time.Time
}

type V1Compatibility struct {
	Architecture  string    `json:"architecture"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	ID            string    `json:"id"`
	Os            string    `json:"os"`
}

type ImageData struct {
	Name    string
	Created time.Time
	Tag     string
}

type AllImages struct {
	Images []ImageData
}

func (images *AllImages) AddImage(image ImageData) {
	images.Images = append(images.Images, image)
}
