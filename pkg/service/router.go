package service

import (
	"github.com/ant0ine/go-json-rest/rest"

	"kubesphere.io/devops/pkg/logger"
)

func Router(s *Server) (app rest.App) {

	app, err := rest.MakeRouter()
	if err != nil {
		logger.Critical("%v", err)
		return
	}
	return
}
