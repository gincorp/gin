package workflow

import (
	"fmt"
	"testing"
)

func cleanContext() map[string]string {
	return map[string]string{"simple_string": "abcd", "templated": "{{.Hello}}"}
}

func TestStep_Compile(t *testing.T) {
	v1 := map[string]interface{}{"Hello": "World!"}
	v2 := map[string]interface{}{}

	tests := []struct {
		name               string
		context            map[string]string
		vars               map[string]interface{}
		wantErr            bool
		wantTemplatedValue string
	}{
		{"compile simple template", cleanContext(), v1, false, "World!"},
		{"fail to compile simple template (no error)", cleanContext(), v2, false, "<no value>"},
	}
	for _, tt := range tests {
		s := Step{Context: tt.context}
		err := s.Compile(tt.vars)

		t.Run(tt.name, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("Step.Compile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		t.Run(fmt.Sprintf("%s value", tt.name), func(t *testing.T) {
			if s.Context["templated"] != tt.wantTemplatedValue {
				t.Errorf("Step.Compile()[templated] = %q, expected %q", s.Context["templated"], tt.wantTemplatedValue)
			}
		})
	}
}

func TestStep_SetStatus(t *testing.T) {
	tests := []struct {
		name  string
		error bool
	}{
		{"set step status on success", false},
		{"set step status on failure", true},
	}
	for _, tt := range tests {
		data := map[string]interface{}{
			"Duration":     "1 ms",
			"Start":        "some-timestamp",
			"End":          "another-timestamp",
			"Failed":       tt.error,
			"ErrorMessage": "an error occurred",
		}
		s := Step{}
		s.SetStatus(data)

		t.Run(fmt.Sprintf("%s - Duration", tt.name), func(t *testing.T) {
			if s.Duration != data["Duration"] {
				t.Errorf("Step.SetStatus().Duration = %q, expected %q", s.Duration, data["Duration"].(string))
			}
		})

		t.Run(fmt.Sprintf("%s - Start", tt.name), func(t *testing.T) {
			if s.Start != data["Start"] {
				t.Errorf("Step.SetStatus().Start = %q, expected %q", s.Start, data["Start"].(string))
			}
		})

		t.Run(fmt.Sprintf("%s - End", tt.name), func(t *testing.T) {
			if s.End != data["End"] {
				t.Errorf("Step.SetStatus().End = %q, expected %q", s.End, data["End"].(string))
			}
		})

		t.Run(fmt.Sprintf("%s - Failed", tt.name), func(t *testing.T) {
			if s.Failed != data["Failed"] {
				t.Errorf("Step.SetStatus().Failed = %v, expected %v", s.Duration, data["Failed"].(bool))
			}
		})

		if data["Failed"].(bool) {
			t.Run(fmt.Sprintf("%s - ErrorMessage", tt.name), func(t *testing.T) {
				if s.ErrorMessage != data["ErrorMessage"] {
					t.Errorf("Step.SetStatus().ErrorMessage = %q, expected %q", s.Duration, data["ErrorMessage"].(string))
				}
			})
		}
	}
}
