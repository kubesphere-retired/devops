package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/koding/multiconfig"
	"kubesphere.io/devops/pkg/logger"
)

type Config struct {
	Log     LogConfig
	Mysql   MysqlConfig
	Jenkins JenkinsConfig
}

type LogConfig struct {
	Level string `default:"info"` // debug, info, warn, error, fatal
}

type MysqlConfig struct {
	Host     string `default:"kubesphere-db"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Database string `default:"kubesphere"`
}
type JenkinsConfig struct {
	Address  string `default:"http://jenkins.kubesphere.com/"`
	User     string `default:"magicsong"`
	Password string `default:"devops"`
}

func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func PrintUsage() {
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader("devopsphere")
	e.PrintEnvs(new(Config))
	fmt.Println("")
}

func GetFlagSet() *flag.FlagSet {
	flag.CommandLine.Usage = PrintUsage
	return flag.CommandLine
}

func ParseFlag() {
	GetFlagSet().Parse(os.Args[1:])
}

var profilingServerStarted = false

func LoadConf() *Config {
	ParseFlag()

	config := new(Config)
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader("devopsphere"))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		panic(err)
	}
	logger.SetLevelByString(config.Log.Level)
	logger.Info("LoadConf: %+v", config)

	return config
}
