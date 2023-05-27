package main

import (
	"flag"
	"github.com/peakedshout/cfc-proxyNet/config"
	"github.com/peakedshout/go-CFC/loger"
	"github.com/peakedshout/go-CFC/tool"
	"github.com/peakedshout/go-socks/relay"
	"github.com/peakedshout/go-socks/share"
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

	sc := &relay.ServerConfig{
		Addr:   c.ProxyServerHost.ProxyServerAddr,
		RawKey: c.ProxyServerHost.LinkProxyKey,
		CMDSwitch: share.SocksCMDSwitch{
			SwitchCMDCONNECT:      true,
			SwitchCMDBIND:         true,
			SwitchCMDUDPASSOCIATE: true,
		},
		ConnTimeout: 0,
		DialTimeout: 0,
		BindTimeout: 0,
	}
	err := tool.ReRun(c.Setting.ReLinkTime, func() bool {
		rs, err := relay.NewServer(sc)
		if err != nil {
			loger.SetLogWarn(err)
		}
		err = rs.Wait()
		if err != nil {
			loger.SetLogWarn(err)
		}
		return true
	})
	if err != nil {
		loger.SetLogError(err)
	}
}
