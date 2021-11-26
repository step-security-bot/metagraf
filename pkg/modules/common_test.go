package modules

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
