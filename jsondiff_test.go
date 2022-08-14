package jsondiff

import (
	"fmt"
	"testing"
)

var diffCases = []struct {
	a      string
	b      string
	match  bool
	result []JsonDiffItem
}{
	{`{"a": 1}`, `{"a": 1}`, true, []JsonDiffItem{}},
	{`{"a": 1}`, `{"a": 2}`, false, []JsonDiffItem{}},
	{`{"a": 1}`, `{"a": 2, "b": 3}`, false, []JsonDiffItem{}},
	{`{"a": 1}`, `{"a": true}`, false, []JsonDiffItem{}},
	{`{"a": [1,2,3]}`, `{"a": [1,2,3]}`, true, []JsonDiffItem{}},
	{`{"a": [1,2,3]}`, `{"a": [2,3,4]}`, false, []JsonDiffItem{}},
	{`{"a": {"c":1}}`, `{"a": {"c":1}}`, true, []JsonDiffItem{}},
	{`{"a": {"c":1}}`, `{"a": {"c":2}}`, false, []JsonDiffItem{}},
	{`{"a": {"c":[1,2,3]}}`, `{"a": {"c":[1,2,3]}}`, true, []JsonDiffItem{}},
	{`{"a": {"c":[1,2,3]}}`, `{"a": {"c":[1,2,4]}}`, false, []JsonDiffItem{}},
	{`{"a": {"c":[1,2,3]}, "d": 1}`, `{"a": {"c":[1,2,4]}}`, false, []JsonDiffItem{}},
	{`[1,2,3]`, `[1,2,3]`, true, []JsonDiffItem{}},
	{`[1,2,3]`, `[2,3,4]`, false, []JsonDiffItem{}},
	{`[1,2,3]`, `[1,2,3,4]`, false, []JsonDiffItem{}},
	{`[1,{"a": 1}]`, `[1,{"a": 1}]`, true, []JsonDiffItem{}},
	{`[1,{"a": 1}]`, `[1,{"a": 2}]`, false, []JsonDiffItem{}},
}

func TestJsonDiff(t *testing.T) {
	for caseIndex, caseOne := range diffCases {
		changes, _ := JsonDiff(caseOne.a, caseOne.b)
		isMatched := true
		if len(changes) > 0 {
			isMatched = false
		}

		for changeIndex, changeOne := range changes {
			fmt.Printf("Case[%d] change[%d]:(%s)%s\n", caseIndex, changeIndex+1, changeOne.Type, changeOne.Path)
		}

		if isMatched != caseOne.match {
			t.Errorf("Case[%d] failed, got: %v, expected: %v\n", caseIndex, isMatched, caseOne.match)
		}
	}
}
