package main

import (
	"flag"
	"github.com/peakedshout/cfc-proxyNet/config"
	"github.com/peakedshout/go-CFC/loger"
	"github.com/peakedshout/go-CFC/server"
	"github.com/peakedshout/go-CFC/tool"
)

func main() {
	p := flag.String("c", "./config.json", "config file path , default is ./config.json")
	flag.Parse()
	c := config.ReadConfig(*p)
	if c.Setting.LogLevel == 0 {
		c.Setting.LogLevel = loger.LogLevelWarn
	}
	loger.SetLoggerLevel(c.Setting.LogLevel)
	loger.SetLoggerStack(c.Setting.LogStack)

	sc := &server.Config{
		RawKey:           c.ProxyServerHost.LinkProxyKey,
		LnAddr:           c.ProxyServerHost.ProxyServerAddr,
		SwitchVPNProxy:   true,
		SwitchLinkClient: false,
	}

	err := tool.ReRun(c.Setting.ReLinkTime, func() bool {
		server.NewProxyServer2(sc).Wait()
		return true
	})
	if err != nil {
		loger.SetLogError(err)
	}
}
