package node

import (
	"reflect"
	"testing"
)

func TestNewProducer(t *testing.T) {
	type args struct {
		uri string
		key string
	}
	tests := []struct {
		name string
		args args
		want *Producer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProducer(tt.args.uri, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProducer_Send(t *testing.T) {
	type fields struct {
		exch string
		key  string
		uri  string
	}
	type args struct {
		body []byte
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
			p := &Producer{
				exch: tt.fields.exch,
				key:  tt.fields.key,
				uri:  tt.fields.uri,
			}
			if err := p.Send(tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Producer.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
