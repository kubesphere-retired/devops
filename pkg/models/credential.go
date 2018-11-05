package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

const (
	ProjectCredentialTableName    = "project_credential"
	ProjectCredentialIdColumn     = "credential_id"
	ProjectCredentialDomainColumn = "domain"
)

type ProjectCredential struct {
	ProjectId    string    `json:"project_id"`
	CredentialId string    `json:"credential_id"`
	Domain       string    `json:"domain"`
	Creator      string    `json:"creator"`
	CreateTime   time.Time `json:"create_time"`
}

var ProjectCredentialColumns = GetColumnsFromStruct(&ProjectCredential{})

func NewProjectCredential(projectId, credentialId, domain, creator string) *ProjectCredential {
	if govalidator.IsNull(domain) {
		domain = "_"
	}
	return &ProjectCredential{
		ProjectId:    projectId,
		CredentialId: credentialId,
		Domain:       domain,
		Creator:      creator,
		CreateTime:   time.Now(),
	}
}
