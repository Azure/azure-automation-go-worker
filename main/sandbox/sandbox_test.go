// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"github.com/Azure/azure-automation-go-worker/main/sandbox/job"
	"testing"
	"time"
)

var (
	sandboxId      = "ccb7bc90-20e9-4e5c-bbf0-f265c1de7000"
	jobId          = "ccb7bc90-20e9-4e5c-bbf0-f265c1de7111"
	subscriptionId = "ccb7bc90-20e9-4e5c-bbf0-f265c1de7222"
)

type jrdsMock struct {
	unloadJob_f func(subscriptionId string, sandboxId string, jobId string, isTest bool, startTime time.Time, executionTimeInSeconds int) error
}

func (jrds *jrdsMock) GetJobActions(sandboxId string, jobData *jrds.JobActions) error {
	panic("implement me")
}

func (jrds *jrdsMock) GetJobData(jobId string, jobData *jrds.JobData) error {
	panic("implement me")
}

func (jrds *jrdsMock) GetUpdatableJobData(jobId string, jobData *jrds.JobUpdatableData) error {
	panic("implement me")
}

func (jrds *jrdsMock) GetRunbookData(runbookVersionId string, runbookData *jrds.RunbookData) error {
	panic("implement me")
}

func (jrds *jrdsMock) AcknowledgeJobAction(sandboxId string, messageMetadata jrds.MessageMetadatas) error {
	panic("implement me")
}

func (jrds *jrdsMock) SetJobStatus(sandboxId string, jobId string, status int, isTermial bool, exception *string) error {
	panic("implement me")
}

func (jrds *jrdsMock) SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error {
	panic("implement me")
}

func (jrds *jrdsMock) SetLog(eventId int, activityId string, logType int, args ...string) error {
	panic("implement me")
}

func (jrds *jrdsMock) UnloadJob(subscriptionId string, sandboxId string, jobId string, isTest bool, startTime time.Time, executionTimeInSeconds int) error {
	return jrds.unloadJob_f(subscriptionId, sandboxId, jobId, isTest, startTime, executionTimeInSeconds)
}

func Test_CleanCompletedJobs_DoesNotCleanRunningJobs(t *testing.T) {
	// create sandbox
	sbx := NewSandbox(sandboxId, nil)

	//create job
	job := job.NewJob(sandboxId, jrds.JobData{JobId: &jobId}, nil)
	sbx.jobs[jobId] = &job

	stopTrackingCompletedJobs(&sbx)
	if sbx.jobs[jobId] != &job {
		t.Fatal("unexpected error : job is not tracked by sandbox")
	}
}

func Test_CleanCompletedJobs_CleansCompletedJobs(t *testing.T) {
	// create jrdsmock
	unloadCalled := false
	jrdsMock := jrdsMock{}
	jrdsMock.unloadJob_f = func(subscriptionId string, sandboxId string, jobId string, isTest bool, startTime time.Time, executionTimeInSeconds int) error {
		unloadCalled = true
		return nil
	}

	// create job
	job := job.NewJob(sandboxId, jrds.JobData{JobId: &jobId, SubscriptionId: &subscriptionId}, &jrdsMock)
	job.Completed = true
	job.StartTime = time.Now()

	// create sandbox
	sbx := NewSandbox(sandboxId, &jrdsMock)
	sbx.jobs[jobId] = &job

	stopTrackingCompletedJobs(&sbx)

	j := sbx.jobs[jobId]
	if j != nil {
		t.Fatal("unexpected error : job is still tracked by sandbox")
	}

	if !unloadCalled {
		t.Fatal("unexpected error : job was not unloaded")
	}
}
