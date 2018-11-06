// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/internal/tracer"
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
	Id               string
	jobData          jrds.JobData
	jobUpdatableData jrds.JobUpdatableData
	runbookData      jrds.RunbookData

	sandboxId  string
	jrdsClient jrdsClient

	// runtime
	StartTime time.Time
	Completed bool

	// channels
	PendingActions chan int
	Exceptions     chan string
}

func NewJob(sandboxId string, jobData jrds.JobData, jrdsClient jrdsClient) Job {
	return Job{
		Id:             *jobData.JobId,
		jobData:        jobData,
		sandboxId:      sandboxId,
		jrdsClient:     jrdsClient,
		StartTime:      time.Now(),
		Completed:      false,
		PendingActions: make(chan int),
		Exceptions:     make(chan string)}
}

func (job *Job) Run() {

	err := loadJob(job)
	panicOnError(fmt.Sprintf("error loading job %v", err), err)

	err = initializeRuntime(job)
	panicOnError(fmt.Sprintf("error initializing runtime %v", err), err)

	err = executeRunbook(job)
	panicOnError(fmt.Sprintf("error executing runbook %v", err), err)

	err = unloadJob(job)
	panicOnError(fmt.Sprintf("error unloading job %v", err), err)

	job.Completed = true
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

var initializeRuntime = func(job *Job) error {
	return nil // TODO
}

var executeRunbook = func(job *Job) error {
	err := job.jrdsClient.SetJobStatus(job.sandboxId, job.Id, enum_statusRunning, false, nil)
	panicOnError(fmt.Sprintf("error setting job status %v", err), err)

	// temporary
	// running job
	time.Sleep(time.Second * 10)

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
