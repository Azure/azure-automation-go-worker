// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
	"github.com/Azure/azure-automation-go-worker/main/sandbox/job"
	"github.com/Azure/azure-extension-foundation/httputil"
	"os"
	"time"
)

type Sandbox struct {
	id      string
	isAlive bool

	jrdsClient           jrdsClient
	jrdsPollingFrequency time.Duration

	jobs map[string]*job.Job
}

func NewSandbox(sandboxId string, jrdsClient jrdsClient) Sandbox {
	return Sandbox{id: sandboxId,
		isAlive:              true,
		jrdsClient:           jrdsClient,
		jrdsPollingFrequency: time.Duration(int64(time.Second) * configuration.GetJrdsPollingFrequencyInSeconds()),
		jobs:                 make(map[string]*job.Job, 1)}
}

type jrdsClient interface {
	GetJobActions(sandboxId string, jobData *jrds.JobActions) error
	GetJobData(jobId string, jobData *jrds.JobData) error
	GetUpdatableJobData(jobId string, jobData *jrds.JobUpdatableData) error
	GetRunbookData(runbookVersionId string, runbookData *jrds.RunbookData) error
	AcknowledgeJobAction(sandboxId string, messageMetadata jrds.MessageMetadatas) error
	SetJobStatus(sandboxId string, jobId string, status int, isTermial bool, exception *string) error
	SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error
	SetLog(eventId int, activityId string, logType int, args ...string) error
	UnloadJob(subscriptionId string, sandboxId string, jobId string, isTest bool, startTime time.Time, executionTimeInSeconds int) error
}

func (sandbox *Sandbox) Start() {
	for sandbox.isAlive {
		routine(sandbox)
		time.Sleep(sandbox.jrdsPollingFrequency)
	}
}

var routine = func(sandbox *Sandbox) {
	jobActions := jrds.JobActions{}
	err := sandbox.jrdsClient.GetJobActions(sandbox.id, &jobActions)
	if err != nil {
		sandbox.isAlive = false

		switch err.(type) {
		case *jrds.RequestAuthorizationError:
			tracer.LogSandboxJrdsClosureRequest(sandbox.id)
			return
		default:
			tracer.LogErrorTrace(err.Error())
			return
		}
	}

	tracer.LogDebugTrace(fmt.Sprintf("Get job action. Found %v new action(s).", len(jobActions.Value)))
	for _, action := range jobActions.Value {
		tracer.LogSandboxGetJobActions(&jobActions)

		jobData := jrds.JobData{}
		err := sandbox.jrdsClient.GetJobData(*action.JobId, &jobData)
		if err != nil {
			fmt.Printf("error getting jobData %v", err)
		}

		if (jobData.PendingAction != nil && *jobData.PendingAction == 1) ||
			(jobData.PendingAction == nil && *jobData.JobStatus == 1) ||
			(jobData.PendingAction == nil && *jobData.JobStatus == 2) {
			// new job
			job := job.NewJob(sandbox.id, jobData, sandbox.jrdsClient)
			sandbox.jobs[job.Id] = &job

			go job.Run()
		} else if jobData.PendingAction != nil {
			pendingAction := job.GetPendingAction(*jobData.PendingAction)
			// stop pending action
			if job, ok := sandbox.jobs[*jobData.JobId]; ok {
				job.PendingActions <- pendingAction
			}
		} else if jobData.PendingAction == nil {
			// no pending action
			tracer.LogDebugTrace("no pending action")
		} else {
			//unsupported pending action
			tracer.LogDebugTrace("unsupported pending action")
		}
	}

	stopTrackingCompletedJobs(sandbox)
}

var stopTrackingCompletedJobs = func(sandbox *Sandbox) {
	completedJob := make([]string, 1)
	for jobId, runningJob := range sandbox.jobs {
		if runningJob.Completed {
			completedJob = append(completedJob, jobId)
		}
	}

	// stop tracking jobs
	for _, jobId := range completedJob {
		delete(sandbox.jobs, jobId)
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("missing sandbox id parameter")
	}
	sandboxId := os.Args[1]

	httpClient := httputil.NewSecureHttpClient(configuration.GetJrdsCertificatePath(), configuration.GetJrdsKeyPath(), httputil.LinearRetryThrice)
	jrdsClient := jrds.NewJrdsClient(&httpClient, configuration.GetJrdsBaseUri(), configuration.GetAccountId(), configuration.GetHybridWorkerGroupName())
	tracer.InitializeTracer(&jrdsClient)

	tracer.LogSandboxStarting(sandboxId)
	sandbox := NewSandbox(sandboxId, &jrdsClient)
	sandbox.Start()
}
