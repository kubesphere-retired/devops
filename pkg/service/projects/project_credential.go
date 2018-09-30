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
	"kubesphere.io/devops/pkg/gojenkins"
)

var CredentialTypeMap = map[string]string{
	"SSH Username with private key": CredentialTypeSsh,
	"Username with password":        CredentialTypeUsernamePassword,
	"Secret text":                   CredentialTypeSecretText,
}

func formatCredentialResponse(jenkinsCredentialResponse *gojenkins.CredentialResponse) *CredentialResponse {
	response := &CredentialResponse{}
	response.Id = jenkinsCredentialResponse.Id
	response.Description = jenkinsCredentialResponse.Description
	response.DisplayName = jenkinsCredentialResponse.DisplayName
	if jenkinsCredentialResponse.Fingerprint != nil && jenkinsCredentialResponse.Fingerprint.Hash != "" {
		response.Fingerprint = &struct {
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
		}{}
		response.Fingerprint.FileName = jenkinsCredentialResponse.Fingerprint.FileName
		response.Fingerprint.Hash = jenkinsCredentialResponse.Fingerprint.Hash
		for _, usage := range jenkinsCredentialResponse.Fingerprint.Usage {
			response.Fingerprint.Usage = append(response.Fingerprint.Usage, usage)
		}
	}
	response.Domain = jenkinsCredentialResponse.Domain
	credentialType, ok := CredentialTypeMap[jenkinsCredentialResponse.TypeName]
	if ok {
		response.Type = credentialType
		return response
	}
	response.Type = jenkinsCredentialResponse.TypeName
	return response
}

func formatCredentialsResponse(jenkinsCredentialsResponse []*gojenkins.CredentialResponse) []*CredentialResponse {
	responseSlice := make([]*CredentialResponse, 0)
	for _, jenkinsCredential := range jenkinsCredentialsResponse {
		responseSlice = append(responseSlice, formatCredentialResponse(jenkinsCredential))
	}
	return responseSlice
}
