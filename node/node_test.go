package node

import (
	"reflect"
	"testing"

	"github.com/gincorp/gin/taskmanager"
	"github.com/streadway/amqp"
)

func TestNewNode(t *testing.T) {
	type args struct {
		uri      string
		redisURI string
		nodeMode string
	}
	tests := []struct {
		name  string
		args  args
		wantN Node
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := NewNode(tt.args.uri, tt.args.redisURI, tt.args.nodeMode); !reflect.DeepEqual(gotN, tt.wantN) {
				t.Errorf("NewNode() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

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

func TestNode_Consume(t *testing.T) {
	type fields struct {
		Consumer    *Consumer
		Producer    *Producer
		TaskManager taskmanager.TaskManager
	}
	type args struct {
		deliveries <-chan amqp.Delivery
		done       chan error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
			n.Consume(tt.args.deliveries, tt.args.done)
		})
	}
}

func TestNode_Deliver(t *testing.T) {
	type fields struct {
		Consumer    *Consumer
		Producer    *Producer
		TaskManager taskmanager.TaskManager
	}
	type args struct {
		message interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := n.Deliver(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Node.Deliver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
