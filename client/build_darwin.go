//go:build darwin

package main

import (
	"os"
)

func init() {
	err := os.RemoveAll("./asset/complete_resources/darwin")
	if err != nil {
		panic(err)
	}
}

func run() error {
	err := os.Rename("./asset/cfc-proxyNet/deploy/darwin", "./asset/complete_resources/darwin")
	if err != nil {
		return err
	}
	err = CopyFile("./asset/box", "./asset/complete_resources/darwin/cfc-proxyNet.app/Contents/Resources/box")
	if err != nil {
		return err
	}
	err = CopyFile("./asset/box/cpnlogo.icns", "./asset/complete_resources/darwin/cfc-proxyNet.app/Contents/Resources/cfc-proxyNet.icns")
	if err != nil {
		return err
	}
	return nil
}
func del() error {
	return nil
}
