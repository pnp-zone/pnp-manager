package common

import (
	"errors"
	"fmt"
	"github.com/myOmikron/echotools/color"
	"github.com/pelletier/go-toml"
	"github.com/pnp-zone/pnp-manager/conf"
	"io/fs"
	"io/ioutil"
	"os"
)

func GetConfig(configPath string) *conf.Config {
	config := &conf.Config{}

	if configBytes, err := ioutil.ReadFile(configPath); errors.Is(err, fs.ErrNotExist) {
		color.Printf(color.RED, "Config was not found at %s\n", configPath)
		b, _ := toml.Marshal(config)
		fmt.Print(string(b))
		os.Exit(1)
	} else {
		if err := toml.Unmarshal(configBytes, config); err != nil {
			panic(err)
		}
	}

	config.Check()

	return config
}
