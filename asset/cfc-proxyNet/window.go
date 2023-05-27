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
)

const (
	appName      = "cfc-proxyNet"
	appVersion   = "v0.2.0"
	appUrl       = "https://github.com/peakedshout/cfc-proxyNet"
	appAuthor    = "peakedshout"
	appAuthorUrl = "https://github.com/peakedshout"
)

type Ctx struct {
	onceLn net.Listener
	app    *widgets.QApplication
	wd     *widgets.QMainWindow
	ti     *widgets.QSystemTrayIcon

	logo *gui.QIcon

	serverCtx *serverCtx

	active bool
	uuid   string
	acfg   *AppConfig
	logCtx *logCtx

	wc wCtx
}

type wCtx struct {
	ctx *Ctx

	c1, c2                 *widgets.QComboBox
	box1, box2, box3, box4 *widgets.QLineEdit
	s1                     *widgets.QLabel
	s2, s3                 *widgets.QPushButton
	s4                     *widgets.QComboBox
	s5                     *widgets.QCheckBox
	i1, i2                 *widgets.QAction
}

func (ctx *Ctx) buildPre() {
	ctx.app.SetApplicationName(appName)
	ctx.app.SetApplicationVersion(appVersion)
	ctx.app.SetDesktopFileName(appName)
	ctx.logo = gui.NewQIcon5(ctx.getBoxPath("cfcproxynet_logo.png"))
	ctx.app.SetWindowIcon(ctx.logo)

	ctx.app.SetQuitOnLastWindowClosed(false)
	ctx.app.ConnectAboutToQuit(func() {
		if ctx.serverCtx != nil {
			settings.Close()
		}
	})
	ctx.app.ConnectEvent(func(e *core.QEvent) bool {
		if e.Type() == core.QEvent__ApplicationActivate {
			ctx.wd.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
			ctx.wd.Show()
		}
		return false
	})

	newFLog(ctx.getBoxPath(appLogFileName))
	uuid, err := ctx.getUUID()
	if err != nil {
		loger.SetLogError(err)
	}
	ctx.uuid = uuid
	ctx.acfg = newACfg(ctx.getBoxPath(appConfigFileName), ctx.uuid)
	err = ctx.acfg.read()
	if err != nil {
		loger.SetLogError(err)
	}

	ctx.onceLn, err = net.Listen("tcp", "127.0.0.1:19999")
	if err != nil {
		loger.SetLogDebug(errOnlyOnceStartFailurePleaseCheckAppStateOrNetPort, err)
		loger.SetLogError(errOnlyOnceStartFailurePleaseCheckAppStateOrNetPort)
		return
	}
	return
}
func (ctx *Ctx) delWork() {
	ctx.onceLn.Close()
}

func (ctx *Ctx) newApp() {
	ctx.app = widgets.NewQApplication(len(os.Args), os.Args)
	ctx.buildPre()
	defer ctx.delWork()

	ctx.wc.ctx = ctx
	ctx.buildWindow()
	ctx.wc.i1, ctx.wc.i2 = ctx.buildTrayIcon()

	ctx.wc.setFn()

	ctx.app.Exec()
}

func (ctx *Ctx) buildWindow() {
	//window	-----------------------------------------------------------------------------------------
	ctx.wd = widgets.NewQMainWindow(nil, 0)
	ctx.wd.SetWindowTitle(appName)
	ctx.wd.SetMaximumWidth(300)
	ctx.wd.SetMinimumWidth(300)

	//l1
	l1 := widgets.NewQWidget(nil, 0)
	l1.SetLayout(widgets.NewQVBoxLayout())

	ctx.wc.c1, ctx.wc.c2 = ctx.buildConfig(l1)
	ctx.wc.box1, ctx.wc.box2, ctx.wc.box3, ctx.wc.box4 = ctx.buildBox(l1)
	ctx.wc.s1, ctx.wc.s2, ctx.wc.s3, ctx.wc.s4, ctx.wc.s5 = ctx.buildSwitch(l1)
	ctx.buildInfo(l1)
	ctx.wd.SetCentralWidget(l1)
	ctx.wd.Show()
}

var scdCfgSl = []string{"saveCfg", "copy&&saveCfg", "deleteCfg"}
var scdCfgSl2 = []string{"saveCfg"}

