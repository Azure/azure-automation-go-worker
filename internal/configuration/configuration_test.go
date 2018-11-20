// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

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
		t.Fatalf("unexpected error while loading configuration : %v", err)
	}

	if GetWorkerVersion() != testConfig.WorkerVersion {
		t.Fatal("unexpected configuration value")
	}
}

func TestSetConfiguration(t *testing.T) {
	clearConfiguration()
	config := GetConfiguration()
	config.Component = Component_worker
	SetConfiguration(&config)

	updatedConfig := GetConfiguration()
	if updatedConfig.Component != Component_worker {
		t.Fatal("unexpected configuration value")
	}

	updatedConfig.Component = Component_sandbox
	SetConfiguration(&updatedConfig)
	updatedConfig = GetConfiguration()
	if updatedConfig.Component != Component_sandbox {
		t.Fatal("unexpected configuration value")
	}
}
