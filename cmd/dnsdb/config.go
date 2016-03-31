package main

import (
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var CONF = config{}

type config struct {
	APIKEY string
}

func getDefaultConfPath() (path string) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	dir := usr.HomeDir
	return filepath.Join(dir, ".dnsdb-query.conf")
}

func loadConfig(path string) error {
	_, err := toml.DecodeFile(path, &CONF)
	if err != nil {
		return err
	}
	return nil
}
