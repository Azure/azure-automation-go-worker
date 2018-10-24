package configuration

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	testValue = "value"
	testPath  = ""
)

func TestLoadConfiguration_LoadsConfigurationAndGetsValue(t *testing.T) {
	clearConfiguration()

	testConfig := Configuration{JrdsCertificatePath: testValue}
	readDiskConfiguration = func(path string) ([]byte, error) {
		content, _ := json.Marshal(testConfig)
		return content, nil
	}

	err := LoadConfiguration(testPath)
	if err != nil {
		t.Fatal("unexpected error while loading configuration")
	}
	if GetJrdsCertificatePath() != testConfig.JrdsCertificatePath {
		t.Fatal("unexpected value")
	}
}

func TestLoadConfiguration_ReturnsErrorOnInvalidPath(t *testing.T) {
	clearConfiguration()
	readDiskConfiguration = func(path string) ([]byte, error) {
		return nil, fmt.Errorf("invalid path")
	}

	err := LoadConfiguration(testPath)
	if err == nil {
		t.Fatal("unexpected missing error from LoadConfiguration on invalid path")
	}
}

func TestLoadConfiguration_OverrideDefaultValues(t *testing.T) {
	clearConfiguration()
	testConfig := getDefaultConfiguration()
	readDiskConfiguration = func(path string) ([]byte, error) {
		testConfig.WorkerVersion = "99999"
		content, _ := json.Marshal(testConfig)
		return content, nil
	}

	err := LoadConfiguration(testPath)
	if err != nil {
		t.Fatal("unexpected error while loading configuration : %v", err)
	}

	if GetWorkerVersion() != testConfig.WorkerVersion {
		t.Fatal("unexpected configuration value")
	}
}