package splitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func GenerateSplits(srcDir string, numberOfSplits int) ([]Split, error) {
	cmd := exec.Command("go", "test", "-json", srcDir)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running tests: %s\n", err.Error())
		os.Exit(1)
	}

	return GetSplits(string(output), numberOfSplits)
}

//func GenerateSplitsFromOutputs(numberOfSplits int, out []string) ([]Split, error) {
//
//}

func GetSplits(testJsonOut string, numberOfSplits int) ([]Split, error) {
	tests, err := parseTestOutput(testJsonOut)
	if err != nil {
		return nil, err
	}

	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Duration > tests[j].Duration
	})

	// Simple greedy algorithm to split the tests into groups
	// You can learn more about it by searching for Multifit algorithm (or identical machines scheduling)
	groups := make([]group, numberOfSplits+1)
	for _, test := range tests {
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Duration < groups[j].Duration
		})
		groups[0].Tests = append(groups[0].Tests, test)
		groups[0].Duration += test.Duration
	}

	splits := make([]Split, numberOfSplits+1)
	for i, group := range groups {
		splits[i] = make(Split, len(group.Tests))
		fmt.Printf("Group %d (Total Duration: %f seconds):\n", i+1, group.Duration)
		for j, test := range group.Tests {
			splits[i][j] = test.PackageName
			//fmt.Printf("- %s: %f seconds\n", test.PackageName, test.Duration)
		}
		fmt.Println()
	}

	return splits, nil
}

func parseTestOutput(output string) ([]test, error) {
	lines := strings.Split(output, "\n")
	var tests []test
	startTimes := map[string]time.Time{}
	for _, line := range lines {
		if len(line) <= 1 {
			continue
		}
		var event testEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			fmt.Printf("Error decoding JSON: %s\n", err)
			return nil, err
		}

		if event.Action == "fail" {
			fmt.Printf("Test `%s` failed. Stopping splitting as times for passing "+
				"tests would differ from failing tests\n", event.TestName)
			return nil, errors.New("test failed")
		}

		if event.Action == "start" {
			startTimes[event.Package] = event.Time
		}

		if event.Action == "pass" && event.TestName == "" {
			duration := event.Time.Sub(startTimes[event.Package]).Seconds()
			tests = append(tests, test{
				PackageName: event.Package,
				Duration:    duration,
			})
		}
	}
	return tests, nil
}

type testEvent struct {
	Time     time.Time `json:"Time"`
	Action   string    `json:"Action"`
	Package  string    `json:"Package"`
	TestName string    `json:"Test"`
	Elapsed  float64   `json:"Elapsed"`
}

type test struct {
	PackageName string
	Duration    float64
}

type group struct {
	Tests    []test
	Duration float64
}
