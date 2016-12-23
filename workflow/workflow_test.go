package workflow

import (
	"reflect"
	"testing"
)

func TestParseWorkflow(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantW   Workflow
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, err := ParseWorkflow(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("ParseWorkflow() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
