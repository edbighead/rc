package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/edbighead/rc/manifest"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	userImage, keepCount string
)

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://registry.private/"
		username := "user"     // anonymous
		password := "password" // anonymous

		hub, err := registry.New(url, username, password)

		if err != nil {
			exit(err)
		}

		repositories, err := hub.Repositories()
		if err != nil {
			exit(err)
		}

		if userImage == "" || keepCount == "" {
			exit("image flag and keep count is required")
		}

		if !imageInRepo(userImage, repositories) {
			exit("No such image")
		}

		tags, err := hub.Tags(userImage)

		if err != nil {
			panic(err)
		}
		images := manifest.AllImages{}

		for _, t := range tags {
			if err != nil {
				panic(err)
			}
			url := fmt.Sprintf("%s/v2/%s/manifests/%s", url, userImage, t)
			m := basicAuth(username, password, url)

			manifests := manifest.Manifest{}
			v1Compatibility := manifest.V1Compatibility{}

			json.Unmarshal(m, &manifests)
			json.Unmarshal([]byte(manifests.History[0].V1Compatibility), &v1Compatibility)

			image := manifest.ImageData{
				Name:    userImage,
				Created: v1Compatibility.Created,
				Tag:     manifests.Tag,
			}
			images.AddImage(image)
		}

		// sort array by time
		sort.Slice(images.Images, func(i, j int) bool {
			return images.Images[i].Created.After(images.Images[j].Created)
		})

		k, err := strconv.Atoi(keepCount)
		if err == nil {
			if k < len(tags) {
				fmt.Println("Some images will be deleted")
			}
		}

		data := [][]string{}

		for _, i := range images.Images {
			data = append(data, []string{i.Name, i.Created.Format("2 Jan 2006 15:04:05"), i.Tag, ""})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Image", "Date", "Tag", "Delete"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render()

	},
}

func imageInRepo(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func exit(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func basicAuth(user, password, url string) []byte {
	var username string = user
	var passwd string = password
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	return bodyText
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().StringVarP(&userImage, "image", "i", "", "image name to run cleanup")
	cleanupCmd.Flags().StringVarP(&keepCount, "keep", "k", "", "number of images to keep")
}
