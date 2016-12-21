package node

import (
	"reflect"
	"testing"

	"github.com/streadway/amqp"
)

func TestNewConsumer(t *testing.T) {
	type args struct {
		uri string
		key string
	}
	tests := []struct {
		name string
		args args
		want *Consumer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConsumer(tt.args.uri, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsumer_Shutdown(t *testing.T) {
	type fields struct {
		channel *amqp.Channel
		conn    *amqp.Connection
		done    chan error
		exch    string
		key     string
		queue   string
		tag     string
		uri     string
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
			c := &Consumer{
				channel: tt.fields.channel,
				conn:    tt.fields.conn,
				done:    tt.fields.done,
				exch:    tt.fields.exch,
				key:     tt.fields.key,
				queue:   tt.fields.queue,
				tag:     tt.fields.tag,
				uri:     tt.fields.uri,
			}
			if err := c.Shutdown(); (err != nil) != tt.wantErr {
				t.Errorf("Consumer.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