func (ctx *Ctx) buildConfig(l1 *widgets.QWidget) (q1, q2 *widgets.QComboBox) {
	l11 := widgets.NewQWidget(nil, 0)
	l11.SetLayout(widgets.NewQHBoxLayout())
	l111 := widgets.NewQLabel2("Cfg:", nil, 0)
	l112 := widgets.NewQComboBox(nil)
	l112.Clear()
	acfgIdSl := ctx.acfg.getIdList()
	l112.AddItems(acfgIdSl)
	l112.SetSizeAdjustPolicy(widgets.QComboBox__AdjustToContents)
	l113 := widgets.NewQComboBox(nil)
	l113.Clear()

	l113.AddItems(scdCfgSl2)
	l11.Layout().AddWidget(l111)
	l11.Layout().AddWidget(l112)
	l11.Layout().AddWidget(l113)

	l1.Layout().AddWidget(l11)

	return l112, l113
}

func (ctx *Ctx) buildBox(l1 *widgets.QWidget) (q1, q2, q3, q4 *widgets.QLineEdit) {
	//l11
	l11 := widgets.NewQWidget(nil, 0)
	l11.SetLayout(widgets.NewQVBoxLayout())
	t1 := widgets.NewQLabel2("proxy host addr", nil, core.Qt__Widget)
	e1 := widgets.NewQLineEdit(nil)
	e1.SetPlaceholderText("eg:127.0.0.1:8080")
	l11.Layout().AddWidget(t1)
	l11.Layout().AddWidget(e1)

	t2 := widgets.NewQLabel2("proxy host key", nil, core.Qt__Widget)
	e2 := widgets.NewQLineEdit(nil)
	e2.SetPlaceholderText("key len must 32 bytes")
	e2.SetMaxLength(32)
	e2.SetEchoMode(widgets.QLineEdit__PasswordEchoOnEdit)
	l11.Layout().AddWidget(t2)
	l11.Layout().AddWidget(e2)

	t3 := widgets.NewQLabel2("proxy http/https port", nil, core.Qt__Widget)
	e3 := widgets.NewQLineEdit(nil)
	e3.SetPlaceholderText("listen http/https port")
	e3.SetValidator(gui.NewQIntValidator2(0, 65535, nil))
	l11.Layout().AddWidget(t3)
	l11.Layout().AddWidget(e3)

	t4 := widgets.NewQLabel2("proxy socks port", nil, core.Qt__Widget)
	e4 := widgets.NewQLineEdit(nil)
	e4.SetPlaceholderText("listen socks port")
	e4.SetValidator(gui.NewQIntValidator2(0, 65535, nil))
	l11.Layout().AddWidget(t4)
	l11.Layout().AddWidget(e4)

	l1.Layout().AddWidget(l11)

	return e1, e2, e3, e4
}

const (
	active   = "<font color = green>now: active</font>"
	inactive = "<font color = red>now: inactive</font>"
	running  = "<font color = grey>now: running...</font>"

	start = "start"
	stop  = "stop"

	testAndPing = "test&&ping"
)

func (ctx *Ctx) buildSwitch(l1 *widgets.QWidget) (q1 *widgets.QLabel, q2, q3 *widgets.QPushButton, q4 *widgets.QComboBox, q5 *widgets.QCheckBox) {

	logLevelSl := []string{"all", "trace", "debug", "info", "warn", "log", "error", "fatal", "off", "must"}
	l11 := widgets.NewQWidget(nil, 0)
	l11.SetLayout(widgets.NewQHBoxLayout())

	t1 := widgets.NewQLabel2(inactive, nil, 0)
	b1 := widgets.NewQPushButton2(start, nil)
	b2 := widgets.NewQPushButton2(testAndPing, nil)
	l11.Layout().AddWidget(t1)
	l11.Layout().AddWidget(b1)
	l11.Layout().AddWidget(b2)

	l1.Layout().AddWidget(l11)

	l12 := widgets.NewQWidget(nil, 0)
	l12.SetLayout(widgets.NewQHBoxLayout())

	t2 := widgets.NewQLabel2("log settings :", nil, 0)
	b3 := widgets.NewQComboBox(nil)
	b3.Clear()
	b3.AddItems(logLevelSl)
	b3.ConnectActivated(func(index int) {
		loger.SetLoggerLevel(uint8(index))
	})
	b4 := widgets.NewQCheckBox2("stack", nil)
	sr := false
	b4.ConnectClicked(func(checked bool) {
		sr = !sr
		loger.SetLoggerStack(sr)
	})
	l12.Layout().AddWidget(t2)
	l12.Layout().AddWidget(b3)
	l12.Layout().AddWidget(b4)

	l1.Layout().AddWidget(l12)

	return t1, b1, b2, b3, b4
}

