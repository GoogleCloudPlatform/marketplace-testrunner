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

package conditions // import "github.com/GoogleCloudPlatform/marketplace-testrunner/conditions"

import (
	"github.com/GoogleCloudPlatform/marketplace-testrunner/asserts"
	"github.com/GoogleCloudPlatform/marketplace-testrunner/specs"
)

type TestStatus interface {
	FailuresSoFarCount() int
}

// Evaluate evaluates if a condition is true for the current test status.
// If the condition evaluates to false, the second return contains the reason.
func Evaluate(condition *specs.Condition, testStatus TestStatus) (bool, string) {
	if condition.FailuresSoFar != nil {
		if msg := asserts.DoAssert(testStatus.FailuresSoFarCount(), condition.FailuresSoFar); msg != "" {
			return false, asserts.MessageWithContext(msg, "Count of failures so far")
		}
	}
	return true, ""
}
