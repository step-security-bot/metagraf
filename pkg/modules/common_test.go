package modules

import (
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func generateMg(required bool) (mg metagraf.MetaGraf ) {
	mg.Kind = "metagraf"
	mg.Metadata.Name = "test"
	mg.Spec.Version = "1.0.0"

	envVar := generateEnvironmentVar("FOO", required)
	mg.Spec.Environment.Local = append(mg.Spec.Environment.Local, envVar)

	return mg
}

func generateEnvironmentVar(name string, required bool) (envVar metagraf.EnvironmentVar) {
	envVar = metagraf.EnvironmentVar{
		Name: name,
		Required: required,
		Type: "string",
		Description: "test",
	}
	return envVar
}

var emptyProps = metagraf.MGProperties{}
var singleProps = metagraf.MGProperties{
	"FOO": metagraf.MGProperty{
		Source: "local",
		Key: "FOO",
		Value: "BAR",
	},
}
var singleWrongProps = metagraf.MGProperties{
	"FOO2": metagraf.MGProperty{
		Source: "local",
		Key: "FOO2",
		Value: "BAZ",
	},
}
var doubleProps = metagraf.MGProperties{
	"FOO": metagraf.MGProperty{
		Source: "local",
		Key: "FOO",
		Value: "BAR",
	},
	"FOO2": metagraf.MGProperty{
		Source: "local",
		Key: "FOO2",
		Value: "BAZ",
	},
}

func Test_MergeLabels(t *testing.T) {
	first := make(map[string]string)
	first["label1"] = "value1"

	second := make(map[string]string)
	second["label2"] = "value2"

	actualResult := MergeLabels(first, second)

	expectedValue := make(map[string]string)
	expectedValue["label1"] = "value1"
	expectedValue["label2"] = "value2"

	assert.Equal(t, actualResult, expectedValue)
}

func Test_MergeLabelsOverride(t *testing.T) {
	first := make(map[string]string)
	first["label1"] = "value1"

	second := make(map[string]string)
	second["label1"] = "override"
	second["label2"] = "value2"

	actualResult := MergeLabels(first, second)

	expectedValue := make(map[string]string)
	expectedValue["label1"] = "override"
	expectedValue["label2"] = "value2"

	assert.Equal(t, actualResult, expectedValue)
}

// env required, value set = output
// env not required, value not set = don't output
// env not required, value set = output
// env required, value not set = error

func Test_RequiredEnvVarWithValue(t *testing.T) {
	// env required, value set = output
	mg := generateMg(true)

	actualResult, err := GetEnvVars(&mg, singleProps)
	expectedResult := []corev1.EnvVar{
		{
			Name: "MG_APP_NAME",
			Value: "testv1",
		},
		{
			Name: "MG_API_VERSION",
			Value: "1.0.0",
		},
		{
			Name: "FOO",
			Value: "BAR",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_RequiredEnvVarWithMissingValue(t *testing.T) {
	// env required, value not set = error
	mg := generateMg(true)

	actualResult, err := GetEnvVars(&mg, emptyProps)
	var expectedResult []corev1.EnvVar = nil

	assert.Error(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_MultipleRequiredWithValues(t *testing.T) {
	// env required, value set = output
	mg := generateMg(true)
	secondRequired := generateEnvironmentVar("FOO2", true)
	mg.Spec.Environment.Local = append(mg.Spec.Environment.Local, secondRequired)

	actualResult, err := GetEnvVars(&mg, doubleProps)
	expectedResult := []corev1.EnvVar{
		{
			Name: "MG_APP_NAME",
			Value: "testv1",
		},
		{
			Name: "MG_API_VERSION",
			Value: "1.0.0",
		},
		{
			Name: "FOO",
			Value: "BAR",
		},
		{
			Name: "FOO2",
			Value: "BAZ",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_MultipleRequiredWithOneMissingValue(t *testing.T) {
	// env required, value set = output
	mg := generateMg(true)
	secondRequired := generateEnvironmentVar("FOO2", true)
	mg.Spec.Environment.Local = append(mg.Spec.Environment.Local, secondRequired)

	actualResult, err := GetEnvVars(&mg, singleProps)
	var expectedResult []corev1.EnvVar = nil

	assert.Error(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_NotRequiredEnvVarWithValue(t *testing.T) {
	// env not required, value set = output
	mg := generateMg(false)

	actualResult, err := GetEnvVars(&mg, singleProps)
	expectedResult := []corev1.EnvVar{
		{
			Name: "MG_APP_NAME",
			Value: "testv1",
		},
		{
			Name: "MG_API_VERSION",
			Value: "1.0.0",
		},
		{
			Name: "FOO",
			Value: "BAR",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_NotRequiredEnvVarWithoutValue(t *testing.T) {
	// env not required, value not set = don't output
	mg := generateMg(false)

	actualResult, err := GetEnvVars(&mg, emptyProps)
	expectedResult := []corev1.EnvVar{
		{
			Name: "MG_APP_NAME",
			Value: "testv1",
		},
		{
			Name: "MG_API_VERSION",
			Value: "1.0.0",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_CombinedRequiredAndNotRequiredWithValues(t *testing.T) {
	// env required, value set = output
	// env not required, value set = output
	mg := generateMg(true)
	secondRequired := generateEnvironmentVar("FOO2", false)
	mg.Spec.Environment.Local = append(mg.Spec.Environment.Local, secondRequired)

	actualResult, err := GetEnvVars(&mg, doubleProps)
	expectedResult := []corev1.EnvVar{
		{
			Name: "MG_APP_NAME",
			Value: "testv1",
		},
		{
			Name: "MG_API_VERSION",
			Value: "1.0.0",
		},
		{
			Name: "FOO",
			Value: "BAR",
		},
		{
			Name: "FOO2",
			Value: "BAZ",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_CombinedRequiredAndNotRequiredWithOnlyRequiredValue(t *testing.T) {
	// env required, value set = output
	// env not required, value not set = don't output
	mg := generateMg(true)
	secondRequired := generateEnvironmentVar("FOO2", false)
	mg.Spec.Environment.Local = append(mg.Spec.Environment.Local, secondRequired)

	actualResult, err := GetEnvVars(&mg, singleProps)
	expectedResult := []corev1.EnvVar{
		{
			Name: "MG_APP_NAME",
			Value: "testv1",
		},
		{
			Name: "MG_API_VERSION",
			Value: "1.0.0",
		},
		{
			Name: "FOO",
			Value: "BAR",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func Test_CombinedRequiredAndNotRequiredWithMissingRequiredValue(t *testing.T) {
	// env required, value not set = error
	// env not required, value set = output
	mg := generateMg(true)
	secondRequired := generateEnvironmentVar("FOO2", false)
	mg.Spec.Environment.Local = append(mg.Spec.Environment.Local, secondRequired)

	actualResult, err := GetEnvVars(&mg, singleWrongProps)
	var expectedResult []corev1.EnvVar = nil

	assert.Error(t, err)
	assert.Equal(t, expectedResult, actualResult)
}