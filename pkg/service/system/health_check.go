package system

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"kubesphere.io/devops/pkg/constants"
	"kubesphere.io/devops/pkg/ds"
	"kubesphere.io/devops/pkg/logger"
)

type SystemService struct {
	Ds *ds.Ds
}

func (s *SystemService) HealthCheck(w rest.ResponseWriter, r *rest.Request) {
	globalRole, err := s.Ds.Jenkins.GetGlobalRole(constants.JenkinsAllUserRoleName)
	if err != nil {
		logger.Warn("failed to get jenkins role, jenkins is not ready")
		rest.Error(w, "jenkins is not ready", http.StatusServiceUnavailable)
	}

	if globalRole == nil {
		logger.Error("jenkins role has been modified by user")
		rest.Error(w, "jenkins role has been modified by user", http.StatusServiceUnavailable)
	}
	w.WriteJson(struct {
		Status string `json:"status"`
	}{"ok"})
}
