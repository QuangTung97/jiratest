.PHONY: test jira-test

test:
	go test -v ./...

jira-test:
	go test ./... -jira_pwd=${PWD}; cat testrun.tmp.json | jq -s "." > testrun.json
