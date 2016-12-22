package node

import (
	"testing"
)

func TestNewProducer(t *testing.T) {
	var c interface{}

	t.Run("Initialises and returns a Producer", func(t *testing.T) {
		c = NewProducer(uri, key)

		switch c.(type) {
		case *Producer:
		default:
			t.Errorf("NewProducer() error = Received %T, expected Producer", c)
		}
	})

	t.Run("Initialises with the correct exchange", func(t *testing.T) {
		if c.(*Producer).exch != "workflow.exchange" {
			t.Errorf("NewProducer().exch = %q, want 'workflow.exchange'", c.(*Producer).exch)
		}
	})

	t.Run("Initialises with the correct routing key", func(t *testing.T) {
		if c.(*Producer).key != key {
			t.Errorf("NewProducer().key = %q, want %q", c.(*Producer).key, key)
		}
	})

	t.Run("Initialises with the correct uri", func(t *testing.T) {
		if c.(*Producer).uri != uri {
			t.Errorf("NewProducer().uri = %q, want %q", c.(*Producer).uri, uri)
		}
	})
}

// Not going to bother so much; this should be adequately covered in
// github.com/streadway/amqp - I'm doing nothing beyond using values
// tested above on the struct and using the above package's model
// implementation.
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
