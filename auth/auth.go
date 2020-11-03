package auth

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"
)

var cfgFile string

type config struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *config) getConf(path string) *config {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		if strings.Contains(path, ".rc.yaml") {
			fmt.Printf("Warning: %s/.rc.yaml not found, checking environment variables\n", usr.HomeDir)
		} else {
			log.Fatal(err)
		}
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		exit(err)
	}

	return c
}

func Init(cfgFile string) (url, username, password string) {

	var file string
	var c config

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// If config file is defined, it takes priority
	if cfgFile != "" {
		file = cfgFile
		c.getConf(file)
		if credentialsPresent(c.URL, c.Username, c.Password) {
			return c.URL, c.Username, c.Password
		}

		exit("Please check your credentials file")

	}

	// Seecond priority are env variables
	if credentialsPresent(os.Getenv("RC_URL"), os.Getenv("RC_USERNAME"), os.Getenv("RC_PASSWORD")) {
		return os.Getenv("RC_URL"), os.Getenv("RC_USERNAME"), os.Getenv("RC_PASSWORD")
	}

	// Third priority is config file
	file = fmt.Sprintf("%s/.rc.yaml", usr.HomeDir)
	c.getConf(file)

	if credentialsPresent(c.URL, c.Username, c.Password) {
		return c.URL, c.Username, c.Password
	}

	exit("Credentials not found, please check help command")

	return url, username, password
}

func exit(msg interface{}) {
	log.Fatal(msg)
	os.Exit(1)
}

func credentialsPresent(u, n, p string) bool {
	if len(u) == 0 || len(n) == 0 || len(p) == 0 {
		return false
	}
	return true
}
