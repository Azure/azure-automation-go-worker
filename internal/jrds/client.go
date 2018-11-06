// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package jrds

import (
	"encoding/json"
	"fmt"
	"time"
)

type JrdsClient struct {
	baseUri         string
	accountId       string
	workerGroupName string
	workerVersion   string
	protocolVersion string
	client          httpClient
}

const (
	accept_headerKey      = "Accept"
	contenttype_headerKey = "Content-Type"
	conection_headerKey   = "Connection"
	useragent_headerKey   = "User-Agent"

	appjson_headerValue   = "application/json"
	keepalive_headerValue = "keep-alive"

	datetimeFormat = "2006-01-02T15:04:05.999999"
)

type httpClient interface {
	Get(url string, headers map[string]string) (responseCode int, body []byte, err error)
	Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
}

func NewJrdsClient(client httpClient, baseUri string, accountId string, workerGroupName string) JrdsClient {
	return JrdsClient{baseUri: baseUri, client: client, accountId: accountId, workerGroupName: workerGroupName, protocolVersion: "1.0", workerVersion: "2.0.0.0"}
}

func (jrds *JrdsClient) GetSandboxActions(sandboxAction *SandboxActions) error {
	url := fmt.Sprintf("%s/automationAccounts/%s/Sandboxes/GetSandboxActions?HybridWorkerGroupName=%s&api-version=%s", jrds.baseUri, jrds.accountId, jrds.workerGroupName, jrds.protocolVersion)
	err := jrds.issueGetRequest(url, sandboxAction)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetJobActions(sandboxId string, jobActions *JobActions) error {
	url := fmt.Sprintf("%s/automationAccounts/%s/Sandboxes/%s/jobs/getJobActions?api-version=%s", jrds.baseUri, jrds.accountId, sandboxId, jrds.protocolVersion)
	err := jrds.issueGetRequest(url, jobActions)
	if err != nil {
		return err
	}

	actions := (*jobActions).Value
	if len(actions) > 0 {
		var arr []MessageMetadata
		for _, jobaction := range actions {
			arr = append(arr, *jobaction.MessageMetadata)
		}

		metadatas := MessageMetadatas{arr}
		err = jrds.AcknowledgeJobAction(sandboxId, metadatas)
		if err != nil {
			fmt.Printf("error getting messageMetadata %v", err)
		}
	}

	return nil
}

func (jrds *JrdsClient) GetJobData(jobId string, jobData *JobData) error {
	url := fmt.Sprintf("%s/automationAccounts/%s/jobs/%s?api-version=%s", jrds.baseUri, jrds.accountId, jobId, jrds.protocolVersion)
	err := jrds.issueGetRequest(url, jobData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetUpdatableJobData(jobId string, jobData *JobUpdatableData) error {
	url := fmt.Sprintf("%s/automationAccounts/%s/jobs/%s?api-version=%s", jrds.baseUri, jrds.accountId, jobId, jrds.protocolVersion)
	err := jrds.issueGetRequest(url, jobData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetRunbookData(runbookVersionId string, runbookData *RunbookData) error {
	url := fmt.Sprintf("%s/automationAccounts/%s/runbooks/%s?api-version=%s", jrds.baseUri, jrds.accountId, runbookVersionId, jrds.protocolVersion)
	err := jrds.issueGetRequest(url, runbookData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) AcknowledgeJobAction(sandboxId string, messageMetadata MessageMetadatas) error {
	url := fmt.Sprintf("%s/automationAccounts/%s/Sandboxes/%s/jobs/AcknowledgeJobActions?api-version=%s", jrds.baseUri, jrds.accountId, sandboxId, jrds.protocolVersion)
	err := jrds.issuePostRequest(url, messageMetadata, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) SetJobStatus(sandboxId string, jobId string, status int, isTermial bool, exception *string) error {
	jobStatus := JobStatus{JobStatus: &status, Exception: exception, IsFinalStatus: &isTermial}
	url := fmt.Sprintf("%s/automationAccounts/%s/Sandboxes/%s/jobs/%s/ChangeStatus?api-version=%s", jrds.baseUri, jrds.accountId, sandboxId, jobId, jrds.protocolVersion)
	err := jrds.issuePostRequest(url, jobStatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error {
	recordTime := time.Now().Format(datetimeFormat)
	stream := Stream{AccountId: &jrds.accountId, JobId: &jobId, RecordTime: &recordTime, RunbookVersionId: &runbookVersionId, SequenceNumber: &sequence, StreamRecord: nil, StreamRecordText: &text, Type: &streamType} // Todo : datetime
	url := fmt.Sprintf("%s/automationAccounts/%s/jobs/%s/postJobStream?api-version=%s", jrds.baseUri, jrds.accountId, jobId, jrds.protocolVersion)
	err := jrds.issuePostRequest(url, stream, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) SetLog(eventId int, activityId string, logType int, args ...string) error {
	log := Log{EventId: &eventId, Arguments: &args, LogType: &logType, ActivityId: &activityId}
	url := fmt.Sprintf("%s/automationAccounts/%s/logs?api-version=%s", jrds.baseUri, jrds.accountId, jrds.protocolVersion)
	err := jrds.issuePostRequest(url, log, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) UnloadJob(subscriptionId string, sandboxId string, jobId string, isTest bool, startTime time.Time, executionTimeInSeconds int) error {
	jobStartTime := startTime.Format(datetimeFormat)
	payload := UnloadJob{JobId: &jobId, IsTest: &isTest, StartTime: &jobStartTime, SubscriptionId: &subscriptionId, ExecutionTimeInSeconds: &executionTimeInSeconds}
	url := fmt.Sprintf("%s/automationAccounts/%s/Sandboxes/%s/jobs/%s/unload?api-version=%s", jrds.baseUri, jrds.accountId, sandboxId, jobId, jrds.protocolVersion)
	err := jrds.issuePostRequest(url, payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds JrdsClient) getDefaultHeaders() map[string]string {
	return map[string]string{accept_headerKey: appjson_headerValue,
		conection_headerKey: keepalive_headerValue,
		useragent_headerKey: fmt.Sprintf("AzureAutomationHybridWorker/%s", jrds.workerVersion)}
}

func (jrds *JrdsClient) issuePostRequest(url string, payload interface{}, out interface{}) error {
	headers := jrds.getDefaultHeaders()
	headers[contenttype_headerKey] = appjson_headerValue

	var body []byte
	if payload != nil {
		out, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response: %v", err)
		}
		body = out
	}

	code, _, err := jrds.client.Post(url, headers, body)

	if err != nil {
		return NewRequestError(fmt.Sprintf("request error %v : %v\n", url, code))
	}

	if code == 401 {
		return NewRequestAuthorizationError(fmt.Sprintf("authorization error %v : %v\n", url, code))
	}

	if code != 200 {
		return NewRequestInvalidStatusError(fmt.Sprintf("invalid return code for %v : %v\n", url, code))
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("failed to unmarshal request response: %v", err)
		}
	}

	return err
}

func (jrds *JrdsClient) issueGetRequest(url string, out interface{}) error {
	code, body, err := jrds.client.Get(url, jrds.getDefaultHeaders())

	if err != nil {
		return NewRequestError(fmt.Sprintf("request error %v : %v\n", url, code))
	}

	if code == 401 {
		return NewRequestAuthorizationError(fmt.Sprintf("authorization error %v : %v\n", url, code))
	}

	if code != 200 {
		return NewRequestInvalidStatusError(fmt.Sprintf("invalid return code for %v : %v\n", url, code))
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("failed to unmarshal request response: %v", err)
		}
	}

	return err
}
