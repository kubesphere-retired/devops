package projects

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_NoScmPipelineConfig(t *testing.T) {
	inputs := []*Pipeline{
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
		},
		&Pipeline{
			Name:        "",
			Description: "",
			Jenkinsfile: "node{echo 'hello'}",
		},
		&Pipeline{
			Name:              "",
			Description:       "",
			Jenkinsfile:       "node{echo 'hello'}",
			DisableConcurrent: true,
		},
	}
	for _, input := range inputs {
		outputString, err := createPipelineConfigXml(input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parsePipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_NoScmPipelineConfig_Discarder(t *testing.T) {
	inputs := []*Pipeline{
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			Discarder: &DiscarderProperty{
				"3", "5",
			},
		},
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			Discarder: &DiscarderProperty{
				"3", "",
			},
		},
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			Discarder: &DiscarderProperty{
				"", "21321",
			},
		},
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			Discarder: &DiscarderProperty{
				"", "",
			},
		},
	}
	for _, input := range inputs {
		outputString, err := createPipelineConfigXml(input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parsePipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_NoScmPipelineConfig_Param(t *testing.T) {
	inputs := []*Pipeline{
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			Parameters: []*Parameter{
				&Parameter{
					Name:         "d",
					DefaultValue: "a\nb",
					Type:         "choice",
					Description:  "fortest",
				},
			},
		},
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			Parameters: []*Parameter{
				&Parameter{
					Name:         "a",
					DefaultValue: "abc",
					Type:         "string",
					Description:  "fortest",
				},
				&Parameter{
					Name:         "b",
					DefaultValue: "false",
					Type:         "boolean",
					Description:  "fortest",
				},
				&Parameter{
					Name:         "c",
					DefaultValue: "password \n aaa",
					Type:         "text",
					Description:  "fortest",
				},
				&Parameter{
					Name:         "d",
					DefaultValue: "a\nb",
					Type:         "choice",
					Description:  "fortest",
				},
			},
		},
	}
	for _, input := range inputs {
		outputString, err := createPipelineConfigXml(input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parsePipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_NoScmPipelineConfig_Trigger(t *testing.T) {
	inputs := []*Pipeline{
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			TimerTrigger: &TimerTrigger{
				Cron: "1 1 1 * * *",
			},
		},

		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			RemoteTrigger: &RemoteTrigger{
				Token: "abc",
			},
		},
		&Pipeline{
			Name:        "",
			Description: "for test",
			Jenkinsfile: "node{echo 'hello'}",
			TimerTrigger: &TimerTrigger{
				Cron: "1 1 1 * * *",
			},
			RemoteTrigger: &RemoteTrigger{
				Token: "abc",
			},
		},
	}

	for _, input := range inputs {
		outputString, err := createPipelineConfigXml(input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parsePipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_MultiBranchPipelineConfig(t *testing.T) {

	inputs := []*MultiBranchPipeline{
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "git",
			},
		},
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "github",
			},
		},
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "single_svn",
			},
		},
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "svn",
			},
		},
	}
	for _, input := range inputs {
		outputString, err := createMultiBranchPipelineConfigXml("", input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parseMultiBranchPipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_MultiBranchPipelineConfig_Discarder(t *testing.T) {

	inputs := []*MultiBranchPipeline{
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "git",
			},
			Discarder: &DiscarderProperty{
				DaysToKeep: "1",
				NumToKeep:  "2",
			},
		},
	}
	for _, input := range inputs {
		outputString, err := createMultiBranchPipelineConfigXml("", input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parseMultiBranchPipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_MultiBranchPipelineConfig_TimerTrigger(t *testing.T) {
	inputs := []*MultiBranchPipeline{
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "git",
			},
			TimerTrigger: &TimerTrigger{
				Interval: "12345566",
			},
		},
	}
	for _, input := range inputs {
		outputString, err := createMultiBranchPipelineConfigXml("", input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parseMultiBranchPipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_MultiBranchPipelineConfig_Source(t *testing.T) {

	inputs := []*MultiBranchPipeline{
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "git",
			},
			TimerTrigger: &TimerTrigger{
				Interval: "12345566",
			},
		},
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "github",
			},
			TimerTrigger: &TimerTrigger{
				Interval: "12345566",
			},
		},

		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "svn",
			},
			TimerTrigger: &TimerTrigger{
				Interval: "12345566",
			},
		},
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "single_svn",
			},
			TimerTrigger: &TimerTrigger{
				Interval: "12345566",
			},
		},
	}
	jsonByte, _ := json.Marshal(&GitSource{
		Url:              "https://github.com/kubesphere/devops",
		CredentialId:     "git",
		DiscoverBranches: true,
	})
	json.Unmarshal(jsonByte, &inputs[0].Source.Define)

	jsonByte, _ = json.Marshal(&GithubSource{
		Owner:                "kubesphere",
		Repo:                 "devops",
		CredentialId:         "github",
		ApiUri:               "https://api.github.com",
		DiscoverBranches:     1,
		DiscoverPRFromOrigin: 2,
		DiscoverPRFromForks: &GithubDiscoverPRFromForks{
			Strategy: 1,
			Trust:    1,
		},
	})
	json.Unmarshal(jsonByte, &inputs[1].Source.Define)

	jsonByte, _ = json.Marshal(&SvnSource{
		Remote:       "https://api.svn.com/bcd",
		CredentialId: "svn",
		Excludes:     "truck",
		Includes:     "tag/*",
	})
	json.Unmarshal(jsonByte, &inputs[2].Source.Define)

	jsonByte, _ = json.Marshal(&SingleSvnSource{
		Remote:       "https://api.svn.com/bcd",
		CredentialId: "svn",
	})
	json.Unmarshal(jsonByte, &inputs[3].Source.Define)
	for _, input := range inputs {
		outputString, err := createMultiBranchPipelineConfigXml("", input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parseMultiBranchPipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}
}

func Test_MultiBranchPipelineCloneConfig(t *testing.T) {

	inputs := []*MultiBranchPipeline{
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "git",
			},
		},
		&MultiBranchPipeline{
			Name:        "",
			Description: "for test",
			ScriptPath:  "Jenkinsfile",
			Source: &Source{
				Type: "github",
			},
		},
	}
	jsonByte, _ := json.Marshal(&GitSource{
		Url:              "https://github.com/kubesphere/devops",
		CredentialId:     "git",
		DiscoverBranches: true,
		CloneOption: &GitCloneOption{
			Shallow: false,
			Depth:   3,
			Timeout: 20,
		},
	})
	json.Unmarshal(jsonByte, &inputs[0].Source.Define)

	jsonByte, _ = json.Marshal(&GithubSource{
		Owner:                "kubesphere",
		Repo:                 "devops",
		CredentialId:         "github",
		ApiUri:               "https://api.github.com",
		DiscoverBranches:     1,
		DiscoverPRFromOrigin: 2,
		DiscoverPRFromForks: &GithubDiscoverPRFromForks{
			Strategy: 1,
			Trust:    1,
		},
		CloneOption: &GitCloneOption{
			Shallow: false,
			Depth:   3,
			Timeout: 20,
		},
	})
	json.Unmarshal(jsonByte, &inputs[1].Source.Define)

	for _, input := range inputs {
		outputString, err := createMultiBranchPipelineConfigXml("", input)
		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		output, err := parseMultiBranchPipelineConfigXml(outputString)

		if err != nil {
			t.Fatalf("should not get error %+v", err)
		}
		if !reflect.DeepEqual(input, output) {
			t.Fatalf("input [%+v] output [%+v] should equal ", input, output)
		}
	}

}
