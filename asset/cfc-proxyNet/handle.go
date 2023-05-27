package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/peakedshout/cfc-proxyNet/settings"
	"github.com/peakedshout/go-CFC/loger"
	"github.com/peakedshout/go-CFC/tool"
	"github.com/peakedshout/go-socks/client"
	"github.com/peakedshout/go-socks/server"
	"github.com/peakedshout/go-socks/share"
	"net"
	"net/http"
	"strings"
	"time"
)

type serverCtx struct {
	httpServer  net.Listener
	socksServer *server.SocksServer

	lnHAddr, lnSAddr    string
	proxyAddr, proxyKey string

	connMap *share.ConnMap

	*tool.CloseWaiter
}

func checkParameters(lnHAddr, lnSAddr, proxyAddr, proxyKey string) error {
	if lnHAddr == "" {
		return errors.New("http/https port is nill")
	}
	if lnSAddr == "" {
		return errors.New("socks port is nill")
	}
	if proxyAddr == "" {
		return errors.New("proxy addr is nill")
	}
	if len(proxyKey) != 32 {
		return errors.New("proxy key is not 32 bytes")
	}
	return nil
}

func newServerCtx(lnHAddr, lnSAddr, proxyAddr, proxyKey string) (*serverCtx, error) {
	err := checkParameters(lnHAddr, lnSAddr, proxyAddr, proxyKey)
	if err != nil {
		return nil, err
	}

	sc := &serverCtx{
		httpServer:  nil,
		socksServer: nil,
		lnHAddr:     lnHAddr,
		lnSAddr:     lnSAddr,
		proxyAddr:   proxyAddr,
		proxyKey:    proxyKey,
		connMap:     share.NewConnMap(),
		CloseWaiter: tool.NewCloseWaiter(),
	}
	sc.AddCloseFn(func() {
		if sc.httpServer != nil {
			sc.httpServer.Close()
		}
		if sc.socksServer != nil {
			sc.socksServer.Close(nil)
		}
		sc.connMap.Disable()
		sc.connMap.RangeConn(func(conn net.Conn) bool {
			conn.Close()
			return true
		})
	})
	err = sc.newSocksServer()
	if err != nil {
		return nil, err
	}
	err = sc.newHttpServer()
	if err != nil {
		return nil, err
	}
	return sc, nil
}
func (sc *serverCtx) newSocksServer() error {
	config := &server.SocksServerConfig{
		TlnAddr: sc.lnSAddr,
		SocksAuthCb: server.SocksAuthCb{
			Socks5AuthNOAUTH: true,
		},
		RelayConfig: &server.SocksRelayConfig{
			Addr:        sc.proxyAddr,
			RawKey:      sc.proxyKey,
			KeepEncrypt: true,
		},
		VersionSwitch: share.SocksVersionSwitch{
			SwitchSocksVersion4: true,
			SwitchSocksVersion5: true,
		},
		CMDSwitch: share.SocksCMDSwitch{
			SwitchCMDCONNECT:      true,
			SwitchCMDBIND:         true,
			SwitchCMDUDPASSOCIATE: true,
		},
	}
	ss, err := server.NewSocksServer(config)
	if err != nil {
		return err
	}
	sc.socksServer = ss
	return nil
}

func (sc *serverCtx) newHttpServer() error {
	lAddr, err := net.ResolveTCPAddr("", sc.lnHAddr)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		return err
	}
	sc.httpServer = ln
	err = settings.Init(sc.lnHAddr)
	if err != nil {
		return err
	}
	go func() {
		defer settings.Close()
		for {
			conn, err := sc.httpServer.Accept()
			if err != nil {
				loger.SetLogWarn(err)
				return
			}
			go sc.handleHttpTcpConn(conn)
		}
	}()
	return nil
}

func (sc *serverCtx) handleHttpTcpConn(conn net.Conn) {
	defer conn.Close()
	if !sc.connMap.SetConn(conn) {
		return
	}
	defer sc.connMap.DelConn(conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		loger.SetLogInfo(err)
		return
	}
	loger.SetLogInfo("http:", req.Host, "~Link~")
	defer loger.SetLogInfo("http:", req.Host, "~Over~")
	otherMethod := false
	bs := new(bytes.Buffer)
	if req.Method != http.MethodConnect {
		err = req.Write(bs)
		if err != nil {
			loger.SetLogInfo("http:", req.Host, err)
			return
		}
		otherMethod = true
	}
	host := req.Host
	_, _, err = net.SplitHostPort(host)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			host += ":80"
		} else {
			loger.SetLogInfo("http:", req.Host, err)
			return
		}
	}

	dr, err := client.NewSocks5ConnCONNECT(sc.lnSAddr, &client.Socks5Auth{Socks5AuthNOAUTH: true}, nil)
	if err != nil {
		loger.SetLogInfo(err)
		return
	}
	conn2, err := dr.Dial("tcp", host)
	if err != nil {
		loger.SetLogInfo(err)
		return
	}
	if otherMethod {
		_, err = conn2.Write(bs.Bytes())
		if err != nil {
			loger.SetLogWarn("http:", req.Host, err)
			return
		}
	} else {
		resp := http.Response{
			Status:     "200 Connection Established",
			StatusCode: 200,
			Proto:      req.Proto,
			ProtoMajor: req.ProtoMajor,
			ProtoMinor: req.ProtoMinor,
		}
		resp.Write(conn)
	}
	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				loger.SetLogInfo(err)
				return
			}
			_, err = conn2.Write(buf[:n])
			if err != nil {
				loger.SetLogInfo(err)
				return
			}
		}
	}()
	for {
		buf := make([]byte, 4096)
		n, err := conn2.Read(buf)
		if err != nil {
			loger.SetLogInfo(err)
			return
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			loger.SetLogInfo(err)
			return
		}
	}
}

func testP(lnAddr string) (string, error) {
	testAddrList := []string{"https://www.google.com", "https://baidu.com"}
	var tL []string
	for i, one := range testAddrList {
		t := time.Now()
		err := testP2(one, lnAddr)
		if err != nil {
			return "", err
		}
		t1 := time.Now().Sub(t)
		tr := fmt.Sprintf("test %d:%s--%s", i+1, one, t1.String())
		tL = append(tL, tr)
	}
	return strings.Join(tL, "\n"), nil
}
func testP2(rAddr, lnAddr string) error {
	req, err := http.NewRequest(http.MethodConnect, rAddr, nil)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", lnAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = req.Write(conn)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errProxyPingAndTestIsFailure
	}
	req, err = http.NewRequest(http.MethodGet, rAddr, nil)
	if err != nil {
		return err
	}
	err = req.Write(conn)
	if err != nil {
		return err
	}
	resp, err = http.ReadResponse(reader, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errProxyPingAndTestIsFailure
	}
	return nil
}
