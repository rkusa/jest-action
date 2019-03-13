package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v24/github"
	"github.com/ldez/ghactions"
)

type Report struct {
	NumFailedTests      int
	NumPassedTests      int
	NumTotalTests       int
	NumFailedTestSuites int
	NumPassedTestSuites int
	NumTotalTestSuites  int
	Success             bool
	TestResults         []*TestResult
}

type TestResult struct {
	AssertionResults []*AssertionResult
	Message          string
	FilePath         string `json:"name"`
	Status           string
	Summary          string
}

type AssertionResult struct {
	AncestorTitles  []string
	FailureMessages []string
	FullName        string
	Location        Location
	Status          string
	Title           string
}

type Location struct {
	Column int
	Line   int
}

func main() {
	var report Report
	err := json.NewDecoder(os.Stdin).Decode(&report)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	action := ghactions.NewAction(ctx)

	action.OnPush(func(client *github.Client, event *github.PushEvent) error {
		return handlePush(ctx, client, event, report)
	})

	if err := action.Run(); err != nil {
		log.Fatal(err)
	}
}

func handlePush(ctx context.Context, client *github.Client, event *github.PushEvent, report Report) error {
	if report.Success {
		return nil
	}

	head := os.Getenv(ghactions.GithubSha)
	owner, repoName := ghactions.GetRepoInfo()

	// find the action's checkrun
	checkName := os.Getenv(ghactions.GithubAction)
	result, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repoName, head, &github.ListCheckRunsOptions{
		CheckName: github.String(checkName),
		Status:    github.String("in_progress"),
	})
	if err != nil {
		return err
	}

	if len(result.CheckRuns) == 0 {
		return fmt.Errorf("Unable to find check run for action: %s", checkName)
	}
	checkRun := result.CheckRuns[0]

	// add annotations for test failures
	workspacePath := os.Getenv(ghactions.GithubWorkspace) + "/"
	var annotations []*github.CheckRunAnnotation
	for _, t := range report.TestResults {
		if t.Status == "passed" {
			continue
		}

		path := strings.TrimPrefix(t.FilePath, workspacePath)

		if len(t.AssertionResults) > 0 {
			for _, a := range t.AssertionResults {
				if a.Status == "passed" {
					continue
				}

				annotations = append(annotations, &github.CheckRunAnnotation{
					Path:            github.String(path),
					StartLine:       github.Int(a.Location.Line),
					EndLine:         github.Int(a.Location.Line),
					AnnotationLevel: github.String("failure"),
					Title:           github.String(a.FullName),
					Message:         github.String(strings.Join(a.FailureMessages, "\n\n")),
				})
			}
		} else {
			// usually the case for failed test suites
			annotations = append(annotations, &github.CheckRunAnnotation{
				Path:            github.String(path),
				StartLine:       github.Int(1),
				EndLine:         github.Int(1),
				AnnotationLevel: github.String("failure"),
				Title:           github.String("Test Suite Error"),
				Message:         github.String(t.Message),
			})
		}
	}

	summary := fmt.Sprintf(
		"Test Suites: %d failed, %d passed, %d total\n",
		report.NumFailedTests,
		report.NumPassedTests,
		report.NumTotalTests,
	)
	summary += fmt.Sprintf(
		"Tests: %d failed, %d passed, %d total",
		report.NumFailedTestSuites,
		report.NumPassedTestSuites,
		report.NumTotalTestSuites,
	)

	// add annotations in #50 chunks
	for i := 0; i < len(annotations); i += 50 {
		end := i + 50

		if end > len(annotations) {
			end = len(annotations)
		}

		_, _, err = client.Checks.UpdateCheckRun(ctx, owner, repoName, checkRun.GetID(), github.UpdateCheckRunOptions{
			Name:    checkName,
			HeadSHA: github.String(head),
			Output: &github.CheckRunOutput{
				Title:       github.String("Result"),
				Summary:     github.String(summary),
				Annotations: annotations[i:end],
			},
		})
		if err != nil {
			return err
		}
	}

	return fmt.Errorf(summary)
}
