package search

import (
	"fmt"
	"github.com/myOmikron/echotools/color"
	"github.com/pnp-zone/pkg-manager/task"
	"github.com/pnp-zone/pnp-manager/common"
	"os"
	"strings"
)

func Search(searchPattern string, configPath string) {
	config := common.GetConfig(configPath)
	index, err := common.ReceiveIndex(config)
	if err != nil {
		color.Printf(color.RED, "ERROR:\n%s\n", err.Error())
		os.Exit(1)
	}

	exactMatches := make([]task.Package, 0)
	for _, p := range index {
		if strings.Contains(p.Name, searchPattern) {
			exactMatches = append(exactMatches, p)
		}
	}

	for _, match := range exactMatches {
		color.Print(color.GREEN, match.Name)

		var latest *task.PackageVersion
		for _, pv := range match.PackageVersions {
			if pv.FlagLatest {
				latest = &pv
			}
		}

		if latest == nil {
			continue
		}

		orphaned := ""
		if match.FlagOrphaned {
			orphaned = color.Colorize(color.RED, "(Orphaned) ")
		}
		fmt.Printf(" / %sv%d.%d.%d\n", orphaned, latest.VersionMajor, latest.VersionMinor, latest.VersionPatch)
		fmt.Println("\t" + latest.Description)
	}

}
