// Acceldata Inc. and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// 	Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package libsysd is a simple wrapper on top of go-systemd module
package libsysd

// SystemDEvent represents a single systemd service event
type SystemDEvent struct {
	Timestamp      int64                  // Timestamp of when did we receive the event
	PropertyUpdate map[string]interface{} // Property systemd property name:value/systemd property values map
	UnitName       string                 // UnitName of the systemd service
	Hostname       string                 // Hostname of the current machine
}
