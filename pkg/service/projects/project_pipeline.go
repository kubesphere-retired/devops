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
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/mitchellh/mapstructure"

	"kubesphere.io/devops/pkg/utils/idutils"
)

const (
	JenkinsJobPipeline            = "pipeline"
	JenkinsJobMultiBranchPipeline = "multi-branch-pipeline"
)

var ParameterTypeMap = map[string]string{
	"hudson.model.StringParameterDefinition":   "string",
	"hudson.model.ChoiceParameterDefinition":   "choice",
	"hudson.model.TextParameterDefinition":     "text",
	"hudson.model.BooleanParameterDefinition":  "boolean",
	"hudson.model.FileParameterDefinition":     "file",
	"hudson.model.PasswordParameterDefinition": "password",
}

type JenkinsJobRequest struct {
	Type   string                 `json:"type"`
	Define map[string]interface{} `json:"define"`
}

type Pipeline struct {
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Discarder         *DiscarderProperty `json:"discarder"`
	Parameters        []*Parameter       `json:"parameters"`
	DisableConcurrent bool               `json:"disable_concurrent" mapstructure:"disable_concurrent"`
	TimerTrigger      *TimerTrigger      `json:"timer_trigger" mapstructure:"timer_trigger"`
	RemoteTrigger     *RemoteTrigger     `json:"remote_trigger" mapstructure:"remote_trigger"`
	Jenkinsfile       string             `json:"jenkinsfile"`
}

type MultiBranchPipeline struct {
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Discarder    *DiscarderProperty `json:"discarder"`
	TimerTrigger *TimerTrigger      `json:"timer_trigger" mapstructure:"timer_trigger"`
	Source       Source             `json:"source"`
	ScriptPath   string             `json:"script_path" mapstructure:"script_path"`
}

type Source struct {
	Type   string                 `json:"type"`
	Define map[string]interface{} `json:"define"`
}

type GitSource struct {
	Url              string `json:"url"`
	CredentialId     string `json:"credential_id" mapstructure:"credential_id"`
	DiscoverBranches bool   `json:"discover_branches" mapstructure:"discover_branches"`
}

type GithubSource struct {
	Owner                string                     `json:"owner"`
	Repo                 string                     `json:"repo"`
	CredentialId         string                     `json:"credential_id" mapstructure:"credential_id"`
	ApiUri               string                     `json:"api_uri,omitempty"`
	DiscoverBranches     int                        `json:"discover_branches" mapstructure:"discover_branches"`
	DiscoverPRFromOrigin int                        `json:"discover_pr_from_origin" mapstructure:"discover_pr_from_origin"`
	DiscoverPRFromForks  *GithubDiscoverPRFromForks `json:"discover_pr_from_forks" mapstructure:"discover_pr_from_forks"`
}

type ScmInfo struct {
	Type   string `json:"type"`
	Repo   string `json:"repo"`
	ApiUri string `json:"api_uri,omitempty"`
	Path   string `json:"path"`
}

type GithubDiscoverPRFromForks struct {
	Strategy int `json:"strategy"`
	Trust    int `json:"trust"`
}

type DiscarderProperty struct {
	DaysToKeep string `json:"days_to_keep" mapstructure:"days_to_keep"`
	NumToKeep  string `json:"num_to_keep" mapstructure:"num_to_keep"`
}

type Parameter struct {
	Name         string `json:"name"`
	DefaultValue string `json:"default_value,omitempty" mapstructure:"default_value"`
	Type         string `json:"type"`
	Description  string `json:"description"`
}

type TimerTrigger struct {
	Cron     string `json:"cron,omitempty"`
	Interval string `json:"interval,omitempty"`
}

type RemoteTrigger struct {
	Token string `json:"token"`
}

