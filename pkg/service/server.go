package service

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"kubesphere.io/devops/pkg/config"
	"kubesphere.io/devops/pkg/ds"
	"kubesphere.io/devops/pkg/logger"
	"kubesphere.io/devops/pkg/service/projects"
)

type Server struct {
	Ds       *ds.Ds
	Projects *projects.ProjectService
}

const APIVersion = "/api/v1alpha"

func Serve(cfg *config.Config) {

	s := Server{}
	s.Ds = ds.NewDs(cfg)
	s.Projects = &projects.ProjectService{Ds: s.Ds}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.AuthBasicMiddleware{
		Realm: "temp",
		Authenticator: func(userId string, password string) bool {
			if userId == "" {
				return false
			}
			return true
		},
	})
	api.SetApp(Router(&s))
	http.Handle(APIVersion+"/", http.StripPrefix(APIVersion, api.MakeHandler()))
	logger.Critical("%v", http.ListenAndServe(":8080", nil))
}
