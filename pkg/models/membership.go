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

package models

import "kubesphere.io/devops/pkg/constants"

const (
	ProjectMembershipTableName       = "project_membership"
	ProjectMembershipUsernameColumn  = "username"
	ProjectMembershipProjectIdColumn = "project_id"
	ProjectMembershipRoleColumn      = "role"
)

type ProjectMembership struct {
	Username  string `json:"username"`
	ProjectId string `json:"project_id" db:"project_id"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	GrantBy   string `json:"grand_by,omitempty"`
}

var ProjectMembershipColumns = GetColumnsFromStruct(&ProjectMembership{})

func NewProjectMemberShip(username, projectId, role, grantBy string) *ProjectMembership {
	return &ProjectMembership{
		Username:  username,
		ProjectId: projectId,
		Role:      role,
		Status:    constants.StatusActive,
		GrantBy:   grantBy,
	}
}
