// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests_test

import (
	"errors"
	"testing"

	"github.com/GoogleCloudPlatform/marketplace-testrunner/specs"
	"github.com/GoogleCloudPlatform/marketplace-testrunner/tests"
	"github.com/ghodss/yaml"
)

type SimpleSetupExecutor struct{}

func (e *SimpleSetupExecutor) ExecScript(script string) (int, string, string, error) {
	return 0, "FOO", "BAR", nil
}

func TestSimpleSetup(t *testing.T) {
	shouldPass(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
			Stdout: &specs.StringAssert{
				Exactly: newString("FOO"),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString("BAR"),
			},
		},
	}, &SimpleSetupExecutor{})

	shouldPass(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
			Stdout: &specs.StringAssert{
				Exactly: newString("FOO"),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString("BAR"),
			},
		},
	}, &SimpleSetupExecutor{})
}

type FailingExecutor struct{}

func (e *FailingExecutor) ExecScript(script string) (int, string, string, error) {
	return 0, "FOO", "BAR", errors.New("Error while executing external process")
}

func TestFailingExecutor(t *testing.T) {
	shouldFail(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
			Stdout: &specs.StringAssert{
				Exactly: newString("FOO"),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString("BAR"),
			},
		},
	}, &FailingExecutor{})

	shouldFail(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
			Stdout: &specs.StringAssert{
				Exactly: newString("FOO"),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString("BAR"),
			},
		},
	}, &FailingExecutor{})
}

type ScriptsRequired struct{}

func (e *ScriptsRequired) ExecScript(script string) (int, string, string, error) {
	return 0, "FOO", "BAR", nil
}

func TestScriptsRequired(t *testing.T) {
	shouldPass(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
		},
	}, &ScriptsRequired{})

	shouldPass(t, &specs.BashTest{
		Scripts: []string{
			"# 1st script is not executed",
			"# 2nd script is not executed",
			"# 3rd script is not executed",
		},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
		},
	}, &ScriptsRequired{})

	// Script(s) field is required
	shouldFail(t, &specs.BashTest{
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
		},
	}, &ScriptsRequired{})

	// Cannot use both script and scripts fields at the same time.
	shouldFail(t, &specs.BashTest{
		Script: "# this script is not executed",
		Scripts: []string{
			"# 1st script is not executed",
			"# 2nd script is not executed",
			"# 3rd script is not executed",
		},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
		},
	}, &ScriptsRequired{})
}

type ExitCodeExecutor struct{}

func (e *ExitCodeExecutor) ExecScript(script string) (int, string, string, error) {
	return 42, "", "", nil
}

func TestExitCodeParsing(t *testing.T) {
	shouldPass(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(42),
			},
		},
	}, &ExitCodeExecutor{})

	shouldPass(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(42),
			},
		},
	}, &ExitCodeExecutor{})

	shouldPass(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals:      newInt(42),
				NotEquals:   newInt(41),
				GreaterThan: newInt(41),
			},
			Stdout: &specs.StringAssert{
				Exactly: newString(""),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString(""),
			},
		},
	}, &ExitCodeExecutor{})

	shouldPass(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals:      newInt(42),
				NotEquals:   newInt(41),
				GreaterThan: newInt(41),
			},
			Stdout: &specs.StringAssert{
				Exactly: newString(""),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString(""),
			},
		},
	}, &ExitCodeExecutor{})

	shouldFail(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
		},
	}, &ExitCodeExecutor{})

	shouldFail(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
		},
	}, &ExitCodeExecutor{})

	shouldFail(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				LessThan: newInt(40),
			},
		},
	}, &ExitCodeExecutor{})

	shouldFail(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				LessThan: newInt(40),
			},
		},
	}, &ExitCodeExecutor{})
}

type StdoutExecutor struct{}

func (e *StdoutExecutor) ExecScript(script string) (int, string, string, error) {
	return 0, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nProin nibh augue, suscipit a, scelerisque sed, lacinia in, mi.", "", nil
}

func TestStdoutParsing(t *testing.T) {
	shouldPass(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			Stdout: &specs.StringAssert{
				Exactly: newString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nProin nibh augue, suscipit a, scelerisque sed, lacinia in, mi."),
			},
		},
	}, &StdoutExecutor{})

	shouldPass(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			Stdout: &specs.StringAssert{
				Exactly: newString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nProin nibh augue, suscipit a, scelerisque sed, lacinia in, mi."),
			},
		},
	}, &StdoutExecutor{})

	shouldPass(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
			Stdout: &specs.StringAssert{
				Exactly:  newString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nProin nibh augue, suscipit a, scelerisque sed, lacinia in, mi."),
				Contains: newString("consectetur a"),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString(""),
			},
		},
	}, &StdoutExecutor{})

	shouldPass(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			ExitCode: &specs.IntAssert{
				Equals: newInt(0),
			},
			Stdout: &specs.StringAssert{
				Exactly:  newString("Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nProin nibh augue, suscipit a, scelerisque sed, lacinia in, mi."),
				Contains: newString("consectetur a"),
			},
			Stderr: &specs.StringAssert{
				Exactly: newString(""),
			},
		},
	}, &StdoutExecutor{})

	shouldFail(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			Stdout: &specs.StringAssert{
				Exactly:  newString("Proin nibh augue, suscipit a, scelerisque sed, lacinia in, mi."),
				Contains: newString("consecteturadipiscingadipiscing"),
			},
		},
	}, &StdoutExecutor{})

	shouldFail(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			Stdout: &specs.StringAssert{
				Exactly:  newString("Proin nibh augue, suscipit a, scelerisque sed, lacinia in, mi."),
				Contains: newString("consecteturadipiscingadipiscing"),
			},
		},
	}, &StdoutExecutor{})

	shouldFail(t, &specs.BashTest{
		Script: "# this script is not executed",
		Expect: &specs.CliExpect{
			Stdout: &specs.StringAssert{
				NotContains: newString("Lorem"),
			},
		},
	}, &StdoutExecutor{})

	shouldFail(t, &specs.BashTest{
		Scripts: []string{"# 1st script is not executed"},
		Expect: &specs.CliExpect{
			Stdout: &specs.StringAssert{
				NotContains: newString("Lorem"),
			},
		},
	}, &StdoutExecutor{})
}

func shouldPass(t *testing.T, test *specs.BashTest, executor tests.CommandExecutor) {
	outcome := runBashTest(t, test, executor)
	if len(outcome) > 0 {
		t.Errorf("Expected to pass but was failing: %s", outcome)
	} else {
		t.Logf("Passing OK!")
	}
}

func shouldFail(t *testing.T, test *specs.BashTest, executor tests.CommandExecutor) {
	outcome := runBashTest(t, test, executor)
	if len(outcome) == 0 {
		t.Errorf("Expected to fail but was passing")
	} else {
		t.Logf("Failing OK! Test output: %s", outcome)
	}
}

func runBashTest(t *testing.T, test *specs.BashTest, executor tests.CommandExecutor) string {
	testRule, _ := yaml.Marshal(test)
	t.Logf("---\nTest rule:\n%s", testRule)
	return tests.RunBashTest(test, executor)
}
