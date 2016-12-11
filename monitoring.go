package main

import (
	"os"
)

type Mon struct {
	Hostname            string
	EUID, GID, PID, UID int
	UTF8                string
}

func NewMon() (m Mon) {
	h, err := os.Hostname()
	if err != nil {
		m.Hostname = err.Error()
	} else {
		m.Hostname = h
	}

	m.UTF8 = "✔"

	m.EUID = os.Geteuid()
	m.GID = os.Getgid()
	m.PID = os.Getpid()
	m.UID = os.Getuid()

	return
}
