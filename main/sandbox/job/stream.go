package job

import "github.com/Azure/azure-automation-go-worker/internal/tracer"

type streamHandler struct {
	client           streamClient
	runbookVersionId string
	jobId            string
	sequence         int
}

type streamClient interface {
	SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error
}

func NewStreamHandler(client streamClient, jobId, runbookVersionId string) streamHandler {
	return streamHandler{
		client:           client,
		runbookVersionId: runbookVersionId,
		jobId:            jobId,
		sequence:         0}
}

func (s *streamHandler) SetStream(message string) {
	err := s.client.SetJobStream(s.jobId, s.runbookVersionId, message, "output", s.sequence)
	if err != nil {
		panic(err)
	}
	tracer.LogDebugTrace(message)
	s.sequence += 1
}
