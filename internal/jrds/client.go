package jrds

import (
	"encoding/json"
	"fmt"
	"time"
)

type JrdsClient struct {
	baseuri         string
	accountid       string
	workergroupname string
	workerversion   string
	protocolversion string
	client          httpClient
}

const (
	accept_headerKey      = "Accept"
	contenttype_headerKey = "Content-Type"
	conection_headerKey   = "Connection"
	useragent_headerKey   = "User-Agent"

	appjson_headerValue   = "application/json"
	keepalive_headerValue = "keep-alive"

	datetimeFormat = "2018-10-20T01:00:00.0000000"
)

type httpClient interface {
	Get(url string, headers map[string]string) (responseCode int, body []byte, err error)
	Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
}

func NewJrdsClient(client httpClient, baseuri string, accountid string, workergroupname string) JrdsClient {
	return JrdsClient{baseuri: baseuri, client: client, accountid: accountid, workergroupname: workergroupname, protocolversion: "1.0", workerversion: "2.0.0.0"}
}

func (jrds *JrdsClient) GetSandboxActions(sandboxAction *SandboxActions) error {
	url := fmt.Sprintf("%v/automationAccounts/%v/Sandboxes/GetSandboxActions?HybridWorkerGroupName=%v&api-version=%v", jrds.baseuri, jrds.accountid, jrds.workergroupname, jrds.protocolversion)
	err := jrds.issueGetRequest(url, sandboxAction)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetJobActions(sandboxId string, jobData *[]JobActions) error {
	url := fmt.Sprintf("%v/automationAccounts/%v/Sandboxes/jobs/getJobActions?%v?api-version=%v", jrds.baseuri, jrds.accountid, sandboxId, jrds.protocolversion)
	err := jrds.issueGetRequest(url, jobData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetJobData(jobId string, jobData *[]JobData) error {
	url := fmt.Sprintf("%v/automationAccounts/%v/jobs/%v?api-version=%v", jrds.baseuri, jrds.accountid, jobId, jrds.protocolversion)
	err := jrds.issueGetRequest(url, jobData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetUpdatableJobData(jobId string, jobData *JobUpdatableData) error {
	url := fmt.Sprintf("%v/automationAccounts/%v/jobs/%v?api-version=%v", jrds.baseuri, jrds.accountid, jobId, jrds.protocolversion)
	err := jrds.issueGetRequest(url, jobData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) GetRunbookData(runbookVersionId string, runbookData *RunbookData) error {
	url := fmt.Sprintf("%v/automationAccounts/%v/runbooks/%v?api-version=%v", jrds.baseuri, jrds.accountid, runbookVersionId, jrds.protocolversion)
	err := jrds.issueGetRequest(url, runbookData)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) AcknowledgeJobAction(sandboxId string, messageMetadata MessageMetadata) error {
	url := fmt.Sprintf("%v/automationAccounts/%v/Sandboxes/%v/jobs/AcknowledgeJobActions?api-version=%v", jrds.baseuri, jrds.accountid, sandboxId, jrds.protocolversion)
	err := jrds.issuePostRequest(url, messageMetadata, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) SetJobStatus(sandboxId string, jobId string, status int, isTermial bool, exception string) error {
	jobStatus := JobStatus{JobStatus: &status, Exception: &exception, IsFinalStatus: &isTermial}
	url := fmt.Sprintf("%v/automationAccounts/%v/Sandboxes/%v/jobs/%v/ChangeStatus?api-version=%v", jrds.baseuri, jrds.accountid, sandboxId, jobId, jrds.protocolversion)
	err := jrds.issuePostRequest(url, jobStatus, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error {
	recordTime := time.Now().Format(datetimeFormat)
	stream := Stream{AccountId: &jrds.accountid, JobId: &jobId, RecordTime: &recordTime, RunbookVersionId: &runbookVersionId, SequenceNumber: &sequence, StreamRecord: nil, StreamRecordText: &text, Type: &streamType} // Todo : datetime
	url := fmt.Sprintf("%v/automationAccounts/%v/jobs/%v/postJobStream?api-version=%v", jrds.baseuri, jrds.accountid, jobId, jrds.protocolversion)
	err := jrds.issuePostRequest(url, stream, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds *JrdsClient) SetLog(eventId string, activityId string, logType string, args ...string) error {
	log := Log{EventId: &eventId, Arguments: &args, LogType: &logType, ActivityId: &activityId}
	url := fmt.Sprintf("%v/automationAccounts/%v/logs?api-version=%v", jrds.baseuri, jrds.accountid, jrds.protocolversion)
	err := jrds.issuePostRequest(url, log, nil)
	if err != nil {
		return err
	}

	return nil
}

func (jrds JrdsClient) getDefaultHeaders() map[string]string {
	return map[string]string{accept_headerKey: appjson_headerValue,
		conection_headerKey: keepalive_headerValue,
		useragent_headerKey: fmt.Sprintf("AzureAutomationHybridWorker/%v", jrds.workerversion)}
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
		return NewRequestError(fmt.Sprintf("request error : %v \n", err))
	}

	if code == 401 {
		return NewRequestAuthorizationError(fmt.Sprintf("authorization error : %v\n", code))
	}

	if code != 200 {
		return NewRequestInvalidStatusError(fmt.Sprintf("invalid return code : %v\n", code))
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
		return NewRequestError(fmt.Sprintf("request error : %v \n", err))
	}

	if code == 401 {
		return NewRequestAuthorizationError(fmt.Sprintf("authorization error : %v\n", code))
	}

	if code != 200 {
		return NewRequestInvalidStatusError(fmt.Sprintf("invalid return code : %v\n", code))
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("failed to unmarshal request response: %v", err)
		}
	}

	return err
}