func createPipelineConfigXml(pipeline *Pipeline) (string, error) {
	doc := etree.NewDocument()
	xmlString := `<?xml version='1.0' encoding='UTF-8'?>
<flow-definition plugin="workflow-job">
  <actions>
    <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobAction plugin="pipeline-model-definition"/>
    <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobPropertyTrackerAction plugin="pipeline-model-definition">
      <jobProperties/>
      <triggers/>
      <parameters/>
      <options/>
    </org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobPropertyTrackerAction>
  </actions>
</flow-definition>
`
	doc.ReadFromString(xmlString)
	flow := doc.SelectElement("flow-definition")
	flow.CreateElement("description").SetText(pipeline.Description)
	properties := flow.CreateElement("properties")

	if pipeline.DisableConcurrent {
		properties.CreateElement("org.jenkinsci.plugins.workflow.job.properties.DisableConcurrentBuildsJobProperty")
	}

	if pipeline.Discarder != nil {
		discarder := properties.CreateElement("jenkins.model.BuildDiscarderProperty")
		strategy := discarder.CreateElement("strategy")
		strategy.CreateAttr("class", "hudson.tasks.LogRotator")
		strategy.CreateElement("daysToKeep").SetText(pipeline.Discarder.DaysToKeep)
		strategy.CreateElement("numToKeep").SetText(pipeline.Discarder.NumToKeep)
		strategy.CreateElement("artifactDaysToKeep").SetText("-1")
		strategy.CreateElement("artifactNumToKeep").SetText("-1")
	}
	if pipeline.Parameters != nil {
		parameterDefinitions := properties.CreateElement("hudson.model.ParametersDefinitionProperty").
			CreateElement("parameterDefinitions")
		for _, parameter := range pipeline.Parameters {
			for className, typeName := range ParameterTypeMap {
				if typeName == parameter.Type {
					paramDefine := parameterDefinitions.CreateElement(className)
					paramDefine.CreateElement("name").SetText(parameter.Name)
					paramDefine.CreateElement("description").SetText(parameter.Description)
					switch parameter.Type {
					case "choice":
						choices := paramDefine.CreateElement("choices")
						choices.CreateAttr("class", "java.util.Arrays$ArrayList")
						a := choices.CreateElement("a")
						a.CreateAttr("class", "string-array")
						choiceValues := strings.Split(parameter.DefaultValue, "\n")
						for _, choiceValue := range choiceValues {
							a.CreateElement("string").SetText(choiceValue)
						}
					case "file":
						break
					default:
						paramDefine.CreateElement("defaultValue").SetText(parameter.DefaultValue)
					}
				}
			}
		}
	}

	if pipeline.TimerTrigger != nil {
		triggers := properties.
			CreateElement("org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty").
			CreateElement("triggers")
		triggers.CreateElement("hudson.triggers.TimerTrigger").CreateElement("spec").SetText(pipeline.TimerTrigger.Cron)
	}

	pipelineDefine := flow.CreateElement("definition")
	pipelineDefine.CreateAttr("class", "org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition")
	pipelineDefine.CreateAttr("plugin", "workflow-cps")
	pipelineDefine.CreateElement("script").SetText(pipeline.Jenkinsfile)

	pipelineDefine.CreateElement("sandbox").SetText("true")

	flow.CreateElement("triggers")

	if pipeline.RemoteTrigger != nil {
		flow.CreateElement("authToken").SetText(pipeline.RemoteTrigger.Token)
	}
	flow.CreateElement("disabled").SetText("false")

	doc.Indent(2)
	stringXml, err := doc.WriteToString()
	if err != nil {
		return "", err
	}
	return replaceXmlVersion(stringXml, "1.0", "1.1"), err
}

