package main

import (
	"github.com/peakedshout/go-CFC/loger"
	"time"
)

func main() {
	runMain()
}

func runMain() {
	ctx := Ctx{
		onceLn:    nil,
		app:       nil,
		wd:        nil,
		ti:        nil,
		logo:      nil,
		serverCtx: nil,
		active:    false,
		uuid:      "",
		acfg:      nil,
		logCtx:    nil,
		wc:        wCtx{},
	}
	defer func() {
		err := recover()
		if err != nil {
			loger.SetLogMust("Application death:", err)
			time.Sleep(3 * time.Second)
		}
	}()
	ctx.newApp()
}
