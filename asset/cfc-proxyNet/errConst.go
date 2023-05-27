package main

import (
	"errors"
	"fmt"
)

var errGetUUIDBad = errors.New("get uuid bad")

var errParsingFailure = errors.New("parsing failure")

var errProxyAddrIsNil = errors.New("proxy addr is nil")

var errProxyPingAndTestIsFailure = errors.New("proxy ping and test is failure")

var errOnlyOnceStartFailurePleaseCheckAppStateOrNetPort = errors.New("only once start failure: please check app state or net port")

var errSocksHandshakeBad = errors.New("socks handshake bad")

var errSocksHandshakeBadAny = func(a ...any) error { return errors.New(errSocksHandshakeBad.Error() + fmt.Sprint(a...)) }

var errSocksHandshakeAddrBad = errors.New("socks handshake addr bad")

var errSocksHandshakeAddrBadAny = func(a ...any) error { return errors.New(errSocksHandshakeAddrBad.Error() + fmt.Sprint(a...)) }

var errSocksUdpDataPauseBad = errors.New("socks udp data pause is bad")
