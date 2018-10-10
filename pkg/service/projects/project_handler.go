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
	"net/http"
	"strconv"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/asaskevich/govalidator"
	"github.com/gocraft/dbr"

	"kubesphere.io/devops/pkg/constants"
	"kubesphere.io/devops/pkg/db"
	"kubesphere.io/devops/pkg/logger"
	"kubesphere.io/devops/pkg/models"
	"kubesphere.io/devops/pkg/utils/stringutils"
	"kubesphere.io/devops/pkg/utils/userutils"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Description string `json:"description"`
}

type AddProjectMemberRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UpdateProjectMemberRequest struct {
	Role string `json:"role"`
}

func (s *ProjectService) GetProjectHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId,
		[]string{ProjectOwner, ProjectMaintainer, ProjectReporter, ProjectDeveloper})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	project := &models.Project{}
	err = s.Ds.Db.Select(models.ProjectColumns...).
		From(models.ProjectTableName).
		Where(db.Eq(models.ProjectIdColumn, projectId)).
		LoadOne(project)
	if err != nil && err != dbr.ErrNotFound {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err == dbr.ErrNotFound {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteJson(project)
	return
}

func (s *ProjectService) GetProjectsHandler(w rest.ResponseWriter, r *rest.Request) {
	operator := userutils.GetUserNameFromRequest(r)
	id := r.URL.Query().Get("id")
	query := s.Ds.Db.Select(models.ProjectColumns...).
		From(models.ProjectTableName)
	var conditions []dbr.Builder
	switch operator {
	case constants.KS_ADMIN:
		if !govalidator.IsNull(id) {
			ids := strings.Split(id, ",")
			conditions = append(conditions, db.Eq(models.ProjectIdColumn, ids))
		}
	default:
		var membershipCondition []dbr.Builder
		membershipCondition = append(membershipCondition, db.Eq(models.ProjectMembershipUsernameColumn, operator))
		membershipCondition = append(membershipCondition, db.Eq(constants.StatusColumn, constants.StatusActive))
		if !govalidator.IsNull(id) {
			ids := strings.Split(id, ",")
			membershipCondition = append(membershipCondition, db.Eq(models.ProjectIdColumn, ids))
		}
		projectMemberships := make([]*models.ProjectMembership, 0)
		_, err := s.Ds.Db.Select(models.ProjectMembershipColumns...).
			From(models.ProjectMembershipTableName).
			Where(db.And(membershipCondition...)).
			Load(&projectMemberships)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		projectIdArray := make([]string, 0)
		for _, projectMembership := range projectMemberships {
			projectIdArray = append(projectIdArray, projectMembership.ProjectId)
		}
		conditions = append(conditions, db.Eq(models.ProjectIdColumn, projectIdArray))
	}
	projects := make([]*models.Project, 0)
	if len(conditions) > 0 {
		query.Where(db.And(conditions...))
	}
	_, err := query.Load(&projects)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(projects)
	return
}

func (s *ProjectService) CreateProjectHandler(w rest.ResponseWriter, r *rest.Request) {
	creator := userutils.GetUserNameFromRequest(r)
	request := &CreateProjectRequest{}
	err := r.DecodeJsonPayload(request)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	project := models.NewProject(request.Name, request.Description, creator)
	_, err = s.Ds.Jenkins.CreateFolder(project.ProjectId, project.Description)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	for role, permission := range JenkinsProjectPermissionMap {
		_, err := s.Ds.Jenkins.AddProjectRole(GetProjectRoleName(project.ProjectId, role),
			GetProjectRolePattern(project.ProjectId), permission, true)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
	}
	for role, permission := range JenkinsPipelinePermissionMap {
		_, err := s.Ds.Jenkins.AddProjectRole(GetPipelineRoleName(project.ProjectId, role),
			GetPipelineRolePattern(project.ProjectId), permission, true)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
	}

	globalRole, err := s.Ds.Jenkins.GetGlobalRole(constants.JenkinsAllUserRoleName)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = globalRole.AssignRole(creator)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	projectRole, err := s.Ds.Jenkins.GetProjectRole(GetProjectRoleName(project.ProjectId, ProjectOwner))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = projectRole.AssignRole(creator)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pipelineRole, err := s.Ds.Jenkins.GetProjectRole(GetPipelineRoleName(project.ProjectId, ProjectOwner))
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	err = pipelineRole.AssignRole(creator)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	_, err = s.Ds.Db.InsertInto(models.ProjectTableName).
		Columns(models.ProjectColumns...).Record(project).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectMembership := models.NewProjectMemberShip(creator, project.ProjectId, ProjectOwner, creator)
	_, err = s.Ds.Db.InsertInto(models.ProjectMembershipTableName).
		Columns(models.ProjectMembershipColumns...).Record(projectMembership).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(project)
	return
}

func (s *ProjectService) DeleteProjectHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, err = s.Ds.Jenkins.DeleteJob(projectId)
	if err != nil && err.Error() != strconv.Itoa(http.StatusNotFound) {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	roleNames := make([]string, 0)
	for role := range JenkinsProjectPermissionMap {
		roleNames = append(roleNames, GetProjectRoleName(projectId, role))
		roleNames = append(roleNames, GetPipelineRoleName(projectId, role))
	}
	err = s.Ds.Jenkins.DeleteProjectRoles(roleNames...)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	_, err = s.Ds.Db.DeleteFrom(models.ProjectMembershipTableName).
		Where(db.Eq(models.ProjectMembershipProjectIdColumn, projectId)).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = s.Ds.Db.Update(models.ProjectTableName).
		Set(constants.StatusColumn, constants.StatusDeleted).
		Where(db.Eq(models.ProjectIdColumn, projectId)).Exec()
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	project := &models.Project{}
	err = s.Ds.Db.Select(models.ProjectColumns...).
		From(models.ProjectTableName).
		Where(db.Eq(models.ProjectIdColumn, projectId)).
		LoadOne(project)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(project)
	return
}

func (s *ProjectService) UpdateProjectHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	request := &UpdateProjectRequest{}
	err := r.DecodeJsonPayload(request)
	if err != nil {
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
	if !govalidator.IsNull(request.Description) {
		_, err := s.Ds.Db.Update(models.ProjectTableName).
			Set(models.ProjectDescriptionColumn, request.Description).
			Where(db.Eq(models.ProjectIdColumn, projectId)).Exec()
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	project := &models.Project{}
	err = s.Ds.Db.Select(models.ProjectColumns...).
		From(models.ProjectTableName).
		Where(db.Eq(models.ProjectIdColumn, projectId)).
		LoadOne(project)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(project)
	return
}
