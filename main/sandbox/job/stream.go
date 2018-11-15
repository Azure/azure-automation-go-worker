// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package job

import (
	"fmt"
	"strings"
)

const (
	typeOutput   = "Output"
	typeProgress = "Progress"
	typeWarning  = "Warning"
	typeDebug    = "Debug"
	typeVerbose  = "Verbose"
	typeError    = "Error"
)

var (
	prefixDebug    = fmt.Sprintf("%v:", strings.ToLower(typeDebug))
	prefixVerbose  = fmt.Sprintf("%v:", strings.ToLower(typeVerbose))
	prefixWarning  = fmt.Sprintf("%v:", strings.ToLower(typeWarning))
	prefixError    = fmt.Sprintf("%v:", strings.ToLower(typeError))
	prefixProgress = fmt.Sprintf("%v:", strings.ToLower(typeProgress))
)

type StreamHandler struct {
	client           streamClient
	runbookVersionId string
	jobId            string
	sequence         int
}

type streamClient interface {
	SetJobStream(jobId string, runbookVersionId string, text string, streamType string, sequence int) error
}

func NewStreamHandler(client streamClient, jobId, runbookVersionId string) StreamHandler {
	return StreamHandler{
		client:           client,
		runbookVersionId: runbookVersionId,
		jobId:            jobId,
		sequence:         -1}
}

func (s *StreamHandler) SetStream(message string) {
	streamType := typeOutput

	if strings.HasPrefix(message, prefixDebug) ||
		strings.HasPrefix(message, strings.ToUpper(prefixDebug)) ||
		strings.HasPrefix(message, strings.Title(prefixDebug)) {
		streamType = typeDebug
	} else if strings.HasPrefix(message, prefixError) ||
		strings.HasPrefix(message, strings.ToUpper(prefixError)) ||
		strings.HasPrefix(message, strings.Title(prefixError)) {
		streamType = typeError
	} else if strings.HasPrefix(message, prefixVerbose) ||
		strings.HasPrefix(message, strings.ToUpper(prefixVerbose)) ||
		strings.HasPrefix(message, strings.Title(prefixVerbose)) {
		streamType = typeVerbose
	} else if strings.HasPrefix(message, prefixWarning) ||
		strings.HasPrefix(message, strings.ToUpper(prefixWarning)) ||
		strings.HasPrefix(message, strings.Title(prefixWarning)) {
		streamType = typeWarning
	} else if strings.HasPrefix(message, prefixProgress) ||
		strings.HasPrefix(message, strings.ToUpper(prefixProgress)) ||
		strings.HasPrefix(message, strings.Title(prefixProgress)) {
		streamType = typeProgress
	}

	s.sequence += 1
	err := s.client.SetJobStream(s.jobId, s.runbookVersionId, message, streamType, s.sequence)
	if err != nil {
		panic(err)
	}
}
