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

package projects

import (
	"fmt"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/asaskevich/govalidator"

	"kubesphere.io/devops/pkg/constants"
	"kubesphere.io/devops/pkg/db"
	"kubesphere.io/devops/pkg/logger"
	"kubesphere.io/devops/pkg/models"
	"kubesphere.io/devops/pkg/utils/reflectutils"
	"kubesphere.io/devops/pkg/utils/stringutils"
	"kubesphere.io/devops/pkg/utils/userutils"
)

type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var DefaultRoles = []*Role{
	{
		Name:        ProjectMaintainer,
		Description: "项目的主要维护者，可以进行项目内的凭证配置、pipeline配置等操作",
	},
	{
		Name:        ProjectOwner,
		Description: "项目的所有者，可以进行项目的所有操作",
	},
	{
		Name:        ProjectDeveloper,
		Description: "项目的开发者，可以进行pipeline的触发以及查看",
	},
	{
		Name:        ProjectReporter,
		Description: "项目的观察者，可以查看pipeline的运行情况",
	},
}

func (s *ProjectService) GetMembersHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId, []string{
		ProjectOwner, ProjectMaintainer, ProjectReporter, ProjectDeveloper})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	memberships := make([]*models.ProjectMembership, 0)
	_, err = s.Ds.Db.Select(models.ProjectMembershipColumns...).
		From(models.ProjectMembershipTableName).
		Where(db.Eq(models.ProjectIdColumn, projectId)).
		Load(&memberships)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(memberships)
	return
}

func (s *ProjectService) GetMemberHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	username := r.PathParams["uid"]
	err := s.checkProjectUserInRole(operator, projectId, []string{
		ProjectOwner, ProjectMaintainer, ProjectReporter, ProjectDeveloper})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	memberships := &models.ProjectMembership{}
	err = s.Ds.Db.Select(models.ProjectMembershipColumns...).
		From(models.ProjectMembershipTableName).
		Where(db.And(
			db.Eq(models.ProjectIdColumn, projectId),
			db.Eq(models.ProjectMembershipUsernameColumn, username))).
		LoadOne(&memberships)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(memberships)
	return
}

func (s *ProjectService) AddProjectMemberHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	request := &AddProjectMemberRequest{}
	err := r.DecodeJsonPayload(request)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if govalidator.IsNull(request.Username) {
		err := fmt.Errorf("error need username")
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !reflectutils.In(request.Role, AllRoleSlice) {
		err := fmt.Errorf("err role [%s] not in [%s]", request.Role,
			AllRoleSlice)
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	globalRole, err := s.Ds.Jenkins.GetGlobalRole(constants.JenkinsAllUserRoleName)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = globalRole.AssignRole(request.Username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	projectRole, err := s.Ds.Jenkins.GetProjectRole(GetProjectRoleName(projectId, request.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = projectRole.AssignRole(request.Username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	pipelineRole, err := s.Ds.Jenkins.GetProjectRole(GetPipelineRoleName(projectId, request.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = pipelineRole.AssignRole(request.Username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	projectMembership := models.NewProjectMemberShip(request.Username, projectId, request.Role, operator)
	_, err = s.Ds.Db.
		InsertInto(models.ProjectMembershipTableName).
		Columns(models.ProjectMembershipColumns...).
		Record(projectMembership).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(projectMembership)
	return
}

func (s *ProjectService) UpdateMemberHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	username := r.PathParams["uid"]
	operator := userutils.GetUserNameFromRequest(r)
	request := &UpdateProjectMemberRequest{}
	err := r.DecodeJsonPayload(request)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !reflectutils.In(request.Role, AllRoleSlice) {
		err := fmt.Errorf("err role [%s] not in [%s]", request.Role, AllRoleSlice)
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	oldMembership := &models.ProjectMembership{}
	err = s.Ds.Db.Select(models.ProjectMembershipColumns...).
		From(models.ProjectMembershipTableName).
		Where(db.And(
			db.Eq(models.ProjectMembershipUsernameColumn, username),
			db.Eq(models.ProjectMembershipProjectIdColumn, projectId),
		)).LoadOne(oldMembership)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oldProjectRole, err := s.Ds.Jenkins.GetProjectRole(GetProjectRoleName(projectId, oldMembership.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = oldProjectRole.UnAssignRole(username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	oldPipelineRole, err := s.Ds.Jenkins.GetProjectRole(GetPipelineRoleName(projectId, oldMembership.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = oldPipelineRole.UnAssignRole(username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	projectRole, err := s.Ds.Jenkins.GetProjectRole(GetProjectRoleName(projectId, request.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = projectRole.AssignRole(username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	pipelineRole, err := s.Ds.Jenkins.GetProjectRole(GetPipelineRoleName(projectId, request.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = pipelineRole.AssignRole(username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	_, err = s.Ds.Db.Update(models.ProjectMembershipTableName).
		Set(models.ProjectMembershipRoleColumn, request.Role).
		Where(db.And(
			db.Eq(models.ProjectMembershipProjectIdColumn, projectId),
			db.Eq(models.ProjectMembershipUsernameColumn, username),
		)).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseMembership := &models.ProjectMembership{}
	err = s.Ds.Db.Select(models.ProjectMembershipColumns...).
		From(models.ProjectMembershipTableName).
		Where(db.And(
			db.Eq(models.ProjectMembershipUsernameColumn, username),
			db.Eq(models.ProjectMembershipProjectIdColumn, projectId),
		)).LoadOne(responseMembership)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(responseMembership)
	return

}

func (s *ProjectService) DeleteMemberHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	username := r.PathParams["uid"]
	operator := userutils.GetUserNameFromRequest(r)

	err := s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	oldMembership := &models.ProjectMembership{}
	err = s.Ds.Db.Select(models.ProjectMembershipColumns...).
		From(models.ProjectMembershipTableName).
		Where(db.And(
			db.Eq(models.ProjectMembershipUsernameColumn, username),
			db.Eq(models.ProjectMembershipProjectIdColumn, projectId),
		)).LoadOne(oldMembership)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oldProjectRole, err := s.Ds.Jenkins.GetProjectRole(GetProjectRoleName(projectId, oldMembership.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = oldProjectRole.UnAssignRole(username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	oldPipelineRole, err := s.Ds.Jenkins.GetProjectRole(GetPipelineRoleName(projectId, oldMembership.Role))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = oldPipelineRole.UnAssignRole(username)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	_, err = s.Ds.Db.DeleteFrom(models.ProjectMembershipTableName).
		Where(db.And(
			db.Eq(models.ProjectMembershipProjectIdColumn, projectId),
			db.Eq(models.ProjectMembershipUsernameColumn, username),
		)).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(struct {
		Username string `json:"username"`
	}{Username: username})
	return
}

func (s *ProjectService) GetProjectDefaultRolesHandler(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(DefaultRoles)
	return
}
