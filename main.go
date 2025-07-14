package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type state int

const (
	good state = iota
	bad
)

func main() {

	// Create a new random number generator instance
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	probes := []Probe{
		&TlsProbe{host: "1.1.1.1", port: 443, duration: 20 * time.Second},
		&TlsProbe{host: "9.9.9.9", port: 443, duration: 40 * time.Second},
	}

	var wg sync.WaitGroup
	wg.Add(len(probes))

	for _, probe := range probes {
		p := probe // capture range variable
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(p.Interval())
			defer ticker.Stop()
			log.Printf("[INIT] Starting probe: %s with interval %s", p.Name(), p.Interval())

			prevState := good
			var lastError error
			for {
				select {
				case <-ticker.C:
					// Add jitter to avoid synchronized checks
					jitter := time.Duration(rnd.Intn(500)) * time.Millisecond
					time.Sleep(jitter)

					stateStr := "bad"
					if prevState == good {
						stateStr = "good"
					}
					log.Printf("[RUN] Probe: %s (last state: %s)", p.Name(), stateStr)
					err := p.Run()
					if err != nil {
						if prevState == good {
							log.Printf("[FAIL] Probe %s failed: %v", p.Name(), err)
						} else {
							log.Printf("[FAIL] Probe %s still failing: %v (no additional notification)", p.Name(), err)
						}
						prevState = bad
						lastError = err
						continue
					}

					if prevState == bad {
						log.Printf("[RECOVER] Probe %s recovered from failure: %v", p.Name(), lastError)
						prevState = good
						lastError = nil
					}
				}
			}
		}()
	}

	// Wait for all goroutines to start (they never finish, so block forever)
	wg.Wait()
}
