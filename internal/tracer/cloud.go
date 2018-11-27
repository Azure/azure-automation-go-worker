package tracer

import (
	"fmt"
	"math/rand"
	"strconv"
)

// generateActivityId generates a somewhat unique uuid; this isn't a proper uuid implementation and should only
// be used temporally until we have time to implement a proper uuid generation algorithm
var generateActivityId = func() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

var formatAndIssueTrace = func(trace trace) error {
	// this format matches the cloud etw manifest; do not reorder
	cloudTraceFormat := []string{
		trace.accountId,
		trace.subscriptionId,
		trace.hybridWorkerGroupName,
		trace.machineId,
		trace.component,
		strconv.Itoa(trace.eventId),
		trace.taskName,
		trace.keyword,
		strconv.Itoa(trace.threadId),
		strconv.Itoa(trace.processId),
		trace.activityId,
		trace.hybridWorkerVersion,
		trace.message,
	}

	err := issueJrdsTrace(cloudHybridTraceEventId, activityId, cloudDebugLogType, cloudTraceFormat)
	if err != nil {
		return err
	}

	return nil
}

var issueJrdsTrace = func(eventId int, activityId string, logType int, arg []string) error {
	if jrdsClient == nil {
		return fmt.Errorf("error emitting trace; nil jrds client in tracer package \n")
	}

	err := jrdsClient.SetLog(eventId, activityId, logType, arg...)
	if err != nil {
		return fmt.Errorf("error emitting trace to jrds : %v \n", err)
	}

	return nil
}
