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

const (
	enum_statusActivating = 1
	enum_statusRunning    = 2
	enum_statusCompleted  = 3
	enum_statusFailed     = 4
	enum_statusStopped    = 5
)

type Job struct {
	Id string

	// runtime
	StartTime time.Time
	Completed bool

	// channels
	PendingActions chan int
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
		PendingActions:   make(chan int),
		Exceptions:       make(chan string)}
}

func (job *Job) Run() {
	err := loadJob(job)
	panicOnError(fmt.Sprintf("error loading job %v", err), err)

	jobRuntime, err := initializeRuntime(job)
	panicOnError(fmt.Sprintf("error initializing jobRuntime %v", err), err)

	err = executeRunbook(jobRuntime, job)
	panicOnError(fmt.Sprintf("error executing runbook %v", err), err)
}

var loadJob = func(job *Job) error {
	err := job.jrdsClient.SetJobStatus(job.sandboxId, job.Id, enum_statusActivating, false, nil)
	if err != nil {
		return err
	}

	jobUpdatableData := jrds.JobUpdatableData{}
	err = job.jrdsClient.GetUpdatableJobData(job.Id, &jobUpdatableData)
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

	// test if is the runtime supported by the os
	supp := runtime.IsSupported()
	if !supp {
		tracer.LogErrorTrace("Runbook definition kind not supported")
	}

	return &runtime, nil
}

var executeRunbook = func(runtime *runtime.Runtime, job *Job) error {
	err := job.jrdsClient.SetJobStatus(job.sandboxId, job.Id, enum_statusRunning, false, nil)
	panicOnError(fmt.Sprintf("error setting job status %v", err), err)

	// temporary
	// running job
	runtime.StartRunbook()

	// temporary
	// check pending action while job is running
	if action, found := getPendingActions(job); found {
		if action == 5 {
			err = job.jrdsClient.SetJobStatus(job.sandboxId, job.Id, enum_statusStopped, true, nil)
			panicOnError(fmt.Sprintf("error stopping job %v", err), err)
			return nil
		}
	}

	err = job.jrdsClient.SetJobStatus(job.sandboxId, job.Id, enum_statusCompleted, true, nil)
	panicOnError(fmt.Sprintf("error setting job status %v", err), err)

	err = unloadJob(job)
	panicOnError(fmt.Sprintf("error unloading job %v", err), err)

	job.Completed = true
	return nil
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

var panicOnError = func(message string, err error) {
	if err != nil {
		panic(err)
	}
}

var getPendingActions = func(job *Job) (pendingAction int, found bool) {
	select {
	case action := <-job.PendingActions:
		return action, true
	default:
	}

	return -1, false
}
