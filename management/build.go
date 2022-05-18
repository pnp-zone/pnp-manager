package management

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/myOmikron/echotools/color"
	"github.com/pelletier/go-toml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func checkDir(dirname string) {
	if dirname == "" {
		color.Println(color.RED, "error")
		color.Println(color.RED, dirname+" must be an existing directory.")
		os.Exit(1)
	}

	if stat, err := os.Stat(dirname); err != nil {
		if os.IsNotExist(err) {
			color.Println(color.RED, "error")
			color.Println(color.RED, dirname+" must be an existing directory.")
			os.Exit(1)
		} else {
			color.Println(color.RED, "error")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	} else {
		if !stat.IsDir() {
			color.Println(color.RED, "error")
			color.Println(color.RED, dirname+" must be an existing directory.")
			os.Exit(1)
		}
	}
}

func Build() {
	var config Config

	fmt.Print("Checking config ... ")

	if data, err := ioutil.ReadFile("plugin.toml"); err != nil {
		if os.IsNotExist(err) {
			color.Println(color.RED, "error")
			color.Println(color.RED, "plugin.toml was not found")
			os.Exit(1)
		} else {
			color.Println(color.RED, "error")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	} else {
		if err := toml.Unmarshal(data, &config); err != nil {
			color.Println(color.RED, "error")
			color.Println(color.RED, "Could not unmarshal plugin.toml")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	}

	if !strings.HasSuffix(config.Local.OutputDir, "/") {
		config.Local.OutputDir += "/"
	}

	if !strings.HasSuffix(config.Local.StaticDir, "/") {
		config.Local.StaticDir += "/"
	}

	if config.Local.OutputDir == "" {
		color.Println(color.RED, "error")
		color.Println(color.RED, "OutputDir must be empty")
		os.Exit(1)
	}

	if stat, err := os.Stat(config.Local.OutputDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(config.Local.OutputDir, 0700); err != nil {
				color.Println(color.RED, "error")
				color.Println(color.RED, "Could not create OutputDir")
				color.Println(color.RED, err.Error())
				os.Exit(1)
			}
		} else {
			color.Println(color.RED, "error")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	} else {
		if !stat.IsDir() {
			color.Println(color.RED, "error")
			color.Println(color.RED, config.Local.OutputDir+" must be an directory.")
			os.Exit(1)
		}
	}

	checkDir(config.Local.StaticDir)

	if _, err := os.Stat(config.Local.GoPath); err != nil {
		if os.IsNotExist(err) {
			color.Println(color.RED, "error")
			color.Println(color.RED, "GoPath must be an existing go plugin file.")
			os.Exit(1)
		} else {
			color.Println(color.RED, "error")
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	}

	color.Println(color.GREEN, "done")

	filename := fmt.Sprintf(
		"%s%s_v%d-%d-%d.tar",
		config.Local.OutputDir,
		config.Global.Name,
		config.Global.VersionMajor, config.Global.VersionMinor, config.Global.VersionPatch,
	)

	if _, err := os.Stat(filename + ".gz"); err != nil {
		if !os.IsNotExist(err) {
			color.Println(color.RED, err.Error())
			os.Exit(1)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		color.Print(color.YELLOW, "Build with that version already exists. Press [y] to overwrite. ")
		line, _, _ := reader.ReadLine()
		if string(line) != "y" {
			os.Exit(1)
		}
	}

	fmt.Print("Creating archive ... ")
	tarfile, err := os.Create(filename)
	if err != nil {
		color.Println(color.RED, "Could not build file in "+config.Local.OutputDir)
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}
	tarball := tar.NewWriter(tarfile)

	baseDir := filepath.Base(config.Local.StaticDir)

	if err := filepath.Walk(config.Local.StaticDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, config.Local.StaticDir))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		}); err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	info, err := os.Stat(config.Local.GoPath)
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	if err := tarball.WriteHeader(header); err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	file, err := os.Open(config.Local.GoPath)
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	_, err = io.Copy(tarball, file)
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	tarfile.Close()
	tarball.Close()

	color.Println(color.GREEN, "done")

	fmt.Print("Compressing archive ... ")
	reader, err := os.Open(filename)
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, "Could not open tar archive")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	writer, err := os.Create(filename + ".gz")
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	if _, err = io.Copy(archiver, reader); err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	if err := os.Remove(filename); err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, "Could not remove tar archive")
		os.Exit(1)
	}

	color.Println(color.GREEN, "done")

	fmt.Print("Writing config file ... ")
	config.Local = Local{}
	cleanedConfig, err := toml.Marshal(&config)
	if err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	if err := ioutil.WriteFile(filename+".toml", cleanedConfig, 0600); err != nil {
		color.Println(color.RED, "error")
		color.Println(color.RED, err.Error())
		os.Exit(1)
	}

	color.Println(color.GREEN, "done")

	color.Println(color.GREEN, "\nFinished building. You can now sign the build:")
	fmt.Println("\tgpg -u <YOUR_GPG_USER> --detach-sign " + filename + ".gz\n")
}