func (ctx *Ctx) buildInfo(l1 *widgets.QWidget) {
	//log
	ll := widgets.NewQTextBrowser(nil)
	ll.SetLayout(widgets.NewQVBoxLayout())
	ctx.logCtx = newLogCtx(100, func(s string) {
		ll.SetText(s)
	})

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

	l1.Layout().AddWidget(ll)
	l1.Layout().AddWidget(i1)
	l1.Layout().AddWidget(i2)
}

const (
	startService = "Start service"
	stopService  = "Stop service"

	startService2 = appName + "(" + appVersion + ")" + ": Active"
	stopService2  = appName + "(" + appVersion + ")" + ": Inactive"
)

func (ctx *Ctx) buildTrayIcon() (q1, q2 *widgets.QAction) {
	ctx.ti = widgets.NewQSystemTrayIcon2(ctx.logo, nil)
	m1 := widgets.NewQMenu(nil)
	b0 := m1.AddAction(stopService2)
	b0.SetDisabled(true)

	b1 := m1.AddAction("Home")
	b1.ConnectTriggered(func(checked bool) {
		ctx.wd.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
		ctx.wd.Show()
	})
	m1.AddSeparator()

	b2 := m1.AddAction(startService)
	m1.AddSeparator()
	b3 := m1.AddAction("Exit")
	b3.ConnectTriggered(func(checked bool) {
		ctx.app.Exit(0)
	})

	ctx.ti.SetContextMenu(m1)

	ctx.ti.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
		switch reason {
		case widgets.QSystemTrayIcon__Context, widgets.QSystemTrayIcon__Trigger:
		case widgets.QSystemTrayIcon__DoubleClick, widgets.QSystemTrayIcon__MiddleClick:
			ctx.wd.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
			ctx.wd.Show()
		}
	})

	ctx.ti.Show()

	return b0, b2
}

func (wc *wCtx) switchFn(r bool) {
	wc.c1.SetDisabled(r)
	wc.c2.SetDisabled(r)
	wc.box1.SetDisabled(r)
	wc.box2.SetDisabled(r)
	wc.box3.SetDisabled(r)
	wc.box4.SetDisabled(r)
	wc.i1.SetDisabled(r)
}
func (wc *wCtx) switchFn2(r bool) {
	wc.i2.SetDisabled(r)
	wc.s2.SetDisabled(r)
}

