package enviroment

import (
	"os"
	"text/template"
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
}

// CreateTemplateRds creates the terraform files from environment's template folder
func CreateTemplateRds() error {

	RdsTemplate, err := downloadTemplate("https://raw.githubusercontent.com/ministryofjustice/cloud-platform-terraform-rds-instance/add-template/template/rds.tmpl")
	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.New("rds").Parse(RdsTemplate))

	rdsValues, err := templateRdsSetValues()
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(os.Stdout, rdsValues)
	if err != nil {
		panic(err)
	}

	return nil
}

func templateRdsSetValues() (*templateRds, error) {
	values := templateRds{}

	err := validatePath()
	if err != nil {
		outsidePath := promptYesNo{label: "WARNING: You are outside the cloud-platform environment. Do you want to continue and render templates on the screen?", defaultValue: 0}
		err = outsidePath.promptyesNo()
		if err != nil {
			return nil, err
		}
	}

	namespaces, err := GetNamespacesFromGH()
	if err != nil {
		panic(err)
	}

	// spew.Dump(environments)

	namespaceName, err := promptSelectNamespaces(namespaces)
	if err != nil {
		return nil, err
	}

	metadata := MetaDataFromGH{namespace: namespaceName}
	err = metadata.GetEnvironmentsMetadataFromGH()
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
