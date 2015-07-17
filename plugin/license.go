package plugin

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

// Format: kismatic plugin license kubernetes-rbac ~/license.lic
//
func License(c *cli.Context) {
	// TODO: Load settings from config file
	// cfg, err := Load(c.GlobalString(`configuration`))
	// if err != nil {
	// 	log.Fatalln(`unable to load configuration file:`, err)
	// }

	if len(c.Args()) != 2 {
		fmt.Printf("Incorrect usage. Must have 2 arguments\n")
		return
	}

	pluginName := c.Args()[0]
	licenseFilePath := c.Args()[1]

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	destinationDir := filepath.Join(usr.HomeDir, ".kismatic")

	fmt.Println("Licensing plugin", pluginName)

	p := filepath.Join(destinationDir, fmt.Sprintf("license-%s.lic", pluginName))

	// make license dir if need be
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		if os.IsNotExist(err) {
			fmt.Println("Creating directory", destinationDir)
			os.MkdirAll(destinationDir, 0755)
		} else {
			// Some other error
			fmt.Println(err)
			// os.Exit(1)
		}
	}

	reader, err := os.Open(licenseFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	writer, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Installed license for", pluginName)

}
