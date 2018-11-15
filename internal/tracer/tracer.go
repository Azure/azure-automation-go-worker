// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package tracer

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"github.com/Azure/azure-automation-go-worker/internal/jrds"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	empty = ""

	logPrefix        = "Log"
	debugTracePrefix = "[DebugTrace]"

	traceDatetimeFormat = "2006-01-02T15:04:05.00"

	cloudDebugLogType       = 0
	cloudHybridTraceEventId = 16000

	keywordError         = "Error"
	keywordDebug         = "Debug"
	keywordStartup       = "Startup"
	keywordRoutine       = "Routine"
	keywordInformational = "Informational"
	keywordJob           = "Job"

	tasknameTraceError     = "TraceError"
	trasknameSandboxStdout = "SandboxStdout"
	trasknameSandboxStderr = "SandboxStderr"
)

var (
	jrdsClient jrdsTracer

	activityId        = generateActivityId()
	tracerPackageName = reflect.TypeOf(tracer{}).PkgPath()
)

type tracer struct {
}

type trace struct {
	component string

	threadId  int
	processId int

	eventId    int
	taskName   string
	message    string
	keyword    string
	activityId string

	accountId             string
	subscriptionId        string
	machineId             string
	hybridWorkerGroupName string
	hybridWorkerVersion   string
}

type jrdsTracer interface {
	SetLog(eventId int, activityId string, logType int, args ...string) error
}

func InitializeTracer(client jrdsTracer) {
	jrdsClient = client
}

var traceGenericHybridWorkerDebugEvent = func(eventId int, taskName string, message string, keyword string) {
	trace := NewTrace(eventId, taskName, message, keyword)
	go traceGenericHybridWorkerEventRoutine(trace, true, false)
}

var traceGenericHybridWorkerEvent = func(eventId int, taskName string, message string, keyword string) {
	trace := NewTrace(eventId, taskName, message, keyword)
	go traceGenericHybridWorkerEventRoutine(trace, false, false)
}

var traceGenericHybridWorkerEventRoutine = func(trace trace, debug bool, localonly bool) {
	// do not log debug traces based on configuration
	if !configuration.GetDebugTraces() && debug {
		return
	}

	// local stdout
	traceLocally(trace)

	if localonly {
		return
	}

	// cloud stdout
	err := formatAndIssueTrace(trace)
	if err != nil {
		traceErrorLocally(fmt.Sprintf("error while calling formatAndIssueTrace : %v \n", err))
	}
}

var NewTrace = func(eventId int, taskName string, message string, keyword string) trace {
	return trace{
		component:             configuration.GetComponent(),
		threadId:              1, // thread id are not exposed in go but this field is require for cloud traces
		processId:             os.Getpid(),
		eventId:               eventId,
		taskName:              taskName,
		message:               message,
		keyword:               keyword,
		activityId:            activityId,
		accountId:             configuration.GetAccountId(),
		subscriptionId:        empty,
		machineId:             empty,
		hybridWorkerGroupName: configuration.GetHybridWorkerGroupName(),
		hybridWorkerVersion:   configuration.GetWorkerVersion()}
}

var traceErrorLocally = func(message string) {
	errorTrace := NewTrace(0, tasknameTraceError, message, keywordError)
	traceLocally(errorTrace)
}

var traceLocally = func(trace trace) {
	traceOutput := ""

	if configuration.GetComponent() == configuration.Component_worker &&
		trace.component == configuration.Component_sandbox {
		const format = "%v \n"
		traceOutput = fmt.Sprintf(format, trace.message)
	} else {
		const format = "%s (%v)[%s] : [%s] %v \n"
		var now = time.Now().Format(traceDatetimeFormat)
		traceOutput = fmt.Sprintf(format, now, trace.processId, trace.component, trace.taskName, trace.message)
	}

	fmt.Print(traceOutput)
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

// generateActivityId generates a somewhat unique uuid; this isn't a proper uuid implementation and should only
// be used temporally until we have time to implement a proper uuid generation algorithm
var generateActivityId = func() string {
	b := make([]byte, 16)

	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

var getTraceName = func() string {
	pc := make([]uintptr, 10)

	// skip 2 frames (Callers() and getTraceName()); assumes this is called from the tracing function directly only
	callers := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:callers])
	frame, _ := frames.Next()

	replacer := strings.NewReplacer(logPrefix, empty, fmt.Sprintf("%s.", tracerPackageName), empty)
	return replacer.Replace(frame.Function)
}

