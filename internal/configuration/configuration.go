// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	EnvironmentConfigurationKey = "WORKERCONF"

	DEFAULT_empty                         = ""
	DEFAULT_workerVersion                 = "2.0.0"
	DEFAULT_sandboxExecutableName         = "sandbox"
	DEFAULT_jrdsPollingFrequencyInSeconds = 10
	DEFAULT_component                     = Component_worker
	DEFAULT_debugTraces                   = false

	Component_sandbox = "sandbox"
	Component_worker  = "worker"
)

type Configuration struct {
	JrdsCertificatePath string `json:"jrds_cert_path"`
	JrdsKeyPath         string `json:"jrds_key_path"`
	JrdsBaseUri         string `json:"jrds_base_uri"`

	AccountId              string `json:"account_id"`
	MachineId              string `json:"machine_id"`
	HybridWorkerGroupName  string `json:"hybrid_worker_group_name"`
	WorkerVersion          string `json:"worker_version"`
	WorkerWorkingDirectory string `json:"working_directory_path"`
	SandboxExecutablePath  string `json:"sandbox_executable_path"`

	JrdsPollingFrequency int  `json:"jrds_polling_frequency"`
	DebugTraces          bool `json:"debug_traces"`

	// runtime configuration
	Component string `json:"component"`
}

func LoadConfiguration(path string) error {
	configuration := getDefaultConfiguration()
	content, err := readDiskConfiguration(path)
	if err != nil {
		return err
	}
	err = DeserializeConfiguration(content, &configuration)
	if err != nil {
		return err
	}

	setConfiguration(&configuration)
	return nil
}

func SetConfiguration(configuration *Configuration) {
	setConfiguration(configuration)
}

func GetConfiguration() Configuration {
	return getEnvironmentConfiguration()
}

var readDiskConfiguration = func(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

var setConfiguration = func(config *Configuration) {
	configuration, err := SerializeConfiguration(config)
	if err != nil {
		panic("unable to serialize configuration from environment")
	}

	err = os.Setenv(EnvironmentConfigurationKey, string(configuration))
	if err != nil {
		panic("unable to set configuration to environment")
	}
}

var clearConfiguration = func() {
	os.Unsetenv(EnvironmentConfigurationKey)
}

var getEnvironmentConfiguration = func() Configuration {
	value, exists := os.LookupEnv(EnvironmentConfigurationKey)

	configuration := Configuration{}
	if exists {
		err := DeserializeConfiguration([]byte(value), &configuration)
		if err != nil {
			panic("unable to deserialize configuration from environment")
		}
	}
	return configuration
}

var SerializeConfiguration = func(configuration *Configuration) ([]byte, error) {
	return json.Marshal(configuration)
}

var DeserializeConfiguration = func(data []byte, configuration *Configuration) error {
	return json.Unmarshal(data, &configuration)
}

var getDefaultConfiguration = func() Configuration {
	return Configuration{
		JrdsCertificatePath:    DEFAULT_empty,
		JrdsKeyPath:            DEFAULT_empty,
		JrdsBaseUri:            DEFAULT_empty,
		AccountId:              DEFAULT_empty,
		MachineId:              DEFAULT_empty,
		HybridWorkerGroupName:  DEFAULT_empty,
		WorkerVersion:          DEFAULT_workerVersion,
		WorkerWorkingDirectory: DEFAULT_empty,
		SandboxExecutablePath:  DEFAULT_sandboxExecutableName,
		Component:              DEFAULT_component,
		DebugTraces:            DEFAULT_debugTraces}
}

var GetJrdsCertificatePath = func() string {
	config := getEnvironmentConfiguration()
	return config.JrdsCertificatePath
}

var GetJrdsKeyPath = func() string {
	config := getEnvironmentConfiguration()
	return config.JrdsKeyPath
}

var GetJrdsBaseUri = func() string {
	config := getEnvironmentConfiguration()
	return config.JrdsBaseUri
}

var GetAccountId = func() string {
	config := getEnvironmentConfiguration()
	return config.AccountId
}

var GetHybridWorkerGroupName = func() string {
	config := getEnvironmentConfiguration()
	return config.HybridWorkerGroupName
}

var GetWorkingDirectory = func() string {
	config := getEnvironmentConfiguration()
	return config.WorkerWorkingDirectory
}

var GetSandboxExecutablePath = func() string {
	config := getEnvironmentConfiguration()
	return config.SandboxExecutablePath
}

var GetWorkerVersion = func() string {
	config := getEnvironmentConfiguration()
	return config.WorkerVersion
}

var GetJrdsPollingFrequencyInSeconds = func() int64 {
	config := getEnvironmentConfiguration()
	return int64(config.JrdsPollingFrequency)
}

var GetComponent = func() string {
	config := getEnvironmentConfiguration()
	return config.Component
}

var GetDebugTraces = func() bool {
	config := getEnvironmentConfiguration()
	return config.DebugTraces
}
