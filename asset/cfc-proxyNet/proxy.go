package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/peakedshout/cfc-proxyNet/settings"
	"github.com/peakedshout/go-CFC/client"
	"github.com/peakedshout/go-CFC/loger"
	"github.com/peakedshout/go-CFC/tool"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type proxyCtx struct {
	lnAddr    string
	proxyAddr string
	proxyKey  string

	ln      net.Listener
	closer  sync.Once
	stop    chan error
	connMap sync.Map
	connId  int64
}

func (ctx *proxyCtx) close(err error) {
	ctx.closer.Do(func() {
		ctx.ln.Close()
		settings.Close()
		ctx.rangeClose()
		ctx.stop <- err
	})
}
func (ctx *proxyCtx) wait() error {
	return <-ctx.stop
}

func runProxy(lnAddr, proxyAddr, proxyKey string) (*proxyCtx, error) {
	c := &proxyCtx{
		lnAddr:    lnAddr,
		proxyAddr: proxyAddr,
		proxyKey:  proxyKey,
		ln:        nil,
		closer:    sync.Once{},
		stop:      make(chan error, 1),
		connMap:   sync.Map{},
		connId:    0,
	}
	_, err := net.ResolveTCPAddr("", c.proxyAddr)
	if err != nil {
		return nil, err
	}
	lAddr, err := net.ResolveTCPAddr("", c.lnAddr)
	if err != nil {
		return nil, err
	}
	if len(c.proxyKey) != 32 {
		return nil, tool.ErrKeyIsNot32Bytes
	}

	err = settings.Init(lnAddr)
	if err != nil {
		return nil, err
	}

	ln, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		return nil, err
	}
	c.ln = ln

	go func(ctx2 *proxyCtx) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				ctx2.close(err)
				return
			}
			go ctx2.handleConn(conn)
		}
	}(c)
	return c, err
}
func (ctx *proxyCtx) handleConn(conn net.Conn) {
	defer conn.Close()

	id := ctx.setConn(conn)
	defer ctx.delConn(id)

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		loger.SetLogWarn(err)
		return
	}
	loger.SetLogXY(req.Host, "vpn:", "~Link~")
	defer loger.SetLogXY(req.Host, "vpn:", "~Over~")
	otherMethod := false
	bs := new(bytes.Buffer)
	if req.Method == http.MethodConnect {
		loger.SetLogXY(req.Host, "vpn:", "x the first is", req.Method, ", is", http.MethodConnect, "can start!")
	} else {
		loger.SetLogXY(req.Host, "vpn:", "x the first is", req.Method, ", not", http.MethodConnect, "cant start! try go it!")
		err = req.Write(bs)
		if err != nil {
			loger.SetLogXY(req.Host, "vpn:", err)
			return
		}
		otherMethod = true
	}
	host := req.Host
	sl := strings.Split(req.Host, ":")
	if len(sl) == 1 {
		host = req.Host + ":80"
	} else if len(sl) == 2 {
		host = req.Host
	} else {
		loger.SetLogXY(req.Host, "vpn:", req.Host, "is bad addr")
		return
	}

	lc, err := client.LinkOtherConn(client.LinkConnReq{
		CopyConn:  conn,
		ConnType:  tool.LinkConnTypeTCP,
		ConnAddr:  host,
		ProxyAddr: ctx.proxyAddr,
		ProxyKey:  ctx.proxyKey,
	})
	if err != nil {
		loger.SetLogXY(req.Host, "vpn:", err)
		return
	}
	if otherMethod {
		err = lc.WriteToLinkConn(bs.Bytes())
		if err != nil {
			loger.SetLogXY(req.Host, "vpn:", err)
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
	err = lc.Wait()
	if err != nil {
		loger.SetLogXY(req.Host, "vpn:", err)
		return
	}
	return
}

func (ctx *proxyCtx) setConn(conn net.Conn) int64 {
	id := atomic.AddInt64(&ctx.connId, 1)
	ctx.connMap.Store(id, conn)
	return id
}
func (ctx *proxyCtx) delConn(id int64) {
	ctx.connMap.Delete(id)
}
func (ctx *proxyCtx) rangeClose() {
	ctx.connMap.Range(func(key, value any) bool {
		conn := value.(net.Conn)
		conn.Close()
		return true
	})
}

func testProxy(lnAddr string) (string, error) {
	testAddrList := []string{"https://www.google.com", "https://baidu.com"}
	var tL []string
	for i, one := range testAddrList {
		t := time.Now()
		err := testProxy2(one, lnAddr)
		if err != nil {
			return "", err
		}
		t1 := time.Now().Sub(t)
		tr := fmt.Sprintf("test %d:%s--%s", i+1, one, t1.String())
		tL = append(tL, tr)
	}
	return strings.Join(tL, "\n"), nil
}
func testProxy2(rAddr, lnAddr string) error {
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
