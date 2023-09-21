package main

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/usama28232/fusion/fuse"
)

// a sample struct to implement send function
type SampleRequestSender struct {
	URL    string
	status bool // a casual field to toggle return type
}

func (s SampleRequestSender) Send() bool {
	return s.status
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestBreakerSendBestCase(t *testing.T) {
	s := SampleRequestSender{
		URL:    "address",
		status: true,
	}
	b := fuse.Breaker{}
	// prints detailed information in console
	b = *b.Debug()

	for i := 0; i < 5; i++ {
		b.Send(s)
	}
}

func TestBreakerSendDeadCase(t *testing.T) {
	s := SampleRequestSender{
		URL:    "address",
		status: false,
	}
	b := fuse.Breaker{}
	b = *b.Debug()

	for i := 0; i < 5; i++ {
		b.Send(s)
	}
}

func TestBreakerResetState(t *testing.T) {
	var wg sync.WaitGroup
	b := fuse.Breaker{}
	b = *b.Debug()
	wg.Add(1)
	resetChannel := b.Reset()

	go func() {
		defer wg.Done()
		for f := range resetChannel {
			fmt.Println(f)
		}
	}()
	wg.Wait()
}
