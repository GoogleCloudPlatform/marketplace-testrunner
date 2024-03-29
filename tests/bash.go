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

package tests // import "github.com/GoogleCloudPlatform/marketplace-testrunner/tests"

import (
	"bytes"
	"errors"

	"fmt"
	"os/exec"
	"syscall"

	"github.com/GoogleCloudPlatform/marketplace-testrunner/asserts"
	"github.com/GoogleCloudPlatform/marketplace-testrunner/specs"
)

type CommandExecutor interface {
	RunTest(test *specs.BashTest) (int, string, string, error)
}

type RealExecutor struct{}

func (e *RealExecutor) RunTest(test *specs.BashTest) (int, string, string, error) {
	cmd := exec.Command("bash", "-c", test.Script)
	err, stdout, stderr := e.executeProcess(cmd)
	var exitErr *exec.ExitError
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exitErr = e
		} else {
			return 0, "", "", errors.New(fmt.Sprintf("Process running error: %v", err))
		}
	}
	status, err := e.transformExitErrorToCode(exitErr)
	if err != nil {
		return 0, "", "", errors.New(fmt.Sprintf("Status code extraction error: %v", err))
	}
	return status, stdout, stderr, nil
}

func (e *RealExecutor) executeProcess(cmd *exec.Cmd) (error, string, string) {
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer
	err := cmd.Run()
	stdout := stdoutBuffer.String()
	stderr := stderrBuffer.String()
	return err, stdout, stderr
}

func (e *RealExecutor) transformExitErrorToCode(exitErr *exec.ExitError) (int, error) {
	if exitErr == nil {
		return 0, nil
	}
	if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
		return status.ExitStatus(), nil
	}
	return 0, errors.New("Cannot extract status code from process")
}

// RunBashTest executes a BashTest rule, returning an empty string if
// the test passes, otherwise the error message.
func RunBashTest(test *specs.BashTest, executor CommandExecutor) string {
	status, stdout, stderr, err := executor.RunTest(test)
	if err != nil {
		return asserts.MessageWithContext(fmt.Sprintf("%v", err), "Failed to execute bash script")
	}
	return validateResult(status, stdout, stderr, test)
}

func validateResult(status int, stdout, stderr string, test *specs.BashTest) string {
	if test.Expect == nil {
		return ""
	}
	if msg := validateErrorCode(status, test.Expect.ExitCode); msg != "" {
		return asserts.BashMessageWithContext(msg, "Unexpected exit status code", status, stdout, stderr)
	}
	if msg := validateBufferedOutput(stdout, test.Expect.Stdout); msg != "" {
		return asserts.BashMessageWithContext(msg, "Unexpected standard output stream", status, stdout, stderr)
	}
	if msg := validateBufferedOutput(stderr, test.Expect.Stderr); msg != "" {
		return asserts.BashMessageWithContext(msg, "Unexpected standard error stream", status, stdout, stderr)
	}
	return ""
}

func validateErrorCode(status int, expect *specs.IntAssert) string {
	if expect != nil {
		result := asserts.DoAssert(status, expect)
		if result != "" {
			return result
		}
	}
	return ""
}

func validateBufferedOutput(out string, expect *specs.StringAssert) string {
	if expect != nil {
		result := asserts.DoAssert(out, expect)
		if result != "" {
			return result
		}
	}
	return ""
}
