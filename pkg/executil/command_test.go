package executil

import (
	"reflect"
	"testing"
)

const (
	cmd  = "cmdA"
	arg1 = "arg1"
	arg2 = "arg2"
)

func TestNewCommand(t *testing.T) {
	c := NewCommand(cmd, arg1, arg2)

	if c.Name != cmd {
		t.Errorf("invalid command name")
	}
	if reflect.DeepEqual(c.Arguments, []string{arg1, arg2}) == false {
		t.Errorf("invalid command arguments")
	}
	if c.IsSuccessful == true {
		t.Errorf("invalid command IsSuccessful property")
	}
}
