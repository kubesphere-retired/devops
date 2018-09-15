package main

import (
	"kubesphere.io/devops/pkg/config"
	"kubesphere.io/devops/pkg/service"
)

func main() {
	cfg := config.LoadConf()
	service.Serve(cfg)
}
