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
	"fmt"
	"time"

	"github.com/acceldata-io/goutils/netutils"
)

func (w *watcher) poll() {
	if len(w.watchList) < 1 {
		ErrCh <- fmt.Errorf("no systemd services were provided")
		return
	}
	for _, unit := range w.watchList {
		unitList, err := w.systemD.ListUnitsByPattern(states, []string{unit})
		if err != nil {
			ErrCh <- err
			return
		}
		if len(unitList) < 1 {
			ErrCh <- fmt.Errorf("%s unit listed cannot be found", unit)
			return
		}
	}
	pollTicker := time.NewTicker(time.Duration(w.pollInterval) * time.Second)
	for ; true; <-pollTicker.C {
		for _, unit := range w.watchList {
			event, err := w.systemD.GetPropertiesForUnit(unit)
			if err != nil {
				ErrCh <- err
			}
			hostName, err := netutils.GetHostName(w.hostnameMethod, 20)
			if err != nil {
				hostName = "localhost"
			}
			e := &SystemDEvent{
				Timestamp:      time.Now().UnixMilli(),
				PropertyUpdate: event,
				UnitName:       unit,
				Hostname:       hostName,
			}
			EventsOut <- e
		}
	}
}
