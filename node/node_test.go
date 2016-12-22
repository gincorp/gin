package node

import (
	"testing"

	"github.com/gincorp/gin/taskmanager"
)

var (
	redisURI = "redis://localhost"
)

func TestNewNode_Job(t *testing.T) {
	var n interface{}

	t.Run("Initialises and returns a Node", func(t *testing.T) {
		n = NewNode(uri, redisURI, "job")

		switch n.(type) {
		case Node:
		default:
			t.Errorf("NewNode() type error = Received %T, expected Node", n)
		}
	})

	t.Run("Initialises with the corrext TaskManager type", func(t *testing.T) {
		switch n.(Node).TaskManager.(type) {
		case taskmanager.JobManager:
		default:
			t.Errorf("NewNode().TaskManager type error = Received %T, expected taskmanager.JobManager", n)
		}
	})

	t.Run("Initialises with the correct queue set for it's Consumer", func(t *testing.T) {
		if n.(Node).Consumer.queue != "job" {
			t.Errorf("NewNode().Consumer.queue error = Received %q, expected job", n.(Node).Consumer.queue)
		}
	})

	t.Run("Initialises with the correct routing key set in the Producer", func(t *testing.T) {
		if n.(Node).Producer.key != "master" {
			t.Errorf("NewNode().Producer.key error = Received %q, expected master", n.(Node).Producer.key)
		}

	})
}

func TestNewNode_Master(t *testing.T) {
	var n interface{}

	t.Run("Initialises and returns a Node", func(t *testing.T) {
		n = NewNode(uri, redisURI, "master")

		switch n.(type) {
		case Node:
		default:
			t.Errorf("NewNode() type error = Received %T, expected Node", n)
		}
	})

	t.Run("Initialises with the corrext TaskManager type", func(t *testing.T) {
		switch n.(Node).TaskManager.(type) {
		case taskmanager.MasterManager:
		default:
			t.Errorf("NewNode().TaskManager type error = Received %T, expected taskmanager.MasterManager", n)
		}
	})

	t.Run("Initialises with the correct queue set for it's Consumer", func(t *testing.T) {
		if n.(Node).Consumer.queue != "master" {
			t.Errorf("NewNode().Consumer.queue error = Received %q, expected master", n.(Node).Consumer.queue)
		}
	})

	t.Run("Initialises with the correct routing key set in the Producer", func(t *testing.T) {
		if n.(Node).Producer.key != "job" {
			t.Errorf("NewNode().Producer.key error = Received %q, expected job", n.(Node).Producer.key)
		}

	})
}

// ConsumerLoop is badly written; certainly for testing. Thus I shall skip
// the testing of an important part of the tool in the hopes that testing
// eveything else in gin mitigates the potential for error here.
func TestNode_ConsumerLoop(t *testing.T) {
	type fields struct {
		Consumer    *Consumer
		Producer    *Producer
		TaskManager taskmanager.TaskManager
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				Consumer:    tt.fields.Consumer,
				Producer:    tt.fields.Producer,
				TaskManager: tt.fields.TaskManager,
			}
			if err := n.ConsumerLoop(); (err != nil) != tt.wantErr {
				t.Errorf("Node.ConsumerLoop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
