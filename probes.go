package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
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

// Structure: DnsProbe
type DnsProbe struct {
	host     string
	resolver string // Optional resolver address
	duration time.Duration
}

func (d *DnsProbe) Run() error {
	if d.resolver != "" {
		// Use custom resolver if specified
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout: time.Millisecond * 3000,
				}
				return dialer.DialContext(ctx, network, d.resolver+":53")
			},
		}
		net.DefaultResolver = resolver
	}

	ips, err := net.LookupIP(d.host)
	if err != nil {
		return fmt.Errorf("DNS lookup failed for %s: %w", d.host, err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("no IP addresses found for %s", d.host)
	}

	fmt.Printf("DNS lookup successful for %s: %v\n", d.host, ips)
	return nil
}

func (d *DnsProbe) Interval() time.Duration {
	return d.duration
}

func (d *DnsProbe) Name() string {
	return fmt.Sprintf("DNS %s", d.host)
}
