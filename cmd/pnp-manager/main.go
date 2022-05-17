package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/pnp-zone/pnp-manager/management"
	"os"
)

func main() {
	parser := argparse.NewParser("pnp-manager", "", &argparse.ParserConfig{})

	configPath := parser.String("", "config-path", &argparse.Option{
		Default:     "/etc/pnp-manager/config.toml",
		Inheritable: true,
		Help:        "Set the path of the configuration file. Defaults to /etc/pnp-manager/config.toml",
	})

	installParser := parser.AddCommand("install", "Install plugins", &argparse.ParserConfig{})
	removeParser := parser.AddCommand("remove", "Remove plugins", &argparse.ParserConfig{})
	upgradeParser := parser.AddCommand("upgrade", "Upgrade plugins", &argparse.ParserConfig{})
	initParser := parser.AddCommand("init", "Initialize plugin structure", &argparse.ParserConfig{})
	initName := initParser.String("", "name", &argparse.Option{
		Positional: true,
		Required:   true,
	})
	buildParser := parser.AddCommand("build", "Build plugin", &argparse.ParserConfig{})
	uploadParser := parser.AddCommand("upload", "Upload plugin to ", &argparse.ParserConfig{})
	uploadPath := uploadParser.String("", "path", &argparse.Option{
		Positional: true,
		Required:   true,
	})

	searchParser := parser.AddCommand("search", "Search for plugins", &argparse.ParserConfig{})
	searchPattern := searchParser.String("", "search-pattern", &argparse.Option{
		Positional: true,
		Required:   true,
	})

	if err := parser.Parse(nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case installParser.Invoked:
	case removeParser.Invoked:
	case upgradeParser.Invoked:
	case initParser.Invoked:
		management.Init(*initName)
	case buildParser.Invoked:
	case uploadParser.Invoked:
		management.Upload(*uploadPath)
	case searchParser.Invoked:
		management.Search(*searchPattern, *configPath)
	}
}
