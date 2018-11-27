package tracer

import (
	"fmt"
	"github.com/Azure/azure-automation-go-worker/internal/configuration"
	"math"
	"os"
	"path"
	"sync"
	"time"
)

var diskMutex = &sync.Mutex{}

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
	writeToDisk(traceOutput)
}

var writeToDisk = func(msg string) {
	diskMutex.Lock()
	defer diskMutex.Unlock()

	// open file
	logPath := path.Join(configuration.GetWorkingDirectory(), localLogFilename)
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		fmt.Printf("Unable to open worker.log, error : %v\n", err.Error())
		return
	}

	// get file size
	info, err := file.Stat()
	if err != nil {
		fmt.Printf("Unable to get worker.log file informations, error : %v\n", err.Error())
		return
	}

	// rotate if needed; only keep 2 iteration of log file
	if info.Size() > int64(math.Pow(1, 7)) {
		// rotate
		file.Close()
		os.Rename(logPath, fmt.Sprintf("%v.1", logPath))

		file, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			fmt.Printf("Unable to open worker.log after rotation, error : %v\n", err.Error())
			return
		}
	}

	file.WriteString(msg)
	file.Close()
}
