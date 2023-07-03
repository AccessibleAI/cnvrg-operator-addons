package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/assert"
)

var (
	DeployPath = "../"
	cnvrgApp   = CnvrgAppSpec{}
	values     = map[string]string{"clusterDomain": "aws.dilerous.cloud",
		"controlPlane.image":         "cnvrg/app:v4.7.85",
		"networking.ingress.type":    "ingress",
		"controlPlane.hyper.enabled": "false"}
)

func TestChartTemplateRenders(t *testing.T) {
	// Path to the helm chart we will test
	helmChartPath := DeployPath

	// Setup the args.
	// For this test, we will set the following input values:
	// equivalent of passing --set on helm template comman
	options := &helm.Options{
		SetValues: values,
	}

	render := helm.RenderTemplate(
		t, options, helmChartPath, "cnvrg",
		[]string{"templates/cap.yml", "templates/hooks.yml", "templates/operator.yml"})

	helm.UnmarshalK8SYaml(t, render, &cnvrgApp)

	assert.Equal(t, "aws.dilerous.cloud", cnvrgApp.Spec.ClusterDomain)
	assert.Equal(t, "cnvrg/app:v4.7.85", cnvrgApp.Spec.ControlPlane.Image)
	assert.Equal(t, "ingress", cnvrgApp.Spec.Networking.Ingress.Type)
	assert.False(t, false, cnvrgApp.Spec.ControlPlane.Hyper.Enabled)
}
