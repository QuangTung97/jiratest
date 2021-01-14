package sample

import (
	"jiratest"
	"testing"
)

func TestAdd(t *testing.T) {
	detail := jiratest.Detail{
		IssueLinks: []string{"PES-566"},
		Objective:  "add 2 numbers",
		Folder:     "Some folder",
	}

	table := []struct {
		name      string
		objective string
		a         int
		b         int
		result    int
	}{
		{
			name:      "test-1",
			objective: "objective 1",
			a:         4,
			b:         5,
			result:    9,
		},
		{
			name:      "test-2",
			objective: "objective 2",
			a:         4,
			b:         6,
			result:    9,
		},
		{
			name:      "test-3",
			objective: "objective 3",
			a:         2,
			b:         3,
			result:    5,
		},
	}

	for _, e := range table {
		t.Run(e.name, func(t *testing.T) {
			detail.Objective = e.objective
			defer jiratest.Setup(t, detail)()

			if e.a+e.b != e.result {
				t.Error(e.a, "+", e.b, "!=", e.result)
			}
		})
	}
}
