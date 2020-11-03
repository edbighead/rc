package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/edbighead/rc/auth"
	"github.com/edbighead/rc/manifest"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	userImage, keepCount, cfgFile string
	dryRun                        bool

	// cleanupCmd represents the cleanup command
	cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "Remove old images from repository",
		Long: `
Remove images from repository, based on creation date, keeping
only a certain number of images`,
		Run: func(cmd *cobra.Command, args []string) {
			url, username, password := auth.Init(cfgFile)

			hub, err := registry.New(url, username, password)

			if err != nil {
				exit(err)
			}

			repositories, err := hub.Repositories()
			if err != nil {
				exit(err)
			}

			if userImage == "" || keepCount == "" {
				exit("image and keep count flags are required")
			}

			if !manifest.ImageInRepo(userImage, repositories) {
				exit("No such image")
			}

			tags, err := hub.Tags(userImage)

			if err != nil {
				exit(err)
			}

			images := manifest.AllImages{}

			for _, t := range tags {
				if err != nil {
					panic(err)
				}
				url := fmt.Sprintf("%s/v2/%s/manifests/%s", url, userImage, t)
				body, _, _ := manifest.RegistryCall(username, password, url, "GET", "application/json")

				manifests := manifest.Manifest{}
				v1Compatibility := manifest.V1Compatibility{}

				json.Unmarshal(body, &manifests)
				json.Unmarshal([]byte(manifests.History[0].V1Compatibility), &v1Compatibility)

				image := manifest.ImageData{
					Name:    userImage,
					Created: v1Compatibility.Created,
					Tag:     manifests.Tag,
				}
				images.AddImage(image)
			}

			// sort images by time
			sort.Slice(images.Images, func(i, j int) bool {
				return images.Images[i].Created.After(images.Images[j].Created)
			})

			k, err := strconv.Atoi(keepCount)

			if err != nil {
				exit("number of images to keep should be an integer")
			}

			if k < len(tags) {
				log.Println("Some images will be deleted")
			} else {
				log.Println("No images will be deleted")
			}

			data := [][]string{}
			deleteTags := []string{}

			count := 1
			for _, i := range images.Images {
				deleteMark := "no"
				if count > k {
					deleteMark = "yes"
					deleteTags = append(deleteTags, i.Tag)
				}
				data = append(data, []string{strconv.Itoa(count), i.Name, i.Created.Format("2 Jan 2006 15:04:05"), i.Tag, deleteMark})
				count++
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"№", "Image", "Created", "Tag", "Delete"})

			for _, v := range data {
				table.Append(v)
			}
			table.Render()

			if !dryRun {
				for _, dt := range deleteTags {
					fullURL := fmt.Sprintf("%s/v2/%s/manifests/%s", url, userImage, dt)
					_, sha, _ := manifest.RegistryCall(username, password, fullURL, "HEAD", "application/vnd.docker.distribution.manifest.v2+json")
					deleteUrlfmt := fmt.Sprintf("%s/v2/%s/manifests/%s", url, userImage, sha)
					_, _, status := manifest.RegistryCall(username, password, deleteUrlfmt, "DELETE", "application/vnd.docker.distribution.manifest.v2+json")
					if status == 202 {
						log.Printf("%s:%s successfully deleted!\n", userImage, dt)
					}

				}
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rc.yaml)")
	cleanupCmd.Flags().StringVarP(&userImage, "image", "i", "", "image name")
	cleanupCmd.Flags().StringVarP(&keepCount, "keep", "k", "", "number of image tags to keep")
	cleanupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "only output image tags")
}

func exit(msg interface{}) {
	log.Fatal(msg)
	os.Exit(1)
}
