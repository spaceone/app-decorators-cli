package main

import (
	"fmt"
	"github.com/serkansipahi/app-decorators-cli/helper"
	"github.com/serkansipahi/app-decorators-cli/install"
	"github.com/serkansipahi/app-decorators-cli/options"
	utilOs "github.com/serkansipahi/app-decorators-cli/util/os"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	CLI_NAME     = "appdec"
	AUTHOR_NAME  = "Serkan Sipahi"
	AUTHOR_EMAIL = "serkan.sipahi@yahoo.de"
	APP_VERSION  = "0.8.221"
	COPYRIGHT    = "(c) 2017"
)

func main() {

	_, err := Init(utilOs.Os{})
	if err != nil {
		log.Fatalln("Failed while initializing...", err)
	}

	rootPath, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed while getting root path")
	}
	app := cli.NewApp()
	app.Name = CLI_NAME
	app.Usage = "command line tool"
	app.Version = APP_VERSION
	app.Copyright = COPYRIGHT + " " + CLI_NAME
	app.Authors = []cli.Author{
		cli.Author{
			Name:  AUTHOR_NAME,
			Email: AUTHOR_EMAIL,
		},
	}

	/**
	 * Setting up allowed commands
	 */
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create new component",
			Flags: []cli.Flag{
				options.Name,
				options.Debug,
			},
			Action: func(c *cli.Context) error {

				// assign passed args
				name := strings.ToLower(c.String("name"))
				debug := c.Bool("debug")

				// create installer
				installer := install.New(
					name,
					rootPath,
					APP_VERSION,
					CLI_NAME,
					debug,
				)

				// start installer
				if err := installer.Run(); err != nil {
					log.Fatalln("Failed while installing...", err)
				}

				return nil
			},
		},
		/*
			{
				Name:    "recreate",
				Aliases: []string{"rc"},
				Usage:   "reinit an existing component (without deleting ./src files)",
				Flags: []cli.Flag{
					options.Name,
					options.Debug,
				},
				Action: func(c *cli.Context) error {

					return nil
				},
			},
		*/
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "starting workflow, see: run --help",
			Flags: []cli.Flag{
				options.Name,
				options.Watch,
				options.Server,
				options.Production,
				options.Minify,
				//options.Format,
				//options.Port,
				//options.Browser,
			},
			Action: func(c *cli.Context) error {

				var (
					name       string    = strings.ToLower(c.String("name"))
					watch      bool      = c.Bool("watch")
					format     string    = "default"
					server     bool      = c.Bool("server")
					production bool      = c.Bool("production")
					minify     bool      = c.Bool("minify")
					ch         chan bool = make(chan bool)
					//format   = c.String("format")
					//port     = c.String("port")
				)

				if name == "" {
					log.Fatalln("\nFailed: please pass component name e.g. --name=component")
				}

				// component has appdec.json
				if err := helper.ModuleExists(name); err != nil {
					log.Fatalln("\nComponent: " + name + " does not exists!")
				}

				// change to component directory
				if err = os.Chdir(name); err != nil {
					log.Fatalln("\nCant change to: "+name, err)
				}

				// compile files
				var cmdBuild *exec.Cmd
				cmdCompile := compile("src", "lib", watch, func() {

					if !production {
						return
					}

					cmdBuild = build("src/index.js", "lib/index.js", format, minify, true, true)
					err = cmdBuild.Run()
					fmt.Println("Build done!")
					if err != nil {
						log.Fatalln(err)
					}

					// @todo: lösche alle files außer lib/index.js
					// ... code ...
					// ... code ...

				})

				//@issue: ./appdec run --name=collapsible --production --watch --server when cmdCompile is .Start
				err = cmdCompile.Run()
				if err != nil {
					log.Fatalln(err)
				}

				if server {
					webserver("3000")
				}

				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete existing component",
			Flags: []cli.Flag{
				options.Name,
			},
			Action: func(c *cli.Context) error {

				name := strings.ToLower(c.String("name"))

				if name == "" {
					log.Fatalln("Failed: please pass module-name with --name=mymodule")
				}

				err := Delete(rootPath, name)
				if err != nil {
					log.Fatalln("Failed while deleting...", err)
				}

				return nil
			},
		},
		/*
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "list usage",
				Flags: []cli.Flag{
					options.Name,
				},
				Action: func(c *cli.Context) error {
					fmt.Println("test component")
					return nil
				},
			},
			{
				Name:    "publish",
				Aliases: []string{"p"},
				Usage:   "publish component on npm",
				Flags: []cli.Flag{
					options.Name,
				},
				Action: func(c *cli.Context) error {
					// use lerna (internal)
					fmt.Println("publish component")
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all available modules of app-decorators",
				Action: func(c *cli.Context) error {
					fmt.Println("list all modules")
					return nil
				},
			},
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "install usage",
				Flags: []cli.Flag{
					options.Name,
				},
				Action: func(c *cli.Context) error {

					// When installing an existing app-dec or vendor module
					// it will store the name of module and the type(existing or vendor)
					// in a file. This is important if we make a bundle/codesplitting
					// for current developed module.

					// Some other ideas:
					//
					fmt.Println("list existing app-dec or vendor module")
					return nil
				},
			},
		*/
	}

	app.Run(os.Args)
}