func parsePipelineConfigXml(config string) (*Pipeline, error) {
	pipeline := &Pipeline{}
	config = replaceXmlVersion(config, "1.1", "1.0")
	doc := etree.NewDocument()
	err := doc.ReadFromString(config)
	if err != nil {
		return nil, err
	}
	flow := doc.SelectElement("flow-definition")
	if flow == nil {
		return nil, fmt.Errorf("can not find pipeline definition")
	}
	pipeline.Description = flow.SelectElement("description").Text()

	properties := flow.SelectElement("properties")
	if properties.
		SelectElement(
			"org.jenkinsci.plugins.workflow.job.properties.DisableConcurrentBuildsJobProperty") != nil {
		pipeline.DisableConcurrent = true
	}
	if properties.SelectElement("jenkins.model.BuildDiscarderProperty") != nil {
		strategy := properties.
			SelectElement("jenkins.model.BuildDiscarderProperty").
			SelectElement("strategy")
		pipeline.Discarder = &DiscarderProperty{
			DaysToKeep: strategy.SelectElement("daysToKeep").Text(),
			NumToKeep:  strategy.SelectElement("numToKeep").Text(),
		}
	}
	if parametersProperty := properties.SelectElement("hudson.model.ParametersDefinitionProperty"); parametersProperty != nil {
		params := parametersProperty.SelectElement("parameterDefinitions").ChildElements()
		for _, param := range params {
			switch param.Tag {
			case "hudson.model.StringParameterDefinition":
				pipeline.Parameters = append(pipeline.Parameters, &Parameter{
					Name:         param.SelectElement("name").Text(),
					Description:  param.SelectElement("description").Text(),
					DefaultValue: param.SelectElement("defaultValue").Text(),
					Type:         ParameterTypeMap["hudson.model.StringParameterDefinition"],
				})
			case "hudson.model.BooleanParameterDefinition":
				pipeline.Parameters = append(pipeline.Parameters, &Parameter{
					Name:         param.SelectElement("name").Text(),
					Description:  param.SelectElement("description").Text(),
					DefaultValue: param.SelectElement("defaultValue").Text(),
					Type:         ParameterTypeMap["hudson.model.BooleanParameterDefinition"],
				})
			case "hudson.model.TextParameterDefinition":
				pipeline.Parameters = append(pipeline.Parameters, &Parameter{
					Name:         param.SelectElement("name").Text(),
					Description:  param.SelectElement("description").Text(),
					DefaultValue: param.SelectElement("defaultValue").Text(),
					Type:         ParameterTypeMap["hudson.model.TextParameterDefinition"],
				})
			case "hudson.model.FileParameterDefinition":
				pipeline.Parameters = append(pipeline.Parameters, &Parameter{
					Name:        param.SelectElement("name").Text(),
					Description: param.SelectElement("description").Text(),
					Type:        ParameterTypeMap["hudson.model.FileParameterDefinition"],
				})
			case "hudson.model.PasswordParameterDefinition":
				pipeline.Parameters = append(pipeline.Parameters, &Parameter{
					Name:         param.SelectElement("name").Text(),
					Description:  param.SelectElement("description").Text(),
					DefaultValue: param.SelectElement("name").Text(),
					Type:         ParameterTypeMap["hudson.model.PasswordParameterDefinition"],
				})
			case "hudson.model.ChoiceParameterDefinition":
				choiceParameter := &Parameter{
					Name:        param.SelectElement("name").Text(),
					Description: param.SelectElement("description").Text(),
					Type:        ParameterTypeMap["hudson.model.ChoiceParameterDefinition"],
				}
				choices := param.SelectElement("choices").SelectElement("a").SelectElements("string")
				for _, choice := range choices {
					choiceParameter.DefaultValue += fmt.Sprintf("%s\n", choice.Text())
				}
				pipeline.Parameters = append(pipeline.Parameters, choiceParameter)
			default:
				pipeline.Parameters = append(pipeline.Parameters, &Parameter{
					Name:         param.SelectElement("name").Text(),
					Description:  param.SelectElement("description").Text(),
					DefaultValue: "unknown",
					Type:         param.Tag,
				})
			}
		}
	}

	if triggerProperty := properties.
		SelectElement(
			"org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty"); triggerProperty != nil {
		triggers := triggerProperty.SelectElement("triggers")
		if timerTrigger := triggers.SelectElement("hudson.triggers.TimerTrigger"); timerTrigger != nil {
			pipeline.TimerTrigger = &TimerTrigger{
				Cron: timerTrigger.SelectElement("spec").Text(),
			}
		}
	}
	if authToken := flow.SelectElement("authToken"); authToken != nil {
		pipeline.RemoteTrigger = &RemoteTrigger{
			Token: authToken.Text(),
		}
	}
	if definition := flow.SelectElement("definition"); definition != nil {
		if script := definition.SelectElement("script"); script != nil {
			pipeline.Jenkinsfile = script.Text()
		}
	}
	return pipeline, nil
}
func parseMultiBranchPipelineConfigXml(config string) (*MultiBranchPipeline, error) {
	pipeline := &MultiBranchPipeline{}
	config = replaceXmlVersion(config, "1.1", "1.0")
	doc := etree.NewDocument()
	err := doc.ReadFromString(config)
	if err != nil {
		return nil, err
	}
	project := doc.SelectElement("org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject")
	if project == nil {
		return nil, fmt.Errorf("can not parse mutibranch pipeline config")
	}
	pipeline.Description = project.SelectElement("description").Text()

	if discarder := project.SelectElement("orphanedItemStrategy"); discarder != nil {
		pipeline.Discarder = &DiscarderProperty{
			DaysToKeep: discarder.SelectElement("daysToKeep").Text(),
			NumToKeep:  discarder.SelectElement("numToKeep").Text(),
		}
	}
	if triggers := project.SelectElement("triggers"); triggers != nil {
		if timerTrigger := triggers.SelectElement(
			"com.cloudbees.hudson.plugins.folder.computed.PeriodicFolderTrigger"); timerTrigger != nil {
			pipeline.TimerTrigger = &TimerTrigger{
				Interval: timerTrigger.SelectElement("interval").Text(),
			}
		}
	}

	if sources := project.SelectElement("sources"); sources != nil {
		if sourcesData := sources.SelectElement("data"); sourcesData != nil {
			if branchSource := sourcesData.SelectElement("jenkins.branch.BranchSource"); branchSource != nil {
				source := branchSource.SelectElement("source")
				switch source.SelectAttr("class").Value {
				case "org.jenkinsci.plugins.github_branch_source.GitHubSCMSource":
					githubSource := &GithubSource{}
					if credential := source.SelectElement("credentialsId"); credential != nil {
						githubSource.CredentialId = credential.Text()
					}
					if repoOwner := source.SelectElement("repoOwner"); repoOwner != nil {
						githubSource.Owner = repoOwner.Text()
					}
					if repository := source.SelectElement("repository"); repository != nil {
						githubSource.Repo = repository.Text()
					}
					if apiUri := source.SelectElement("apiUri"); apiUri != nil {
						githubSource.ApiUri = apiUri.Text()
					}
					traits := source.SelectElement("traits")
					if branchDiscoverTrait := traits.SelectElement(
						"org.jenkinsci.plugins.github__branch__source.BranchDiscoveryTrait"); branchDiscoverTrait != nil {
						strategyId, err := strconv.Atoi(branchDiscoverTrait.SelectElement("strategyId").Text())
						if err != nil {
							return nil, err
						}
						githubSource.DiscoverBranches = strategyId
					}
					if originPRDiscoverTrait := traits.SelectElement(
						"org.jenkinsci.plugins.github__branch__source.OriginPullRequestDiscoveryTrait"); originPRDiscoverTrait != nil {
						strategyId, err := strconv.Atoi(originPRDiscoverTrait.SelectElement("strategyId").Text())
						if err != nil {
							return nil, err
						}
						githubSource.DiscoverPRFromOrigin = strategyId
					}
					if forkPRDiscoverTrait := traits.SelectElement(
						"org.jenkinsci.plugins.github__branch__source.ForkPullRequestDiscoveryTrait"); forkPRDiscoverTrait != nil {
						strategyId, err := strconv.Atoi(forkPRDiscoverTrait.SelectElement("strategyId").Text())
						if err != nil {
							return nil, err
						}
						trustClass := forkPRDiscoverTrait.SelectElement("trust").SelectAttr("class").Value
						trust := strings.Split(trustClass, "$")
						switch trust[1] {
						case "TrustContributors":
							githubSource.DiscoverPRFromForks = &GithubDiscoverPRFromForks{
								Strategy: strategyId,
								Trust:    1,
							}
						case "TrustEveryone":
							githubSource.DiscoverPRFromForks = &GithubDiscoverPRFromForks{
								Strategy: strategyId,
								Trust:    2,
							}
						case "TrustPermission":
							githubSource.DiscoverPRFromForks = &GithubDiscoverPRFromForks{
								Strategy: strategyId,
								Trust:    3,
							}
						case "TrustNobody":
							githubSource.DiscoverPRFromForks = &GithubDiscoverPRFromForks{
								Strategy: strategyId,
								Trust:    4,
							}
						}

					}
					scmSource := Source{
						Type: "github",
					}
					jsonByte, err := json.Marshal(githubSource)
					if err != nil {
						return nil, err
					}
					err = json.Unmarshal(jsonByte, &scmSource.Define)
					if err != nil {
						return nil, err
					}
					pipeline.Source = scmSource
				case "jenkins.plugins.git.GitSCMSource":
					gitSource := &GitSource{}
					if credential := source.SelectElement("credentialsId"); credential != nil {
						gitSource.CredentialId = credential.Text()
					}
					if remote := source.SelectElement("remote"); remote != nil {
						gitSource.Url = remote.Text()
					}

					traits := source.SelectElement("traits")
					if branchDiscoverTrait := traits.SelectElement(
						"org.jenkinsci.plugins.github__branch__source.BranchDiscoveryTrait"); branchDiscoverTrait != nil {
						gitSource.DiscoverBranches = true
					}
					scmSource := Source{
						Type: "git",
					}
					jsonByte, err := json.Marshal(gitSource)
					if err != nil {
						return nil, err
					}
					err = json.Unmarshal(jsonByte, &scmSource.Define)
					if err != nil {
						return nil, err
					}
					pipeline.Source = scmSource
				}

			}
		}
	}

	pipeline.ScriptPath = project.SelectElement("factory").SelectElement("scriptPath").Text()
	return pipeline, nil
}

