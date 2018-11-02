package jrds

type SandboxActions struct {
	Value []SandboxAction `json:"value"`
}

type SandboxAction struct {
	SandboxId *string `json:"sandboxId"`
}

type JobActions struct {
	Value []JobAction `json:"value"`
}

type JobAction struct {
	MessageMetadata *MessageMetadata `json:"MessageMetadata"`
	MessageSource   *string          `json:"MessageSource"`
	LockToken       *string          `json:"LockToken"`
	JobId           *string          `json:"JobId"`
}

type MessageMetadatas struct {
	MessageMetadatas []MessageMetadata `json:"MessageMetadatas"`
}

type JobData struct {
	RunbookVersionId *string `json:"runbookVersionId"`
}

type JobUpdatableData struct {
	LogActivityTrace       *int    `json:"logActivityTrace"`
	TriggerSource          *int    `json:"triggerSource"`
	JobId                  *string `json:"jobId"`
	LogProgress            *bool   `json:"logProgress"`
	JobStatus              *int    `json:"jobStatus"`
	AccountName            *string `json:"accountName"`
	PartitionId            *int    `json:"partitionId"`
	LogDebug               *bool   `json:"logDebug"`
	IsDraft                *bool   `json:"isDraft"`
	WorkflowInstanceId     *string `json:"workflowInstanceId"`
	JobKey                 *int    `json:"JobKey"`
	SubscriptionId         *string `json:"subscriptionId"`
	NoPersistEvictionCount *int    `json:"noPersistEvictionsCount"`
	AccountKey             *int    `json:"AccountKey"`
	JobStartedBy           *string `json:"jobStartedBy"`
	JobDestination         *string `json:"jobRunDestination"`
	TierName               *string `json:"tierName"`
	PendingActionData      *string `json:"pendingActionData"`
	JobStatusDetails       *int    `json:"jobStatusDetails"`
	ResourceGroupName      *string `json:"resourceGroupName"`
	PendingAction          *int    `json:"pendingAction"`
	LogVerbose             *bool   `json:"logVerbose"`
	RunbookKey             *int    `json:"RunbookKey"`
}

type RunbookData struct {
	Name                  *string `json:"name"`
	AccountId             *string `json:"accountId"`
	RunbookId             *string `json:"runbookId"`
	Definition            *string `json:"definition"`
	RunbookDefinitionKind *int    `json:"runbookDefinitionKind"`
	RunbookVersionId      *int    `json:"runbookVersionId"`
	Parameters            *bool   `json:"parameters"`
}

type MessageMetadata struct {
	PopReceipt *string `json:"PopReceipt"`
	MessageId  *string `json:"MessageId"`
}

type AcknowledgeJobAction struct {
	MessageMetadata MessageMetadata `json:"MessageMetadata"`
}

type JobStatus struct {
	Exception     *string `json:"exception"`
	IsFinalStatus *bool   `json:"isFinalStatus"`
	JobStatus     *int    `json:"jobStatus"`
}

type Stream struct {
	AccountId        *string `json:"AccountId"`
	JobId            *string `json:"JobId"`
	RecordTime       *string `json:"RecordTime"`
	RunbookVersionId *string `json:"RunbookVersionId"`
	SequenceNumber   *int    `json:"SequenceNumber"`
	StreamRecord     *string `json:"StreamRecord"`
	StreamRecordText *string `json:"StreamRecordText"`
	Type             *string `json:"Type"`
}

type Log struct {
	ActivityId *string   `json:"activityId"`
	Arguments  *[]string `json:"args"`
	EventId    *int      `json:"eventId"`
	LogType    *int      `json:"logtype"`
}

type UnloadJob struct {
	IsTest                 *bool   `json:"isTest"`
	JobId                  *string `json:"jobId"`
	StartTime              *string `json:"startTime"`
	SubscriptionId         *string `json:"subscriptionId"`
	ExecutionTimeInSeconds *int    `json:"executionTimeInSeconds"`
}
