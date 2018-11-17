// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
	"github.com/Azure/azure-automation-go-worker/main/worker/sandbox"
	"github.com/Azure/azure-extension-foundation/httputil"
	"os"
	"time"
)

type Worker struct {
	jrdsPollingFrequency time.Duration
	jrdsClient           JrdsClient
	sandboxCollection    map[string]*sandbox.Sandbox
}

type JrdsClient interface {
	GetSandboxActions(sandboxAction *jrds.SandboxActions) error
}

// NewWorker creates a new hybrid worker
func NewWorker(client JrdsClient) Worker {
	return Worker{jrdsClient: client,
		jrdsPollingFrequency: time.Duration(int64(time.Second) * configuration.GetJrdsPollingFrequencyInSeconds()),
		sandboxCollection:    make(map[string]*sandbox.Sandbox)}
}

// Start starts the main loop of the hybrid worker which polls JRDS for sandbox actions
func (worker *Worker) Start() {
	for {
		worker.routine()
		time.Sleep(worker.jrdsPollingFrequency)
	}
}

// routine defines the main polling logic and actions to perform when a new sandbox has to be created
func (worker *Worker) routine() {
	// get sandbox actions
	actions := jrds.SandboxActions{}
	err := worker.jrdsClient.GetSandboxActions(&actions)
	if err != nil {
		tracer.LogWorkerErrorGettingSandboxActions(err)
		return
	}

	// start a new sandbox for each actions returned by jrds
	tracer.LogDebugTrace(fmt.Sprintf("Get sandbox action. Found %v action(s).", len(actions.Value)))
	if len(actions.Value) > 0 {
		tracer.LogWorkerSandboxActionsFound(actions)
		for _, action := range actions.Value {
			sandboxId := *action.SandboxId
			if _, tracked := worker.sandboxCollection[sandboxId]; tracked {
				// sandbox already tracked, skip sandbox creation
				continue
			}

			sandbox := sandbox.NewSandbox(sandboxId)
			worker.sandboxCollection[sandbox.Id] = &sandbox
			err := createAndStartSandbox(&sandbox)
			if err != nil {
				tracer.LogWorkerFailedToCreateSandbox(err)
			}
		}
	}
}

var createAndStartSandbox = func(sandbox *sandbox.Sandbox) error {
	err := sandbox.CreateBaseDirectory()
	if err != nil {
		return err
	}

	err = sandbox.Start()
	if err != nil {
		return err
	}

	go monitorSandbox(sandbox)
	return nil
}

var monitorSandbox = func(sandbox *sandbox.Sandbox) {
	for sandbox.IsAlive() {
		time.Sleep(time.Millisecond * 100) // TODO: temporary until async output is implemented
	}

	sandbox.Cleanup()
}

func main() {
	// always load configuration and initialize tracer before anything else
	err := configuration.LoadConfiguration(os.Args[1])
	if err != nil {
		panic(err)
	}

	httpClient := httputil.NewSecureHttpClientWithCertificates(configuration.GetJrdsCertificatePath(), configuration.GetJrdsKeyPath(), httputil.LinearRetryThrice)
	jrdsClient := jrds.NewJrdsClient(httpClient, configuration.GetJrdsBaseUri(), configuration.GetAccountId(), configuration.GetHybridWorkerGroupName())
	tracer.InitializeTracer(&jrdsClient)

	tracer.LogWorkerStarting()
	worker := NewWorker(&jrdsClient)
	worker.Start()
}
