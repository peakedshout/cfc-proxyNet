package main

import (
	"github.com/peakedshout/cfc-proxyNet/settings"
	"github.com/peakedshout/go-CFC/loger"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"net"
	"os"
	"strconv"
	"sync"
)

const (
	appName      = "cfc-proxyNet"
	appVersion   = "v0.1.0"
	appUrl       = "https://github.com/peakedshout/cfc-proxyNet"
	appAuthor    = "peakedshout"
	appAuthorUrl = "https://github.com/peakedshout"
)

type appCtx struct {
	app *widgets.QApplication

	wd *widgets.QMainWindow
	ti *widgets.QSystemTrayIcon

	logo *gui.QIcon

	lock sync.Mutex
}

func (rc *runCtx) newAppCtx() {
	rc.ac = &appCtx{}

	rc.ac.app = widgets.NewQApplication(len(os.Args), os.Args)

	rc.buildPre()
	defer rc.delWork()

	rc.ac.app.SetApplicationName(appName)
	rc.ac.app.SetApplicationVersion(appVersion)
	rc.ac.app.SetDesktopFileName(appName)

	rc.ac.logo = gui.NewQIcon5(rc.getBoxPath("cfcproxynet_logo.png"))

	rc.ac.app.SetWindowIcon(rc.ac.logo)

	rc.buildWindow()

	rc.ac.app.SetQuitOnLastWindowClosed(false)

	rc.ac.app.ConnectAboutToQuit(func() {
		if rc.pc != nil {
			settings.Close()
		}
	})

	rc.ac.app.ConnectEvent(func(e *core.QEvent) bool {
		if e.Type() == core.QEvent__ApplicationActivate {
			rc.ac.wd.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
			rc.ac.wd.Show()
		}
		return false
	})

	rc.ac.app.Exec()
}

func (rc *runCtx) buildPre() {
	newFLog(rc.getBoxPath(appLogFileName))
	uuid, err := rc.getUUID()
	if err != nil {
		loger.SetLogError(err)
	}
	rc.uuid = uuid
	rc.acfg = newACfg(rc.getBoxPath(appConfigFileName), rc.uuid)
	err = rc.acfg.read()
	if err != nil {
		loger.SetLogError(err)
	}

	rc.onceLn, err = net.Listen("tcp", "127.0.0.1:19999")
	if err != nil {
		loger.SetLogDebug(errOnlyOnceStartFailurePleaseCheckAppStateOrNetPort, err)
		loger.SetLogError(errOnlyOnceStartFailurePleaseCheckAppStateOrNetPort)
		return
	}
	return
}
func (rc *runCtx) delWork() {
	rc.onceLn.Close()
}

