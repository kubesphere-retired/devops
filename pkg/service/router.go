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
	"github.com/ant0ine/go-json-rest/rest"

	"kubesphere.io/devops/pkg/logger"
)

func Router(s *Server) (app rest.App) {

	app, err := rest.MakeRouter(
		rest.Get("/projects", s.Projects.GetProjectsHandler),
		rest.Get("/projects/:id", s.Projects.GetProjectHandler),
		rest.Post("/projects", s.Projects.CreateProjectHandler),
		rest.Patch("/projects/:id", s.Projects.UpdateProjectHandler),
		rest.Delete("/projects/:id", s.Projects.DeleteProjectHandler),
		rest.Get("/projects/:id/members", s.Projects.GetMembersHandler),
		rest.Get("/projects/:id/members/:uid", s.Projects.GetMemberHandler),
		rest.Post("/projects/:id/members", s.Projects.AddProjectMemberHandler),
		rest.Patch("/projects/:id/members/:uid", s.Projects.UpdateMemberHandler),
		rest.Delete("/projects/:id/members/:uid", s.Projects.DeleteMemberHandler))
	if err != nil {
		logger.Critical("%v", err)
		return
	}
	return
}
