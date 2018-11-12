package test_config

import (
	"os"
	"testing"

	"kubesphere.io/devops/pkg/config"
	"kubesphere.io/devops/pkg/db"
	"kubesphere.io/devops/pkg/logger"
)

type DbTestConfig struct {
	openDbUnitTests string
	EnvConfig       *config.Config
}

func NewDbTestConfig() DbTestConfig {
	tc := DbTestConfig{
		openDbUnitTests: os.Getenv("KS_DEVOPS_DB_UNIT_TEST"),
		EnvConfig:       config.LoadConf(),
	}
	return tc
}

func (tc DbTestConfig) GetDatabaseConn() *db.Database {
	if tc.openDbUnitTests == "1" {
		d, err := db.OpenDatabase(tc.EnvConfig.Mysql)
		if err != nil {
			logger.Critical( "failed to open database %+v", tc.EnvConfig.Mysql)
		}
		return d
	}
	return nil
}

func (tc DbTestConfig) CheckDbUnitTest(t *testing.T) {
	if tc.openDbUnitTests != "1" {
		t.Skipf("if you want run unit tests with db,set KS_DEVOPS_DB_UNIT_TEST=1")
	}
}


type JenkinsTestConfig struct {
	openJenkinsUnitTests string
	EnvConfig       *config.Config
}

func NewJenkinsTestConfig() JenkinsTestConfig {
	tc := JenkinsTestConfig{
		openJenkinsUnitTests: os.Getenv("KS_DEVOPS_JK_UNIT_TEST"),
		EnvConfig:       config.LoadConf(),
	}
	return tc
}


func (tc JenkinsTestConfig) CheckJenkinsUnitTest(t *testing.T) {
	if tc.openJenkinsUnitTests != "1" {
		t.Skipf("if you want run unit tests with db,set KS_DEVOPS_JK_UNIT_TEST=1")
	}
}

