// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package job

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
	"github.com/Azure/azure-automation-go-worker/main/sandbox/runtime"
	"os"
	"path/filepath"
	"time"
)

type Job struct {
	Id string

	// runtime
	StartTime time.Time
	Completed bool

	// channels
	PendingActions chan PendingAction
	Exceptions     chan string

	jobData          jrds.JobData
	jobUpdatableData jrds.JobUpdatableData
	runbookData      jrds.RunbookData

	sandboxId        string
	workingDirectory string
	jrdsClient       jrdsClient
}

type jrdsClient interface {
	GetJobActions(sandboxId string, jobData *jrds.JobActions) error
	GetJobData(jobId string, jobData *jrds.JobData) error
	GetUpdatableJobData(jobId string, jobData *jrds.JobUpdatableData) error
	GetRunbookData(runbookVersionId string, runbookData *jrds.RunbookData) error
	AcknowledgeJobAction(sandboxId string, messageMetadata jrds.MessageMetadatas) error
	SetJobStatus(sandboxId string, jobId string, status int, isTermial bool, exception *string) error
	SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error
	UnloadJob(subscriptionId string, sandboxId string, jobId string, isTest bool, startTime time.Time, executionTimeInSeconds int) error
}

func NewJob(sandboxId string, jobData jrds.JobData, jrdsClient jrdsClient) Job {
	workingDirectory := filepath.Join(configuration.GetWorkingDirectory(), *jobData.JobId)
	err := os.MkdirAll(workingDirectory, 0750)
	panicOnError("Unable to create job working directory", err)

	return Job{
		Id:               *jobData.JobId,
		jobData:          jobData,
		sandboxId:        sandboxId,
		workingDirectory: workingDirectory,
		jrdsClient:       jrdsClient,
		StartTime:        time.Now(),
		Completed:        false,
		PendingActions:   make(chan PendingAction),
		Exceptions:       make(chan string)}
}

func (job *Job) Run() {
	err := loadJob(job)
	panicOnError(fmt.Sprintf("error loading job : %v", err), err)

	jobRuntime, err := initializeRuntime(job)
	panicOnError(fmt.Sprintf("error initializing jobRuntime %v", err), err)

	executeRunbook(jobRuntime, job)

	err = unloadJob(job)
	panicOnError(fmt.Sprintf("error unloading job : %v", err), err)
}

var loadJob = func(job *Job) error {
	setStatus(job, getActivatingStatus())

	jobUpdatableData := jrds.JobUpdatableData{}
	err := job.jrdsClient.GetUpdatableJobData(job.Id, &jobUpdatableData)
	if err != nil {
		return err
	}

	runbookData := jrds.RunbookData{}
	err = job.jrdsClient.GetRunbookData(*job.jobData.RunbookVersionId, &runbookData)
	if err != nil {
		return err
	}

	job.jobUpdatableData = jobUpdatableData
	job.runbookData = runbookData

	tracer.LogSandboxJobLoaded(job.sandboxId, job.Id)
	return nil
}

var initializeRuntime = func(job *Job) (*runtime.Runtime, error) {
	// create runbook
	runbook, err := runtime.NewRunbook(
		*job.runbookData.Name,
		*job.runbookData.RunbookVersionId,
		runtime.DefinitionKind(*job.runbookData.RunbookDefinitionKind),
		*job.runbookData.Definition)
	if err != nil {
		return nil, err
	}

	// create language; failed the job if the language isn't supported by the worker
	language, err := runtime.GetLanguage(runbook.Kind)
	if err != nil {
		return nil, err
	}

	// create runtime
	runtime := runtime.NewRuntime(language, runbook, job.jobData, job.workingDirectory)
	err = runtime.Initialize()
	if err != nil {
		return nil, err
	}

	return &runtime, nil
}

var executeRunbook = func(runtime *runtime.Runtime, job *Job) {
	// test if is the runtime supported by the os
	supp := runtime.IsSupported()
	if !supp {
		tracer.LogSandboxJobUnsupportedRunbookType(job.sandboxId, job.Id, fmt.Sprintf("Runtime not supported"))
		setStatus(job, getFailedStatus("Language not supported on this host."))
	}

	setStatus(job, getRunningStatus())

	streamHandler := NewStreamHandler(job.jrdsClient, job.Id, *job.jobData.RunbookVersionId)
	runtime.StartRunbookAsync(streamHandler.SetStream)

	// check pending action while job is running
	stopped := false
	for runtime.IsRunbookRunning() {
		if action, found := getPendingActions(job); found {
			if action.Enum == Stop {
				runtime.StopRunbook()
				stopped = true
				break
			}
		}
		time.Sleep(time.Millisecond * 10)
	}

	if stopped {
		setStatus(job, getStoppedStatus())
	} else {
		setStatus(job, getCompletedStatus())
	}

	job.Completed = true
}

var unloadJob = func(job *Job) error {
	executionTimeInSeconds := int((time.Now().Sub(job.StartTime)).Seconds())
	err := job.jrdsClient.UnloadJob(*job.jobData.SubscriptionId, job.sandboxId, job.Id, false, job.StartTime, executionTimeInSeconds)
	if err != nil {
		return err
	}

	tracer.LogSandboxJobUnloaded(job.sandboxId, job.Id)
	return nil
}

var getPendingActions = func(job *Job) (pendingAction PendingAction, found bool) {
	// read from channel without blocking
	select {
	case action := <-job.PendingActions:
		return action, true
	default:
	}

	return PendingAction{}, false
}

var setStatus = func(job *Job, jobstatus status) {
	err := job.jrdsClient.SetJobStatus(job.sandboxId, job.Id, jobstatus.enum, jobstatus.isTerminal, jobstatus.exception)
	panicOnError(fmt.Sprintf("error setting job status : %v", err), err)
}

var panicOnError = func(message string, err error) {
	if err != nil {
		panic(err)
	}
}
