package plugin

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Install(c *cli.Context) {
	// TODO: Load settings from config file
	// cfg, err := Load(c.GlobalString(`configuration`))
	// if err != nil {
	// 	log.Fatalln(`unable to load configuration file:`, err)
	// }

	if len(c.Args()) != 1 {
		fmt.Printf("Incorrect usage. Must have 1 plugin argument\n")
		return
	}

	version := c.String("version")

	pluginName := c.Args()[0]

	// TODO: Make destinationDir read from a config value
	destinationDir := "./"

	owner := "kismatic"
	repo := ""

	switch pluginName {
	case "ldap":
		repo = "kubernetes-ldap"
		// TODO: Remove test case
		// case "test":
		// 	owner = "bcbroussard"
		// 	repo = "rkt"
	}

	if repo == "" {
		log.Fatalf("Plugin not found")
	}

	fmt.Println("installing plugin", pluginName)

	client := github.NewClient(nil)

	if version == "" {
		release, _, err := client.Repositories.GetLatestRelease(owner, repo)
		if err != nil {
			log.Fatalf("Could not download plugin")
		}

		saveBinaryFromRelease(release, destinationDir)
	} else {

		fmt.Printf("using version", version)
		release, _, err := client.Repositories.GetReleaseByTag(owner, repo, version)
		if err != nil {
			log.Fatalf("Could not download plugin")
		}

		saveBinaryFromRelease(release, destinationDir)
	}
}

func saveBinaryFromRelease(release *github.RepositoryRelease, destinationDir string) {
	if len(release.Assets) == 0 {
		log.Fatalf("no releases found")
	}

	url := release.Assets[0].BrowserDownloadURL

	fmt.Println("Getting url ", *url)
	response, err := http.Get(*url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	zipFile, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	unzip(zipFile, destinationDir)
}

func unzip(zipFile []byte, destinationDir string) {

	reader, err := zip.NewReader(
		bytes.NewReader(zipFile), int64(len(zipFile)))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// defer reader.Close()

	for _, f := range reader.File {

		zipped, err := f.Open()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer zipped.Close()

		path := filepath.Join(destinationDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			fmt.Println("Creating directory", path)
		} else {
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			defer writer.Close()

			if _, err = io.Copy(writer, zipped); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Decompressing : ", path)
		}

	}

}
