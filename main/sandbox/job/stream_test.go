package job

import (
	"fmt"
	"testing"
)

type clientMock struct {
	setStream_f func(jobId string, runbookVersionId string, text string, streamType string, sequence int) error
}

func (c *clientMock) SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error {
	return c.setStream_f(jobId, runbookVersionId, text, streamType, sequence)
}

func TestStreamHandler_SetStream_Debug(t *testing.T) {
	jrds := clientMock{}
	streamClient := NewStreamHandler(&jrds, "", "")

	streamType := ""
	jrds.setStream_f = func(jobId string, runbookVersionId string, text string, streamType string, sequence int) error {
		streamType = streamType
		return nil
	}

	format := "%v helloworld"

	prefix := "debug:"
	streamClient.SetStream(fmt.Sprintf(format, prefix))
	if streamType != typeDebug {
		t.Fatalf("unexpected stream type for prefix : %v", prefix)
	}

	prefix = "Debug:"
	streamClient.SetStream(fmt.Sprintf(format, prefix))
	if streamType != typeDebug {
		t.Fatalf("unexpected stream type for prefix : %v", prefix)
	}

	prefix = "DEBUG:"
	streamClient.SetStream(fmt.Sprintf(format, prefix))
	if streamType != typeDebug {
		t.Fatalf("unexpected stream type for prefix : %v", prefix)
	}

}
