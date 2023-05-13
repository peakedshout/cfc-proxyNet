//go:build darwin

package main

import (
	"github.com/peakedshout/go-CFC/loger"
	"os/exec"
	"path"
	"strings"
)

func (rc *runCtx) getBoxPath(fileName string) string {
	return path.Join(rc.ac.app.ApplicationDirPath(), "../Resources/box/", fileName)
}
func (rc *runCtx) getUUID() (string, error) {
	b, err := exec.Command("system_profiler", "SPHardwareDataType").CombinedOutput()
	if err != nil {
		loger.SetLogDebug(errGetUUIDBad, err)
		return "", err
	}
	sl := strings.Split(string(b), "\n")
	str := ""
	for _, one := range sl {
		if strings.Contains(one, "Hardware UUID: ") {
			sx := strings.Split(one, "UUID: ")
			if len(sx) != 2 {
				break
			}
			str = sx[1]
			break
		}
	}
	if str == "" {
		loger.SetLogDebug(errGetUUIDBad, string(b))
		return "", errGetUUIDBad
	}
	out := strings.Replace(str, "-", "", -1)
	if len(out) != 32 {
		loger.SetLogDebug(errGetUUIDBad, string(b))
		return "", errGetUUIDBad
	}
	return out, nil
}
