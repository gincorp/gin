package workflow

import (
	"fmt"
	"testing"
)

func TestNewRunner(t *testing.T) {
	uuid := "some-uuid"

	varsWF := Workflow{"Some workflow", []Step{Step{Name: "a-step", Type: "none", Context: map[string]string{}}}, map[string]string{"a": "b"}}
	noVarsWF := Workflow{"Some workflow", []Step{Step{Name: "a-step", Type: "none", Context: map[string]string{}}}, map[string]string{}}

	type args struct {
		uuid string
		wf   Workflow
	}
	tests := []struct {
		name    string
		args    args
		hasVars bool
	}{
		{"creates a runner with vars", args{uuid, varsWF}, true},
		{"creates a runner with no vars", args{uuid, noVarsWF}, false},
	}
	for _, tt := range tests {
		var r interface{}
		r = NewRunner(tt.args.uuid, tt.args.wf)

		t.Run(tt.name, func(t *testing.T) {
			switch r.(type) {
			case Runner:
			default:
				t.Errorf("NewRunner() error: expected Runner, got %T", r)
			}
		})

		t.Run(fmt.Sprintf("%s - UUID", tt.name), func(t *testing.T) {
			if r.(Runner).UUID != uuid {
				t.Errorf("NewRunner() UUID error: expected %q, got %q", uuid, r.(Runner).UUID)
			}
		})

		t.Run(fmt.Sprintf("%s - Variables", tt.name), func(t *testing.T) {
			if _, ok := r.(Runner).Variables["Defaults"]; !ok {
				t.Errorf("NewRunner() Variables error: expected key Defaults, got %V", r.(Runner).Variables)
			}
		})

		if tt.hasVars {
			t.Run(fmt.Sprintf("%s - vars", tt.name), func(t *testing.T) {
				if _, ok := r.(Runner).Variables["Defaults"].(map[string]string)["a"]; !ok {
					t.Errorf("NewRunner() Default -> Variables error: expected key 'a', got %V", r.(Runner).Variables["Defaults"])
				}
			})
		}
	}
}

func TestParseRunner(t *testing.T) {
	wfrJSON := `
{
  "UUID": "some-uuid",
  "Last": "a step",
  "Variables": {
    "Defaults": {
      "foo": "bar"
    },
    "Something": {
      "baz": "quux"
    }
  },
  "Workflow": {
    "Name": "A Workflow",
    "Steps": [
      {
        "Name": "a step",
        "Type": "none"
      }
    ]
  }
}
`
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple runner", args{wfrJSON}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRunner(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRunner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRunner_Next(t *testing.T) {
	wf := Workflow{"Some workflow", []Step{Step{Name: "a step", Type: "none"}, Step{Name: "final step", Type: "none"}}, map[string]string{}}

	type fields struct {
		Last     string
		Workflow Workflow
	}
	tests := []struct {
		name     string
		fields   fields
		stepName string
		wantDone bool
	}{
		{"workflow with non-final step", fields{"a step", wf}, "final step", false},
		{"workflow with final step", fields{"final step", wf}, "final step", true},
	}
	for _, tt := range tests {
		wfr := &Runner{
			Last:     tt.fields.Last,
			Workflow: tt.fields.Workflow,
		}
		gotStep, gotDone := wfr.Next()

		t.Run(fmt.Sprintf("%s - done status", tt.name), func(t *testing.T) {
			if gotDone != tt.wantDone {
				t.Errorf("Runner.Next() gotDone = %v, want %v", gotDone, tt.wantDone)
			}
		})

		t.Run(fmt.Sprintf("%s - next step name", tt.name), func(t *testing.T) {
			if gotStep.Name != tt.stepName {
				t.Errorf("Runner.Next() gotStep.Name = %v, want %v", gotStep.Name, tt.stepName)
			}
		})

	}
}

func TestRunner_Current(t *testing.T) {
	wf := Workflow{"Some workflow", []Step{Step{Name: "a step", Type: "none"}, Step{Name: "final step", Type: "none"}}, map[string]string{}}

	type fields struct {
		Last     string
		Workflow Workflow
	}
	tests := []struct {
		name     string
		fields   fields
		stepName string
	}{
		{"No set step", fields{"", wf}, "a step"},
		{"First step", fields{"a step", wf}, "a step"},
		{"Last step", fields{"final step", wf}, "final step"},
	}
	for _, tt := range tests {
		wfr := &Runner{
			Last:     tt.fields.Last,
			Workflow: tt.fields.Workflow,
		}
		_, s := wfr.Current()

		t.Run(tt.name, func(t *testing.T) {
			if s.Name != tt.stepName {
				t.Errorf("Runner.Current() Step name = %v, want %v", s.Name, tt.stepName)
			}
		})
	}
}

func TestRunner_End(t *testing.T) {
	wfr := &Runner{}
	t.Run("set end state", func(t *testing.T) {
		if wfr.State == "ended" {
			t.Errorf("Runner.End() error: step already set to ended, %V", wfr)
		}

		wfr.End()

		if wfr.State != "ended" {
			t.Errorf("Runner.End() state change did not persist to ended, %V", wfr)
		}
	})
}

func TestRunner_Fail(t *testing.T) {
	wfr := &Runner{}
	failMsg := "a failure occurred"

	t.Run("set fail state", func(t *testing.T) {
		if wfr.State == "failed" {
			t.Errorf("Runner.Fail() error: step already set to failed, %V", wfr)
		}

		wfr.Fail(failMsg)

		if wfr.State != "failed" {
			t.Errorf("Runner.Fail() state change did not persist to failed, %V", wfr)
		}

		if wfr.ErrorMessage != failMsg {
			t.Errorf("Runner.Fail() state change did not persist error message, %V", wfr)
		}
	})
}

func TestRunner_Start(t *testing.T) {
	wfr := &Runner{}
	t.Run("set start state", func(t *testing.T) {
		if wfr.State == "started" {
			t.Errorf("Runner.Start() error: step already set to started, %V", wfr)
		}

		wfr.Start()

		if wfr.State != "started" {
			t.Errorf("Runner.Start() state change did not persist to started, %V", wfr)
		}
	})
}
