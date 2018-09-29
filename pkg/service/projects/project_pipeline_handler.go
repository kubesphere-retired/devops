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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mitchellh/mapstructure"

	"kubesphere.io/devops/pkg/logger"
	"kubesphere.io/devops/pkg/utils/stringutils"
	"kubesphere.io/devops/pkg/utils/userutils"
)

func (s *ProjectService) CreatePipelineHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	request := &JenkinsJobRequest{}

	err := r.DecodeJsonPayload(request)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner, ProjectMaintainer})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	switch request.Type {
	case JenkinsJobPipeline:
		pipeline := &Pipeline{}
		err := mapstructure.Decode(request.Define, pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		config, err := createPipelineConfigXml(pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = s.Ds.Jenkins.CreateJobInFolder(config, pipeline.Name, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Name string `json:"name"`
		}{Name: pipeline.Name})
		return
	case JenkinsJobMultiBranchPipeline:
		pipeline := &MultiBranchPipeline{}
		err := mapstructure.Decode(request.Define, pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		config, err := createMultiBranchPipelineConfigXml(pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = s.Ds.Jenkins.CreateJobInFolder(config, pipeline.Name, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Name string `json:"name"`
		}{Name: pipeline.Name})
		return

	default:
		err := fmt.Errorf("error unsupport job type")
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *ProjectService) DeletePipelineHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	pipelineId := r.PathParams["pid"]
	err := s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner, ProjectMaintainer})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, err = s.Ds.Jenkins.DeleteJob(pipelineId, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	w.WriteJson(struct {
		Name string `json:"name"`
	}{Name: pipelineId})
}

func (s *ProjectService) UpdatePipelineHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	pipelineId := r.PathParams["pid"]
	operator := userutils.GetUserNameFromRequest(r)
	request := &JenkinsJobRequest{}

	err := r.DecodeJsonPayload(request)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner, ProjectMaintainer})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	switch request.Type {
	case JenkinsJobPipeline:
		pipeline := &Pipeline{}
		err := mapstructure.Decode(request.Define, pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		config, err := createPipelineConfigXml(pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		job, err := s.Ds.Jenkins.GetJob(pipelineId, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		err = job.UpdateConfig(config)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Name string `json:"name"`
		}{Name: pipeline.Name})
		return
	case JenkinsJobMultiBranchPipeline:
		multiBranchPipeline := &MultiBranchPipeline{}
		err := mapstructure.Decode(request.Define, multiBranchPipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		config, err := createMultiBranchPipelineConfigXml(multiBranchPipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		job, err := s.Ds.Jenkins.GetJob(pipelineId, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		err = job.UpdateConfig(config)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Name string `json:"name"`
		}{Name: multiBranchPipeline.Name})
		return
	default:
		err := fmt.Errorf("error unsupport job type")
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
}

func (s *ProjectService) GetPipelineHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	pipelineId := r.PathParams["pid"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner, ProjectMaintainer})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	job, err := s.Ds.Jenkins.GetJob(pipelineId, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	switch job.Raw.Class {
	case "org.jenkinsci.plugins.workflow.job.WorkflowJob":
		config, err := job.GetConfig()
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		pipeline, err := parsePipelineConfigXml(config)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pipeline.Name = pipelineId
		jobRequest := JenkinsJobRequest{
			Type: "pipeline",
		}
		jsonByte, err := json.Marshal(pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(jsonByte, &jobRequest.Define)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteJson(jobRequest)
		return
		return
	case "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject":
		config, err := job.GetConfig()
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		pipeline, err := parseMultiBranchPipelineConfigXml(config)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pipeline.Name = pipelineId
		jobRequest := JenkinsJobRequest{
			Type: "multi-branch-pipeline",
		}
		jsonByte, err := json.Marshal(pipeline)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(jsonByte, &jobRequest.Define)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteJson(jobRequest)
		return

	default:
		err := fmt.Errorf("error unsupport job type")
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *ProjectService) GetPipelineScmHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	pipelineId := r.PathParams["pid"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId, AllRoleSlice)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	job, err := s.Ds.Jenkins.GetJob(pipelineId, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	switch job.Raw.Class {
	case "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject":
		config, err := job.GetConfig()
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		scm, err := parseMultiBranchPipelineScm(config)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteJson(scm)
		return

	default:
		err := fmt.Errorf("error unsupport job type")
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
