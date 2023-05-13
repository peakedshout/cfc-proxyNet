package main

import (
	"github.com/peakedshout/go-CFC/loger"
	"os/exec"
)

func errCheck(err error, b []byte) {
	if err != nil {
		loger.SetLogError(err, string(b))
	}
}

func main() {
	var err error
	var cmd *exec.Cmd
	var b []byte
	cmd = exec.Command("go", "run", "./client")
	b, err = cmd.CombinedOutput()
	errCheck(err, b)
	cmd = exec.Command("go", "run", "./server")
	b, err = cmd.CombinedOutput()
	errCheck(err, b)
}
