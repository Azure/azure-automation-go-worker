// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package runtime

import "fmt"

const (
	pythonExtension     Extension = "py"
	powershellExtension Extension = "ps1"
	bashExtension       Extension = "sh"
)

type Extension string
type Name string

type Language struct {
	extension   Extension
	interpreter Interpreter
}

func (l *Language) GetExtension() Extension {
	return l.extension
}

func (l *Language) GetInterpreter() Interpreter {
	return l.interpreter
}

var GetLanguage = func(definitionKind DefinitionKind) (Language, error) {
	var language Language
	switch definitionKind {
	case PowerShell:
		language = Language{extension: powershellExtension, interpreter: getPowerShellInterpreter()}
		break
	case Python2:
		language = Language{extension: pythonExtension, interpreter: getPython2Interpreter()}
		break
	case Python3:
		language = Language{extension: pythonExtension, interpreter: getPython3Interpreter()}
		break
	case Bash:
		language = Language{extension: bashExtension, interpreter: getBashInterpreter()}
		break
	default:
		return Language{}, fmt.Errorf("unsupported language")
	}
	return language, nil
}
