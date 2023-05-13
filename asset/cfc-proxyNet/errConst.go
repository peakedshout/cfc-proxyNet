package main

import "errors"

var errGetUUIDBad = errors.New("get uuid bad")

var errParsingFailure = errors.New("parsing failure")

var errProxyAddrIsNil = errors.New("proxy addr is nil")

var errProxyPingAndTestIsFailure = errors.New("proxy ping and test is failure")

var errOnlyOnceStartFailurePleaseCheckAppStateOrNetPort = errors.New("only once start failure: please check app state or net port")
