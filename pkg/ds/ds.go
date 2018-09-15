package ds

import (
	"kubesphere.io/devops/pkg/config"
	"kubesphere.io/devops/pkg/constants"
	"kubesphere.io/devops/pkg/db"
	"kubesphere.io/devops/pkg/gojenkins"
	"kubesphere.io/devops/pkg/logger"
)

type Ds struct {
	cfg     *config.Config
	Db      *db.Database
	Jenkins *gojenkins.Jenkins
}

func NewDs(cfg *config.Config) *Ds {
	s := &Ds{cfg: cfg}
	s.openDatabase()
	s.connectJenkins()
	return s
}

func (p *Ds) openDatabase() *Ds {
	db, err := db.OpenDatabase(p.cfg.Mysql)
	if err != nil {
		logger.Critical("failed to connect mysql")
		panic(err)
	}
	p.Db = db
	return p
}

func (p *Ds) connectJenkins() {
	jenkins := gojenkins.CreateJenkins(nil, p.cfg.Jenkins.Address, p.cfg.Jenkins.User, p.cfg.Jenkins.Password)
	jenkins, err := jenkins.Init()
	if err != nil {
		logger.Critical("failed to connect jenkins")
		panic(err)
	}
	p.Jenkins = jenkins
	globalRole, err := jenkins.GetGlobalRole(constants.JenkinsAllUserRoleName)
	if err != nil {
		logger.Critical("failed to get jenkins role")
		panic(err)
	}
	if globalRole == nil {
		_, err := jenkins.AddGlobalRole(constants.JenkinsAllUserRoleName, gojenkins.GlobalPermissionIds{
			GlobalRead: true,
		}, true)
		if err != nil {
			logger.Critical("failed to create jenkins global role")
			panic(err)
		}
	}

}
