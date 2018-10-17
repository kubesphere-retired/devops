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

package service

import (
	"net/http"
	"time"

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

	// func to connect jenkins solve https://issues.jenkins-ci.org/browse/JENKINS-2489
	go func() {
		for {
			s.Ds.Jenkins.Info()
			time.Sleep(time.Second * 30)
		}
	}()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.SetApp(Router(&s))
	http.Handle(APIVersion+"/", http.StripPrefix(APIVersion, api.MakeHandler()))
	logger.Critical("%v", http.ListenAndServe(":8080", nil))
}
