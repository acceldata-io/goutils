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

package main

import (
	"fmt"
	"log"

	"github.com/acceldata-io/goutils/libsysd"

	"github.com/integrii/flaggy"
)

var (
	sub               *flaggy.Subcommand
	poll              *flaggy.Subcommand
	watchList         []string
	metricBufferLimit int64 = 10000
	pollInterval      int64 = 10
	hostnameMethod          = "CMD"
)

func main() {
	flaggy.SetName("sysd")

	sub = flaggy.NewSubcommand("sub")
	sub.Description = "Subscribe to a list of systemd services, and wait for any events"
	sub.StringSlice(&watchList, "w", "watchlist", "systemd services")
	sub.String(&hostnameMethod, "n", "hostnamemethod", "Host name method to use")
	sub.Int64(&metricBufferLimit, "m", "metricbufferlimit", "Metric Buffer limit")
	sub.Int64(&pollInterval, "p", "interval", "Poll interval in seconds")
	flaggy.AttachSubcommand(sub, 1)

	poll = flaggy.NewSubcommand("poll")
	poll.Description = "Poll the metrics from systemd for a list of systemd services, in a given interval"
	poll.StringSlice(&watchList, "w", "watchlist", "systemd services")
	poll.String(&hostnameMethod, "n", "hostnamemethod", "Host name method to use")
	poll.Int64(&metricBufferLimit, "m", "metricbufferlimit", "Metric Buffer limit")
	poll.Int64(&pollInterval, "p", "interval", "Poll interval in seconds")
	flaggy.AttachSubcommand(poll, 1)
	flaggy.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case a := <-libsysd.EventsOut:
			fmt.Println(a)
		case err := <-libsysd.ErrCh:
			log.Fatal(err)
		}
	}
}

func run() error {
	switch {
	case sub.Used:
		SysWatcher := libsysd.New(watchList)
		SysWatcher.Sub(libsysd.WithMetricsBufferLimit(metricBufferLimit), libsysd.WithPollInterval(pollInterval))
	case poll.Used:
		SysWatcher := libsysd.New(watchList)
		SysWatcher.Poll(libsysd.WithMetricsBufferLimit(metricBufferLimit), libsysd.WithPollInterval(pollInterval))
	}
	return nil
}
