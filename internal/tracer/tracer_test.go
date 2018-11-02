package tracer

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type jrdsMock struct {
	setLog_f func(eventId int, activityId string, logType int, args ...string) error
}

func (j *jrdsMock) SetLog(eventId int, activityId string, logType int, args ...string) error {
	return j.setLog_f(eventId, activityId, logType, args...)
}

func TestLogWorkerStarting(t *testing.T) {
	mock := jrdsMock{}
	var (
		setLogCalled       = false
		traceLocallyCalled = false
	)
	mock.setLog_f = func(eventId int, activityId string, logType int, args ...string) error {
		setLogCalled = true
		return nil
	}

	traceLocally = func(trace trace) {
		traceLocallyCalled = true
	}

	InitializeTracer(&mock)
	LogWorkerStarting()

	// tracing are background routine
	time.Sleep(10 * time.Millisecond)

	if !setLogCalled || !traceLocallyCalled {
		t.Fatal("unexpected missing call to local or cloud trace")
	}
}

func Test_LocalTraceOnJrdsTraceError(t *testing.T) {
	mock := jrdsMock{}

	mock.setLog_f = func(eventId int, activityId string, logType int, args ...string) error {
		return fmt.Errorf(empty)
	}

	localTraceCount := 0
	traceLocally = func(trace trace) {
		localTraceCount += 1
	}

	InitializeTracer(&mock)
	LogWorkerStarting()

	// tracing are background routine
	time.Sleep(10 * time.Millisecond)

	if localTraceCount != 2 {
		t.Fatal("missing local trace on jrds trace exception")
	}
}

func Test_CloudTraceFormat(t *testing.T) {
	// we need to follow a strict contract for cloud trace format;
	// this test is simply to ensure that this format is enforced; do not change this format

	trace := NewTrace(0, empty, empty, empty)

	issueJrdsTrace = func(eventId int, activityId string, logType int, arg []string) error {
		if eventId != 16000 ||
			logType != 0 ||
			arg[0] != trace.accountId ||
			arg[1] != trace.subscriptionId ||
			arg[2] != trace.hybridWorkerGroupName ||
			arg[3] != trace.machineId ||
			arg[4] != trace.component ||
			arg[5] != strconv.Itoa(trace.eventId) ||
			arg[6] != trace.taskName ||
			arg[7] != trace.keyword ||
			arg[8] != strconv.Itoa(trace.threadId) ||
			arg[9] != strconv.Itoa(trace.processId) ||
			arg[10] != trace.activityId ||
			arg[11] != trace.hybridWorkerVersion ||
			arg[12] != trace.message {
			t.Fatal("unexpected cloud trace format")
		}

		return nil
	}

	formatAndIssueTrace(trace)

	// tracing are background routine
	time.Sleep(10 * time.Millisecond)
}

func LogWorkerTestMethod() string {
	return getTraceName()
}

func Test_GetTraceName(t *testing.T) {
	expectedTraceName := "WorkerTestMethod"
	traceName := LogWorkerTestMethod()

	if traceName != expectedTraceName {
		t.Fatal("unexpected trace name")
	}
}
