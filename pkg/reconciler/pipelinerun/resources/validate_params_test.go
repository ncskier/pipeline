/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

import (
	"testing"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestValidateParamTypesMatching_Valid(t *testing.T) {
	tcs := []struct {
		name string
		p    *v1alpha1.Pipeline
		pr   *v1alpha1.PipelineRun
	}{{
		name: "proper param types",
		p: tb.Pipeline("a-pipeline", namespace, tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1alpha1.ParamTypeString),
			tb.PipelineParamSpec("correct-type-2", v1alpha1.ParamTypeString),
			tb.PipelineParamSpec("correct-type-3", v1alpha1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun", namespace, tb.PipelineRunSpec(
			"test-pipeline",
			tb.PipelineRunParam("correct-type-1", "somestring"),
			tb.PipelineRunParam("correct-type-2", "astring"),
			tb.PipelineRunParam("correct-type-3", "another", "array"))),
	}, {
		name: "empty param type",
		p: tb.Pipeline("a-pipeline", namespace, tb.PipelineSpec(
			tb.PipelineParamSpec("default-correct-type", ""),
		)),
		pr: tb.PipelineRun("a-pipelinerun", namespace, tb.PipelineRunSpec(
			"test-pipeline",
			tb.PipelineRunParam("default-correct-type", "astring"),
		)),
	}, {
		name: "no params to get wrong",
		p:    tb.Pipeline("a-pipeline", namespace),
		pr:   tb.PipelineRun("a-pipelinerun", namespace),
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateParamTypesMatching(tc.p, tc.pr); err != nil {
				t.Errorf("Pipeline.Validate() returned error: %v", err)
			}
		})
	}
}

func TestValidateParamTypesMatching_Invalid(t *testing.T) {
	tcs := []struct {
		name string
		p    *v1alpha1.Pipeline
		pr   *v1alpha1.PipelineRun
	}{{
		name: "string-array mismatch",
		p: tb.Pipeline("a-pipeline", namespace, tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1alpha1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1alpha1.ParamTypeString),
			tb.PipelineParamSpec("correct-type-2", v1alpha1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun", namespace,
			tb.PipelineRunSpec("test-pipeline",
				tb.PipelineRunParam("correct-type-1", "somestring"),
				tb.PipelineRunParam("mismatching-type", "an", "array"),
				tb.PipelineRunParam("correct-type-2", "another", "array"))),
	}, {
		name: "array-string mismatch",
		p: tb.Pipeline("a-pipeline", namespace, tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1alpha1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1alpha1.ParamTypeArray),
			tb.PipelineParamSpec("correct-type-2", v1alpha1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun", namespace,
			tb.PipelineRunSpec("test-pipeline",
				tb.PipelineRunParam("correct-type-1", "somestring"),
				tb.PipelineRunParam("mismatching-type", "astring"),
				tb.PipelineRunParam("correct-type-2", "another", "array"))),
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateParamTypesMatching(tc.p, tc.pr); err == nil {
				t.Errorf("Expected to see error when validating PipelineRun/Pipeline param types but saw none")
			}
		})
	}
}
