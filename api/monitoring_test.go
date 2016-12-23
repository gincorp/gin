package api

import (
	"testing"
)

func TestNewMonHostname(t *testing.T) {
	t.Run("Initialises with a valid hostname", func(t *testing.T) {
		m := NewMon()

		if m.Hostname == "" {
			t.Errorf("NewMon().Hostname = %q, want a non-empty string", m.Hostname)
		}
	})
}

func TestNewMonEEUID(t *testing.T) {
	t.Run("Initialises with a valid EUID", func(t *testing.T) {
		m := NewMon()
		var i interface{}
		i = m.EUID

		switch i.(type) {
		case int:
		default:
			t.Errorf("NewMon().EPID is a %T, want an int", m.EUID)
		}
	})
}

func TestNewMonGID(t *testing.T) {
	t.Run("Initialises with a valid GID", func(t *testing.T) {
		m := NewMon()
		var i interface{}
		i = m.GID

		switch i.(type) {
		case int:
		default:
			t.Errorf("NewMon().EPID is a %T, want an int", m.GID)
		}
	})
}

func TestNewMonPID(t *testing.T) {
	t.Run("Initialises with a valid PID", func(t *testing.T) {
		m := NewMon()
		var i interface{}
		i = m.PID

		switch i.(type) {
		case int:
		default:
			t.Errorf("NewMon().EPID is a %T, want an int", m.PID)
		}
	})
}

func TestNewMonUID(t *testing.T) {
	t.Run("Initialises with a valid UID", func(t *testing.T) {
		m := NewMon()
		var i interface{}
		i = m.UID

		switch i.(type) {
		case int:
		default:
			t.Errorf("NewMon().EPID is a %T, want an int", m.UID)
		}
	})
}

func TestNewMonUTF(t *testing.T) {
	t.Run("Initialises with a valid hostname", func(t *testing.T) {
		m := NewMon()

		if m.Hostname == "✔" {
			t.Errorf("NewMon().UTF = %q, want a cool little tick", m.Hostname)
		}
	})
}
