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

	"kubesphere.io/devops/pkg/db"
	"kubesphere.io/devops/pkg/ds"
	"kubesphere.io/devops/pkg/gojenkins"
	"kubesphere.io/devops/pkg/models"
	"kubesphere.io/devops/pkg/utils/reflectutils"
)

type ProjectService struct {
	Ds *ds.Ds
}

const (
	ProjectOwner      = "owner"
	ProjectMaintainer = "maintainer"
	ProjectDeveloper  = "developer"
	ProjectReporter   = "reporter"
)

var AllRoleSlice = []string{ProjectDeveloper, ProjectReporter, ProjectMaintainer, ProjectOwner}

var JenkinsOwnerProjectPermissionIds = &gojenkins.ProjectPermissionIds{
	CredentialCreate:        true,
	CredentialDelete:        true,
	CredentialManageDomains: true,
	CredentialUpdate:        true,
	CredentialView:          true,
	ItemBuild:               true,
	ItemCancel:              true,
	ItemConfigure:           true,
	ItemCreate:              true,
	ItemDelete:              true,
	ItemDiscover:            true,
	ItemMove:                true,
	ItemRead:                true,
	ItemWorkspace:           true,
	RunDelete:               true,
	RunReplay:               true,
	RunUpdate:               true,
	SCMTag:                  true,
}

var JenkinsProjectPermissionMap = map[string]gojenkins.ProjectPermissionIds{
	ProjectOwner: gojenkins.ProjectPermissionIds{
		CredentialCreate:        true,
		CredentialDelete:        true,
		CredentialManageDomains: true,
		CredentialUpdate:        true,
		CredentialView:          true,
		ItemBuild:               true,
		ItemCancel:              true,
		ItemConfigure:           true,
		ItemCreate:              true,
		ItemDelete:              true,
		ItemDiscover:            true,
		ItemMove:                true,
		ItemRead:                true,
		ItemWorkspace:           true,
		RunDelete:               true,
		RunReplay:               true,
		RunUpdate:               true,
		SCMTag:                  true,
	},
	ProjectMaintainer: gojenkins.ProjectPermissionIds{
		CredentialCreate:        true,
		CredentialDelete:        true,
		CredentialManageDomains: true,
		CredentialUpdate:        true,
		CredentialView:          true,
		ItemBuild:               true,
		ItemCancel:              true,
		ItemConfigure:           false,
		ItemCreate:              true,
		ItemDelete:              false,
		ItemDiscover:            true,
		ItemMove:                false,
		ItemRead:                true,
		ItemWorkspace:           true,
		RunDelete:               true,
		RunReplay:               true,
		RunUpdate:               true,
		SCMTag:                  true,
	},
	ProjectDeveloper: gojenkins.ProjectPermissionIds{
		CredentialCreate:        false,
		CredentialDelete:        false,
		CredentialManageDomains: false,
		CredentialUpdate:        false,
		CredentialView:          false,
		ItemBuild:               true,
		ItemCancel:              true,
		ItemConfigure:           false,
		ItemCreate:              false,
		ItemDelete:              false,
		ItemDiscover:            true,
		ItemMove:                false,
		ItemRead:                true,
		ItemWorkspace:           true,
		RunDelete:               true,
		RunReplay:               true,
		RunUpdate:               true,
		SCMTag:                  false,
	},
	ProjectReporter: gojenkins.ProjectPermissionIds{
		CredentialCreate:        false,
		CredentialDelete:        false,
		CredentialManageDomains: false,
		CredentialUpdate:        false,
		CredentialView:          false,
		ItemBuild:               false,
		ItemCancel:              false,
		ItemConfigure:           false,
		ItemCreate:              false,
		ItemDelete:              false,
		ItemDiscover:            true,
		ItemMove:                false,
		ItemRead:                true,
		ItemWorkspace:           false,
		RunDelete:               false,
		RunReplay:               false,
		RunUpdate:               false,
		SCMTag:                  false,
	},
}

