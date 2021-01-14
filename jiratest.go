package jiratest

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

// TestRunScript ...
type TestRunScript struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Detail ...
type Detail struct {
	Name               string
	IssueLinks         []string
	Objective          string
	Precondition       string
	WebLinks           []string
	Folder             string
	ConfluenceLinks    []string
	TestScript         TestRunScript
	TestRunEnvironment string
}

type testResult struct {
	Name               string        `json:"name"`
	IssueLinks         []string      `json:"issueLinks"`
	Objective          string        `json:"objective"`
	Precondition       string        `json:"precondition"`
	WebLinks           []string      `json:"webLinks"`
	Folder             string        `json:"folder"`
	ConfluenceLinks    []string      `json:"confluenceLinks"`
	TestScript         TestRunScript `json:"testScript"`
	TestRunStatus      string        `json:"testrun_status"`
	TestRunEnvironment string        `json:"testrun_environment"`
	TestRunComment     string        `json:"testrun_comment"`
	TestRunDuration    float64       `json:"testrun_duration"`
	TestRunDate        string        `json:"testrun_date"`
}

var results []testResult
var resultMut sync.Mutex

var outputFile *os.File
var outputFileMut sync.Mutex

var directory = flag.String("jira_pwd", "", "go test -v ./... -jira_pwd=$PWD")

func getFile() *os.File {
	outputFileMut.Lock()
	defer outputFileMut.Unlock()
	if outputFile != nil {
		return outputFile
	}

	file, err := os.Create(path.Join(*directory, "testrun.json"))
	if err != nil {
		panic(err)
	}
	outputFile = file

	return outputFile
}

func writeResultsToFile(writer io.Writer) {
	resultMut.Lock()
	defer resultMut.Unlock()

	err := json.NewEncoder(writer).Encode(results)
	if err != nil {
		panic(err)
	}
}

func writeResults() {
	if *directory != "" {
		file := getFile()
		_, err := file.Seek(0, os.SEEK_SET)
		if err != nil {
			panic(err)
		}

		writeResultsToFile(file)
	}
}

// Setup set up and tear down a Functional Test Case
func Setup(t *testing.T, detail Detail) func() {
	start := time.Now()

	return func() {
		name := detail.Name
		if name == "" {
			name = t.Name()
		}

		d := time.Since(start)

		status := "Pass"
		if t.Failed() {
			status = "Fail"
		}

		result := testResult{
			Name:               name,
			IssueLinks:         detail.IssueLinks,
			Objective:          detail.Objective,
			Precondition:       detail.Precondition,
			WebLinks:           detail.WebLinks,
			Folder:             detail.Folder,
			ConfluenceLinks:    detail.ConfluenceLinks,
			TestScript:         detail.TestScript,
			TestRunStatus:      status,
			TestRunEnvironment: detail.TestRunEnvironment,
			TestRunDuration:    float64(d.Milliseconds()) / 1000.0,
			TestRunDate:        start.Format(time.RFC3339),
		}

		resultMut.Lock()
		results = append(results, result)
		resultMut.Unlock()

		writeResults()
	}
}