func parseMultiBranchPipelineScm(config string) (*ScmInfo, error) {
	config = replaceXmlVersion(config, "1.1", "1.0")
	doc := etree.NewDocument()
	err := doc.ReadFromString(config)
	if err != nil {
		return nil, err
	}
	project := doc.SelectElement("org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject")
	if project == nil {
		return nil, fmt.Errorf("can not parse mutibranch pipeline config")
	}
	scmInfo := &ScmInfo{}
	scmInfo.Path = project.SelectElement("factory").SelectElement("scriptPath").Text()
	if sources := project.SelectElement("sources"); sources != nil {
		if sourcesData := sources.SelectElement("data"); sourcesData != nil {
			if branchSource := sourcesData.SelectElement("jenkins.branch.BranchSource"); branchSource != nil {
				source := branchSource.SelectElement("source")
				switch source.SelectAttr("class").Value {
				case "org.jenkinsci.plugins.github_branch_source.GitHubSCMSource":
					scmInfo.Type = "github"
					repoOwner := source.SelectElement("repoOwner")
					repository := source.SelectElement("repository")
					if repoOwner != nil && repository != nil {
						scmInfo.Repo = repoOwner.Text() + ":" + repository.Text()
					}
					if apiUri := source.SelectElement("apiUri"); apiUri != nil {
						scmInfo.ApiUri = apiUri.Text()
					}
					return scmInfo, nil
				case "jenkins.plugins.git.GitSCMSource":
					scmInfo.Type = "git"
					if remote := source.SelectElement("remote"); remote != nil {
						scmInfo.Repo = remote.Text()
					}
					return scmInfo, nil
				default:
					return nil, nil
				}
			}
		}
	}
	return nil, nil

}
func createMultiBranchPipelineConfigXml(pipeline *MultiBranchPipeline) (string, error) {
	doc := etree.NewDocument()
	xmlString := `
<?xml version='1.0' encoding='UTF-8'?>
<org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject plugin="workflow-multibranch">
  <actions/>
  <properties>
    <org.jenkinsci.plugins.pipeline.modeldefinition.config.FolderConfig plugin="pipeline-model-definition">
      <dockerLabel></dockerLabel>
      <registry plugin="docker-commons"/>
    </org.jenkinsci.plugins.pipeline.modeldefinition.config.FolderConfig>
  </properties>
  <folderViews class="jenkins.branch.MultiBranchProjectViewHolder" plugin="branch-api">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </folderViews>
  <healthMetrics>
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric plugin="cloudbees-folder">
      <nonRecursive>false</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  </healthMetrics>
  <icon class="jenkins.branch.MetadataActionFolderIcon" plugin="branch-api">
    <owner class="org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject" reference="../.."/>
  </icon>
</org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject>`
	err := doc.ReadFromString(xmlString)
	if err != nil {
		return "", err
	}

	project := doc.SelectElement("org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject")
	project.CreateElement("description").SetText(pipeline.Description)

	if pipeline.Discarder != nil {
		discarder := project.CreateElement("orphanedItemStrategy")
		discarder.CreateAttr("class", "com.cloudbees.hudson.plugins.folder.computed.DefaultOrphanedItemStrategy")
		discarder.CreateAttr("plugin", "cloudbees-folder")
		discarder.CreateElement("pruneDeadBranches").SetText("true")
		discarder.CreateElement("daysToKeep").SetText(pipeline.Discarder.DaysToKeep)
		discarder.CreateElement("numToKeep").SetText(pipeline.Discarder.NumToKeep)
	}

	triggers := project.CreateElement("triggers")
	if pipeline.TimerTrigger != nil {
		timeTrigger := triggers.CreateElement(
			"com.cloudbees.hudson.plugins.folder.computed.PeriodicFolderTrigger")
		timeTrigger.CreateAttr("plugin", "cloudbees-folder")
		timeTrigger.CreateElement("spec")
		timeTrigger.CreateElement("interval").SetText(pipeline.TimerTrigger.Interval)

		triggers.CreateElement("disabled").SetText("false")
	}

	sources := project.CreateElement("sources")
	sources.CreateAttr("class", "jenkins.branch.MultiBranchProject$BranchSourceList")
	sources.CreateAttr("plugin", "branch-api")
	sourcesOwner := sources.CreateElement("owner")
	sourcesOwner.CreateAttr("class", "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject")
	sourcesOwner.CreateAttr("reference", "../..")

	branchSource := sources.CreateElement("data").CreateElement("jenkins.branch.BranchSource")
	branchSourceStrategy := branchSource.CreateElement("strategy")
	branchSourceStrategy.CreateAttr("class", "jenkins.branch.NamedExceptionsBranchPropertyStrategy")
	branchSourceStrategy.CreateElement("defaultProperties").CreateAttr("class", "empty-list")
	branchSourceStrategy.CreateElement("namedExceptions").CreateAttr("class", "empty-list")

	switch pipeline.Source.Type {
	case "git":
		gitDefine := &GitSource{}
		err := mapstructure.Decode(pipeline.Source.Define, gitDefine)
		if err != nil {
			return "", err
		}
		gitSource := branchSource.CreateElement("source")
		gitSource.CreateAttr("class", "jenkins.plugins.git.GitSCMSource")
		gitSource.CreateAttr("plugin", "git")
		gitSource.CreateElement("id").SetText(idutils.GetUuid("git-"))
		gitSource.CreateElement("remote").SetText(gitDefine.Url)
		if gitDefine.CredentialId != "" {
			gitSource.CreateElement("credentialsId").SetText(gitDefine.CredentialId)
		}
		traits := gitSource.CreateElement("traits")
		if gitDefine.DiscoverBranches {
			traits.CreateElement("jenkins.plugins.git.traits.BranchDiscoveryTrait")
		}

	case "github":
		githubDefine := &GithubSource{}
		err := mapstructure.Decode(pipeline.Source.Define, githubDefine)
		if err != nil {
			return "", err
		}
		githubSource := branchSource.CreateElement("source")
		githubSource.CreateAttr("class", "org.jenkinsci.plugins.github_branch_source.GitHubSCMSource")
		githubSource.CreateAttr("plugin", "github-branch-source")
		githubSource.CreateElement("id").SetText(idutils.GetUuid("github-"))
		githubSource.CreateElement("credentialsId").SetText(githubDefine.CredentialId)
		githubSource.CreateElement("repoOwner").SetText(githubDefine.Owner)
		githubSource.CreateElement("repository").SetText(githubDefine.Repo)
		if githubDefine.ApiUri != "" {
			githubSource.CreateElement("apiUri").SetText(githubDefine.ApiUri)
		}
		traits := githubSource.CreateElement("traits")
		if githubDefine.DiscoverBranches != 0 {
			traits.CreateElement("org.jenkinsci.plugins.github__branch__source.BranchDiscoveryTrait").
				CreateElement("strategyId").SetText(strconv.Itoa(githubDefine.DiscoverBranches))
		}
		if githubDefine.DiscoverPRFromOrigin != 0 {
			traits.CreateElement("org.jenkinsci.plugins.github__branch__source.OriginPullRequestDiscoveryTrait").
				CreateElement("strategyId").SetText(strconv.Itoa(githubDefine.DiscoverPRFromOrigin))
		}
		if githubDefine.DiscoverPRFromForks != nil {
			forkTrait := traits.CreateElement("org.jenkinsci.plugins.github__branch__source.ForkPullRequestDiscoveryTrait")
			forkTrait.CreateElement("strategyId").SetText(strconv.Itoa(githubDefine.DiscoverPRFromForks.Strategy))
			trustClass := "org.jenkinsci.plugins.github_branch_source.ForkPullRequestDiscoveryTrait$"
			switch githubDefine.DiscoverPRFromForks.Trust {
			case 1:
				trustClass += "TrustContributors"
			case 2:
				trustClass += "TrustEveryone"
			case 3:
				trustClass += "TrustPermission"
			case 4:
				trustClass += "TrustNobody"
			default:
				return "", fmt.Errorf("unsupport trust choice")
			}
			forkTrait.CreateElement("trust").CreateAttr("class", trustClass)
		}
	default:
		return "", fmt.Errorf("unsupport source type")
	}
	factory := project.CreateElement("factory")
	factory.CreateAttr("class", "org.jenkinsci.plugins.workflow.multibranch.WorkflowBranchProjectFactory")

	factoryOwner := factory.CreateElement("owner")
	factoryOwner.CreateAttr("class", "org.jenkinsci.plugins.workflow.multibranch.WorkflowMultiBranchProject")
	factoryOwner.CreateAttr("reference", "../..")
	factory.CreateElement("scriptPath").SetText(pipeline.ScriptPath)

	doc.Indent(2)
	stringXml, err := doc.WriteToString()
	return replaceXmlVersion(stringXml, "1.0", "1.1"), err
}

func replaceXmlVersion(config, oldVersion, targetVersion string) string {
	lines := strings.Split(string(config), "\n")
	lines[0] = strings.Replace(lines[0], oldVersion, targetVersion, -1)
	output := strings.Join(lines, "\n")
	return output
}