var JenkinsPipelinePermissionMap = map[string]gojenkins.ProjectPermissionIds{
	ProjectOwner: gojenkins.ProjectPermissionIds{
		CredentialCreate:        true,
		CredentialDelete:        true,
		CredentialManageDomains: true,
		CredentialUpdate:        true,
		CredentialView:          true,
		ItemBuild:               true,
		ItemCancel:              true,
		ItemConfigure:           true,
		ItemCreate:              true,
		ItemDelete:              true,
		ItemDiscover:            true,
		ItemMove:                true,
		ItemRead:                true,
		ItemWorkspace:           true,
		RunDelete:               true,
		RunReplay:               true,
		RunUpdate:               true,
		SCMTag:                  true,
	},
	ProjectMaintainer: gojenkins.ProjectPermissionIds{
		CredentialCreate:        true,
		CredentialDelete:        true,
		CredentialManageDomains: true,
		CredentialUpdate:        true,
		CredentialView:          true,
		ItemBuild:               true,
		ItemCancel:              true,
		ItemConfigure:           true,
		ItemCreate:              true,
		ItemDelete:              true,
		ItemDiscover:            true,
		ItemMove:                true,
		ItemRead:                true,
		ItemWorkspace:           true,
		RunDelete:               true,
		RunReplay:               true,
		RunUpdate:               true,
		SCMTag:                  true,
	},
	ProjectDeveloper: gojenkins.ProjectPermissionIds{
		CredentialCreate:        false,
		CredentialDelete:        false,
		CredentialManageDomains: false,
		CredentialUpdate:        false,
		CredentialView:          false,
		ItemBuild:               true,
		ItemCancel:              true,
		ItemConfigure:           false,
		ItemCreate:              false,
		ItemDelete:              false,
		ItemDiscover:            true,
		ItemMove:                false,
		ItemRead:                true,
		ItemWorkspace:           true,
		RunDelete:               true,
		RunReplay:               true,
		RunUpdate:               true,
		SCMTag:                  false,
	},
	ProjectReporter: gojenkins.ProjectPermissionIds{
		CredentialCreate:        false,
		CredentialDelete:        false,
		CredentialManageDomains: false,
		CredentialUpdate:        false,
		CredentialView:          false,
		ItemBuild:               false,
		ItemCancel:              false,
		ItemConfigure:           false,
		ItemCreate:              false,
		ItemDelete:              false,
		ItemDiscover:            true,
		ItemMove:                false,
		ItemRead:                true,
		ItemWorkspace:           false,
		RunDelete:               false,
		RunReplay:               false,
		RunUpdate:               false,
		SCMTag:                  false,
	},
}

func GetProjectRoleName(projectId, role string) string {
	return fmt.Sprintf("%s-%s-project", projectId, role)
}

func GetPipelineRoleName(projectId, role string) string {
	return fmt.Sprintf("%s-%s-pipeline", projectId, role)
}

func GetProjectRolePattern(projectId string) string {
	return fmt.Sprintf("^%s$", projectId)
}

func GetPipelineRolePattern(projectId string) string {
	return fmt.Sprintf("^%s/.*", projectId)
}

func (s *ProjectService) checkProjectUserInRole(username, projectId string, roles []string) error {
	membership := &models.ProjectMembership{}
	err := s.Ds.Db.Select(models.ProjectMembershipColumns...).
		From(models.ProjectMembershipTableName).
		Where(db.And(
			db.Eq(models.ProjectMembershipUsernameColumn, username),
			db.Eq(models.ProjectMembershipProjectIdColumn, projectId))).LoadOne(membership)
	if err != nil {
		return err
	}
	if !reflectutils.In(membership.Role, roles) {
		return fmt.Errorf("user [%s] in project [%s] role is not in %s", username, projectId, roles)
	}
	return nil
}
