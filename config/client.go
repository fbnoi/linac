package config

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync/atomic"
	"time"

	"linac"
	"linac/config/env"
)

var (
	_codeOk          = 0
	_codeModified    = -301
	_codeNotModified = -201

	_retryInterval        = time.Second
	_httpTimeout          = time.Second * 60
	_unknownVersion int64 = -1

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
	Version   string // 当前应用版本
	Since     string // 应用挂载时间
	Addr      string // 配置中心主机
	Path      string // 配置文件路径
	Token     string // 获取配置请求需要的token
	Customize string // 配置自定义字段

	// 应用环境
	AppName   string // 应用名称
	AppID     string // 应用ID
	DeployEnv string // 运行环境
	Hostname  string // 挂载主机名
	Zone      string // 挂载可用区
	Region    string // 挂载地区
}

// Value config value
type Value struct {
	CID    int64  `json:"cid"`
	Name   string `json:"name"`
	Config string `json:"config"`
}

// Client is config client.
type Client struct {
	ver   *version
	confs atomic.Value
	event chan string

	httpCli *http.Client

	local atomic.Value // NOTE: struct: map[string]interface{}
}

func init() {
	// os env
	conf.Version = os.Getenv("CONF_VERSION")
	conf.Addr = os.Getenv("CONF_HOST")
	conf.Hostname = os.Getenv("CONF_HOSTNAME")
	conf.Path = os.Getenv("CONF_PATH")
	conf.DeployEnv = os.Getenv("CONF_ENV")
	conf.Token = os.Getenv("CONF_TOKEN")
	conf.Region = os.Getenv("REGION")
	conf.Zone = os.Getenv("ZONE")
	conf.AppID = os.Getenv("APP_ID")

	// flags
	flag.StringVar(&conf.Version, "conf_version", conf.Version, `app version.`)
	flag.StringVar(&conf.Addr, "conf_host", conf.Addr, `config center api host.`)
	flag.StringVar(&conf.Path, "conf_path", conf.Path, `config file path.`)
	flag.StringVar(&conf.DeployEnv, "conf_env", conf.DeployEnv, `config Env.`)
	flag.StringVar(&conf.Token, "conf_token", conf.Token, `config Token.`)

	// env set
	conf.Region = env.Region
	conf.Zone = env.Zone
	conf.AppID = env.AppID
	conf.AppName = env.AppName
	conf.DeployEnv = env.DeployEnv
	conf.Hostname = env.Hostname

	conf.Since = time.Now().String()
}

// New 创建并返回一个新的配置客户端
func New() (cli *Client, err error) {
	cli = &Client{
		httpCli: &http.Client{Timeout: _httpTimeout},
		event:   make(chan string, 10),
	}

	if conf.AppName != "" && conf.Hostname != "" && conf.Path != "" && conf.Addr != "" && conf.Token != "" &&
		conf.Version != "" && conf.AppID != "" && conf.Since != "" && conf.DeployEnv != "" && conf.Zone != "" && conf.Region != "" {
		if err = cli.init(); err != nil {
			return nil, err
		}
		go cli.updateProc()
		return
	}
	err = fmt.Errorf("at least one params is empty. app=%s, version=%s, hostname=%s, path=%s, Token=%s, DeployEnv=%s, appID=%s, since=%s, zone=%s, region=%s",
		conf.AppName, conf.Version, conf.Addr, conf.Path, conf.Token, conf.DeployEnv, conf.AppID, conf.Since, conf.Zone, conf.Region)
	return
}

// Value 返回 config 值
func (c *Client) Value(key string) (val string, ok bool) {
	var (
		m map[string]*Value
		n *Value
	)
	if m, ok = c.confs.Load().(map[string]*Value); !ok {
		return
	}
	if n, ok = m[key]; !ok {
		return
	}
	val = n.Config
	return
}

// Configs 获取所有的配置
func (c *Client) Configs() (confs []*Value, ok bool) {
	var (
		m map[string]*Value
	)
	if m, ok = c.confs.Load().(map[string]*Value); !ok {
		return
	}
	for _, v := range m {
		if v.CID == 0 {
			continue
		}
		confs = append(confs, v)
	}
	return
}

// Event 获取配置更新事件
func (c *Client) Event() <-chan string {
	return c.event
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

// init 初始化
func (c *Client) init() (err error) {
	var v *version
	c.confs.Store(make(map[string]*Value))
	if v, err = c.checkVersion(&version{ver: _unknownVersion}); err != nil {
		fmt.Printf("get remote version error(%v)\n", err)
		return
	}
	for i := 0; i < 3; i++ {
		if v.ver == _unknownVersion {
			fmt.Println("get null version")
			return
		}
		if err = c.synchro(v); err == nil {
			return
		}
		fmt.Printf("retry times: %d, c.synchro(%d) error(%v)\n", v, i, err)
		time.Sleep(_retryInterval)
	}
	return
}

// updateProc 更新配置进程
func (c *Client) updateProc() (err error) {
	var ver *version
	for {
		time.Sleep(_retryInterval)
		if ver, err = c.checkVersion(c.ver); err != nil {
			log.Fatalf("c.checkVersion(%d) error(%v)", c.ver, err)
			continue
		} else if ver.ver == c.ver.ver {
			continue
		}
		if err = c.synchro(ver); err != nil {
			log.Fatalf("c.synchro(%d) error(%s)", ver, err)
			continue
		}
	}
}

// synchro 同步远端配置并保存到本地文件中
func (c *Client) synchro(ver *version) (err error) {
	var (
		bs             []byte
		tmp            []*Value
		oConfs, nConfs map[string]*Value
		ok             bool
	)
	if bs, err = c.download(ver); err != nil {
		return
	}

	// 合并配置并写入配置文件
	if err = json.Unmarshal(bs, &tmp); err != nil {
		return
	}
	nConfs = make(map[string]*Value)
	if oConfs, ok = c.confs.Load().(map[string]*Value); ok {
		for key, value := range oConfs {
			nConfs[key] = value
		}
	}
	for _, v := range tmp {
		if err = ioutil.WriteFile(path.Join(conf.Path, v.Name), linac.StringToBytes(v.Config), 0644); err != nil {
			return
		}
		nConfs[v.Name] = v
	}

	// 更新当前配置版本号，保存配置
	c.ver = ver
	c.confs.Store(nConfs)

	// 发布事件
	for _, v := range tmp {
		c.event <- v.Name
	}
	return
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

// download 下载远端配置
func (c *Client) download(reqVer *version) (confs []byte, err error) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		rbs  []byte
	)

	if url, err = c.makeURL(_apiCheck, reqVer); err != nil {
		err = fmt.Errorf("download(): make url error:" + err.Error())
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
		err = fmt.Errorf("download(): http error url(%s) status: %d", url, resp.StatusCode)
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
			err = fmt.Errorf("download(): response error: %v", v)
			return
		}
		mh := md5.Sum(linac.StringToBytes(v.Data.Content))
		m5bs := hex.EncodeToString(mh[:])
		if m5bs != v.Data.MD5 {
			err = fmt.Errorf("download(): md5 mismatch, local: %s, remote: %s", m5bs, v.Data.MD5)
		}
		confs = linac.StringToBytes(v.Data.Content)
	default:
		err = fmt.Errorf("download(): response error: %v", v)
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
