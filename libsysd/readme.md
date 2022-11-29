
# libsysd

A simple wrapper on top `go-systemd` module to poll or subscribe to systemd events for a list of systemd units / services.

---

## Interface

```go
type Watcher interface {

    // Poll is used for polling all the systemd metrics for a list of systemd services in every n interval.
    Poll(opts ...WatcherOpts)
    // Sub is used for making a subscription to a list of systemd services. This results in a subscription model, where if there is a change in the subscribed systemd service, then it will send to the buffer channel.
    Sub(opts ...WatcherOpts)
}


// Options available for configuration
WatcherOpts struct {
    HostNameMethod string
    PollTimeOut int
    MetricsBufferLimit int
}

// A systemdEvent structure 
SystemdEvent struct {
    Timestamp int64   // Timestamp of when did we receive the event
    PropertyUpdate map[string]interface{} // Property systemd property name:value/systemd property values map  

    UnitName       string                 // UnitName  
    Hostname       string
}
```

---

## Usage

`example/main.go`

---
