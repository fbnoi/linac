package env

import (
	"flag"
	"os"
)

// 几种部署环境.
const (
	DeployEnvDev  = "dev"
	DeployEnvTEST = "test"
	DeployEnvProd = "prod"
)

// env 默认值
const (
	_defaultRegion    = "wuhan"
	_defaultZone      = "wh001"
	_defaultDeployEnv = DeployEnvDev
)

// env 环境配置.
var (
	// Region 地区
	Region string
	// Zone 可用域
	Zone string
	// Hostname 主机名
	Hostname string
	// DeployEnv 部署环境
	DeployEnv string
	// IP 服务IP
	IP string
	// AppID 服务ID
	AppID string
	// AppID 服务名
	AppName string
)

// app 默认值
var (
	_defaultHTTPPort = "8080"
)

// app 配置项.
var (
	// HTTPPort app listen http port.
	HTTPPort string
)

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}
	IP = os.Getenv("PODIP")
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&Region, "region", defaultString("REGION", _defaultRegion), "avaliable region. or use REGION env variable, value: wuhan etc.")
	fs.StringVar(&Zone, "zone", defaultString("ZONE", _defaultZone), "avaliable zone. or use ZONE env variable, value: wh001/wh002 etc.")
	fs.StringVar(&DeployEnv, "deploy.env", defaultString("DEPLOY_ENV", _defaultDeployEnv), "deploy env. or use DEPLOY_ENV env variable, value: dev/test/prod etc.")
	fs.StringVar(&AppID, "appid", os.Getenv("APP_ID"), "appid is global unique application id, register by service tree. or use APP_ID env variable.")
	fs.StringVar(&AppName, "appname", os.Getenv("APP_NAME"), "AppName is the service name.")
	// app
	fs.StringVar(&HTTPPort, "http.port", defaultString("DISCOVERY_HTTP_PORT", _defaultHTTPPort), "app listen http port, default: 8080")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}
