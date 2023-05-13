package main

import (
	"encoding/json"
	"github.com/peakedshout/cfc-proxyNet/config"
	"github.com/peakedshout/go-CFC/loger"
	"github.com/peakedshout/go-CFC/tool"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"
)

const appConfigFileName = "CacheData.cpn"

const appConfigNone = "<none>"

type AppConfig struct {
	filePath string
	key      []byte
	rawKey   string
	lock     sync.Mutex

	AppName string
	Version string
	RawKey  string

	LastConfig string
	ConfigList map[string]config.Config
}

func newACfg(p, k string) *AppConfig {
	return &AppConfig{
		filePath: p,
		key:      []byte(k),
		rawKey:   k,
		lock:     sync.Mutex{},
	}
}

func (acfg *AppConfig) setFPath(p string) {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	acfg.filePath = p
}
func (acfg *AppConfig) setKey(key string) {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	acfg.key = []byte(key)
	acfg.rawKey = key
}

func (acfg *AppConfig) readData() ([]byte, error) {
	f, err := os.OpenFile(acfg.filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b []byte
	for {
		buf := make([]byte, 4096)
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		b = append(b, buf[:n]...)
	}
	return b, err
}

func (acfg *AppConfig) read() error {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	b, err := acfg.readData()
	if err != nil {
		loger.SetLogDebug(err)
		return errParsingFailure

	}
	if len(b) == 0 {
		acfg.Version = appVersion
		acfg.AppName = appName
		acfg.RawKey = acfg.rawKey
		acfg.LastConfig = ""
		acfg.ConfigList = make(map[string]config.Config)
		err = acfg.save()
		if err != nil {
			loger.SetLogDebug(err)
			return errParsingFailure
		}
		return nil
	}
	bs, err := tool.Decrypt(b, acfg.key)
	if err != nil {
		loger.SetLogDebug(err)
		return errParsingFailure
	}

	err = json.Unmarshal(bs, acfg)
	if err != nil {
		loger.SetLogDebug(err)
		return errParsingFailure
	}

	if acfg.RawKey != acfg.rawKey || acfg.AppName != appName || acfg.Version > appVersion {
		err = errParsingFailure
		loger.SetLogDebug(err)
		return errParsingFailure
	}
	return nil
}

func (acfg *AppConfig) save() error {

	b, err := json.Marshal(acfg)
	if err != nil {
		return err
	}

	bs, err := tool.Encrypt(b, acfg.key)
	if err != nil {
		return err
	}

	f, err := os.Create(acfg.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (acfg *AppConfig) cutConfig(nowName, paddr, pkey, phport, psport string, ll uint8, ls bool) (string, error) {
	if paddr == "" {
		return "", errProxyAddrIsNil
	}
	php1, err := strconv.Atoi(phport)
	if err != nil {
		return "", err
	}
	psp1, err := strconv.Atoi(psport)
	if err != nil {
		return "", err
	}
	if nowName == appConfigNone || nowName == "" {
		for i := 0; ; i++ {
			id := paddr
			if i != 0 {
				id += "(" + strconv.Itoa(i) + ")"
			}
			_, ok := acfg.ConfigList[id]
			if ok {
				continue
			} else {
				nowName = id
				break
			}
		}
	}

	c := acfg.ConfigList[nowName]

	c.ProxyServerHost.ProxyServerAddr = paddr
	c.ProxyServerHost.LinkProxyKey = pkey
	c.ProxyMethod.Http.Host = "127.0.0.1"
	c.ProxyMethod.Http.Port = php1
	c.ProxyMethod.Https.Host = "127.0.0.1"
	c.ProxyMethod.Https.Port = php1
	c.ProxyMethod.Socks.Host = "127.0.0.1"
	c.ProxyMethod.Socks.Port = psp1
	acfg.ConfigList[nowName] = c
	c.Setting.LogLevel = ll
	c.Setting.LogStack = ls
	err = acfg.save()
	if err != nil {
		return "", err
	}
	return nowName, nil
}

func (acfg *AppConfig) setConfig2(paddr, pkey, phport, psport string, ll uint8, ls bool) (string, error) {
	if paddr == "" {
		return "", errProxyAddrIsNil
	}
	php1, err := strconv.Atoi(phport)
	if err != nil {
		return "", err
	}
	psp1, err := strconv.Atoi(psport)
	if err != nil {
		return "", err
	}

	c := config.Config{
		ProxyServerHost: config.ProxyServerHostConfig{
			ProxyServerAddr: paddr,
			LinkProxyKey:    pkey,
		},
		ProxyMethod: config.ProxyMethodConfig{
			Http: struct {
				Host string `json:"Host"`
				Port int    `json:"Port"`
			}{
				Host: "127.0.0.1",
				Port: php1,
			},
			Https: struct {
				Host string `json:"Host"`
				Port int    `json:"Port"`
			}{
				Host: "127.0.0.1",
				Port: php1,
			},
			Socks: struct {
				Host string `json:"Host"`
				Port int    `json:"Port"`
			}{
				Host: "127.0.0.1",
				Port: psp1,
			},
		},
		Setting: config.SettingConfig{
			ReLinkTime: "",
			LogLevel:   ll,
			LogStack:   ls,
		},
	}
	return acfg.setConfig(c)
}

func (acfg *AppConfig) setConfig(c config.Config) (string, error) {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()

	var id string
	for i := 0; ; i++ {
		id = c.ProxyServerHost.ProxyServerAddr
		if i != 0 {
			id += "(" + strconv.Itoa(i) + ")"
		}
		_, ok := acfg.ConfigList[id]
		if ok {
			continue
		} else {
			acfg.ConfigList[id] = c
			acfg.LastConfig = id
			break
		}
	}
	err := acfg.save()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (acfg *AppConfig) getIdList() []string {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	var sl []string
	for k := range acfg.ConfigList {
		sl = append(sl, k)
	}
	sort.Strings(sl)
	sl = append([]string{appConfigNone}, sl...)
	return sl
}

func (acfg *AppConfig) getConfig(id string) config.Config {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	if id == appConfigNone {
		return config.Config{
			Setting: config.SettingConfig{
				LogLevel: loger.LogLevelWarn,
				LogStack: false,
			},
		}
	}
	return acfg.ConfigList[id]
}
func (acfg *AppConfig) delConfig(id string) error {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	delete(acfg.ConfigList, id)
	return acfg.save()
}
func (acfg *AppConfig) setLastConfig(id string) error {
	acfg.lock.Lock()
	defer acfg.lock.Unlock()
	acfg.LastConfig = id
	return acfg.save()
}
