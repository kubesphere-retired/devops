/*
Copyright 2018 The KubeSphere Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
