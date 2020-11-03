package manifest

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

func ImageInRepo(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RegistryCall(user, password, url, method, reqHeader string) (respBody []byte, header string, statusCode int) {
	var username string = user
	var passwd string = password
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(username, passwd)
	req.Header.Add("Accept", reqHeader)
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	return bodyText, resp.Header.Get("docker-content-digest"), resp.StatusCode
}
