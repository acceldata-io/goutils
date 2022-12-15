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

package libsysd

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gobwas/glob"
)

var states = []string{"active", "activating", "failed", "inactive", "deactivating", "maintenance", "reloading"}

const versionProperty = "Version"

// Adapter implements a systemd adapter
type Adapter interface {
	ListUnitsByPattern(states, patterns []string) ([]dbus.UnitStatus, error)
	GetPropertiesForUnit(unit string) (map[string]interface{}, error)
	GetPropertiesForAUnitType(unit, unitType string) (map[string]interface{}, error)
	GetPropertyForService(unitName, propertyName string) (*dbus.Property, error)
	RestartService(serviceName string) (*dbus.UnitStatus, error)
	StartService(serviceName string) error
	StopService(serviceName string) error
	ReloadService(serviceName string) error
	SubscribeToUnitProperties(sysEventCh chan *dbus.PropertiesUpdate, errCh chan error) error
	GetVersion() (int, error)
	ReloadDaemon() error
	Close()
}

type systemDAdapter struct {
	conn           *dbus.Conn
	systemDVersion int
	mutex          *sync.Mutex
}

func (s *systemDAdapter) Close() {
	if s.conn != nil {
		s.mutex.Lock()
		s.conn.Close()
		s.mutex.Unlock()
	}
}

// NewSystemDAdapter provides a new systemd adapter
func NewSystemDAdapter() Adapter {
	return &systemDAdapter{
		conn:           nil,
		systemDVersion: 0,
		mutex:          &sync.Mutex{},
	}
}

var reVersion = regexp.MustCompile(`\d\d\d`)

func (s *systemDAdapter) GetPropertiesForUnit(unit string) (map[string]interface{}, error) {
	err := s.getConnection()
	if err != nil {
		return nil, err
	}
	return s.conn.GetAllPropertiesContext(context.Background(), unit)
}

func (s *systemDAdapter) GetPropertiesForAUnitType(unit, unitType string) (map[string]interface{}, error) {
	err := s.getConnection()
	if err != nil {
		return nil, err
	}
	return s.conn.GetUnitTypePropertiesContext(context.Background(), unit, unitType)
}

func (s *systemDAdapter) GetPropertyForService(unitName, propertyName string) (*dbus.Property, error) {
	err := s.getConnection()
	if err != nil {
		return nil, err
	}
	return s.conn.GetServicePropertyContext(context.Background(), unitName, propertyName)
}

func (s *systemDAdapter) SubscribeToUnitProperties(sysEvent chan *dbus.PropertiesUpdate, errCh chan error) error {
	err := s.getConnection()
	if err != nil {
		return err
	}
	s.conn.SetPropertiesSubscriber(sysEvent, errCh)
	return nil
}

func (s *systemDAdapter) ListUnitsByPattern(states, patterns []string) ([]dbus.UnitStatus, error) {
	err := s.getConnection()
	if err != nil {
		return nil, err
	}
	version, err := s.getVersion()
	if err != nil {
		return nil, err
	}
	if version >= 230 {
		return s.conn.ListUnitsByPatternsContext(context.Background(), states, patterns)
	}
	return s.listUnitsAndFilterPatterns(states, patterns)
}

func (s *systemDAdapter) getConnection() error {
	if s.conn == nil {
		s.mutex.Lock()
		var err error
		s.conn, err = dbus.NewSystemdConnectionContext(context.Background())
		s.mutex.Unlock()
		return err
	}
	return nil
}

func (s *systemDAdapter) listUnitsAndFilterPatterns(states, patterns []string) ([]dbus.UnitStatus, error) {
	units, err := s.conn.ListUnitsContext(context.Background())
	if err != nil {
		return nil, err
	}

	compiledStates := getCompiledMapGlob(states)
	compiledPatterns := getCompiledMapGlob(patterns)

	matchedUnits := make([]dbus.UnitStatus, 0)

	for _, unit := range units {
		stateMatched := false
		patternMatched := false
		for _, compiledState := range compiledStates {
			if compiledState.Match(unit.ActiveState) || compiledState.Match(unit.SubState) {
				stateMatched = true
			}
		}
		for _, compiledPattern := range compiledPatterns {
			if compiledPattern.Match(unit.Name) {
				patternMatched = true
			}
		}
		if stateMatched && patternMatched {
			matchedUnits = append(matchedUnits, unit)
		}
	}
	return matchedUnits, nil
}

func getCompiledMapGlob(arr []string) map[string]glob.Glob {
	compiledMap := make(map[string]glob.Glob)
	for _, str := range arr {
		compiledMap[str] = glob.MustCompile(str)
	}
	return compiledMap
}

func (s *systemDAdapter) getVersion() (int, error) {
	if s.systemDVersion != 0 {
		return s.systemDVersion, nil
	}
	s.mutex.Lock()
	version, err := s.conn.GetManagerProperty(versionProperty)
	if err != nil {
		s.mutex.Unlock()
		return 0, err
	}

	major := reVersion.FindString(version)
	if major == "" {
		s.mutex.Unlock()
		return 0, fmt.Errorf("couldn't parse systemd version string '%s'", version)
	}

	ver, err := strconv.Atoi(major)
	if err != nil {
		s.mutex.Unlock()
		return 0, fmt.Errorf("couldn't parse systemd version string '%s': %v", version, err)
	}
	s.systemDVersion = ver
	s.mutex.Unlock()
	return s.systemDVersion, nil
}

func (s *systemDAdapter) GetVersion() (int, error) {
	err := s.getConnection()
	if err != nil {
		return -1, err
	}

	return s.getVersion()
}

func (s *systemDAdapter) RestartService(serviceName string) (*dbus.UnitStatus, error) {
	err := s.getConnection()
	if err != nil {
		return nil, err
	}

	wait := make(chan string)
	_, err = s.conn.RestartUnitContext(context.Background(), serviceName, "replace", wait)
	if err != nil {
		return nil, err
	}
	<-wait
	dbusStatus, err := s.ListUnitsByPattern(states, []string{serviceName})
	return &dbusStatus[0], err
}

func (s *systemDAdapter) ReloadDaemon() error {
	err := s.getConnection()
	if err != nil {
		return err
	}

	err = s.conn.ReloadContext(context.Background())
	return err
}

func (s *systemDAdapter) StartService(serviceName string) error {
	err := s.getConnection()
	if err != nil {
		return err
	}

	wait := make(chan string)
	_, err = s.conn.StartUnitContext(context.Background(), serviceName, "replace", wait)
	if err != nil {
		return err
	}
	<-wait
	return nil
}

func (s *systemDAdapter) StopService(serviceName string) error {
	err := s.getConnection()
	if err != nil {
		return err
	}

	wait := make(chan string)
	_, err = s.conn.StopUnitContext(context.Background(), serviceName, "replace", wait)
	if err != nil {
		return err
	}
	<-wait
	return nil
}

func (s *systemDAdapter) ReloadService(serviceName string) error {
	err := s.getConnection()
	if err != nil {
		return err
	}

	wait := make(chan string)
	_, err = s.conn.ReloadUnitContext(context.Background(), serviceName, "replace", wait)
	if err != nil {
		return err
	}
	<-wait
	return nil
}
