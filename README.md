# fusion

[![GoDoc](https://godoc.org/pault.ag/go/piv?status.svg)](https://pkg.go.dev/github.com/usama28232/fusion) 

A lightweight Circuit Breaker pattern implementation in GO to handle failures and prevent an application from making repeated requests to a service or resource that is likely to fail.

## Installation

Import the package using the go get command

```
go get github.com/usama28232/fusion
```

## Usage

In order to utilize `fusion` *ciruit-breaker* we require a struct that implements send interface

**Note:** Breaker performs internal operations with `sync.Mutex`, making it thread-safe

```
// a sample struct to implement send function
type SampleRequestSender struct {
	URL    string
	status bool // a casual field to toggle return type
}

func (s SampleRequestSender) Send() bool {
	return s.status
}
```

Note: This struct can have any number of fields, we are just concerned with Send function which returns a boolean

```
// sender interface exposes a send method
//
// Breaker uses this pattern to trigger send call - see test example
type Sender interface {
	Send() bool
}
```

Observe the following test cases in `main_test.go`

```
func TestBreakerSendBestCase(t *testing.T) {
	s := SampleRequestSender{
		URL:    "address",
		status: true, // this is what our send function will be returning
	}
	b := fuse.Breaker{}
	// prints detailed information in console
	b = *b.Debug()

	for i := 0; i < 5; i++ {
		b.Send(s)
	}
    ...
```

For the failure case:

```
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
    ...
```

To subscribe to Reset State:

```
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
    ...
```

You can also use following methods to provide custom `max-retries` and `reset-timeouts` by:

```
...
// sets breaker max failure count
func (b *Breaker) SetMaxFailureCount(value int) {
	b.maxFailureCount = value
}

// sets breaker reset timeout
func (b *Breaker) SetResetTimeout(value int) {
	b.resetTimeout = value
}
...
```

## Roadmap

* `Peeks` (a mechanism to send test requests to check availability) of a service will be added in later version.
* `Half-Closed` states will be added in later versions.



### Feel free to edit/expand/explore this repository

For feedback and queries, reach me on LinkedIn [here](https://www.linkedin.com/in/usama28232/?original_referer=)
