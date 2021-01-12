package config

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"linac"
	"log"
	"net/http"
	"strconv"
	"sync"
	"testing"
)

var once sync.Once
var (
	_addr = "127.0.0.1:9527"
	_port = ":9527"
)

func TestClientNew(t *testing.T) {
	once.Do(server)
	initConf()
	if _, err := New(); err != nil {
		t.Errorf("client.New() error(%v)", err)
		t.FailNow()
	}
}

func TestCheckVersion(t *testing.T) {
	once.Do(server)
	c := initConf()
	ver, err := c.checkVersion(&version{ver: _unknownVersion})
	if err != nil && ver.ver == _unknownVersion {
		t.Errorf("client.checkVersion() error(%v) ver(%d)", err, ver)
		t.FailNow()
	}
}

func TestSynchro(t *testing.T) {
	once.Do(server)
	c := initConf()
	ver := &version{ver: _unknownVersion}
	if err := c.synchro(ver); err != nil {
		t.Errorf("client.downloda() error(%v) ", err)
		t.FailNow()
	}
}

func TestValue(t *testing.T) {
	once.Do(server)
	c := initConf()
	ver := &version{ver: _unknownVersion}
	if err := c.synchro(ver); err != nil {
		t.Errorf("client.downloda() error(%v) ", err)
		t.FailNow()
	}
	if conf, ok := c.Value("linac.micro.users.timeout"); !ok {
		t.Errorf("client.Value() error ")
		t.FailNow()
	} else if conf != "10" {
		t.Errorf("client.Value() error, 10!=%s ", conf)
		t.FailNow()
	}
}

func TestConfigs(t *testing.T) {
	once.Do(server)
	c := initConf()
	ver := &version{ver: _unknownVersion}
	if err := c.synchro(ver); err != nil {
		t.Errorf("client.downloda() error(%v)", err)
		t.FailNow()
	}

	if confs, ok := c.Configs(); !ok {
		t.Errorf("client.Configs() error")
		t.FailNow()
	} else if len(confs) == 0 {
		t.Errorf("client.Configs() error")
		t.Error(confs)
		t.FailNow()
	}
}

func initConf() (c *Client) {
	conf.Addr = _addr
	conf.Hostname = "testHost"
	conf.Path = "./test_data"
	conf.AppID = "linac.micro"
	conf.AppName = "users-service"
	conf.Version = "1.0.1"
	conf.DeployEnv = "dev"
	conf.Zone = "wh001"
	conf.Region = "wh"
	conf.Token = "45338e440bdc11e880ce02420a0a0204"
	c = &Client{
		httpCli: &http.Client{Timeout: _httpTimeout},
		event:   make(chan string, 10),
	}
	return
}

func server() {
	mux := handler()
	go func() {
		log.Fatal(http.ListenAndServe(_port, mux))
	}()
}

func handler() *http.ServeMux {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("test")
	})
	mux.HandleFunc("/config/check", checkConfig)
	mux.HandleFunc("/config/get", getConfig)
	return mux
}

func checkConfig(w http.ResponseWriter, r *http.Request) {
	sver := r.URL.Query().Get("version")
	if sver == "" {
		log.Print(r.URL.Query())
	}
	if ver, err := strconv.Atoi(sver); err != nil {
		log.Print(err)
	} else if int64(ver) < NowVer.ver {
		json.NewEncoder(w).Encode(cVerResp{
			Code:    _codeShouldUpdate,
			Message: "",
			Data:    NowVer,
		})
	} else {
		json.NewEncoder(w).Encode(cVerResp{
			Code:    _codeNoNeedUpdate,
			Message: "",
			Data:    NowVer,
		})
	}
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	if content, err := json.Marshal(Vers); err != nil {
		log.Print(err)
	} else {
		mh := md5.Sum(content)
		m5bs := hex.EncodeToString(mh[:])
		json.NewEncoder(w).Encode(dVerResp{
			Code:    0,
			Message: "",
			Data: &struct {
				Version int    `json:"version"`
				Content string `json:"content"`
				MD5     string `json:"md5"`
			}{
				Version: 1024,
				Content: linac.BytesToString(content),
				MD5:     m5bs,
			},
		})
	}
}
