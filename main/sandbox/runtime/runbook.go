// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package runtime

import (
	"fmt"
)

const (
	PowerShell DefinitionKind = 5
	Python2    DefinitionKind = 9
	Python3    DefinitionKind = 10
	Bash       DefinitionKind = 11
)

type Runbook struct {
	Name       string
	Kind       DefinitionKind
	Definition string

	FileName string
}

type DefinitionKind int

var NewRunbook = func(Name string, versionId string, kind DefinitionKind, definition string) (Runbook, error) {
	language, err := GetLanguage(kind)
	if err != nil {
		return Runbook{}, err
	}

	runbook := Runbook{
		Name:       Name,
		Kind:       kind,
		Definition: definition,
		FileName:   fmt.Sprintf("%v-%v.%v", Name, versionId, language.extension)}

	return runbook, nil
}
