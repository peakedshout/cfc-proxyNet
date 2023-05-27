//go:build windows

package main

import (
	"github.com/peakedshout/go-CFC/loger"
	"os/exec"
	"path"
	"strings"
)

func (ctx *Ctx) getBoxPath(fileName string) string {
	return path.Join(ctx.app.ApplicationDirPath(), "./box/", fileName)
}

func (ctx *Ctx) getUUID() (string, error) {
	b, err := exec.Command("wmic", "csproduct", "get", "uuid").CombinedOutput()
	if err != nil {
		loger.SetLogDebug(errGetUUIDBad, err)
		return "", err
	}
	sl := strings.Split(string(b), "\n")
	if len(sl) < 2 || strings.TrimSpace(sl[0]) != "UUID" {
		loger.SetLogDebug(errGetUUIDBad, string(b))
		return "", errGetUUIDBad
	}
	str := strings.TrimSpace(sl[1])
	out := strings.Replace(str, "-", "", -1)
	if len(out) != 32 {
		loger.SetLogDebug(errGetUUIDBad, string(b))
		return "", errGetUUIDBad
	}
	return out, nil
}
