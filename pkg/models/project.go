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

import (
	"time"

	"kubesphere.io/devops/pkg/constants"
	"kubesphere.io/devops/pkg/utils/idutils"
)

var ProjectColumns = GetColumnsFromStruct(&Project{})

const (
	ProjectTableName         = "project"
	ProjectPrefix            = "project-"
	ProjectDescriptionColumn = "description"
	ProjectIdColumn          = "project_id"
)

type Project struct {
	ProjectId   string    `json:"project_id" db:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Creator     string    `json:"creator"`
	CreateTime  time.Time `json:"create_time"`
	Status      string    `json:"status"`
	Visibility  string    `json:"visibility"`
}

func NewProject(name, description, creator string) *Project {
	return &Project{
		ProjectId:   idutils.GetUuid(ProjectPrefix),
		Name:        name,
		Description: description,
		Creator:     creator,
		CreateTime:  time.Now(),
		Status:      constants.StatusActive,
		Visibility:  constants.VisibilityPrivate,
	}
}
