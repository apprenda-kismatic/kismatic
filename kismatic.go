// Copyright 2015 Kismatic Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package kismatic contains common functions for Kismatic clusters.
package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/kismatic/kismatic/plugin"
)

func main() {
	app := cli.NewApp()
	app.Name = "kismatic"
	app.Usage = "Install and manage Kismatic plugins for your Kubernetes cluster"
	app.Action = func(c *cli.Context) {
		// run help
		println("Run 'kismatic help' for help")
	}

	// Global flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "configuration, c",
			Value: "",
			Usage: "path to config file",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "plugin",
			Aliases: []string{"plugins"},
			Usage:   "Installs and manages Kismatic plugins",
			Subcommands: []cli.Command{
				{
					Name:  "install",
					Usage: "installs a Kismatc plugin",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "version",
							Value: "",
							Usage: "plugin version to install",
						},
					},
					Action: func(c *cli.Context) {
						plugin.Install(c)
					},
				},
				{
					Name:  "license",
					Usage: "licenses a Kismatc plugin",
					Action: func(c *cli.Context) {
						plugin.License(c)
					},
				},
			},
		},
		{
			Name:  "login",
			Usage: "Login",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username, u",
					Value: "",
					Usage: "Username for authenticating",
				},
				cli.StringFlag{
					Name:  "password, p",
					Value: "",
					Usage: "Password for authenticating",
				},
			},
			Action: func(c *cli.Context) {
				plugin.Login(c)
			},
		},
	}

	app.Run(os.Args)
}