func LogSandboxStdout(message string) {
	trace := NewTrace(0, trasknameSandboxStdout, message, keywordInformational)
	trace.component = configuration.Component_sandbox
	go traceGenericHybridWorkerEventRoutine(trace, strings.Contains(message, debugTracePrefix), true)
}

func LogSandboxStderr(message string) {
	trace := NewTrace(0, trasknameSandboxStderr, message, keywordInformational)
	trace.component = configuration.Component_sandbox
	go traceGenericHybridWorkerEventRoutine(trace, strings.Contains(message, debugTracePrefix), true)
}

func LogWorkerTraceError(message string) {
	traceGenericHybridWorkerEvent(20000, getTraceName(), message, keywordStartup)
}

func LogDebugTrace(message string) {
	traceGenericHybridWorkerDebugEvent(20001, getTraceName(), message, keywordDebug)
}

func LogErrorTrace(error string) {
	message := fmt.Sprintf("Error : %v", error)
	traceGenericHybridWorkerDebugEvent(20001, getTraceName(), message, keywordDebug)
}

func LogWorkerStarting() {
	message := "Worker starting."
	traceGenericHybridWorkerEvent(20020, getTraceName(), message, keywordStartup)
}

func LogWorkerSandboxActionsFound(actions jrds.SandboxActions) {
	message := fmt.Sprintf("Get sandbox actions found %v new action(s).", len(actions.Value))
	traceGenericHybridWorkerEvent(20100, getTraceName(), message, keywordRoutine)
}

func LogWorkerErrorGettingSandboxActions(err error) {
	message := fmt.Sprintf("Error getting sandbox actions. [error=%v]", err.Error())
	traceGenericHybridWorkerEvent(20101, getTraceName(), message, keywordRoutine)
}

func LogWorkerFailedToCreateSandbox(err error) {
	message := fmt.Sprintf("Error creating sandbox. [error=%v]", err.Error())
	traceGenericHybridWorkerEvent(20102, getTraceName(), message, keywordRoutine)
}

func LogWorkerSandboxProcessExited(sandboxId string, pid, exitCode int) {
	message := fmt.Sprintf("Sandbox process exited. [sandboxId=%v][pId=%v][exitCode=%v]", sandboxId, pid, exitCode)
	traceGenericHybridWorkerEvent(20102, getTraceName(), message, keywordRoutine)
}

func LogSandboxStarting(id string) {
	message := fmt.Sprintf("Sandbox starting [sandboxId=%v]", id)
	traceGenericHybridWorkerEvent(25000, getTraceName(), message, keywordStartup)
}

func LogSandboxGetJobActions(actions *jrds.JobActions) {
	message := fmt.Sprintf("Get job actions. Found %v new action(s).", len(actions.Value))
	traceGenericHybridWorkerEvent(25001, getTraceName(), message, keywordRoutine)
}

func LogSandboxJobLoaded(sandboxId, jobId string) {
	message := fmt.Sprintf("Job loaded. [sandboxId=%v][jobId=%v]", sandboxId, jobId)
	traceGenericHybridWorkerEvent(25010, getTraceName(), message, keywordJob)
}

func LogSandboxJobUnloaded(sandboxId, jobId string) {
	message := fmt.Sprintf("Job unloaded. [sandboxId=%v][jobId=%v]", sandboxId, jobId)
	traceGenericHybridWorkerEvent(25013, getTraceName(), message, keywordJob)
}

func LogSandboxJobUnsupportedRunbookType(sandboxId, jobId, exception string) {
	message := fmt.Sprintf("Unsupported runbook type. [sandboxId=%v][jobId=%v][exception=%v]", sandboxId, jobId, exception)
	traceGenericHybridWorkerEvent(25014, getTraceName(), message, keywordJob)
}
