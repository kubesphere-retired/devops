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
	"github.com/mitchellh/mapstructure"

	"kubesphere.io/devops/pkg/logger"
	"kubesphere.io/devops/pkg/utils/stringutils"
	"kubesphere.io/devops/pkg/utils/userutils"
)

const (
	CredentialTypeUsernamePassword = "username_password"
	CredentialTypeSsh              = "ssh"
	CredentialTypeSecretText       = "secret_text"
)

type CredentialRequest struct {
	Type    string                 `json:"type"`
	Domain  string                 `json:"domain"`
	Content map[string]interface{} `json:"content"`
}

type UsernamePasswordCredentialRequest struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type SshCredentialRequest struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Passphrase  string `json:"passphrase"`
	PrivateKey  string `json:"private_key" mapstructure:"private_key"`
	Description string `json:"description"`
}

type SecretTextCredentialRequest struct {
	Id          string `json:"id"`
	Secret      string `json:"secret"`
	Description string `json:"description"`
}

type DeleteCredentialRequest struct {
	Domain string `json:"domain"`
}

type CopySshCredentialRequest struct {
	Id string `json:"id"`
}

type CredentialResponse struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	DisplayName string `json:"display_name"`
	Fingerprint *struct {
		FileName string `json:"file_name,omitempty"`
		Hash     string `json:"hash,omitempty"`
		Usage    []*struct {
			Name   string `json:"name,omitempty"`
			Ranges struct {
				Ranges []*struct {
					Start int `json:"start"`
					End   int `json:"end"`
				} `json:"ranges"`
			} `json:"ranges"`
		} `json:"usage,omitempty"`
	} `json:"fingerprint,omitempty"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
}

func (s *ProjectService) CreateCredentialHandler(w rest.ResponseWriter, r *rest.Request) {
	request := &CredentialRequest{}
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)

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
	case CredentialTypeUsernamePassword:
		UPRequest := &UsernamePasswordCredentialRequest{}
		err := mapstructure.Decode(request.Content, UPRequest)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credentialId, err := s.Ds.Jenkins.CreateUsernamePasswordCredentialInFolder(request.Domain, UPRequest.Id,
			UPRequest.Username, UPRequest.Password, UPRequest.Description, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Id string `json:"id"`
		}{Id: *credentialId})
		return

	case CredentialTypeSsh:
		SshRequest := &SshCredentialRequest{}
		err := mapstructure.Decode(request.Content, SshRequest)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credentialId, err := s.Ds.Jenkins.CreateSshCredentialInFolder(request.Domain, SshRequest.Id,
			SshRequest.Username, SshRequest.Passphrase, SshRequest.PrivateKey, SshRequest.Description, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Id string `json:"id"`
		}{Id: *credentialId})
		return

	case CredentialTypeSecretText:
		TextRequest := &SecretTextCredentialRequest{}
		err := mapstructure.Decode(request.Content, TextRequest)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credentialId, err := s.Ds.Jenkins.CreateSecretTextCredentialInFolder(request.Domain, TextRequest.Id,
			TextRequest.Secret, TextRequest.Description, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Id string `json:"id"`
		}{Id: *credentialId})
		return
	default:
		err := fmt.Errorf("error unsupport  credential type")
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
}

func (s *ProjectService) DeleteCredentialHandler(w rest.ResponseWriter, r *rest.Request) {
	request := &DeleteCredentialRequest{}
	projectId := r.PathParams["id"]
	credentialId := r.PathParams["cid"]
	operator := userutils.GetUserNameFromRequest(r)
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
	id, err := s.Ds.Jenkins.DeleteCredentialInFolder(request.Domain, credentialId, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	w.WriteJson(struct {
		Id string `json:"id"`
	}{Id: *id})
	return
}

func (s *ProjectService) UpdateCredentialHandler(w rest.ResponseWriter, r *rest.Request) {
	request := &CredentialRequest{}
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	credentialId := r.PathParams["cid"]
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
	jenkinsCredential, err := s.Ds.Jenkins.GetCredentialInFolder(request.Domain, credentialId, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	credentialType := CredentialTypeMap[jenkinsCredential.TypeName]
	switch credentialType {
	case CredentialTypeUsernamePassword:
		UPRequest := &UsernamePasswordCredentialRequest{}
		UPRequest.Id = credentialId
		err := mapstructure.Decode(request.Content, UPRequest)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credentialId, err := s.Ds.Jenkins.UpdateUsernamePasswordCredentialInFolder(request.Domain, UPRequest.Id,
			UPRequest.Username, UPRequest.Password, UPRequest.Description, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Id string `json:"id"`
		}{Id: *credentialId})
		return

	case CredentialTypeSsh:
		SshRequest := &SshCredentialRequest{}
		SshRequest.Id = credentialId
		err := mapstructure.Decode(request.Content, SshRequest)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credentialId, err := s.Ds.Jenkins.UpdateSshCredentialInFolder(request.Domain, SshRequest.Id,
			SshRequest.Username, SshRequest.Passphrase, SshRequest.PrivateKey, SshRequest.Description, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteJson(struct {
			Id string `json:"id"`
		}{Id: *credentialId})
		return

	case CredentialTypeSecretText:
		TextRequest := &SecretTextCredentialRequest{}
		TextRequest.Id = credentialId
		err := mapstructure.Decode(request.Content, TextRequest)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		credentialId, err := s.Ds.Jenkins.UpdateSecretTextCredentialInFolder(request.Domain, TextRequest.Id,
			TextRequest.Secret, TextRequest.Description, projectId)
		if err != nil {
			logger.Error("%v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		w.WriteJson(struct {
			Id string `json:"id"`
		}{Id: *credentialId})
		return
	default:
		err := fmt.Errorf("error unsupport credential type %s", credentialType)
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *ProjectService) GetCredentialHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	credentialId := r.PathParams["cid"]
	domain := r.URL.Query().Get("domain")
	err := s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner, ProjectMaintainer})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	credentialResponse, err := s.Ds.Jenkins.GetCredentialInFolder(domain, credentialId, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	response := formatCredentialResponse(credentialResponse)
	w.WriteJson(response)
	return
}

func (s *ProjectService) GetCredentialsHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	operator := userutils.GetUserNameFromRequest(r)
	domain := r.URL.Query().Get("domain")
	err := s.checkProjectUserInRole(operator, projectId, []string{ProjectOwner, ProjectMaintainer})
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	jenkinsCredentialResponse, err := s.Ds.Jenkins.GetCredentialsInFolder(domain, projectId)
	if err != nil {
		logger.Error("%v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	response := formatCredentialsResponse(jenkinsCredentialResponse)
	w.WriteJson(response)
	return
}
