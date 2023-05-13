package main

import (
	"github.com/peakedshout/go-CFC/loger"
	"os"
	"strings"
	"sync"
)

const appLogFileName = appName + "-" + appVersion + ".log"

type fLogCtx struct {
	wLock sync.Mutex
	once  int

	fName string
}

func (flc *fLogCtx) Write(b []byte) (n int, err error) {
	flc.wLock.Lock()
	defer flc.wLock.Unlock()
	var flag int
	if flc.once == 0 {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
		flc.once++
	} else {
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	}
	f, err := os.OpenFile(flc.fName, flag, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(b)
}

func newFLog(f string) {
	loger.SetLoggerLevel(loger.LogLevelError)
	loger.SetLoggerColor(false)
	flc := &fLogCtx{fName: f}
	loger.SetLoggerCopy(flc)
}

type logCtx struct {
	wLock sync.Mutex
	sl    []string
	ml    int
	cbFn  func(string)
}

func newLogCtx(ml int, fn func(string)) *logCtx {
	lc := &logCtx{
		wLock: sync.Mutex{},
		sl:    []string{},
		ml:    ml,
		cbFn:  fn,
	}
	loger.SetLoggerCopy(lc)
	return lc
}

func (lc *logCtx) Write(b []byte) (n int, err error) {
	lc.wLock.Lock()
	defer lc.wLock.Unlock()
	lc.sl = append([]string{string(b)}, lc.sl...)
	if len(lc.sl) > lc.ml {
		lc.sl = lc.sl[:lc.ml]
	}
	lc.cbFn(lc.String())
	return len(b), nil
}
func (lc *logCtx) String() string {
	return strings.Join(lc.sl, "\n")
}