func (wc *wCtx) setFn() {
	proxyFn := func() {
		wc.switchFn2(true)
		defer wc.switchFn2(false)
		if wc.ctx.active {
			if wc.ctx.serverCtx != nil {
				wc.ctx.serverCtx.Close(nil)
				wc.i1.SetText(stopService2)
				wc.i2.SetText(startService)
				wc.s1.SetText(inactive)
				wc.s2.SetText(start)

				wc.ctx.active = false
				loger.SetLogXY("now: inactive service")

				wc.switchFn(false)
				wc.s3.SetDisabled(true)
			}
		} else {
			wc.s1.SetText(running)
			sc, err := newServerCtx("127.0.0.1:"+wc.box3.Text(), "127.0.0.1:"+wc.box4.Text(), wc.box1.Text(), wc.box2.Text())
			if err != nil {
				loger.SetLogWarn(err)
				wc.s1.SetText(inactive)
				return
			}
			wc.ctx.serverCtx = sc
			go func() {
				err = wc.ctx.serverCtx.Wait()
				if err != nil {
					loger.SetLogWarn(err)
				}
			}()
			wc.i1.SetText(startService2)
			wc.i2.SetText(stopService)
			wc.s1.SetText(active)
			wc.s2.SetText(stop)
			wc.ctx.active = true

			wc.s3.SetDisabled(false)
			wc.switchFn(true)

			loger.SetLogXY("now: active service")
		}
	}

	defaultDataFn := func() {
		wc.box1.SetText("")
		wc.box2.SetText("")
		wc.box3.SetText("")
		wc.box4.SetText("")

		wc.s4.SetCurrentIndex(loger.LogLevelWarn)
		loger.SetLoggerLevel(loger.LogLevelWarn)
		wc.s5.SetChecked(false)
		loger.SetLoggerStack(false)
	}
	defaulC2Fn := func(text string) {
		if text == appConfigNone {
			wc.c2.Clear()
			wc.c2.AddItems(scdCfgSl2)
		} else {
			wc.c2.Clear()
			wc.c2.AddItems(scdCfgSl)
		}
	}

	acfgIdSl := wc.ctx.acfg.getIdList()
	wc.c1.Clear()
	wc.c1.AddItems(acfgIdSl)

	if wc.ctx.acfg.LastConfig != "" {
		c := wc.ctx.acfg.getConfig(wc.ctx.acfg.LastConfig)
		wc.box1.SetText(c.ProxyServerHost.ProxyServerAddr)
		wc.box2.SetText(c.ProxyServerHost.LinkProxyKey)
		wc.box3.SetText(strconv.Itoa(c.ProxyMethod.Http.Port))
		wc.box4.SetText(strconv.Itoa(c.ProxyMethod.Socks.Port))

		wc.s4.SetCurrentIndex(int(c.Setting.LogLevel))
		wc.s5.SetChecked(c.Setting.LogStack)
		loger.SetLoggerLevel(c.Setting.LogLevel)
		loger.SetLoggerStack(c.Setting.LogStack)
		for i, one := range acfgIdSl {
			if one == wc.ctx.acfg.LastConfig {
				wc.c1.SetCurrentIndex(i)
				break
			}
		}
	} else {
		defaultDataFn()
		wc.c1.SetCurrentIndex(0)
	}

	wc.s2.ConnectClicked(func(checked bool) {
		proxyFn()
	})
	wc.i2.ConnectTriggered(func(checked bool) {
		proxyFn()
	})
	wc.c1.ConnectActivated2(func(text string) {
		wc.switchFn2(true)
		defer wc.switchFn2(false)
		wc.switchFn(true)
		defer wc.switchFn(false)
		c := wc.ctx.acfg.getConfig(text)
		wc.box1.SetText(c.ProxyServerHost.ProxyServerAddr)
		wc.box2.SetText(c.ProxyServerHost.LinkProxyKey)
		wc.box3.SetText(strconv.Itoa(c.ProxyMethod.Http.Port))
		wc.box4.SetText(strconv.Itoa(c.ProxyMethod.Socks.Port))
		wc.s4.SetCurrentIndex(int(c.Setting.LogLevel))
		wc.s5.SetChecked(c.Setting.LogStack)
		loger.SetLoggerLevel(c.Setting.LogLevel)
		loger.SetLoggerStack(c.Setting.LogStack)
		for i, one := range acfgIdSl {
			if one == wc.ctx.acfg.LastConfig {
				wc.c1.SetCurrentIndex(i)
				break
			}
		}
		err := wc.ctx.acfg.setLastConfig(text)
		if err != nil {
			loger.SetLogWarn(err)
		}
		wc.c1.SetCurrentText(text)
		defaulC2Fn(text)
	})
	wc.c2.ConnectActivated2(func(text string) {
		wc.switchFn2(true)
		defer wc.switchFn2(false)
		wc.switchFn(true)
		defer wc.switchFn(false)
		nowN := ""
		switch text {
		case scdCfgSl[0]:
			nowName, err := wc.ctx.acfg.cutConfig(wc.c1.CurrentText(), wc.box1.Text(), wc.box2.Text(), wc.box3.Text(), wc.box4.Text(), uint8(wc.s4.CurrentIndex()), wc.s5.IsChecked())
			if err != nil {
				loger.SetLogXY(err)
				return
			}
			nowN = nowName
		case scdCfgSl[1]:
			nowName, err := wc.ctx.acfg.setConfig2(wc.box1.Text(), wc.box2.Text(), wc.box3.Text(), wc.box4.Text(), uint8(wc.s4.CurrentIndex()), wc.s5.IsChecked())
			if err != nil {
				loger.SetLogXY(err)
				return
			}
			nowN = nowName
		case scdCfgSl[2]:
			err := wc.ctx.acfg.delConfig(wc.c1.CurrentText())
			if err != nil {
				loger.SetLogXY(err)
				return
			}
			nowN = appConfigNone
			defaultDataFn()
		}
		wc.c1.Clear()
		acfgIdSl = wc.ctx.acfg.getIdList()
		wc.c1.AddItems(acfgIdSl)
		wc.c1.SetCurrentText(nowN)
		defaulC2Fn(nowN)
	})

	testCh := make(chan uint8, 1)
	wc.s3.ConnectClicked(func(checked bool) {
		wc.s3.SetDisabled(true)
		defer wc.s3.SetDisabled(false)
		wc.switchFn2(true)
		defer wc.switchFn2(false)
		select {
		case testCh <- 1:
			go func() {
				rx, err := testP("127.0.0.1:" + wc.box3.Text())
				if err != nil {
					loger.SetLogXY("Test failure:", err)
				} else {
					loger.SetLogXY("Test succeed:", rx)
				}
				<-testCh
			}()
			return
		default:
			return
		}
	})
}
