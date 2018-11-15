package job

const (
	None      = 0
	Activate  = 1
	Abort     = 2
	Suspend   = 3
	Resume    = 4
	Stop      = 5
	Terminate = 6
	Remove    = 7
)

type PendingAction struct {
	Enum int
	Name string
}

func GetPendingAction(enum int) PendingAction {
	switch enum {
	case 1:
		return PendingAction{Enum: Activate, Name: "Activate"}
	case 2:
		return PendingAction{Enum: Abort, Name: "Abort"}
	case 3:
		return PendingAction{Enum: Suspend, Name: "Suspend"}
	case 4:
		return PendingAction{Enum: Resume, Name: "Resume"}
	case 5:
		return PendingAction{Enum: Stop, Name: "Stop"}
	case 6:
		return PendingAction{Enum: Terminate, Name: "Terminate"}
	case 7:
		return PendingAction{Enum: Remove, Name: "Remove"}
	}

	return PendingAction{Enum: None, Name: "None"}
}
