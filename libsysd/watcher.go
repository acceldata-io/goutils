package libsysd

import (
	"strings"
)

// Watcher implements a watch mechanism with poll and sub functions
type Watcher interface {
	// Poll method will poll for the systemd properties at a certain interval
	Poll(opts ...WatcherOps)
	// Sub method is an event based method.
	// If and only if there is an event occurred in the systemd managed services these will be captured by this method
	Sub(opts ...WatcherOps)
}

// TODO : Serializers to convert to certain output formats such as JSON, LineProtocol
// TODO : Add more configuration to make the lib more stable
type watcher struct {
	watchList          []string
	systemD            Adapter
	metricsBufferLimit int64
	pollInterval       int64
	hostnameMethod     string
}

var (
	// EventsOut is a channel where systemd events are being pushed
	EventsOut = make(chan *SystemDEvent)
	// eventsIn  = make(chan *SystemDEvent)

	// ErrCh is a channel where any errors are pushed
	ErrCh = make(chan error)
)

// WatcherOps sets optional parameters to a watcher
type WatcherOps func(*watcher)

// WithPollInterval sets poll interval
func WithPollInterval(interval int64) WatcherOps {
	return func(w *watcher) {
		w.pollInterval = interval
	}
}

// WithMetricsBufferLimit set buffer limit for the number of systemd events
func WithMetricsBufferLimit(limit int64) WatcherOps {
	return func(w *watcher) {
		w.metricsBufferLimit = limit
	}
}

// WithHostNameMethod sets the hostname method used to get the machine hostname
// Valid methods are "RFQDN", "FQDN", "OS" and "CMD"
// Uses: github.com/acceldata-io/goutils/netutils
func WithHostNameMethod(method string) WatcherOps {
	return func(w *watcher) {
		w.hostnameMethod = method
	}
}

// New returns a new watcher
func New(watcherList []string, opts ...WatcherOps) Watcher {
	sys := NewSystemDAdapter()
	w := &watcher{
		watchList: convertUnitType(watcherList),
		systemD:   sys,
	}
	for _, opt := range opts {
		opt(w)
	}
	EventsOut = make(chan *SystemDEvent)
	return w
}

func (w *watcher) Sub(opts ...WatcherOps) {
	for _, opt := range opts {
		opt(w)
	}
	EventsOut = make(chan *SystemDEvent)
	go w.sub()
}

func (w *watcher) Poll(opts ...WatcherOps) {
	for _, opt := range opts {
		opt(w)
	}
	EventsOut = make(chan *SystemDEvent)
	go w.poll()
}

func convertUnitType(unitList []string) []string {
	properUnitName := []string{}
	for _, u := range unitList {
		if !stringInSlice(u, properUnitName) {
			if !strings.ContainsRune(u, '.') {
				properUnitName = append(properUnitName, u+".service")
			}
		}
	}
	return properUnitName
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
