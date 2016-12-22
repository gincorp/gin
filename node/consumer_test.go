package node

import (
	"testing"
)

var (
	uri = "amqp://guest:guest@localhost/vhost"
	key = "testkey"
)

func TestNewConsumer(t *testing.T) {
	var c interface{}

	t.Run("Initialises and returns a Consumer", func(t *testing.T) {
		c = NewConsumer(uri, key)

		switch c.(type) {
		case *Consumer:
		default:
			t.Errorf("NewConsumer() error = Received %T, expected Consumer", c)
		}
	})
}

func TestNewConsumer_Exch(t *testing.T) {
	t.Run("Initialises with the correct exchange", func(t *testing.T) {
		c := NewConsumer(uri, key)

		if c.exch != "workflow.exchange" {
			t.Errorf("NewConsumer().exch = %q, want 'workflow.exchange'", c.exch)
		}
	})
}

func TestNewConsumer_Key(t *testing.T) {
	t.Run("Initialises with the correct routing key", func(t *testing.T) {
		c := NewConsumer(uri, key)

		if c.key != key {
			t.Errorf("NewConsumer().key = %q, want %q", c.key, key)
		}
	})
}

func TestNewConsumer_Queue(t *testing.T) {
	t.Run("Initialises with the correct queue name", func(t *testing.T) {
		c := NewConsumer(uri, key)

		if c.queue != key {
			t.Errorf("NewConsumer().queue = %q, want %q", c.queue, key)
		}
	})
}

func TestNewConsumer_Tag(t *testing.T) {
	t.Run("Initialises with the correct tag name", func(t *testing.T) {
		c := NewConsumer(uri, key)

		if c.tag != key {
			t.Errorf("NewConsumer().tag = %q, want %q", c.tag, key)
		}
	})
}

func TestNewConsumer_URI(t *testing.T) {
	t.Run("Initialises with the correct uri", func(t *testing.T) {
		c := NewConsumer(uri, key)

		if c.uri != uri {
			t.Errorf("NewConsumer().uri = %q, want %q", c.uri, uri)
		}
	})
}
