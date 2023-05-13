//go:build windows

package main

import (
	"os"
)

func init() {
	err := os.RemoveAll("./asset/complete_resources/windows")
	if err != nil {
		panic(err)
	}

}
func run() error {
	err := os.Rename("./asset/cfc-proxyNet/deploy/windows", "./asset/complete_resources/windows")
	if err != nil {
		return err
	}
	err = CopyFile("./asset/box", "./asset/complete_resources/windows/box")
	if err != nil {
		return err
	}
	return nil
}