func (rc *runCtx) buildWindow() {
	//window	-----------------------------------------------------------------------------------------
	rc.ac.wd = widgets.NewQMainWindow(nil, 0)
	rc.ac.wd.SetWindowTitle("cfc-proxy")
	rc.ac.wd.SetMaximumWidth(300)
	rc.ac.wd.SetMinimumWidth(300)

	//l1
	l1 := widgets.NewQWidget(nil, 0)
	l1.SetLayout(widgets.NewQVBoxLayout())

	//l10
	l10 := widgets.NewQWidget(nil, 0)
	l10.SetLayout(widgets.NewQHBoxLayout())

	l101 := widgets.NewQLabel(nil, 0)
	l101.SetText("Cfg:")
	l102 := widgets.NewQComboBox(nil)
	l102.Clear()
	acfgIdSl := rc.acfg.getIdList()
	l102.AddItems(acfgIdSl)
	l102.SetSizeAdjustPolicy(widgets.QComboBox__AdjustToContents)
	l103 := widgets.NewQComboBox(nil)
	l103.Clear()
	scdCfgSl := []string{"saveCfg", "copy&&saveCfg", "deleteCfg"}
	scdCfgSl2 := []string{"saveCfg"}
	l103.AddItems(scdCfgSl2)

	l10.Layout().AddWidget(l101)
	l10.Layout().AddWidget(l102)
	l10.Layout().AddWidget(l103)

	//l11
	l11 := widgets.NewQWidget(nil, 0)
	l11.SetLayout(widgets.NewQVBoxLayout())

	l111 := widgets.NewQLabel2("proxy host addr", nil, core.Qt__Widget)
	l112 := widgets.NewQLineEdit(nil)
	l112.SetPlaceholderText("eg:127.0.0.1:8080")
	l11.Layout().AddWidget(l111)
	l11.Layout().AddWidget(l112)

	//l12
	l12 := widgets.NewQWidget(nil, 0)
	l12.SetLayout(widgets.NewQVBoxLayout())

	l121 := widgets.NewQLabel2("proxy host key", nil, core.Qt__Widget)
	l122 := widgets.NewQLineEdit(nil)
	l122.SetPlaceholderText("key len must 32 bytes")
	l122.SetMaxLength(32)
	l122.SetEchoMode(widgets.QLineEdit__PasswordEchoOnEdit)
	l12.Layout().AddWidget(l121)
	l12.Layout().AddWidget(l122)

	//l13
	l13 := widgets.NewQWidget(nil, 0)
	l13.SetLayout(widgets.NewQVBoxLayout())

	l131 := widgets.NewQLabel2("proxy http/https port", nil, core.Qt__Widget)
	l132 := widgets.NewQLineEdit(nil)
	l132.SetPlaceholderText("listen http/https port")
	l132.SetValidator(gui.NewQIntValidator2(0, 65535, nil))
	l13.Layout().AddWidget(l131)
	l13.Layout().AddWidget(l132)

	const (
		active   = "<font color = green>now: active</font>"
		inactive = "<font color = red>now: inactive</font>"
		running  = "<font color = grey>now: running...</font>"

		start = "start"
		stop  = "stop"

		testAndPing = "test&&ping"
	)

	//l14
	l14 := widgets.NewQWidget(nil, 0)
	l14.SetLayout(widgets.NewQHBoxLayout())

	l141 := widgets.NewQLabel(nil, 0)
	l141.SetLayout(widgets.NewQVBoxLayout())
	l141.SetText(inactive)
	l142 := widgets.NewQPushButton(nil)
	l142.SetText(start)
	l143 := widgets.NewQPushButton(nil)
	l143.SetText(testAndPing)
	l143.SetDisabled(true)
	l14.Layout().AddWidget(l141)
	l14.Layout().AddWidget(l142)
	l14.Layout().AddWidget(l143)

	//l15
	l15 := widgets.NewQWidget(nil, 0)
	l15.SetLayout(widgets.NewQHBoxLayout())

	l151 := widgets.NewQLabel(nil, 0)
	l151.SetLayout(widgets.NewQVBoxLayout())
	l151.SetText("log settings :")
	l152 := widgets.NewQComboBox(nil)
	l152.SetLayout(widgets.NewQVBoxLayout())
	//l152.Clear()
	logLevelSl := []string{"all", "trace", "debug", "info", "warn", "log", "error", "fatal", "off", "must"}
	l152.AddItems(logLevelSl)
	l152.ConnectActivated(func(index int) {
		loger.SetLoggerLevel(uint8(index))
	})
	l153 := widgets.NewQCheckBox2("stack", nil)
	l153.SetLayout(widgets.NewQVBoxLayout())
	sr := false
	l153.ConnectClicked(func(checked bool) {
		sr = !sr
		loger.SetLoggerStack(sr)
	})
	l15.Layout().AddWidget(l151)
	l15.Layout().AddWidget(l152)
	l15.Layout().AddWidget(l153)

	//log
	ll := widgets.NewQTextBrowser(nil)
	ll.SetLayout(widgets.NewQVBoxLayout())
	rc.lc = newLogCtx(100, func(s string) {
		ll.SetText(s)
	})

	//info
	i1 := widgets.NewQLabel(nil, 0)
	i1.SetText(`Application name:		<a style='color : cyan;' href=` + appUrl + `>` + appName + ` ` + appVersion + `</a>`)
	i1.SetOpenExternalLinks(true)
	i1.ConnectLinkActivated(func(link string) {
		gui.QDesktopServices_OpenUrl(core.NewQUrl3(link, core.QUrl__TolerantMode))
	})
	i2 := widgets.NewQLabel(nil, 0)
	i2.SetText(`Application author:		<a style='color : cyan;' href=` + appAuthorUrl + `>` + appAuthor + `</a>`)
	i2.SetOpenExternalLinks(true)
	i2.ConnectLinkActivated(func(link string) {
		loger.SetLogWarn(link)
		gui.QDesktopServices_OpenUrl(core.NewQUrl3(link, core.QUrl__TolerantMode))
	})

	//over
	l1.Layout().AddWidget(l10)
	l1.Layout().AddWidget(l11)
	l1.Layout().AddWidget(l12)
	l1.Layout().AddWidget(l13)
	l1.Layout().AddWidget(l14)
	l1.Layout().AddWidget(l15)
	l1.Layout().AddWidget(ll)
	l1.Layout().AddWidget(i1)
	l1.Layout().AddWidget(i2)

	rc.ac.wd.SetCentralWidget(l1)
	rc.ac.wd.Show()
	//window	-----------------------------------------------------------------------------------------

	//TrayIcon	-----------------------------------------------------------------------------------------
	rc.ac.ti = widgets.NewQSystemTrayIcon2(rc.ac.logo, nil)

	rc.ac.ti.SetIcon(rc.ac.logo)

	const (
		startService = "Start service"
		stopService  = "Stop service"
	)
	getStateFn := func(r bool) string {
		s := "Active"
		if !r {
			s = "Inactive"
		}
		return appName + "(" + appVersion + ")" + ": " + s
	}
	//m1
	m1 := widgets.NewQMenu(nil)
	b0 := m1.AddAction(getStateFn(false))
	b0.SetDisabled(true)

	b1 := m1.AddAction("Home")
	b1.ConnectTriggered(func(checked bool) {
		rc.ac.wd.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
		rc.ac.wd.Show()
	})
	m1.AddSeparator()

	b2 := m1.AddAction(startService)
	m1.AddSeparator()
	b3 := m1.AddAction("Exit")
	b3.ConnectTriggered(func(checked bool) {
		rc.ac.app.Exit(0)
	})
	//over

	rc.ac.ti.SetContextMenu(m1)

	rc.ac.ti.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
		switch reason {
		case widgets.QSystemTrayIcon__Context, widgets.QSystemTrayIcon__Trigger:
		case widgets.QSystemTrayIcon__DoubleClick, widgets.QSystemTrayIcon__MiddleClick:
			rc.ac.wd.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
			rc.ac.wd.Show()
		}
	})

	rc.ac.ti.Show()
	//TrayIcon	-----------------------------------------------------------------------------------------

	del1Fn := func(r bool) {
		l112.SetDisabled(r)
		l122.SetDisabled(r)
		l132.SetDisabled(r)
		l102.SetDisabled(r)
		l103.SetDisabled(r)
	}
	delsFn := func(r bool) {
		l142.SetDisabled(r)
		b2.SetDisabled(r)
	}
	proxyFn := func() {
		rc.ac.lock.Lock()
		defer rc.ac.lock.Unlock()
		delsFn(true)
		defer delsFn(false)
		if rc.active {
			if rc.pc != nil {
				rc.pc.close(nil)
			}
			b2.SetText(startService)
			l142.SetText(start)
			l141.SetText(inactive)

			del1Fn(false)

			rc.active = false
			b0.SetText(getStateFn(rc.active))
			l143.SetDisabled(true)
			loger.SetLogXY("now: inactive service")
		} else {
			l141.SetText(running)

			del1Fn(true)

			pc, err := runProxy("127.0.0.1:"+l132.Text(), l112.Text(), l122.Text())
			if err != nil {
				loger.SetLogWarn(err)
				l141.SetText(inactive)
				return
			}
			rc.pc = pc
			go func() {
				err := rc.pc.wait()
				if err != nil {
					loger.SetLogWarn(err)
				}
			}()

			b2.SetText(stopService)
			l142.SetText(stop)
			l141.SetText(active)
			rc.active = true
			b0.SetText(getStateFn(rc.active))
			l143.SetDisabled(false)
			loger.SetLogXY("now: active service")
		}
		return
	}

	defaultDataFn := func() {
		l112.SetText("sv2.peakedshout.top:9988")
		l122.SetText("6a647c0bf889419c84e461486f83d776")
		l132.SetText("9988")

		l152.SetCurrentIndex(loger.LogLevelWarn)
		loger.SetLoggerLevel(loger.LogLevelWarn)
		l153.SetChecked(sr)
		loger.SetLoggerStack(sr)
	}
	defaulL103Fn := func(text string) {
		if text == appConfigNone {
			l103.Clear()
			l103.AddItems(scdCfgSl2)
		} else {
			l103.Clear()
			l103.AddItems(scdCfgSl)
		}
	}

	if rc.acfg.LastConfig != "" {
		c := rc.acfg.getConfig(rc.acfg.LastConfig)
		l112.SetText(c.ProxyServerHost.ProxyServerAddr)
		l122.SetText(c.ProxyServerHost.LinkProxyKey)
		l132.SetText(strconv.Itoa(c.ProxyMethod.Http.Port))
		l152.SetCurrentIndex(int(c.Setting.LogLevel))
		l153.SetChecked(c.Setting.LogStack)
		sr = c.Setting.LogStack
		loger.SetLoggerLevel(c.Setting.LogLevel)
		loger.SetLoggerStack(sr)

		for i, one := range acfgIdSl {
			if one == rc.acfg.LastConfig {
				l102.SetCurrentIndex(i)
				break
			}
		}
	} else {
		defaultDataFn()
		l102.SetCurrentIndex(0)
	}

	l142.ConnectClicked(func(checked bool) {
		proxyFn()
	})
	b2.ConnectTriggered(func(checked bool) {
		proxyFn()
	})
	l102.ConnectActivated2(func(text string) {
		rc.ac.lock.Lock()
		defer rc.ac.lock.Unlock()
		delsFn(true)
		defer delsFn(false)
		del1Fn(true)
		defer del1Fn(false)
		c := rc.acfg.getConfig(text)
		l112.SetText(c.ProxyServerHost.ProxyServerAddr)
		l122.SetText(c.ProxyServerHost.LinkProxyKey)
		l132.SetText(strconv.Itoa(c.ProxyMethod.Http.Port))
		l152.SetCurrentIndex(int(c.Setting.LogLevel))
		l153.SetChecked(c.Setting.LogStack)
		sr = c.Setting.LogStack
		loger.SetLoggerLevel(c.Setting.LogLevel)
		loger.SetLoggerStack(sr)
		for i, one := range acfgIdSl {
			if one == rc.acfg.LastConfig {
				l102.SetCurrentIndex(i)
				break
			}
		}
		err := rc.acfg.setLastConfig(text)
		if err != nil {
			loger.SetLogWarn(err)
		}
		l102.SetCurrentText(text)
		defaulL103Fn(text)
	})
	l103.ConnectActivated2(func(text string) {
		rc.ac.lock.Lock()
		defer rc.ac.lock.Unlock()
		delsFn(true)
		defer delsFn(false)
		del1Fn(true)
		defer del1Fn(false)
		nowN := ""
		switch text {
		case scdCfgSl[0]:
			nowName, err := rc.acfg.cutConfig(l102.CurrentText(), l112.Text(), l122.Text(), l132.Text(), "0", uint8(l152.CurrentIndex()), l153.IsChecked())
			if err != nil {
				loger.SetLogXY(err)
				return
			}
			nowN = nowName
		case scdCfgSl[1]:
			nowName, err := rc.acfg.setConfig2(l112.Text(), l122.Text(), l132.Text(), "0", uint8(l152.CurrentIndex()), l153.IsChecked())
			if err != nil {
				loger.SetLogXY(err)
				return
			}
			nowN = nowName
		case scdCfgSl[2]:
			err := rc.acfg.delConfig(l102.CurrentText())
			if err != nil {
				loger.SetLogXY(err)
				return
			}
			nowN = appConfigNone
			defaultDataFn()
		}
		l102.Clear()
		acfgIdSl = rc.acfg.getIdList()
		l102.AddItems(acfgIdSl)
		l102.SetCurrentText(nowN)
		defaulL103Fn(nowN)
	})
	l143.ConnectClicked(func(checked bool) {
		rc.ac.lock.Lock()
		defer rc.ac.lock.Unlock()
		l143.SetDisabled(true)
		defer l143.SetDisabled(false)
		delsFn(true)
		defer delsFn(false)
		rx, err := testProxy("127.0.0.1:" + l132.Text())
		if err != nil {
			loger.SetLogXY("Test failure:", err)
			return
		}
		loger.SetLogXY("Test succeed:", rx)
	})
}
