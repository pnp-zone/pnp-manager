package main

import (
	"fmt"
	"github.com/hellflame/argparse"
	"github.com/pnp-zone/pnp-manager/search"
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
	mkpluginParser := parser.AddCommand("mkplugin", "Make plugin", &argparse.ParserConfig{})
	uploadParser := parser.AddCommand("upload", "Upload plugin to ", &argparse.ParserConfig{})
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
	case mkpluginParser.Invoked:
	case uploadParser.Invoked:
	case searchParser.Invoked:
		search.Search(*searchPattern, *configPath)
	}
}
