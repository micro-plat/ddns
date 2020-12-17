package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/micro-plat/ddns/static"
)

var archive = "./static.zip"

func init() {
	_, err := os.Stat(archive)
	if err == nil {
		return
	}
	if !os.IsNotExist(err) {
		panic(err)
	}
	for _, v := range static.AssetNames() {
		err := os.MkdirAll(filepath.Dir(v), 0444)
		if err != nil {
			panic(err)
		}
		buff, err := static.Asset(v)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(v, buff, 0777)
		if err != nil {
			panic(err)
		}
	}
	return
}
