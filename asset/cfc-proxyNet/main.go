package main

import (
	"github.com/peakedshout/go-CFC/loger"
	"net"
	"time"
)

func main() {
	runMain()
}

type runCtx struct {
	pc   *proxyCtx
	lc   *logCtx
	ac   *appCtx
	acfg *AppConfig

	onceLn net.Listener

	uuid   string
	active bool
}

func runMain() {
	rc := &runCtx{
		pc:     nil,
		lc:     nil,
		ac:     nil,
		acfg:   nil,
		onceLn: nil,
		uuid:   "",
		active: false,
	}
	defer func() {
		err := recover()
		if err != nil {
			loger.SetLogMust("Application death:", err)
			time.Sleep(3 * time.Second)
		}
	}()
	rc.newAppCtx()
}
