package management

import (
	"bufio"
	"github.com/myOmikron/echotools/color"
	"github.com/pelletier/go-toml"
	"github.com/pnp-zone/pnp-manager/common"
	"os"
)

type Local struct {
	OutputDir string
	StaticDir string
	GoPath    string
}

type Global struct {
	Name            string
	VersionMajor    uint
	VersionMinor    uint
	VersionPatch    uint
	PNPVersionMajor uint
	PNPVersionMinor uint
	PNPVersionPatch uint
	Description     string
	License         string
	SourceURL       string
}

type Config struct {
	Local  Local
	Global Global
}

func Init(name string) {
	if !common.NameRegex.MatchString(name) {
		color.Println(color.RED, "Invalid name specified")
		os.Exit(1)
	}

	if err := os.Mkdir("static", 0700); err != nil {
		if os.IsExist(err) {
			reader := bufio.NewReader(os.Stdin)
			color.Print(color.YELLOW, "static/ already exists. Press [y] to continue. ")
			line, _, _ := reader.ReadLine()
			if string(line) != "y" {
				os.Exit(1)
			}

		} else {
			color.Println(color.RED, "Error: Could not create static directory:")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	}

	if _, err := os.Stat("plugin.toml"); !os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		color.Print(color.YELLOW, "plugin.toml already exists. Press [y] to overwrite. ")
		line, _, _ := reader.ReadLine()
		if string(line) != "y" {
			os.Exit(1)
		}
	}

	defaultConfig := Config{
		Local: Local{
			OutputDir: "upload/",
			StaticDir: "static/",
			GoPath:    name + ".so",
		},
		Global: Global{
			Name:         name,
			Description:  "Short description of your plugin.",
			License:      "",
			SourceURL:    "",
			VersionMajor: 0,
			VersionMinor: 1,
			VersionPatch: 0,
		},
	}
	if marshal, err := toml.Marshal(&defaultConfig); err != nil {
		color.Println(color.RED, "Could not marshal default config")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	} else {
		if err := os.WriteFile("plugin.toml", marshal, 0600); err != nil {
			color.Println(color.RED, "Could not write plugin.toml")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		} else {
			color.Println(color.GREEN, "done")
		}
	}
}
