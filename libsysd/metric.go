// Package libsysd is a simple wrapper on top of go-systemd module
package libsysd

// SystemDEvent represents a single systemd service event
type SystemDEvent struct {
	Timestamp      int64                  // Timestamp of when did we receive the event
	PropertyUpdate map[string]interface{} // Property systemd property name:value/systemd property values map
	UnitName       string                 // UnitName of the systemd service
	Hostname       string                 // Hostname of the current machine
}
