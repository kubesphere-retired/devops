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
		rest.Delete("/projects/:id/members/:uid", s.Projects.DeleteMemberHandler),
		rest.Post("/projects/:id/credentials", s.Projects.CreateCredentialHandler),
		rest.Delete("/projects/:id/credentials/:cid", s.Projects.DeleteCredentialHandler),
		rest.Put("/projects/:id/credentials/:cid", s.Projects.UpdateCredentialHandler),
		rest.Get("/projects/:id/credentials/:cid", s.Projects.GetCredentialHandler),
		rest.Get("/projects/:id/credentials", s.Projects.GetCredentialsHandler),
		rest.Get("/projects/:id/pipelines/:pid/config", s.Projects.GetPipelineHandler),
		rest.Post("/projects/:id/pipelines", s.Projects.CreatePipelineHandler),
		rest.Put("/projects/:id/pipelines/:pid", s.Projects.UpdatePipelineHandler),
		rest.Delete("/projects/:id/pipelines/:pid", s.Projects.DeletePipelineHandler),
		rest.Get("/projects/:id/pipelines/:pid/scm", s.Projects.GetPipelineScmHandler),
		rest.Get("/projects/default_roles/", s.Projects.GetProjectDefaultRolesHandler))
	if err != nil {
		logger.Critical("%v", err)
		return
	}
	return
}
