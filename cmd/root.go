package cmd

import (
	"github.com/hpcugent/gpfsbeat/beater"

	cmd "github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
)

// Name of this beat
var Name = "gpfsbeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name})

/*
// RootCmd to handle beats cli
var RootCmd *cmd.BeatsRootCmd

func init() {
	settings := instance.Settings{
		Name: Name,
		//Processing:    processing.MakeDefaultSupport(true, processing.WithECS, processing.WithAgentMeta()),
		HasDashboards: false,
	}
	RootCmd = cmd.GenRootCmdWithSettings(beater.New, settings)

	// remove dashboard from export commands
	for _, cmd := range RootCmd.ExportCmd.Commands() {
		if cmd.Name() == "dashboard" {
			RootCmd.ExportCmd.RemoveCommand(cmd)
		}
	}

	// only add defined flags to setup command
	setup := RootCmd.SetupCmd
	setup.Short = "Setup Elasticsearch index template and pipelines"
	setup.Long = `This command does initial setup of the environment:
 * Index mapping template in Elasticsearch to ensure fields are mapped.
 * ILM Policy
`
	setup.ResetFlags()
	//setup.Flags().Bool(cmd.IndexManagementKey, false, "Setup all components related to Elasticsearch index management, including template, ilm policy and rollover alias")
	//setup.Flags().MarkDeprecated(cmd.TemplateKey, fmt.Sprintf("use --%s instead", cmd.IndexManagementKey))
	//setup.Flags().MarkDeprecated(cmd.ILMPolicyKey, fmt.Sprintf("use --%s instead", cmd.IndexManagementKey))
	//setup.Flags().Bool(cmd.TemplateKey, false, "Setup index template")
	//setup.Flags().Bool(cmd.ILMPolicyKey, false, "Setup ILM policy")
}
*/
