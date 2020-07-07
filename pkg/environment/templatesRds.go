package enviroment

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

type templateRds struct {
	IsProduction          bool
	EnvironmentName       string
	BusinessUnit          string
	Application           string
	Namespace             string
	InfrastructureSupport string
	RdsModuleName         string
	TeamName              string
	validPath             bool
}

// CreateTemplateRds creates the terraform files from environment's template folder
func CreateTemplateRds(cmd *cobra.Command, args []string) error {

	RdsTemplate, err := downloadTemplate("https://raw.githubusercontent.com/ministryofjustice/cloud-platform-terraform-rds-instance/main/template/rds.tmpl")
	if err != nil {
		return (err)
	}

	rdsValues, err := templateRdsSetValues()
	if err != nil {
		return (err)
	}

	tpl := template.Must(template.New("rds").Parse(RdsTemplate))

	if rdsValues.validPath == true {
		outputPath := fmt.Sprintf("namespaces/live-1.cloud-platform.service.justice.gov.uk/%s/resources/rds.tf", rdsValues.Namespace)
		f, _ := outputFileWriter(outputPath)
		err = tpl.Execute(f, rdsValues)
		if err != nil {
			return (err)
		}
		fmt.Printf("RDS instance created in the file namespaces/live-1.cloud-platform.service.justice.gov.uk/%s/resources/rds.tf\n", rdsValues.Namespace)
		fmt.Printf("Please review the file before raising PR\n")
	} else {
		err = tpl.Execute(os.Stdout, rdsValues)
		if err != nil {
			return (err)
		}
	}

	return nil
}

func templateRdsSetValues() (*templateRds, error) {
	values := templateRds{}

	validPath, err := validPath()
	if err != nil {
		return nil, err
	}
	values.validPath = validPath

	namespaces, err := GetNamespacesFromGH()
	if err != nil {
		return nil, err
	}

	// spew.Dump(environments)

	namespaceName, err := promptSelectNamespaces(namespaces)
	if err != nil {
		return nil, err
	}
	values.Namespace = namespaceName

	metadata := MetaDataFromNamespaceGH{namespace: namespaceName}
	err = metadata.GetEnvironmentsMetadataFromNamespaceGH()
	if err != nil {
		return nil, err
	}

	rdsModuleName := promptString{label: "Module name for RDS?", defaultValue: "rds"}
	rdsModuleName.promptString()
	if err != nil {
		return nil, err
	}

	environmentName := promptString{label: "Environment?", defaultValue: metadata.environmentName}
	environmentName.promptString()
	if err != nil {
		return nil, err
	}

	isProduction := promptYesNo{label: "Is Production?", defaultValue: 0}
	if metadata.isProduction == "false" {
		isProduction.defaultValue = 1
	}

	err = isProduction.promptyesNo()
	if err != nil {
		return nil, err
	}

	application := promptString{label: "Application name?", defaultValue: metadata.application}
	application.promptString()
	if err != nil {
		return nil, err
	}

	businessUnit := promptString{label: "Business Unit?", defaultValue: metadata.businessUnit}
	businessUnit.promptString()
	if err != nil {
		return nil, err
	}

	teamName := promptString{label: "Team's name", defaultValue: ""}
	teamName.promptString()
	if err != nil {
		return nil, err
	}

	email := promptString{label: "Team's email", defaultValue: metadata.ownerEmail}
	email.promptString()
	if err != nil {
		return nil, err
	}

	values.Application = application.value
	values.BusinessUnit = businessUnit.value
	values.EnvironmentName = environmentName.value
	values.IsProduction = isProduction.value
	values.RdsModuleName = rdsModuleName.value
	values.InfrastructureSupport = email.value
	values.TeamName = teamName.value

	return &values, nil
}