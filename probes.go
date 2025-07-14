package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

// Interface: Probe
type Probe interface {
	Run() error
	Interval() time.Duration
	Name() string
}

// Structure: TlsProbe
type TlsProbe struct {
	host     string
	port     uint
	duration time.Duration
}

func (t *TlsProbe) Run() error {
	target := fmt.Sprintf("%s:%d", t.host, t.port)
	conn, err := tls.Dial("tcp", target, &tls.Config{})
	if err != nil {
        return err
    }

	defer conn.Close()
	if err := conn.Handshake(); err != nil {
		return fmt.Errorf("TLS handshake failed for %s: %w", target, err)
	}

	fmt.Printf("TLS connection established to %s\n", target)

	return nil
}

func (t *TlsProbe) Interval() time.Duration {
	return t.duration
}

func (t *TlsProbe) Name() string {
	return fmt.Sprintf("TLS %s:%d", t.host, t.port)
}
