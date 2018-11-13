package job

const (
	activating = 1
	running    = 2
	completed  = 3
	failed     = 4
	stopped    = 5
)

type status struct {
	enum       int
	isTerminal bool
	exception  *string
}

var getActivatingStatus = func() status {
	return status{enum: activating, isTerminal: false, exception: nil}
}

var getRunningStatus = func() status {
	return status{enum: running, isTerminal: false, exception: nil}
}

var getCompletedStatus = func() status {
	return status{enum: completed, isTerminal: true, exception: nil}
}

var getFailedStatus = func(exception string) status {
	return status{enum: failed, isTerminal: true, exception: &exception}
}

var getStoppedStatus = func() status {
	return status{enum: stopped, isTerminal: true, exception: nil}
}
