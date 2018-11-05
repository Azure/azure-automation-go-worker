// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/main/worker/sandbox"
	"testing"
	"time"
)

type jrdsMock struct {
	getSandboxAction_f func(sandboxAction *jrds.SandboxActions) error
}

func (jrds *jrdsMock) GetSandboxActions(sandboxAction *jrds.SandboxActions) error {
	return jrds.getSandboxAction_f(sandboxAction)
}

func TestWorker_Start_CreatesSingleSandboxForMultipleActionForSameSandboxId(t *testing.T) {
	Setup()

	testSandboxId := "b21b3d28-8d20-42d5-8b53-a7cbd0c97886"
	newSbxCount := 0

	jrdsMock := jrdsMock{}
	jrdsMock.getSandboxAction_f = func(sandboxAction *jrds.SandboxActions) error {
		action := jrds.SandboxAction{SandboxId: &testSandboxId}
		sandboxAction.Value = []jrds.SandboxAction{action, action, action}
		return nil
	}

	sandbox.NewSandbox = func(sandboxId string) sandbox.Sandbox {
		if sandboxId != testSandboxId {
			t.Fatal("invalid sandbox id")
		}
		newSbxCount += 1
		return sandbox.Sandbox{Id: sandboxId}
	}
	worker := NewWorker(&jrdsMock)
	worker.routine()

	if _, ok := worker.sandboxCollection[testSandboxId]; !ok {
		t.Fatal("sandbox not tracked")
	}
	if newSbxCount != 1 {
		t.Fatal("unexpected count of sandboxes created")
	}
}

func TestWorker_Start_CreatesMultipleSandboxesForMultipleActionForUniqueSandboxId(t *testing.T) {
	Setup()

	testSandboxIds := []string{"b21b3d28-8d20-42d5-8b53-a7cbd0c97886", "b21b3d28-8d20-42d5-8b53-a7cbd0c97887"}
	newSbxCount := 0

	jrdsMock := jrdsMock{}
	jrdsMock.getSandboxAction_f = func(sandboxAction *jrds.SandboxActions) error {
		actionA := jrds.SandboxAction{SandboxId: &testSandboxIds[0]}
		actionB := jrds.SandboxAction{SandboxId: &testSandboxIds[1]}
		sandboxAction.Value = []jrds.SandboxAction{actionA, actionB}
		return nil
	}

	sandbox.NewSandbox = func(sandboxId string) sandbox.Sandbox {
		newSbxCount += 1
		return sandbox.Sandbox{Id: sandboxId}
	}
	worker := NewWorker(&jrdsMock)
	worker.routine()

	for _, id := range testSandboxIds {
		if _, ok := worker.sandboxCollection[id]; !ok {
			t.Fatal("sandbox %v not tracked", id)
		}
	}

	if newSbxCount != len(testSandboxIds) {
		t.Fatal("unexpected count of sandboxes created")
	}
}

func TestWorker_Start_DoestNotCreateSandboxOnEmptySandboxActions(t *testing.T) {
	Setup()

	newSbxCount := 0
	jrdsMock := jrdsMock{}
	jrdsMock.getSandboxAction_f = func(sandboxAction *jrds.SandboxActions) error {
		sandboxAction.Value = []jrds.SandboxAction{}
		return nil
	}

	sandbox.NewSandbox = func(sandboxId string) sandbox.Sandbox {
		newSbxCount += 1
		return sandbox.Sandbox{Id: sandboxId}
	}
	worker := NewWorker(&jrdsMock)
	worker.routine()

	if newSbxCount != 0 {
		t.Fatal("unexpected count of sandboxes created")
	}
}

func Setup() {
	configuration.GetJrdsPollingFrequencyInSeconds = func() time.Duration {
		return 1
	}
	createAndStartSandbox = func(sandbox *sandbox.Sandbox) error {
		return nil
	}
}
