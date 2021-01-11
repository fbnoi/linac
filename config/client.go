package config

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linac"
	"linac/config/env"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
)

var (
	_codeOk          = 0
	_codeModified    = -301
	_codeNotModified = -201
)

var (
	conf config
)

// config api
var (
	_apiGet     = "http://%s/config/get?%s"
	_apiCheck   = "http://%s/config/check?%s"
	_apiCreate  = "http://%s/config/create"
	_apiUpdate  = "http://%s/config/update"
	_apiConfIng = "http://%s/config/config/ing?%s"
)

// cVerResp 验证配置版本响应
type cVerResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *version `json:"data"`
}

// dVerResp 下载的配置版本响应
type dVerResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Version int    `json:"version"`
		Content string `json:"content"`
		MD5     string `json:"md5"`
	} `json:"ver"`
}

type version struct {
	ver   int64
	diffs []int64
}

// 应用初始配置选项
type config struct {
	Path      string // 配置文件路径
	Token     string // 获取配置请求需要的token
	Version   string // 当前应用版本
	Since     string // 应用挂载时间
	Addr      string // 配置中心主机
	Customize string // 配置自定义字段

	// 应用环境
	AppID     string // 应用ID
	DeployEnv string // 运行环境
	Hostname  string // 挂载主机名
	Zone      string // 挂载可用区
	Region    string // 挂载地区
}

// Value config value
type Value struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// Client is config client.
type Client struct {
	ver   *version
	data  atomic.Value
	event chan string

	httpCli *http.Client

	watchFile map[string]struct{}
	watchAll  bool

	local atomic.Value // NOTE: struct: map[string]interface{}
}

// Value 返回 config 值
func (c *Client) Value(key string) (val string, ok bool) {
	var (
		m map[string]*Value
		n *Value
	)
	if m, ok = c.data.Load().(map[string]*Value); !ok {
		return
	}
	if n, ok = m[key]; !ok {
		return
	}
	val = n.Config
	return
}

// Path 返回文件配置的路径
func (c *Client) Path() string {
	return conf.Path
}

// SetLocal 动态设置本地值
func (c *Client) SetLocal(value map[string]interface{}) {
	c.local.Store(value)
}

// GetLocal 获取本地配置
func (c *Client) GetLocal() (m map[string]interface{}, ok bool) {
	m, ok = c.local.Load().(map[string]interface{})
	return
}

// Watch 观察文件的变化.
func (c *Client) Watch(filename ...string) {
	if c.watchFile == nil {
		c.watchFile = map[string]struct{}{}
	}
	for _, f := range filename {
		c.watchFile[f] = struct{}{}
	}
}

// checkVersion 检查配置版本
func (c *Client) checkVersion(reqVer *version) (ver *version, err error) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		rbs  []byte
	)
	if url, err = c.makeURL(_apiCheck, reqVer); err != nil {
		err = fmt.Errorf("checkVersion(): make url error:" + err.Error())
		return
	}

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}

	if resp, err = c.httpCli.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("checkVersion(): http error url(%s) status: %d", url, resp.StatusCode)
		return
	}

	if rbs, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	v := &cVerResp{}
	if err = json.Unmarshal(rbs, v); err != nil {
		return
	}

	switch v.Code {
	case _codeModified:
		if v.Data == nil {
			err = fmt.Errorf("checkVersion(): response error: %v", v)
			return
		}
		ver = v.Data
	case _codeNotModified:
		ver = reqVer
	default:
		err = fmt.Errorf("checkVersion(): response error: %v", v)
	}
	return
}

func (c *Client) getConfig(reqVer *version) (confs string, err error) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		rbs  []byte
	)

	if url, err = c.makeURL(_apiCheck, reqVer); err != nil {
		err = fmt.Errorf("getConfig(): make url error:" + err.Error())
		return
	}

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}

	if resp, err = c.httpCli.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("getConfig(): http error url(%s) status: %d", url, resp.StatusCode)
		return
	}

	if rbs, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	v := &dVerResp{}
	if err = json.Unmarshal(rbs, v); err != nil {
		return
	}

	switch v.Code {
	case _codeOk:
		if v.Data == nil {
			err = fmt.Errorf("getConfig(): response error: %v", v)
			return
		}
		mh := md5.Sum(linac.StringToBytes(v.Data.Content))
		m5bs := hex.EncodeToString(mh[:])
		if m5bs != v.Data.MD5 {
			err = fmt.Errorf("getConfig(): md5 mismatch, local: %s, remote: %s", m5bs, v.Data.MD5)
		}
		confs = v.Data.Content
	default:
		err = fmt.Errorf("getConfig(): response error: %v", v)
	}
	return
}

func (c *Client) makeURL(api string, ver *version) (query string, err error) {
	var ids []byte
	params := url.Values{}
	// service
	params.Set("service", service())
	params.Set("hostname", env.Hostname)
	params.Set("build", conf.Version)
	params.Set("version", fmt.Sprint(ver.ver))
	if ids, err = json.Marshal(ver.diffs); err != nil {
		return
	}
	params.Set("ids", string(ids))
	params.Set("ip", localIP())
	params.Set("token", conf.Token)
	params.Set("since", conf.Since)
	params.Set("customize", conf.Customize)
	// api
	query = fmt.Sprintf(api, conf.Addr, params.Encode())
	return
}

func service() string {
	return fmt.Sprintf("%s@%s@%s@%s", env.AppID, env.DeployEnv, env.Zone, env.Region)
}

// 获取挂载主机IP
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
