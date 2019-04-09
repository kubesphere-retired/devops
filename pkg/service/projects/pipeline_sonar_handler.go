package projects

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/kubesphere/sonargo/sonar"

	"kubesphere.io/devops/pkg/gojenkins"
	"kubesphere.io/devops/pkg/logger"
	"kubesphere.io/devops/pkg/utils/stringutils"
	"kubesphere.io/devops/pkg/utils/userutils"
)

const (
	SonarAnalysisActionClass = "hudson.plugins.sonar.action.SonarAnalysisAction"
	SonarMetricKeys          = "alert_status,quality_gate_details,bugs,new_bugs,reliability_rating,new_reliability_rating,vulnerabilities,new_vulnerabilities,security_rating,new_security_rating,code_smells,new_code_smells,sqale_rating,new_maintainability_rating,sqale_index,new_technical_debt,coverage,new_coverage,new_lines_to_cover,tests,duplicated_lines_density,new_duplicated_lines_density,duplicated_blocks,ncloc,ncloc_language_distribution,projects,new_lines"
	SonarAdditionalFields    = "metrics,periods"
)

type SonarStatus struct {
	Measures      *sonargo.MeasuresComponentObject `json:"measures,omitempty"`
	Issues        *sonargo.IssuesSearchObject      `json:"issues,omitempty"`
	JenkinsAction *gojenkins.GeneralObj            `json:"jenkinsAction,omitempty"`
	Task          *sonargo.CeTaskObject            `json:"task,omitempty"`
}

func (s *ProjectService) GetPipelineSonarHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	pipelineId := r.PathParams["pid"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId, AllRoleSlice)
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	job, err := s.Ds.Jenkins.GetJob(pipelineId, projectId)
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	build, err := job.GetLastBuild()
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	sonarStatus, err := s.getBuildSonarResults(build)
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(sonarStatus) == 0 {
		build, err := job.GetLastCompletedBuild()
		if err != nil && stringutils.GetJenkinsStatusCode(err) != http.StatusNotFound {
			logger.Error("%+v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		sonarStatus, err = s.getBuildSonarResults(build)
	}

	w.WriteJson(sonarStatus)
	return
}

func (s *ProjectService) GetMultiBranchPipelineSonarHandler(w rest.ResponseWriter, r *rest.Request) {
	projectId := r.PathParams["id"]
	pipelineId := r.PathParams["pid"]
	branchName := r.PathParams["bid"]
	operator := userutils.GetUserNameFromRequest(r)
	err := s.checkProjectUserInRole(operator, projectId, AllRoleSlice)
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	job, err := s.Ds.Jenkins.GetJob(branchName, projectId, pipelineId)
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}
	build, err := job.GetLastBuild()
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
		return
	}

	sonarStatus, err := s.getBuildSonarResults(build)
	if err != nil {
		logger.Error("%+v", err)
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(sonarStatus) == 0 {
		build, err := job.GetLastCompletedBuild()
		if err != nil && stringutils.GetJenkinsStatusCode(err) != http.StatusNotFound {
			logger.Error("%+v", err)
			rest.Error(w, err.Error(), stringutils.GetJenkinsStatusCode(err))
			return
		}
		sonarStatus, err = s.getBuildSonarResults(build)
	}

	w.WriteJson(sonarStatus)
}

func (s *ProjectService) getBuildSonarResults(build *gojenkins.Build) ([]*SonarStatus, error) {
	actions := build.GetActions()
	sonarStatuses := make([]*SonarStatus, 0)
	for _, action := range actions {
		if action.ClassName == SonarAnalysisActionClass {
			sonarStatus := &SonarStatus{}
			taskOptions := &sonargo.CeTaskOption{
				Id: action.SonarTaskId,
			}
			ceTask, _, err := s.Ds.Sonar.Ce.Task(taskOptions)
			if err != nil {
				logger.Error("get sonar task error [%+v]", err)
				continue
			}
			sonarStatus.Task = ceTask
			measuresComponentOption := &sonargo.MeasuresComponentOption{
				Component:        ceTask.Task.ComponentKey,
				AdditionalFields: SonarAdditionalFields,
				MetricKeys:       SonarMetricKeys,
			}
			measures, _, err := s.Ds.Sonar.Measures.Component(measuresComponentOption)
			if err != nil {
				logger.Error("get sonar task error [%+v]", err)
				continue
			}
			sonarStatus.Measures = measures

			issuesSearchOption := &sonargo.IssuesSearchOption{
				AdditionalFields: "_all",
				ComponentKeys:    ceTask.Task.ComponentKey,
				Resolved:         "false",
				Ps:               "10",
				S:                "FILE_LINE",
				Facets:           "severities,types",
			}
			issuesSearch, _, err := s.Ds.Sonar.Issues.Search(issuesSearchOption)
			sonarStatus.Issues = issuesSearch
			jenkinsAction := action
			sonarStatus.JenkinsAction = &jenkinsAction

			sonarStatuses = append(sonarStatuses, sonarStatus)
		}
	}
	return sonarStatuses, nil
}
