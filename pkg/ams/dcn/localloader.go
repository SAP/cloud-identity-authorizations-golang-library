package dcn

import (
	"encoding/json"
	"os"
	"path"
	"strings"
)

type Loader struct {
	dir                string
	DCNChannel         chan DcnContainer
	AssignmentsChannel chan Assignments
	errorHandler       []func(error)
}

func NewLocalLoader(dir string, errorCallback func(error)) *Loader {
	loader := &Loader{
		dir:                dir,
		DCNChannel:         make(chan DcnContainer),
		AssignmentsChannel: make(chan Assignments),
		errorHandler:       []func(error){},
	}
	if errorCallback != nil {
		loader.errorHandler = append(loader.errorHandler, errorCallback)
	}

	go loader.start()
	return loader
}

func (l *Loader) start() {
	dcn, assignments, err := readDirectory(l.dir)
	if err != nil {
		l.handleError(err)
		return
	}
	l.DCNChannel <- dcn
	l.AssignmentsChannel <- assignments.Assignments
}

func (l *Loader) handleError(err error) {
	for _, handler := range l.errorHandler {
		handler(err)
	}
}

func readDirectory(dir string) (DcnContainer, AssignmentsContainer, error) {
	// Read all files in the directory
	// For each file, read the content and parse it
	// Send the parsed content to the channel
	resultDcn := DcnContainer{
		Policies:  []Policy{},
		Functions: []Function{},
		Schemas:   []Schema{},
		Tests:     []Test{},
	}
	resultAssigments := AssignmentsContainer{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return resultDcn, resultAssigments, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subDCN, _, err := readDirectory(path.Join(dir, entry.Name()))
			if err != nil {
				return resultDcn, resultAssigments, err
			}

			resultDcn.Policies = append(resultDcn.Policies, subDCN.Policies...)
			resultDcn.Functions = append(resultDcn.Functions, subDCN.Functions...)
			resultDcn.Schemas = append(resultDcn.Schemas, subDCN.Schemas...)
			resultDcn.Tests = append(resultDcn.Tests, subDCN.Tests...)
		}
		if strings.HasSuffix(entry.Name(), ".dcn") {
			var dcn DcnContainer
			raw, err := os.ReadFile(path.Join(dir, entry.Name()))
			if err != nil {
				return resultDcn, resultAssigments, err
			}
			err = json.Unmarshal(raw, &dcn)
			if err != nil {
				return resultDcn, resultAssigments, err
			}
			resultDcn.Policies = append(resultDcn.Policies, dcn.Policies...)
			resultDcn.Functions = append(resultDcn.Functions, dcn.Functions...)
			resultDcn.Schemas = append(resultDcn.Schemas, dcn.Schemas...)
			resultDcn.Tests = append(resultDcn.Tests, dcn.Tests...)
		}
		if strings.HasSuffix(entry.Name(), "data.json") {
			raw, err := os.ReadFile(path.Join(dir, entry.Name()))
			if err != nil {
				return resultDcn, resultAssigments, err
			}
			err = json.Unmarshal(raw, &resultAssigments)
			if err != nil {
				return resultDcn, resultAssigments, err
			}
		}
	}
	return resultDcn, resultAssigments, nil
}
